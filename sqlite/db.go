package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

// FindOrCreateDB - ищет файл базы данных в папке запуска приложения.
func FindOrCreateDB(todoDB string) (string, error) {

	var dbURL string

	// Если переменная окружения не задана или пуста присвоим адрес пакета проетка
	if todoDB == "" {

		appPath := "./"

		dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
		_, err := os.Stat(dbFile)

		fmt.Println("База тут ", dbFile)

		if err = createDB(dbFile); err != nil {
			fmt.Println("Не удалось создать БД", err)
			return dbURL, err
		}

		dbURL = dbFile

	} else {
		dbURL = todoDB
	}

	return dbURL, nil
}

// createDB - создаёт БД по адресу запуска приложения
func createDB(dbFile string) error {

	// Создаёт файл БД
	_, err := os.Create(dbFile)
	if err != nil {
		fmt.Println("Файл БД не создали", err)
		return err
	}

	// Открывает файл БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println("Файл БД не открывается", err)
		return err
	}
	defer db.Close()

	// Считывает текст SQL комад для создания БД.
	textSQL, err := os.ReadFile("sqlite/scheduler_creator.sql")
	if err != nil {
		fmt.Println("Не удалось прочитать sqlite/scheduler_creator.sql для создания БД ", err)
		return err
	}
	stringSQL := string(textSQL)

	// Создаёт таблицу в БД.
	_, err = db.Exec(stringSQL)
	if err != nil {
		fmt.Println("Не удалось создать таблицу в БД", err)
		return err
	}
	defer db.Close()

	return nil
}
