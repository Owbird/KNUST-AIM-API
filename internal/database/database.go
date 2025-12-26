package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
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
	instance, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	stmt := `
		CREATE TABLE IF NOT EXISTS cache (
			key TEXT PRIMARY KEY NOT NULL,
			data JSONB NOT NULL,
			updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			ttl INTEGER NOT NULL DEFAULT 1
		)
	`

	_, err = instance.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return &Database{DB: instance}, nil
}

func (db *Database) SetCache(key string, data any, ttl int64) error {
	stmt, err := db.Prepare("INSERT OR REPLACE INTO cache(key, data, ttl) VALUES(?, ?, ?)")
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

	if _, err = stmt.Exec(key, string(json), ttl); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *Database) ReadCache(key string, data interface{}) error {
	stmt, err := db.Prepare("select data, ttl, updatedAt from cache where key = ?")
	if err != nil {
		log.Println(err)
		return err
	}

	defer stmt.Close()

	var (
		res       string
		ttl       int64
		updatedAt string
	)
	err = stmt.QueryRow(key).Scan(&res, &ttl, &updatedAt)
	if err != nil {
		log.Println(err)
		return err
	}

	// ttl == 0 does not expire
	if ttl != 0 {
		now := time.Now()
		parsedUpdatedAt, _ := time.Parse(time.RFC3339, updatedAt)

		expirationDuration := time.Duration(ttl) * 24 * time.Hour
		if now.Sub(parsedUpdatedAt) > expirationDuration {
			log.Printf("Data is stale for %v", key)
			return fmt.Errorf("data is stale")
		}
	}

	json.Unmarshal([]byte(res), &data)

	return nil
}
