## SSO сервис на Go
### SSO (Single Sign-On) сервис, реализованный на Go, для централизованной аутентификации и авторизации пользователей.

**Функциональность:**

- Регистрация и вход: Предоставляет API для регистрации новых пользователей и входа в систему.
- Аутентификация: Использует JWT (JSON Web Tokens) для аутентификации пользователей.
- Авторизация: Определяет права доступа пользователей к различным ресурсам.
- Генерация токенов: Выдает JWT-токены для авторизованных пользователей.
- Проверка токенов: Проверяет валидность JWT-токенов, предоставленных пользователями.

**Технологии:**

1. Go: Сервис написан на Go.
2. gRPC: Используется для создания API.
3. PostgreSQL: Используется для хранения информации о пользователях.
4. GOOSE: Используется для миграции базы данных PostgreSQL.
5. cleanenv: Используется для удобной загрузки конфигурации из различных источников (файлы, переменные окружения).
6. JWT: Используется для аутентификации пользователей.

**Требования:**

- Go: Установите Go на вашу систему.
- PostgreSQL: Установите PostgreSQL и создайте базу данных.
- Docker: (Рекомендуется) Установите Docker для простого запуска и развертывания.

**Установка:**

1. Клонируйте репозиторий:
```bash 
git clone https://github.com/your-username/sso-service.git
```
2. Установите зависимости:
```bash
go mod tidy
```

3. Запустите сервис:
```bash 
go run main.go
```
**P.S:
В файле Taskfile миграции и необходимые команды для запуска**
