package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dbFilePath = os.Getenv("sqlite__filepath")

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal("db connect faild.", err)
		return nil, err
	}

	return db, nil
}

func init() {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		file, err := os.Create(dbFilePath)
		defer file.Close()

		if err != nil {
			log.Fatal("Failed to init db:", err)
		}
	}

	err := createTable()
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
}

func createTable() error {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal("db connect failed.", err)
		return err
	}
	defer db.Close()

	createAccessTokenSQL := `CREATE TABLE IF NOT EXISTS AccessToken (
		"Id"			INTEGER				NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Token"		VARCHAR(50)		NOT NULL,
		"IssueAt"	DATETIME			NOT NULL
	);`

	_, err = db.Exec(createAccessTokenSQL)

	if err != nil {
		return err
	}

	return nil
}
