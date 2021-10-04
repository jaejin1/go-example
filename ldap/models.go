package ldap

type LdapConf struct {
	URL               string `json:"ldap_url"`
	SearchDn          string `json:"ldap_search_dn"`
	SearchPassword    string `json:"ldap_search_password"`
	BaseDn            string `json:"ldap_base_dn"`
	Filter            string `json:"ldap_filter"`
	UID               string `json:"ldap_uid"`
	Scope             int    `json:"ldap_scope"`
	ConnectionTimeout int    `json:"ldap_connection_timeout"`
	VerifyCert        bool   `json:"ldap_verify_cert"`
}

// GroupConf holds information about ldap group
type GroupConf struct {
	BaseDN              string `json:"ldap_group_base_dn,omitempty"`
	Filter              string `json:"ldap_group_filter,omitempty"`
	NameAttribute       string `json:"ldap_group_name_attribute,omitempty"`
	SearchScope         int    `json:"ldap_group_search_scope"`
	AdminDN             string `json:"ldap_group_admin_dn,omitempty"`
	MembershipAttribute string `json:"ldap_group_membership_attribute,omitempty"`
}
