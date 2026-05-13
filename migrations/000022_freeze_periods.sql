-- +goose Up
-- Активные периоды заморозок (одобренные запросы)
CREATE TABLE IF NOT EXISTS freeze_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    freeze_request_id UUID REFERENCES freeze_requests(id) ON DELETE SET NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_freeze_period CHECK (end_date >= start_date)
);

CREATE INDEX idx_freeze_periods_student ON freeze_periods(student_id);
CREATE INDEX idx_freeze_periods_active ON freeze_periods(is_active);
CREATE INDEX idx_freeze_periods_dates ON freeze_periods(start_date, end_date);

-- +goose Down
DROP TABLE IF EXISTS freeze_periods;
