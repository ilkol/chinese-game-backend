ALTER TABLE levels ADD COLUMN background_src VARCHAR(255);

-- Обновить существующий уровень
UPDATE levels SET background_src = '/assets/levels/1/bg.jpeg' WHERE id = 1;