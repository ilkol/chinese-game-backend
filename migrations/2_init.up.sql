DROP TABLE IF EXISTS progress;
DROP TABLE IF EXISTS levels;

CREATE TABLE levels (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL, -- заголовок "планеты" (уровня)
	color VARCHAR(100) NOT NULL, -- цвет планеты на карте
	icon VARCHAR(10) NOT NULL, -- иконка, которая будет видна на карте
	order_index INT NOT NULL, -- порядок уровня в цепочке
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- тип этапа внутри уровня
CREATE TYPE step_type AS ENUM ('theory', 'quiz', 'final');

-- этапы в уровне
CREATE TABLE level_steps (
    id SERIAL PRIMARY KEY,
    level_id INTEGER REFERENCES levels(id) ON DELETE CASCADE,
    type step_type NOT NULL,
    title VARCHAR(255),
    content JSONB,
    order_index INT NOT NULL -- порядок этапа в цепочке
);

-- таблица этапов, которые прошел пользователь
CREATE TABLE user_progress (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    step_id INTEGER REFERENCES level_steps(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, step_id)
);