package handlers

import (
	"fmt"
	"net/http"

	"go_final_project/datawork"
)

// NextDateHandler - возвращает значение новой даты, если оно валидно.
func GetNextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	//Проверяем корректность формата входящего времени
	nowTime, err := datawork.DateValidation(now)
	if err != nil {
		fmt.Println("Ошибка конвертации входящего времени nowTime. ", err)
	} else {
		//Получение следующей даты
		res, err := datawork.NextDate(nowTime, date, repeat)
		if err != nil {
			fmt.Println("Ошибка получения NextDate. ", err)
		}

		w.Write([]byte(res))
	}
}
