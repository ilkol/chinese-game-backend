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

- `GET api/levels` - Получение списка всех уровней