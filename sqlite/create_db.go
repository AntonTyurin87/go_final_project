package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

//var AdresDB string

// CreateDB - создаёт БД по адресу запуска приложения
func CreateDB(dbFile string) error {

	//Создаёт файл БД
	_, err := os.Create(dbFile)
	if err != nil {
		fmt.Println("Файл БД не создали", err)
		return err
	}

	//Открывает файл БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println("Файл БД не открывается", err)
		return err
	}
	defer db.Close()

	//Считывает текст SQL комад для создания БД.
	textSQL, err := os.ReadFile("sqlite/scheduler_creator.sql")
	if err != nil {
		fmt.Println("Не удалось прочитать sqlite/scheduler_creator.sql для создания БД ", err)
		return err
	}
	stringSQL := string(textSQL)

	//Создаёт таблицу в БД.
	_, err = db.Exec(stringSQL)
	if err != nil {
		fmt.Println("Не удалось создать таблицу в БД", err)
		return err
	}
	defer db.Close()

	return nil
}
