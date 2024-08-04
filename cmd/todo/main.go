package main

import (
	"fmt"

	"go_final_project/handlers"
	"go_final_project/sqlite"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {

	port := os.Getenv("TODO_PORT")
	todoDB := os.Getenv("TODO_DBFILE")

	if port == "" {
		port = ":7540"
	}

	//Проверить существует ли файл БД. Если его нет, то создать БД.
	dbURL, err := sqlite.FindOrCreateDB(todoDB)
	if err != nil {
		fmt.Println("Ошибка с базой данных ", err)
	}

	//Подключаемся к БД
	db, err := sqlite.InitDB(dbURL)
	if err != nil {
		fmt.Println("Ошибка инициализации БД ", err)
	}

	storage := sqlite.NewStorage(db)

	//Создаём роутер
	r := chi.NewRouter()

	//Запускаем Web интерфейс
	r.Handle("/*", http.FileServer(http.Dir("./web")))

	//Выводим значение новой даты
	r.Get("/api/nextdate", handlers.GetNextDateHandler)

	//Работаем с одной задачей
	r.Post("/api/task", storage.PostOneTaskHandler)
	r.Get("/api/task", storage.GetOneTaskHandler)
	r.Put("/api/task", storage.PutOneTaskHandler)
	r.Post("/api/task/done", storage.DoneOneTaskHandler)
	r.Delete("/api/task", storage.DeleteOneTaskHandler)

	//Работа с группой задач
	r.Get("/api/tasks", storage.GetTasksHandler)

	//Запускаем сервер
	fmt.Printf("Сервер TODO запущен! Порт %s.\n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s.\n", err.Error())
	}
}

//TODO Составить Docker файл и проверить работоспособность в контейнере.
//TODO Очистить main.go от лишних пометок и артефактов.
//TODO * Аутентификация.
//TODO * Цикличность по месяцам.
//TODO * Дополнительная проверка тестами go test -run ^TestNextDate$ ./tests при значении "true" для FullNextDate.
