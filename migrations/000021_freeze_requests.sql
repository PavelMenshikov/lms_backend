-- +goose Up
-- Запросы на заморозку обучения
CREATE TABLE IF NOT EXISTS freeze_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    requested_by UUID NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    -- Статусы: PENDING, APPROVED, REJECTED
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP,
    review_comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_date_range CHECK (end_date >= start_date)
);

CREATE INDEX idx_freeze_requests_student ON freeze_requests(student_id);
CREATE INDEX idx_freeze_requests_status ON freeze_requests(status);
CREATE INDEX idx_freeze_requests_requested_by ON freeze_requests(requested_by);
CREATE INDEX idx_freeze_requests_created_at ON freeze_requests(created_at);

-- +goose Down
DROP TABLE IF EXISTS freeze_requests;
