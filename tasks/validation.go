package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const DateFormat = "20060102"

// TaskDataValidation - проверяет корректность входящих данных
func TaskDataValidation(httpData []byte) (Task, error) {

	var taskData Task
	now := time.Now()

	// Приводим текущее время к формату
	nowString := now.Format(DateFormat)
	now, err := time.Parse(DateFormat, nowString)
	if err != nil {
		fmt.Println("Строковые данные даты не корректны. ", err)
		return taskData, err
	}

	// Проверка распаковки JSON
	err = json.Unmarshal(httpData, &taskData)
	if err != nil {
		fmt.Println("не удачно распаковался JSON запрос ", err)
		return taskData, err
	}

	// Проверка id
	if taskData.ID != "" {
		_, err := IDValidation(taskData.ID)
		if err != nil {
			fmt.Println("не верый формат id ", err)
			return taskData, err
		}
	}

	// Проверка заголовка
	if taskData.Title == "" {
		err0 := errors.New("заголовок задачи отсутствует")
		return taskData, err0
	}

	// Проверка значений повторений
	if RepeatValidation(taskData.Repeat) != nil {
		return taskData, RepeatValidation(taskData.Repeat)
	}

	// Проверка на пустое значение даты
	if taskData.Date == "" {
		taskData.Date = fmt.Sprint(now.Format(DateFormat))
	}

	// Проверка на корректность даты
	date, err := DateValidation(taskData.Date)
	if err != nil {
		fmt.Println(taskData.Date, err)
		return taskData, err
	}
	taskData.Date = fmt.Sprint(date.Format(DateFormat))

	// Если дата меньше сегодняшей
	// Значение переданной даты. Ошибку игнорируем по причине проверки выше.
	dateIn, _ := time.Parse(DateFormat, taskData.Date)

	if now.After(dateIn) {
		if taskData.Repeat == "" {
			taskData.Date = fmt.Sprint(now.Format(DateFormat))
		} else {
			taskData.Date, err = NextDate(now, taskData.Date, taskData.Repeat)
			if err != nil {
				return taskData, err
			}
		}
	}

	return taskData, nil
}
