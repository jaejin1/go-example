package ldap

//func LDAPConf(ctx context.Context) (*cfgModels.LdapConf, error) {
//	mgr := defaultMgr()
//	err := mgr.Load(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return &cfgModels.LdapConf{
//		URL:               mgr.Get(ctx, common.LDAPURL).GetString(),
//		SearchDn:          mgr.Get(ctx, common.LDAPSearchDN).GetString(),
//		SearchPassword:    mgr.Get(ctx, common.LDAPSearchPwd).GetString(),
//		BaseDn:            mgr.Get(ctx, common.LDAPBaseDN).GetString(),
//		UID:               mgr.Get(ctx, common.LDAPUID).GetString(),
//		Filter:            mgr.Get(ctx, common.LDAPFilter).GetString(),
//		Scope:             mgr.Get(ctx, common.LDAPScope).GetInt(),
//		ConnectionTimeout: mgr.Get(ctx, common.LDAPTimeout).GetInt(),
//		VerifyCert:        mgr.Get(ctx, common.LDAPVerifyCert).GetBool(),
//	}, nil
//}
//
//// LDAPGroupConf returns the setting of ldap group search
//func LDAPGroupConf(ctx context.Context) (*cfgModels.GroupConf, error) {
//	mgr := defaultMgr()
//	err := mgr.Load(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return &cfgModels.GroupConf{
//		BaseDN:              mgr.Get(ctx, common.LDAPGroupBaseDN).GetString(),
//		Filter:              mgr.Get(ctx, common.LDAPGroupSearchFilter).GetString(),
//		NameAttribute:       mgr.Get(ctx, common.LDAPGroupAttributeName).GetString(),
//		SearchScope:         mgr.Get(ctx, common.LDAPGroupSearchScope).GetInt(),
//		AdminDN:             mgr.Get(ctx, common.LDAPGroupAdminDn).GetString(),
//		MembershipAttribute: mgr.Get(ctx, common.LDAPGroupMembershipAttribute).GetString(),
//	}, nil
//}