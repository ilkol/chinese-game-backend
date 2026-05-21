# Chinese Game Backend

Бэкэнд для web игры по изучению китайского языка. 


## Стек технологий
- Go
- БД PostgreSQL
- Docker, Docker Compose
- Роутер Chi

## Запуск
1. Создаем .env по шаблону .env.template
2. Запуск образа БД 
```bash
	docker-compose up -d
```
3. Запустите сервер
```bash
   go run cmd/server/main.go
```

## API endpoints
- `GET api/status` - Проверка работоспособности cервера

- `POST api/auth/register` - Запрос на регистрацию пользователя
- `POST api/auth/login` - Запрос на авторизацию пользователя

### запросы с авторизацией

- `GET api/level` - Получение списка всех уровней
- `GET api/level/{id}` - Получение уровень по ID со списком этапов уровней

- `POST api/progress` - Обновляет статус прохождения этапа пользователем

- `POST api/user/join` - Запрос на привязку ученика к учителю по его коду

- `GET api/teacher/students` - Запрос всех студентов, привязанных к учителю

