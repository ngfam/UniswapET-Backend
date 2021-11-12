package pairs

import (
	database "internal/pkg/db/mysql"
	"log"
	"strconv"
)

type TopPair struct {
	ID                  string  `json:"id"`
	Token0              string  `json:"token0"`
	Token1              string  `json:"token1"`
	Icon0               string  `json:"icon0"`
	Icon1               string  `json:"icon1"`
	TotalVolumeRecorded float64 `json:"totalVolumeRecorded"`
	MarketCap           float64 `json:"marketCap"`
}

func TopTradedPairs(pageId int) []TopPair {
	stmt, err := database.Db.Prepare("CALL showTopTradedPairs(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(pageId)
	if err != nil {
		log.Fatal(err)
	}

	var topPairs []TopPair
	id := 10 * (pageId - 1)

	for rows.Next() {
		var pair TopPair
		id = id + 1
		pair.ID = strconv.Itoa(id)

		err = rows.Scan(&pair.TotalVolumeRecorded, &pair.MarketCap, &pair.Token0, &pair.Token1, &pair.Icon0, &pair.Icon1)
		if err != nil {
			log.Fatal(err)
		}
		topPairs = append(topPairs, pair)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return topPairs
}
