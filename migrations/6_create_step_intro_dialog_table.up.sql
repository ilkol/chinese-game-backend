-- Диалоги для этапов уровня
-- Один диалог на один этап (1:1 через step_id как PK)
CREATE TABLE step_dialogs (
    step_id INTEGER PRIMARY KEY REFERENCES level_steps(id) ON DELETE CASCADE,
    steps JSONB NOT NULL DEFAULT '[]',
    -- Пример структуры steps:
    -- [
    --   {
    --     "speaker": "Лун-Лун",
    --     "text": "Привет! Сегодня мы изучим тоны.",
    --     "emotion": "/assets/chars/lun-lun/happy.png",
    --     "bg": "purple"
    --   }
    -- ]
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);