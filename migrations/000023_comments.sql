-- +goose Up
-- Комментарии куратор → учитель
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id),
    recipient_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    parent_comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_student ON comments(student_id);
CREATE INDEX idx_comments_lesson ON comments(lesson_id);
CREATE INDEX idx_comments_author ON comments(author_id);
CREATE INDEX idx_comments_recipient ON comments(recipient_id);
CREATE INDEX idx_comments_is_read ON comments(is_read);
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- +goose Down
DROP TABLE IF EXISTS comments;
