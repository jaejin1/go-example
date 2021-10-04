package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	goldap "github.com/go-ldap/ldap/v3"
	"net"
	"net/url"
	"strings"
	"time"
)

var ErrLDAPServerTimeout = errors.New("ldap server network timeout")
var ErrLDAPPingFail = errors.New("fail to ping LDAP server")
var ErrEmptySearchDN = errors.New("empty search dn")
var ErrInvalidCredential = errors.New("invalid credential")
var ErrInvalidFilter = errors.New("invalid filter syntax")

type Session struct {
	basicCfg LdapConf
	groupCfg GroupConf
	ldapConn *goldap.Conn
}

func NewSession(basicCfg LdapConf, groupCfg GroupConf) *Session {
	return &Session{
		basicCfg: basicCfg,
		groupCfg: groupCfg,
	}
}

func (s *Session) Bind(dn string, password string) error {
	return s.ldapConn.Bind(dn, password)
}

func (s *Session) Open() error {
	ldapURL, err := formatURL(s.basicCfg.URL)
	if err != nil {
		return err
	}
	splitLdapURL := strings.Split(ldapURL, "://")

	protocol, hostport := splitLdapURL[0], splitLdapURL[1]
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return err
	}
	connectionTimeout := s.basicCfg.ConnectionTimeout
	goldap.DefaultTimeout = time.Duration(connectionTimeout) * time.Second

	switch protocol {
	case "ldap":
		ldap, err := goldap.Dial("tcp", hostport)
		if err != nil {
			return err
		}
		s.ldapConn = ldap
	case "ldaps":
		fmt.Println("Start to dial ldaps")
		ldap, err := goldap.DialTLS("tcp", hostport, &tls.Config{ServerName: host, InsecureSkipVerify: !s.basicCfg.VerifyCert})
		if err != nil {
			return err
		}
		s.ldapConn = ldap
	}

	return nil

}

func (s *Session) Close() {
	if s.ldapConn != nil {
		s.ldapConn.Close()
	}
}

func formatURL(ldapURL string) (string, error) {

	var protocol, hostport string
	_, err := url.Parse(ldapURL)
	if err != nil {
		return "", fmt.Errorf("parse Ldap Host ERR: %s", err)
	}

	if strings.Contains(ldapURL, "://") {
		splitLdapURL := strings.Split(ldapURL, "://")
		protocol, hostport = splitLdapURL[0], splitLdapURL[1]
		if !((protocol == "ldap") || (protocol == "ldaps")) {
			return "", fmt.Errorf("unknown ldap protocol")
		}
	} else {
		hostport = ldapURL
		protocol = "ldap"
	}

	if strings.Contains(hostport, ":") {
		_, port, err := net.SplitHostPort(hostport)
		if err != nil {
			return "", fmt.Errorf("illegal ldap url, error: %v", err)
		}
		if port == "636" {
			protocol = "ldaps"
		}

	} else {
		switch protocol {
		case "ldap":
			hostport = hostport + ":389"
		case "ldaps":
			hostport = hostport + ":636"
		}
	}

	fLdapURL := protocol + "://" + hostport

	return fLdapURL, nil

}

func TestConfig(ldapConfig LdapConf) (bool, error) {
	ts := NewSession(ldapConfig, GroupConf{})
	if err := ts.Open(); err != nil {
		if goldap.IsErrorWithCode(err, goldap.ErrorNetwork) {
			return false, ErrLDAPServerTimeout
		}
		return false, ErrLDAPPingFail
	}
	defer ts.Close()

	if ts.basicCfg.SearchDn == "" {
		return false, ErrEmptySearchDN
	}
	if err := ts.Bind(ts.basicCfg.SearchDn, ts.basicCfg.SearchPassword); err != nil {
		if goldap.IsErrorWithCode(err, goldap.LDAPResultInvalidCredentials) {
			return false, ErrInvalidCredential
		}
	}
	return true, nil
}


