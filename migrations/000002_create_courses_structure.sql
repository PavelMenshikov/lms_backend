CREATE TYPE course_status AS ENUM ('active', 'completed', 'archived');
CREATE TYPE assignment_status AS ENUM ('pending_check', 'rejected', 'accepted');


CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY,
    bio TEXT,
    rating DECIMAL(2,1) NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    is_main BOOLEAN NOT NULL DEFAULT FALSE, 
    status course_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

 Модулей
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    order_num INTEGER NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    UNIQUE (course_id, order_num)
);


CREATE TABLE IF NOT EXISTS lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID NOT NULL,
    teacher_id UUID NOT NULL, 
    title VARCHAR(255) NOT NULL,
    lesson_time TIMESTAMP WITH TIME ZONE NOT NULL, 
    duration_min INTEGER NOT NULL DEFAULT 60,
    online_url VARCHAR(255), 
    order_num INTEGER NOT NULL,
    FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE RESTRICT 
);