-- Таблица связей между учителем и учеником
CREATE TABLE IF NOT EXISTS teacher_students (
	teacher_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	student_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (teacher_id, student_id)
);

-- Добавление колонки для кода учителя
ALTER TABLE users ADD COLUMN invite_code VARCHAR(10) UNIQUE;
CREATE INDEX idx_invite_code ON users(invite_code);

CREATE TABLE admin_invites (
	code VARCHAR(10) PRIMARY KEY,
	admin_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	used_at TIMESTAMP WITH TIME ZONE
);