func (s *Session) SearchUser(username string) ([]User, error) {
	var ldapUsers []User
	ldapFilter, err := createUserSearchFilter(s.basicCfg.Filter, s.basicCfg.UID, username)
	if err != nil {
		return nil, err
	}

	result, err := s.SearchLdap(ldapFilter)
	if err != nil {
		return nil, err
	}

	for _, ldapEntry := range result.Entries {
		var u User
		groupDNList := make([]string, 0)
		groupAttr := strings.ToLower(s.groupCfg.MembershipAttribute)
		for _, attr := range ldapEntry.Attributes {
			// OpenLdap sometimes contain leading space in username
			val := strings.TrimSpace(attr.Values[0])
			fmt.Printf("Current ldap entry attr name: %s\n", attr.Name)
			switch strings.ToLower(attr.Name) {
			case strings.ToLower(s.basicCfg.UID):
				u.Username = val
			case "uid":
				u.Realname = val
			case "cn":
				u.Realname = val
			case "mail":
				u.Email = val
			case "email":
				u.Email = val
			case groupAttr:
				for _, dnItem := range attr.Values {
					groupDNList = append(groupDNList, strings.TrimSpace(dnItem))
					fmt.Printf("Found memberof %v", dnItem)
				}
			}
			u.GroupDNList = groupDNList
		}
		u.DN = ldapEntry.DN
		ldapUsers = append(ldapUsers, u)
	}

	return ldapUsers, nil

}

func (s *Session) SearchLdap(filter string) (*goldap.SearchResult, error) {
	attributes := []string{"uid", "cn", "mail", "email"}
	lowerUID := strings.ToLower(s.basicCfg.UID)

	if lowerUID != "uid" && lowerUID != "cn" && lowerUID != "mail" && lowerUID != "email" {
		attributes = append(attributes, s.basicCfg.UID)
	}

	// Add the Group membership attribute
	groupAttr := strings.TrimSpace(s.groupCfg.MembershipAttribute)
	fmt.Printf("Membership attribute: %s\n", groupAttr)
	attributes = append(attributes, groupAttr)

	return s.SearchLdapAttribute(s.basicCfg.BaseDn, filter, attributes)
}

func (s *Session) SearchLdapAttribute(baseDN, filter string, attributes []string) (*goldap.SearchResult, error) {

	if err := s.Bind(s.basicCfg.SearchDn, s.basicCfg.SearchPassword); err != nil {
		return nil, fmt.Errorf("can not bind search dn, error: %v", err)
	}
	filter = normalizeFilter(filter)
	if len(filter) == 0 {
		return nil, ErrInvalidFilter
	}
	if _, err := goldap.CompileFilter(filter); err != nil {

		fmt.Errorf("Wrong filter format, filter:%v", filter)
		return nil, ErrInvalidFilter
	}
	fmt.Printf("Search ldap with filter:%v", filter)
	searchRequest := goldap.NewSearchRequest(
		baseDN,
		//goldap.ScopeWholeSubtree, => 2
		s.basicCfg.Scope,
		goldap.NeverDerefAliases,
		0,     // Unlimited results
		0,     // Search Timeout
		false, // Types only
		filter,
		attributes,
		nil,
	)

	result, err := s.ldapConn.Search(searchRequest)
	if result != nil {
		fmt.Printf("Found entries:%v\n", len(result.Entries))
	} else {
		fmt.Printf("No entries")
	}

	if err != nil {
		fmt.Printf("LDAP search error", err)
		return nil, err
	}

	return result, nil

}

func createUserSearchFilter(origFilter, ldapUID, username string) (string, error) {
	oFilter, err := NewFilterBuilder(origFilter)
	if err != nil {
		return "", err
	}
	var filterTag string
	filterTag = goldap.EscapeFilter(username)
	if len(filterTag) == 0 {
		filterTag = "*"
	}
	uFilterStr := fmt.Sprintf("(%v=%v)", ldapUID, filterTag)
	uFilter, err := NewFilterBuilder(uFilterStr)
	if err != nil {
		return "", err
	}
	filter := oFilter.And(uFilter)
	return filter.String()
}