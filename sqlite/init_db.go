package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// InitDB - создаёт подключение к базе данных
func InitDB(dbURL string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", dbURL)

	if err != nil {
		fmt.Println("Нет подключения к базе данных", err)
	}

	return db, nil
}
