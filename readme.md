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
- `GET api/progress/levels` - Запрос списка доступных пользователю уровней
- `GET api//progress/levels/{level_id}` - Запрос информации пройден ли конкретный уровень или нет

- `POST api/user/join` - Запрос на привязку ученика к учителю по его коду

- `GET api/teacher/students` - Запрос всех студентов, привязанных к учителю
- `GET api//teacher/invite-code` - Запрос кода-приглашения учеников

- `POST api/level/{id}/step` - Запрос создания нового этапа в уровне
- `PUT api/level/{id}/step/{step_id} ` - Запрос на обновление данных этапа уровня
- `DELETE api/level/{id}/step/{step_id}` - Запрос на удаление данных об этапе уровня
- `PUT api/step/{step_id}/dialog` - Запрос на обновление предварительного диалога этапа