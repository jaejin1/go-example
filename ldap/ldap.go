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
