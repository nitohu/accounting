package main

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/nitohu/accounting/server/models"
)

func main() {
	var a models.API
	a.LocalKey = true
	a.GenerateAPIKey()
	k := a.GetAPIKey()
	prefix := strings.Split(k, ".")[0]
	h := sha256.Sum256([]byte(k))
	s := fmt.Sprintf("%x", h)
	query := "INSERT INTO api (active, name, create_date, last_update, api_key, api_prefix, local_key, access_rights) "
	query += "VALUES ('t', 'My API Key', NOW(), NOW(), '" + s + "', '" + prefix + "', 'f', 'ACCESS_RIGHTS');"

	fmt.Println("This is your plain API Key. Make sure you don't lose it and it's stored in a safe location:")
	fmt.Println(k + "\n")
	fmt.Println("You can use the following database query to insert this API key into the database:")
	fmt.Println(query + "\n")
	fmt.Println("You'll need to replace the ACCESS_RIGHTS string with the actual access rights you want to use.")
	fmt.Println("You can choose from the following list:")
	rights := models.GetAllAccessRights()
	for _, r := range rights {
		fmt.Printf("\t- %s\n", r)
	}
	fmt.Println("If you just want to choose 1 access right you can replace the ACCESS_RIGHTS.")
	fmt.Println("If you want to choose multiple access rights, just make a ; seperated list out of it. E.g.:")
	fmt.Println("transaction.read;account.read;statistic.read")
}
