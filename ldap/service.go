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
}

func (s *service) Ping(ctx context.Context, cfg LdapConf) (bool, error) {
	return TestConfig(cfg)
}
