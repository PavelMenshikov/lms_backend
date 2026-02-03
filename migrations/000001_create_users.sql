DROP TABLE IF EXISTS users, teachers, courses, modules, lessons, user_courses, assignments, user_assignments_submission, child_parent_link CASCADE;

CREATE TYPE user_role AS ENUM ('student', 'parent', 'teacher', 'admin');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(120) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'student',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);