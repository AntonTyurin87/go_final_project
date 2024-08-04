package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/datawork"
	"io"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

// Строковый тип для специальной ошибки
type StringError struct {
	StrEr string `json:"error,omitempty"`
}

// Структура для ID
type IDType struct {
	ID int64 `json:"id"`
}

// Структура для группы записей
type TasksType struct {
	Tasks []datawork.TaskData `json:"tasks"`
}

// Структура для хранилища
type Storage struct {
	DB *sql.DB
}

// NewStorage - онструктор для хранилица
func NewStorage(db *sql.DB) Storage {
	return Storage{DB: db}
}

// GetOneTaskHandler - возвращает одну задачу по выданному признаку (по id)
func (s *Storage) GetOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes StringError
	var err error
	var result []byte

	id := r.FormValue("id")

	//Проверяем не пустой ли id
	switch {
	case id == "":
		errRes.StrEr = "Не указан идентификатор"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	default:
		//Проверяем id
		data, err := datawork.IDValidation(id)
		if err != nil {
			fmt.Println("Ошибка конвертации входящего значения api/task. ", err)
		}

		//Идём искать одну задачу по входным данным
		result, err = s.oneTaskDataRead(data)
		if err != nil {
			fmt.Println("Ошибка чтения из БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(result)
}

// PostOneTaskHandler -  записывает в базу одну задачу
func (s *Storage) PostOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var errRes StringError
	var result []byte

	//Читаем сообщение
	httpData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Не прочитано тело запроса api/task. ", err)
	}

	//Проверяем полученную информацию о задаче
	data, err := datawork.TaskDataValidation(httpData)

	switch {
	//Если не смогли валедировать входящие данные
	case err != nil:
		fmt.Println("Ошибка конвертации входящего значения api/task. ", err)
		errRes.StrEr = fmt.Sprint(err)
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	//Если заголовок пуст возвращаем ошибку
	case data.Title == "":
		errRes.StrEr = "Не указан заголовок задачи"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	//Идём записывать задачу в базу
	default:
		result, err = s.oneTaskDataWrite(data)
		if err != nil {
			fmt.Println("Ошибка записи в БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(result)
}

// PutOneTaskHandler - изменяет в базе одну задачу
func (s *Storage) PutOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes StringError
	var result []byte

	//Читаем сообщение
	httpData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Не прочитано тело запроса api/task. ", err)
	}

	//Проверяем полученную информацию о задаче
	data, err := datawork.TaskDataValidation(httpData)

	switch {
	case err != nil:
		fmt.Println("Ошибка конвертации входящего значения api/task. ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
		//Идём записывать задачу в базу
	default:
		result, err = s.oneTaskDataUpdate(data)
		if err != nil {
			fmt.Println("Ошибка записи в БД ", err)
			errRes.StrEr = "Задача не найдена"
			result, err = json.Marshal(errRes)
			if err != nil {
				fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(result)
}

// DoneOneTaskHandler - завершает в базе одну задачу
func (s *Storage) DoneOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes StringError
	var err error
	var result []byte

	id := r.FormValue("id")

	//Проверяем не пустой ли id
	switch {
	case id == "":
		errRes.StrEr = "Не указан идентификатор"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	default:
		//Проверяем id
		data, err := datawork.IDValidation(id)
		if err != nil {
			fmt.Println("Ошибка конвертации входящего значения api/task. ", err)
		}

		//Идём закрывать одну задачу по входным данным
		result, err = s.oneTaskDataDone(data)
		if err != nil {
			fmt.Println("Ошибка чтения из БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(result)
}

// DeleteOneTaskHandler - удаляет из базы одну задачу по выданному признаку (по id)
func (s *Storage) DeleteOneTaskHandler(w http.ResponseWriter, r *http.Request) {

	var errRes StringError
	var err error
	var result []byte

	id := r.FormValue("id")

	//Проверяем не пустой ли id
	switch {
	case id == "":
		errRes.StrEr = "Не указан идентификатор"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
	default:
		//Проверяем id
		data, err := datawork.IDValidation(id)
		if err != nil {
			fmt.Println("Ошибка конвертации входящего значения api/task. ", err)
		}

		//Идём искать одну задачу по входным данным
		result, err = s.oneTaskDataDelete(data)
		if err != nil {
			fmt.Println("Ошибка чтения из БД ", err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(result)
}

// GetTasksHandler - возвращает группу задач
func (s *Storage) GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	var errRes StringError
	var result []byte
	var err error

	search := r.FormValue("search")

	result, err = s.groupTasksDataRead(search)
	if err != nil {
		fmt.Println("Ошибка чтения из БД ", err)
		errRes.StrEr = fmt.Sprint(err)
		result, err = json.Marshal(errRes)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(result)
}

// oneTaskDataRead - возвращает информацию об одной задаче по входным данным
func (s *Storage) oneTaskDataRead(data datawork.TaskData) ([]byte, error) {

	var errRes StringError
	var result []byte
	var returnData datawork.TaskData
	var row *sql.Rows
	var err error

	//Формируем запрос в базу
	qeryToDB := `SELECT id, date, title, comment, repeat
					FROM scheduler
				WHERE id = ?;`

	row, err = s.DB.Query(qeryToDB, data.ID)
	if err != nil {
		fmt.Println("Чтение из БД не состоялась ", err)
		return result, err
	}
	defer row.Close()

	//Укладываем результаты запроса в структуру
	for row.Next() {
		if err := row.Scan(&returnData.ID, &returnData.Date, &returnData.Title, &returnData.Comment, &returnData.Repeate); err != nil {
			return nil, err
		}
		if err != nil {
			fmt.Println("Не удалось записать корректную дату из БД .", err)
			return result, err
		}
	}

	//Если id задачи отсутствует, то формируем сообщение об ошибке
	if returnData.ID == "" {
		errRes.StrEr = "Задача не найдена"
		result, err := json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		}
		return result, err
	}

	//Формируем сообщение с информацией о задаче
	result, err = json.Marshal(returnData)
	if err != nil {
		fmt.Println("Не получилось сформировать вывод в виде JSON ", err)
		return result, err
	}

	return result, nil
}

// oneTaskDataWrite - записывает в БД данные о внесённой задаче
func (s *Storage) oneTaskDataWrite(data datawork.TaskData) ([]byte, error) {

	var err error
	var result []byte
	var returnData IDType

	//Формируем запрос в базу
	qeryToDB := `INSERT INTO
					scheduler (date, title, comment, repeat)
				VALUES (?, ?, ?, ?);`

	res, err := s.DB.Exec(qeryToDB, data.Date, data.Title, data.Comment, data.Repeate)
	if err != nil {
		fmt.Println("Запись в БД не состоялась ", err)
		return result, err
	}

	//Возвращаем id последней записи
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("ID последней записи в БД не удалось получить ", err)
		return result, err
	}

	returnData.ID = id

	//Формируем сообщение с информацией о задаче
	result, err = json.Marshal(returnData)
	if err != nil {
		fmt.Println("Не получилось выдать ID последней записи в виде JSON ", err)
		return result, err
	}

	return result, nil
}

// oneTaskDataUpdate - изменяет в БД данные о внесённой задаче
func (s *Storage) oneTaskDataUpdate(data datawork.TaskData) ([]byte, error) {

	var errRes StringError
	var result []byte
	var err error

	//Формируем запрос в базу
	qeryToDB := `UPDATE
					scheduler SET date = ?, title = ?, comment = ?, repeat = ?
				WHERE id = ? ;`

	res, err := s.DB.Exec(qeryToDB, data.Date, data.Title, data.Comment, data.Repeate, data.ID)
	if err != nil {
		fmt.Println("Запись в БД не состоялась ", err)
		return result, err
	}

	//Возвращаем количество затронутых записей
	num, err := res.RowsAffected()

	if err != nil || num == 0 {
		fmt.Println("ID последней записи в БД не удалось получить ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	//Если ошибок не накопилось, то результат будет {}
	result, err = json.Marshal(errRes)
	if err != nil {
		fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		return result, err
	}

	return result, err
}

// oneTaskDataDelete - удаляет из БД данные об одной задаче
func (s *Storage) oneTaskDataDelete(data datawork.TaskData) ([]byte, error) {

	var errRes StringError
	var err error
	var result []byte

	//Формируем запрос в базу
	qeryToDB := `DELETE FROM
					scheduler 
					WHERE id = ?;`

	res, err := s.DB.Exec(qeryToDB, data.ID)
	if err != nil {
		fmt.Println("Удаление из БД не состоялась ", err)
		return result, err
	}

	//Возвращаем количество затронутых записей
	num, err := res.RowsAffected()

	if err != nil || num == 0 {
		fmt.Println("ID последней записи в БД не удалось получить ", err)
		errRes.StrEr = "Задача не найдена"
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	//Если ошибок не накопилось, то результат будет {}
	result, err = json.Marshal(errRes)
	if err != nil {
		fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		return result, err
	}

	return result, err
}

// oneTaskDataDone - удаляет из БД данные об одной задаче при её выполнении
func (s *Storage) oneTaskDataDone(data datawork.TaskData) ([]byte, error) {

	var errRes StringError
	var result []byte
	var returnData datawork.TaskData
	var row *sql.Rows
	var err error

	now := time.Now()

	//Формируем запрос в базу
	qeryToDB := `SELECT id, date, title, comment, repeat
					FROM scheduler
					WHERE id = ?;`

	row, err = s.DB.Query(qeryToDB, data.ID)
	if err != nil {
		fmt.Println("Чтение из БД не состоялась ", err)
		return result, err
	}
	defer row.Close()

	//Укладываем результаты запроса в структуру
	for row.Next() {
		if err := row.Scan(&returnData.ID, &returnData.Date, &returnData.Title, &returnData.Comment, &returnData.Repeate); err != nil {
			return nil, err
		}

		if err != nil {
			fmt.Println("Не удалось записать корректную дату из БД .", err)
			return result, err
		}
	}

	//Есть ли правило для повторения задачи
	switch {
	case returnData.Repeate == "":
		result, err = s.oneTaskDataDelete(returnData)
		if err != nil {
			fmt.Println("Не удалось удалить задачу из БД .", err)
			return result, err
		}
	default:
		returnData.Date, err = datawork.NextDate(now, returnData.Date, returnData.Repeate)
		if err != nil {
			fmt.Println("Не удалось получить новую дату для задачи .", err)
			return result, err
		}

		result, err = s.oneTaskDataUpdate(returnData)
		if err != nil {
			fmt.Println("Ошибка записи в БД ", err)
			return result, err
		}

	}

	//Если ошибок не накопилось, то результат будет {}
	result, err = json.Marshal(errRes)
	if err != nil {
		fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
		return result, err
	}

	return result, err
}

// groupTasksDataRead - возвращает информацию о группе последних задач
func (s *Storage) groupTasksDataRead(search string) ([]byte, error) {

	var errRes StringError
	var err error
	var task datawork.TaskData
	var rows *sql.Rows
	var queryToDB, searchDate string
	var result []byte
	returnData := TasksType{Tasks: make([]datawork.TaskData, 0, 20)}

	if search != "" && len(search) == 10 {
		searchDate, err = datawork.DateConvert(search)
		if err != nil {
			fmt.Println("На входе не дата ", err)
		}
	}

	switch {
	//Ищем всё подряд
	case search == "":
		queryToDB = `SELECT id, date, title, comment, repeat
						FROM scheduler
					ORDER BY date LIMIT 20;`

		rows, err = s.DB.Query(queryToDB)
		if err != nil {
			fmt.Println("Чтение из БД не состоялась ", err)
			return result, err
		}

	//Ищем по дате
	case searchDate != "":
		queryToDB = `SELECT id, date, title, comment, repeat
							FROM scheduler
							WHERE date = ?
						ORDER BY date LIMIT 20;`

		rows, err = s.DB.Query(queryToDB, searchDate)
		if err != nil {
			fmt.Println("Чтение из БД не состоялась ", err)
			return result, err
		}

	//Ищем по заголовку или комментарию
	default:

		search = fmt.Sprint("%" + search + "%")

		queryToDB = `SELECT id, date, title, comment, repeat
						FROM scheduler 
						WHERE title LIKE ? OR comment LIKE ?
					ORDER BY date LIMIT 20;`

		rows, err = s.DB.Query(queryToDB, search, search)

		if err != nil {
			fmt.Println("Чтение из БД не состоялась ", err)
			return result, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeate); err != nil {
			fmt.Println("Запись в структуру данных из БД не состоялась ", err)
			return result, err
		}

		returnData.Tasks = append(returnData.Tasks, task)
	}

	result, err = json.Marshal(returnData)
	if err != nil {
		fmt.Println("Не получилось сформировать вывод в виде JSON ", err)
		errRes.StrEr = fmt.Sprintln(err)
		result, err = json.Marshal(errRes)
		if err != nil {
			fmt.Println("Не удалось упаковать ошибку в JSON. ", err)
			return result, err
		}
		return result, err
	}

	return result, nil
}
