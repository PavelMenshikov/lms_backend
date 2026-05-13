-- +goose Up
-- Баннеры для отображения на главной странице
CREATE TABLE IF NOT EXISTS banners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'INFO',
    -- Типы: INFO, WARNING, ANNOUNCEMENT, PROMOTION
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    target_roles TEXT[],
    -- Массив ролей: ['student', 'teacher', 'curator']
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_banner_dates CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_banners_is_active ON banners(is_active);
CREATE INDEX idx_banners_priority ON banners(priority DESC);
CREATE INDEX idx_banners_dates ON banners(start_date, end_date);
CREATE INDEX idx_banners_type ON banners(type);

-- +goose Down
DROP TABLE IF EXISTS banners;
