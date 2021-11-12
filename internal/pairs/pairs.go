package pairs

import (
	"context"
	database "internal/pkg/db/mysql"
	"internal/tokens"
	"internal/users"
	"log"
)

type Pair struct {
	ID                  string
	token0              string
	token1              string
	balance0            float64
	balance1            float64
	marketCap           float64
	totalVolumeRecorded float64
}

var whiteListTokens = []string{"btc", "ethereum", "tether", "cardano", "solana", "usd-coin", "terra-luna", "ripple", "polkadot", "dogecoin"}

func getPair(token0 string, token1 string) (Pair, error) {
	stmt, err := database.Db.Prepare("select * from Pair where (token0 = ? and token1 = ?) or (token0 = ? and token1 = ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var pair Pair
	row := stmt.QueryRow(token0, token1, token1, token0)

	err = row.Scan(&pair.ID, &pair.token0, &pair.token1, &pair.balance0, &pair.balance1, &pair.marketCap, &pair.totalVolumeRecorded)
	if err != nil {
		return pair, err
	}
	return pair, nil
}

// always returns token0/token1
func GetExchangeRate(token0 string, token1 string) float64 {
	pair, err := getPair(token0, token1)
	if err != nil {
		return 0
	}
	return pair.balance0 / pair.balance1
}

func calcExactIn(inToken string, amountIn float64, pair Pair) float64 {
	constK := pair.balance0 * pair.balance1
	if inToken == pair.token0 {
		newAmount0 := pair.balance0 + amountIn
		return pair.balance1 - constK/newAmount0
	} else {
		newAmount1 := pair.balance1 + amountIn
		return pair.balance0 - constK/newAmount1
	}
}

func GetBestExchangeRate(token0 string, token1 string, amountIn float64) (float64, string) {
	stmt, err := database.Db.Prepare("select * from Pair where " +
		"((token0 = ? or token0 = ?) and token1 IN ('btc', 'ethereum', 'tether', 'cardano', 'solana', 'usd-coin', 'terra-luna', 'ripple', 'polkadot', 'dogecoin')) or" +
		"((token1 = ? or token1 = ?) and token0 IN ('btc', 'ethereum', 'tether', 'cardano', 'solana', 'usd-coin', 'terra-luna', 'ripple', 'polkadot', 'dogecoin'))")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(token0, token1, token0, token1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	pairMap := make(map[string]Pair)
	for rows.Next() {
		var pair Pair
		err := rows.Scan(&pair.ID, &pair.token0, &pair.token1, &pair.balance0, &pair.balance1, &pair.marketCap, &pair.totalVolumeRecorded)
		if err != nil {
			log.Fatal(err)
		}
		pairMap[pair.token0+"-"+pair.token1] = pair
		pairMap[pair.token1+"-"+pair.token0] = pair
	}

	var bestAns float64
	var bestIntermediaryToken string
	bestAns = 0

	if pairMap[token0+"-"+token1].ID != "" {
		bestIntermediaryToken = ""
		bestAns = calcExactIn(token0, amountIn, pairMap[token0+"-"+token1])
	}

	for _, whiteListed := range whiteListTokens {
		if whiteListed == token0 || whiteListed == token1 {
			continue
		}
		pair0 := pairMap[token0+"-"+whiteListed]
		pair1 := pairMap[token1+"-"+whiteListed]

		if pair0.ID == "" || pair1.ID == "" {
			continue
		}

		rate := calcExactIn(whiteListed, calcExactIn(token0, amountIn, pair0), pair1)
		if rate > bestAns {
			bestAns = rate
			bestIntermediaryToken = whiteListed
		}
	}

	return bestAns, bestIntermediaryToken
}

func Swap(userID int, token0 string, token1 string, inAmount float64, intermediaryToken string) bool {
	if token0 == token1 {
		return false
	}

	tokenMap := make(map[string]tokens.Token)
	tokenMap[token0] = tokens.GetById(token0)
	tokenMap[token1] = tokens.GetById(token1)
	if intermediaryToken != "" {
		tokenMap[intermediaryToken] = tokens.GetById(intermediaryToken)
	}

	bal0, _ := users.GetUserBalance(userID, token0)
	var bal1 users.UserBalance
	bal2, _ := users.GetUserBalance(userID, token1)
	if intermediaryToken != "" {
		bal1, _ = users.GetUserBalance(userID, intermediaryToken)
	}

	var pair0 Pair
	var pair1 Pair

	if intermediaryToken == "" {
		pair0, _ = getPair(token0, token1)
	} else {
		pair0, _ = getPair(token0, intermediaryToken)
		pair1, _ = getPair(intermediaryToken, token1)
	}

	ctx := context.TODO()
	tx, err := database.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal()
	}

	txSwapFunc := func(inToken string, outToken string, inAmount float64, pair *Pair, inBalance *users.UserBalance, outBalance *users.UserBalance) (float64, error) {
		UPDATE_USER_BALANCE := "UPDATE UserBalance set balance = ? where userID = ? and tokenID = ?"
		UPDATE_PAIR_RESERVE := "UPDATE Pair set balance0 = ?, balance1 = ?, totalVolumeRecorded = ? where ID = ?"
		outAmount := calcExactIn(inToken, inAmount, *pair)

		inBalance.Balance = inBalance.Balance - inAmount
		_, err := tx.ExecContext(ctx, UPDATE_USER_BALANCE, inBalance.Balance, userID, inToken)

		if err != nil {
			tx.Rollback()
			log.Println(err)
			return 0, err
		}

		outBalance.Balance = outBalance.Balance + outAmount
		_, err = tx.ExecContext(ctx, UPDATE_USER_BALANCE, outBalance.Balance, userID, outToken)

		if err != nil {
			tx.Rollback()
			return 0, err
		}

		delta0 := inAmount
		delta1 := -outAmount
		if inToken != pair.token0 {
			delta0 = -outAmount
			delta1 = inAmount
		}
		pair.balance0 = pair.balance0 + delta0
		pair.balance1 = pair.balance1 + delta1
		pair.totalVolumeRecorded = pair.totalVolumeRecorded + ((inAmount*tokenMap[inToken].Price)+(outAmount*tokenMap[outToken].Price))/2
		_, err = tx.ExecContext(ctx, UPDATE_PAIR_RESERVE, pair.balance0, pair.balance1, pair.totalVolumeRecorded, pair.ID)

		if err != nil {
			tx.Rollback()
			return 0, err
		}
		return outAmount, nil
	}

	if intermediaryToken != "" {
		inAmount, err = txSwapFunc(token0, intermediaryToken, inAmount, &pair0, &bal0, &bal1)
		if err != nil {
			tx.Rollback()
			return false
		}
		_, err = txSwapFunc(intermediaryToken, token1, inAmount, &pair1, &bal1, &bal2)
		if err != nil {
			tx.Rollback()
			return false
		}
	} else {
		_, err = txSwapFunc(token0, token1, inAmount, &pair0, &bal0, &bal2)
		if err != nil {
			tx.Rollback()
			return false
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return true
}
