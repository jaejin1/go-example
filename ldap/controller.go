package ldap

import "context"

type Controller interface {
	Ping(ctx context.Context, cfg LdapConf) (bool, error)
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{
		service: service,
	}
}

func (c *controller) Ping(ctx context.Context, cfg LdapConf) (bool, error) {
	if len(cfg.SearchPassword) == 0 {
		//pwd, err := defaultPassword(ctx)
	}

	return c.service.Ping(ctx, cfg)
}

//func defaultPassword(ctx context.Context) (string, error) {
//	mod, err := config.AuthMode(ctx)
//	if err != nil {
//		return "", err
//	}
//	if mod == common.LDAPAuth {
//		conf, err := config.LDAPConf(ctx)
//		if err != nil {
//			return "", err
//		}
//		if len(conf.SearchPassword) == 0 {
//			return "", ldap.ErrEmptyPassword
//		}
//		return conf.SearchPassword, nil
//	}
//	return "", ldap.ErrEmptyPassword
//}
