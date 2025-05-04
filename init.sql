-- Создание таблицы Sessions, если она не существует
CREATE TABLE IF NOT EXISTS Sessions (
    session_id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    ip VARCHAR(45),
    user_agent VARCHAR(512),
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expire_at TIMESTAMP NOT NULL
);

ALTER TABLE sessions ADD COLUMN "updated_at" timestamp NULL DEFAULT now();

CREATE OR REPLACE FUNCTION update_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER sessions_updated BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE PROCEDURE update_column();