run:

		GOOS=linux GOARCH=amd64 go run cmd/todo/main.go
		#TODO_DBFILE="/home/anton/go_final_project/scheduler.db" - для проверки с БД в корне проекта

build_win64:
		rm -rf OS_bin/TODO_win64/sqlite
		rm -rf OS_bin/TODO_win64/web
		mkdir OS_bin/TODO_win64/sqlite
		cp sqlite/scheduler_creator.sql OS_bin/TODO_win64/sqlite
		cp -r ./web ./OS_bin/TODO_win64/
		GOOS=windows GOARCH=amd64 TODO_DBFILE="" go build -o OS_bin/TODO_win64/TODO_windows64.exe cmd/todo/main.go

build_lin64:
		rm -rf OS_bin/TODO_lin64/sqlite
		rm -rf OS_bin/TODO_lin64/web
		mkdir OS_bin/TODO_lin64/sqlite
		cp sqlite/scheduler_creator.sql OS_bin/TODO_lin64/sqlite
		cp -r ./web ./OS_bin/TODO_lin64/
		GOOS=linux GOARCH=amd64 TODO_DBFILE="" go build -o OS_bin/TODO_lin64/TODO_linux cmd/todo/main.go

build_mac64:
		rm -rf OS_bin/TODO_mac64/sqlite
		rm -rf OS_bin/TODO_mac64/web
		mkdir OS_bin/TODO_mac64/sqlite
		cp sqlite/scheduler_creator.sql OS_bin/TODO_mac64/sqlite
		cp -r ./web ./OS_bin/TODO_mac64/
		GOOS=darwin GOARCH=arm64 TODO_DBFILE="" go build -o OS_bin/TODO_mac64/TODO_mac cmd/todo/main.go

clean:
		go clean -cache

test1:
		go test -run ^TestApp ./tests

test2:
		go test -run ^TestDB ./tests

test3:
		go test -run ^TestNextDate ./tests

test4:
		go test -run ^TestAddTask ./tests

test5:
		go test -run ^TestTasks ./tests

test6:
		go test -run ^TestEditTask ./tests

test7:
		go test -run ^TestDone ./tests

test8:
		go test -run ^TestDelTask ./tests

test:
		go test -run ^TestApp ./tests
		go test -run ^TestDB ./tests
		go test -run ^TestNextDate ./tests
		go test -run ^TestAddTask ./tests
		go test -run ^TestTasks ./tests
		go test -run ^TestEditTask ./tests
		go test -run ^TestDone ./tests
		go test -run ^TestDelTask ./tests

docker_build:
		docker build --tag todorun:v1 .

docker_run:
		docker run -d -p 7540:7540 todorun:v1

docker_run_db:
		docker run --rm -it -p 7540:7540 -v //$(PWD)/scheduler.db:/app/bin/scheduler.db gfpbig:v1.0.15


#Команды не запускаются из makefile по причине наличия $
#Но пусть будут тут для удобства копирования
#stop:
#		docker stop $(docker ps -a -q)

#s_r:
#		docker rm $(docker ps -a -q)
