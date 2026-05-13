-- +goose Up
-- Статистика ученика (кеш для быстрого доступа)
CREATE TABLE IF NOT EXISTS student_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    total_lessons INTEGER DEFAULT 0,
    attended_lessons INTEGER DEFAULT 0,
    absent_excused INTEGER DEFAULT 0,
    absent_unexcused INTEGER DEFAULT 0,
    freeze_days INTEGER DEFAULT 0,
    remaining_lessons INTEGER DEFAULT 0,
    remaining_excused INTEGER DEFAULT 0,
    last_attendance_date DATE,
    current_freeze_end_date DATE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_student_statistics_student ON student_statistics(student_id);
CREATE INDEX idx_student_statistics_updated_at ON student_statistics(updated_at);

-- Функция для автоматического обновления статистики
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_student_statistics()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO student_statistics (student_id, updated_at)
    VALUES (NEW.student_id, CURRENT_TIMESTAMP)
    ON CONFLICT (student_id)
    DO UPDATE SET updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Триггер на изменение посещаемости
CREATE TRIGGER trigger_update_statistics_on_attendance
AFTER INSERT OR UPDATE ON attendance_records
FOR EACH ROW
EXECUTE FUNCTION update_student_statistics();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_statistics_on_attendance ON attendance_records;
DROP FUNCTION IF EXISTS update_student_statistics();
DROP TABLE IF EXISTS student_statistics;
