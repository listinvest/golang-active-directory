package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/ldap.v3"
)

func main() {
	username := "<username>@<domain>.com"
	password := "<password>"
	ldapServer := "your-server.domain.com"

	// Try to connect to LDAP server
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, 389))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Reconnect with TLS
	err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Bind(username, password)
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		// The base domain name to search
		"dc=<domain>,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		// The filter to apply
		fmt.Sprintf("(&(objectClass=organizationalPerson)(cn=%s))", "<username>"),
		// List of attributes to retrieve
		[]string{"c", "cn", "displayName", "userPrincipalName"},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	if len(searchResult.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}

	searchResult.PrettyPrint(2)
}
