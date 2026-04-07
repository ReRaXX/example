# User API (Go)

RESTful API для управления пользователями на Go с использованием Gin, GORM и PostgreSQL.

## Возможности

- Аутентификация с JWT токенами (access 15 мин, refresh 7 дней)
- Ролевой доступ (admin/user)
- Управление профилем пользователя
- Блокировка/разблокировка аккаунтов (только admin)
- Хеширование паролей bcrypt

## Технологии

- Gin, GORM, PostgreSQL
- JWT, bcrypt, validator

## Быстрый запуск

### Локально

```bash
go mod download
cp .env.example .env
# настройте .env
go run cmd/server/main.go
```

### Docker

```bash
docker-compose up -d
```

Сервер запустится на `http://localhost:8080`

## API Эндпоинты

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/auth/register` | Регистрация |
| POST | `/api/auth/login` | Вход |
| POST | `/api/auth/refresh` | Обновление токена |
| GET | `/api/users/me` | Мой профиль |
| PUT | `/api/users/me` | Обновить профиль |
| GET | `/api/users` | Все пользователи (admin) |
| GET | `/api/users/:id` | Пользователь по ID |
| PATCH | `/api/users/:id/toggle-block` | Блокировка (admin) |
| GET | `/health` | Проверка статуса |

## Пример запроса

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"fullName":"Иван Иванов","birthDate":"1990-01-01","email":"ivan@example.com","password":"password123"}'
```