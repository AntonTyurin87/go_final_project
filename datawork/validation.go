package datawork

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// IDValidation - проверка id на числовое значение
func IDValidation(id string) (TaskData, error) {

	var Task TaskData

	_, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Id не является числом ", err)
		return Task, err
	}

	Task.ID = id

	return Task, nil
}

// DataValidation - проверка корректности переданной даты
func DateValidation(date string) (time.Time, error) {

	var dateTime time.Time

	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		fmt.Println("Строковые данные даты не корректны. ", err)
		return dateTime, err
	}

	return dateTime, nil
}

// RepeatValidation - проверка корректности переданного значениядля повторений
func RepeatValidation(repeat string) error {

	if repeat == "" {
		return nil
	}

	repeatData := strings.Split(repeat, " ")

	switch string(repeatData[0]) {
	//Ежегодно
	case "y":
		if len(repeatData) == 1 {
			return nil

		}
	//Через несколько дней
	case "d":
		if len(repeatData) < 2 {
			err2 := errors.New("не указан интервал для повторений в днях")
			return err2
		}

		if len(repeatData) > 2 {
			err3 := errors.New("не верно указан интервал для повторений в днях")
			return err3
		}

		daysCount, err := strconv.Atoi(repeatData[1])
		if err != nil {
			fmt.Println("Неверный формат дней для повторений. ", err)
			return err
		}

		if daysCount > 400 {
			err4 := errors.New("превышен интервал для повторений в днях")
			return err4
		}
	//По дням недели
	case "w":
		weekDays := strings.Split(repeat, " ")

		//Проверка на наличие дня недели
		if len(weekDays) < 2 {
			err5 := errors.New("не верный день недели")
			fmt.Println(err5)
			return err5
		}

		//Проверка на наличие одного дня недели
		if len(weekDays[1]) == 1 {
			dayNumber, err := strconv.Atoi(weekDays[1])

			if err != nil {
				fmt.Println("не верное значение дня недели", err)
				return err
			}

			if 0 >= dayNumber || dayNumber >= 8 {
				err6 := errors.New("не верный день недели")
				fmt.Println(err6)
				return err6
			}
			//Если дней не один
		} else {

			for _, value := range strings.Split(weekDays[1], ",") {
				day, err := strconv.Atoi(value)

				if err != nil {
					fmt.Println("не верное значение дня недели", err)
					return err
				}

				if 0 >= day || day >= 8 {
					err7 := errors.New("не верный день недели")
					fmt.Println(err7)
					return err7
				}
			}
		}
		/*
			//По месяцам
			case "m": //TODO Дописать проверку по месяцам

				monthDays := strings.Split(repeat, " ")

				months := map[time.Month]int{time.January: 31}

				//Проверка на наличие дней для повторений
				if len(monthDays) < 2 {
					err5 := errors.New("не верный день недели")
					fmt.Println(err5)
					return err5
				}

		*/

	default:
		err8 := errors.New("неизвестное значение для повторений")
		fmt.Println(err8)
		return err8
	}

	return nil
}

type TaskData struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeate string `json:"repeat,omitempty"`
}

func TaskDataValidation(httpData []byte) (TaskData, error) {

	var taskData TaskData
	now := time.Now()

	//Приводим текущее время к формату
	nowString := now.Format("20060102")
	now, err := time.Parse("20060102", nowString)
	if err != nil {
		fmt.Println("Строковые данные даты не корректны. ", err)
		return taskData, err
	}

	//Проверка распаковки JSON
	err = json.Unmarshal(httpData, &taskData)
	if err != nil {
		fmt.Println("не удачно распаковался JSON запрос ", err)
		return taskData, err
	}

	//Проверка id
	if taskData.ID != "" {
		_, err := IDValidation(taskData.ID)
		if err != nil {
			fmt.Println("не верый формат id ", err)
			return taskData, err
		}
	}

	//Проверка заголовка
	if taskData.Title == "" {
		err0 := errors.New("заголовок задачи отсутствует")
		return taskData, err0
	}

	//Проверка значений повторений
	if RepeatValidation(taskData.Repeate) != nil {
		return taskData, RepeatValidation(taskData.Repeate)
	}

	//Проверка на пустое значение даты
	if taskData.Date == "" {
		taskData.Date = fmt.Sprint(now.Format("20060102"))
	}

	//Проверка на корректность даты
	date, err := DateValidation(taskData.Date)
	if err != nil {
		fmt.Println(taskData.Date, err)
		return taskData, err
	}
	taskData.Date = fmt.Sprint(date.Format("20060102"))

	//Если дата меньше сегодняшей
	//Значение переданной даты. Ошибку игнорируем по причине проверки выше.
	dateIn, _ := time.Parse("20060102", taskData.Date)

	if now.After(dateIn) {
		if taskData.Repeate == "" {
			taskData.Date = fmt.Sprint(now.Format("20060102"))
		} else {
			taskData.Date, err = NextDate(now, taskData.Date, taskData.Repeate)
			if err != nil {
				return taskData, err
			}
		}
	}

	return taskData, nil
}
