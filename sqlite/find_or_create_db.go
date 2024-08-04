package sqlite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// FindOrCreateDB - ищет файл базы данных в папке запуска приложения.
func FindOrCreateDB(todoDB string) (string, error) {

	var dbURL string

	//Если переменная окружения не задани или пуста присвоем адрес текущего каталога
	if todoDB == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
		_, err = os.Stat(dbFile)

		fmt.Println("База тут ", dbFile)

		if err = CreateDB(dbFile); err != nil {
			fmt.Println("Не удалось создать БД", err)
			return dbURL, err
		}

		dbURL = dbFile

	} else {
		dbURL = todoDB
	}

	return dbURL, nil
}
