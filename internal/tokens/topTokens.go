package tokens

import (
	database "internal/pkg/db/mysql"
	"log"
)

func TopMarketcapTokens(pageId int) []Token {
	stmt, err := database.Db.Prepare("CALL showTopMarketCapTokens(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(pageId)
	var topTokens []Token

	for rows.Next() {
		var token Token
		err = rows.Scan(&token.ID, &token.Name, &token.TotalSupply, &token.IconURL, &token.Price)
		if err != nil {
			log.Fatal(err)
		}
		topTokens = append(topTokens, token)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return topTokens
}
