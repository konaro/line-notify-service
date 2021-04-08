package accesstoken

import (
	"log"
	"time"

	"github.com/konaro/line-notify-service/data"
	"github.com/konaro/line-notify-service/data/sqlite"
)

func Add(token string, time time.Time) {
	db, err := sqlite.Connect()

	if err != nil {
		return
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO AccessToken(Token, IssueAt) VALUES (?, ?);")

	if err != nil {
		log.Fatal("Insert statement error.", err)
		return
	}

	_, err = statement.Exec(token, time)

	if err != nil {
		log.Fatal("Insert value failed.", err)
		return
	}
}

func GetList(limit, offset int) []data.AccessToken {
	var tokens []data.AccessToken

	db, err := sqlite.Connect()
	if err != nil {
		return tokens
	}

	defer db.Close()

	rows, err := db.Query("SELECT Id, Token, IssueAt FROM AccessToken LIMIT ? OFFSET ?;", limit, offset)

	if err != nil {
		log.Fatal("query faild.", err)
		return tokens
	}

	for rows.Next() {
		var token data.AccessToken
		err = rows.Scan(&token.Id, &token.Token, &token.IssueAt)
		if err != nil {
			log.Fatal("fetch data error.", err)
		} else {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

func GetAllTokens() []string {
	var tokens []string

	db, err := sqlite.Connect()
	if err != nil {
		return tokens
	}

	defer db.Close()

	rows, err := db.Query("SELECT Token FROM AccessToken;")

	if err != nil {
		log.Fatal("query faild.", err)
		return tokens
	}

	for rows.Next() {
		var token string
		err = rows.Scan(&token)
		if err != nil {
			log.Fatal("fetch data error.", err)
		} else {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

func Delete(id int) error {
	db, err := sqlite.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	res, err := db.Exec("DELETE FROM AccessToken WHERE Id = ?", id)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()

	if err != nil {
		return err
	}

	return nil
}
