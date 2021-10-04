package ldap

import (
	"context"
)

func NewService() Service {
	return &service{}
}

type service struct {
}

type Service interface {
	Ping(ctx context.Context, cfg LdapConf) (bool, error)
	SearchUser(ctx context.Context, sess *Session, username string) ([]User, error)
}

func (s *service) Ping(ctx context.Context, cfg LdapConf) (bool, error) {
	return TestConfig(cfg)
}

func (s *service) SearchUser(ctx context.Context, sess *Session, username string) ([]User, error) {
	users := make([]User, 0)
	if err := sess.Open(); err != nil {
		return users, err
	}
	defer sess.Close()

	ldapUsers, err := sess.SearchUser(username)
	if err != nil {
		return users, err
	}
	for _, u := range ldapUsers {
		ldapUser := User{
			Username:    u.Username,
			Realname:    u.Realname,
			GroupDNList: u.GroupDNList,
			Email:       u.Email,
		}
		users = append(users, ldapUser)
	}
	return users, nil
}
