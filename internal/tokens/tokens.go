package tokens

import (
	database "internal/pkg/db/mysql"
	"log"
)

type Token struct {
	ID          string
	Name        string
	TotalSupply float64
	IconURL     string
	Price       float64
}

func GetAll() []Token {
	stmt, err := database.Db.Prepare("select * from Token limit 100")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tokens []Token
	for rows.Next() {
		var token Token
		err := rows.Scan(&token.ID, &token.Name, &token.TotalSupply, &token.IconURL, &token.Price)

		if err != nil {
			log.Fatal(err)
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return tokens
}

func GetByPrefix(prefix string) []Token {
	stmt, err := database.Db.Prepare("select * from Token where LOWER(id) like ? or LOWER(name) like ? limit 30")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(prefix+"%", prefix+"%")
	if err != nil {
		log.Println(stmt)
		log.Fatal(err)
	}

	var tokens []Token
	for rows.Next() {
		var token Token
		err := rows.Scan(&token.ID, &token.Name, &token.TotalSupply, &token.IconURL, &token.Price)

		if err != nil {
			log.Fatal(err)
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return tokens
}

func GetById(id string) Token {
	stmt, err := database.Db.Prepare("select * from Token where ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	var token Token
	err = row.Scan(&token.ID, &token.Name, &token.TotalSupply, &token.IconURL, &token.Price)
	if err != nil {
		log.Fatal(err)
	}
	return token
}
