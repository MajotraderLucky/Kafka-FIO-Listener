# Используйте официальный образ Go как базовый
FROM golang:1.21 as builder

# Установите рабочий каталог
WORKDIR /app

# Копируйте исходный код в контейнер
COPY . .

# Скачивайте зависимости
RUN go mod download

# Собирайте приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o create_topics create_topics.go

# Используйте scratch для минимального образа
FROM scratch

# Копируйте скомпилированный бинарный файл из предыдущего шага
COPY --from=builder /app/create_topics .

# Запускайте программу
ENTRYPOINT ["./create_topics"]
