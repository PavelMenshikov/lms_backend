-- +goose Up
CREATE TYPE assignment_status AS ENUM ('pending_check', 'rejected', 'accepted');

CREATE TABLE IF NOT EXISTS user_lesson_attendance (
    user_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    is_attended BOOLEAN NOT NULL,
    comment_teacher TEXT,
    
    PRIMARY KEY (user_id, lesson_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_assignments_submission (
    user_id UUID NOT NULL,
    assignment_id UUID NOT NULL,
    submission_text TEXT,
    submission_link VARCHAR(500),
    status assignment_status NOT NULL DEFAULT 'pending_check',
    grade INTEGER,
    teacher_comment TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (user_id, assignment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE
);