# Сервис аутентификации

Этот проект представляет собой часть сервиса аутентификации реализованного на Go

## ✅ Основной функционал
- Выдача токенов
- Обновление токенов
- Получение текущего пользователя
- Деавторизация

## 🔐 Безопасность
- Access токен не хранится
- Refresh токены хранятся только в виде bcrypt
- Обновление токенов возможно только с тем же User-Agent
- Отправка уведомления о смене IP

## Установка и запуск
1. Клонируйте репозиторий:
```bash
git clone git@github.com:ivan-karaberov/simple-auth
cd simple-auth
```

2. Установите зависимости:
```bash
go mod tidy
```

3. Заполните `.env` файл:

4. Сгенерируйте сертификаты
```bash
mkdir certs && \
    openssl genrsa -out certs/jwt-private.pem 2048 && \
    openssl rsa -in certs/jwt-private.pem -pubout -out certs/jwt-public.pem
```

5. Запустите сервер:
```bash
go run main.go
```

## Установка и запуск (Docker)
```bash
docker compose up
```

## Документация API
Полная документация, описывающая структуру доступных маршрутов, форматы запросов и ответов, а также возможные коды ошибок, автоматически генерируется с использованием Swagger и доступна по следующему URL:
```bash
http://localhost:3000/docs/index.html
```