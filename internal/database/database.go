package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func GetInstance() (*Database, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(userHome, "kaimd.db")

	instance, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	stmt := `
		CREATE TABLE IF NOT EXISTS cache (
			key TEXT PRIMARY KEY NOT NULL,
			data JSONB NOT NULL,
			updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err = instance.Exec(stmt)
	if err != nil {
		return nil, err

	}

	return &Database{DB: instance}, nil
}

func (db *Database) SetCache(key string, data any) error {

	stmt, err := db.Prepare("INSERT OR REPLACE INTO cache(key, data) VALUES(?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}

	defer stmt.Close()

	json, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = stmt.Exec(key, string(json)); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (db *Database) ReadCache(key string, data interface{}) error {

	stmt, err := db.Prepare("select data from cache where key = ?")
	if err != nil {
		log.Println(err)
		return err
	}

	defer stmt.Close()

	var res string
	err = stmt.QueryRow(key).Scan(&res)
	if err != nil {
		log.Println(err)
		return err
	}

	json.Unmarshal([]byte(res), &data)

	return nil

}
