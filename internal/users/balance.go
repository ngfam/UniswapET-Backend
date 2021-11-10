package users

import (
	"database/sql"
	database "internal/pkg/db/mysql"
	"log"
)

type UserBalance struct {
	User    int     `json:"user"`
	Token   string  `json:"token"`
	Balance float64 `json:"balance"`
}

func CreateUserBalance(userID int, tokenID string) {
	statement, err := database.Db.Prepare("INSERT INTO UserBalance(userID, tokenID, balance) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(userID, tokenID, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserBalance(userID int, tokenID string) (UserBalance, error) {
	statement, err := database.Db.Prepare("SELECT * from UserBalance where userID = ? and tokenID = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	var userBalance UserBalance
	err = (statement.QueryRow(userID, tokenID)).Scan(&userBalance.User, &userBalance.Token, &userBalance.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			CreateUserBalance(userID, tokenID)
			userBalance.User = userID
			userBalance.Token = tokenID
			userBalance.Balance = 0
		} else {
			log.Fatal(err)
		}
	}
	return userBalance, nil
}

func IncreaseBalance(userID int, tokenID string, delta float64) {
	userBalance, _ := GetUserBalance(userID, tokenID)
	if userBalance.Balance+delta < 0 {
		log.Fatal("Payment exceed balance.")
	}
	statement, err := database.Db.Prepare("UPDATE UserBalance SET balance = ? where userID = ? and tokenID = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(userBalance.Balance+delta, userID, tokenID)
	if err != nil {
		log.Fatal(err)
	}
}
