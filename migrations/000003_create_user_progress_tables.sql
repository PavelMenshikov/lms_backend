CREATE TABLE IF NOT EXISTS user_courses (
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    progress_percent INTEGER NOT NULL DEFAULT 0,
    modules_completed INTEGER NOT NULL DEFAULT 0,
    total_assignments_submitted INTEGER NOT NULL DEFAULT 0,
    
    PRIMARY KEY (user_id, course_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE RESTRICT
);


CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    max_score INTEGER NOT NULL DEFAULT 100,
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);