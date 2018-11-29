package main

import (
	"log"

	"github.com/jtblin/go-ldap-client"
)

func main() {
	client := &ldap.LDAPClient{
		Base:         "dc=example,dc=com",
		Host:         "ldap.exapmle.com",
		Port:         636,
		UseSSL:       true,
		BindDN:       "uid=xl,cn=users,cn=accounts,dc=example,dc=com",
		BindPassword: "password",
		UserFilter:   "(uid=%s)",
		GroupFilter: "(memberUid=%s)",
		//GroupFilter: "(gidNumber=*)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	client.ServerName = "ldap.example.com"
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, user, err := client.Authenticate("xl", "password")
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", "username", err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", "username")
	}
	log.Printf("User: %+v", user)
	
	groups, err := client.GetGroupsOfUser("xl")
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", "username", err)
	}
	log.Printf("Groups: %+v", groups) 
}
