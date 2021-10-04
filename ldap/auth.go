package ldap

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func Authenticate(ctx context.Context, m AuthModel, ctl Controller, ldapConf *LdapConf, groupConf *GroupConf) (*User, error) {
	p := m.Principal
	if len(strings.TrimSpace(p)) == 0 {
		fmt.Println("LDAP authentication failed for empty user id.")
		return nil, errors.New("Empty user id")
	}
	ldapSession, err := ctl.Session(ctx, ldapConf, groupConf)
	if err != nil {
		fmt.Println(err)
	}
	if err = ldapSession.Open(); err != nil {
		fmt.Printf("ldap connection fail: %v", err)
		return nil, err
	}
	defer ldapSession.Close()

	ldapUsers, err := ldapSession.SearchUser(p)
	if err != nil {
		fmt.Printf("ldap search fail: %v", err)
		return nil, err
	}
	if len(ldapUsers) == 0 {
		fmt.Printf("Not found an entry.")
		return nil, errors.New("Not found an entry")
	} else if len(ldapUsers) != 1 {
		fmt.Printf("Found more than one entry.")
		return nil, errors.New("Multiple entries found")
	}
	fmt.Printf("found ldap user %+v", ldapUsers[0])

	dn := ldapUsers[0].DN
	if err = ldapSession.Bind(dn, m.Password); err != nil {
		fmt.Printf("Failed to bind user, username %s, dn: %s, error: %v", p, dn, err)
		return nil, errors.New(err.Error())
	}

	u := User{}
	u.Username = ldapUsers[0].Username
	u.Realname = ldapUsers[0].Realname
	u.Email = strings.TrimSpace(ldapUsers[0].Email)


	return &u, nil
}
