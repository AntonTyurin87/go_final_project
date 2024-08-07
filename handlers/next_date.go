package handlers

import (
	"fmt"
	"net/http"

	"go_final_project/tasks"
)

// GetNextDateHandler - возвращает значение новой даты, если оно валидно.
func GetNextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверяем корректность формата входящего времени
	nowTime, err := tasks.DateValidation(now)
	if err != nil {
		fmt.Println("Ошибка конвертации входящего времени nowTime. ", err)
	} else {
		// Получение следующей даты
		res, err := tasks.NextDate(nowTime, date, repeat)
		if err != nil {
			fmt.Println("Ошибка получения NextDate. ", err)
		}

		_, err = w.Write([]byte(res))
		if err != nil {
			fmt.Println("Ошибка формирования ответа. ", err)
		}
	}
}
