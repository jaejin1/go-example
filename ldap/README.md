# ldap

~~~go
package main

import (
	"context"
	"fmt"
	"go-example/ldap"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ldapService := ldap.NewService()
	ldapController := ldap.NewController(ldapService)

	ldapConf := &ldap.LdapConf{
		URL: "",
		SearchDn: "",
		SearchPassword: "",
		BaseDn: "",
		Filter: "",
		UID: "",
		Scope: 2,
		ConnectionTimeout: 10,
		VerifyCert: true,
	}

	ldapGroupConf := &ldap.GroupConf{
		BaseDN: "",
		Filter: "",
		NameAttribute: "",
		AdminDN: "",
		MembershipAttribute: "",
	}

	ping, err := ldapController.Ping(ctx, *ldapConf)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ping)

	users, err := ldapController.SearchUser(ctx, "", ldapConf, ldapGroupConf)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(users)

	userModel := ldap.AuthModel{
		Principal: "example@xxx.xxx",
		Password: "password",
	}
	user, err := ldap.Authenticate(ctx, userModel, ldapController, ldapConf, ldapGroupConf)
	
	fmt.Println(user)
}
~~~