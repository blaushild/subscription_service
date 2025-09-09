FROM golang:1.24-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o subscription ./cmd/main.go

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest
WORKDIR /app

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/subscription .
COPY --from=builder /app/.env .

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./subscription"]
