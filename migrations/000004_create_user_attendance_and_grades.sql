
CREATE TABLE IF NOT EXISTS user_assignments_submission (
    user_id UUID NOT NULL,
    assignment_id UUID NOT NULL,
    submission_text TEXT,
    submission_link VARCHAR(500),
    status assignment_status NOT NULL DEFAULT 'pending_check',
    grade INTEGER,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (user_id, assignment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE
);