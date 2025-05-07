# Используем golang как базовый образ для сборки
FROM golang:1.24.0 AS builder

# Устанавливаем необходимые пакеты
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libssl-dev \
    openssl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем файлы зависимостей и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Создаем директорию для сертификатов и генерируем их
RUN mkdir certs && \
    openssl genrsa -out certs/jwt-private.pem 2048 && \
    openssl rsa -in certs/jwt-private.pem -pubout -out certs/jwt-public.pem

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simpleAuth .

# Создаем финальный образ
FROM alpine:3.20

WORKDIR /app

# Копируем скомпилированное приложение и сертификаты из образа сборки
COPY --from=builder /app/simpleAuth .
COPY --from=builder /app/certs ./certs

# Делаем приложение исполняемым
RUN chmod +x ./simpleAuth

# Указываем команду для запуска приложения
CMD ["./simpleAuth"]
