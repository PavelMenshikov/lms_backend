CREATE TYPE user_role AS ENUM ('student', 'parent', 'teacher', 'curator', 'moderator', 'admin');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(120) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'student',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    phone VARCHAR(20),
    city VARCHAR(100),
    language VARCHAR(10) DEFAULT 'ru',
    gender VARCHAR(10),
    birth_date TIMESTAMP,
    school_name VARCHAR(255),
    experience_years INTEGER,
    whatsapp_link VARCHAR(255),
    telegram_link VARCHAR(255),
    avatar_url VARCHAR(500)
);

CREATE INDEX idx_users_role ON users(role);