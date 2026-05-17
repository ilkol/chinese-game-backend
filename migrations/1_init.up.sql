
-- Роль пользака
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');

-- Пользаки
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
	role user_role DEFAULT 'student',
	coins INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Уровни
CREATE TABLE IF NOT EXISTS levels (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL, -- заголовок "планеты" (уровня)
	color VARCHAR(100) NOT NULL, -- цвет планеты на карте
	icon VARCHAR(10) NOT NULL, -- иконка, которая будет видна на карте
    blocks JSONB DEFAULT '[]', -- этапы внутри уровня
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Прогресс прозождения уровней
CREATE TABLE IF NOT EXISTS progress (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    level_id INT REFERENCES levels(id) ON DELETE CASCADE,
    details JSONB DEFAULT '{"completed_blocks": []}',
    is_completed BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, level_id)
);