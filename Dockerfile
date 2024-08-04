FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./

COPY web/  ./bin/web
COPY sqlite/scheduler_creator.sql ./bin/sqlite

RUN go mod tidy

#Создаёт исполняемый файл в среде, реализованной в контейнере
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 TODO_DBFILE="" go build -o /app/bin/todorun ./cmd/todo/main.go

#Запускает исполняемый файл
CMD ["/app/bin/todorun"]
