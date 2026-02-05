CREATE TYPE lesson_type AS ENUM ('group', 'trial', 'individual');
ALTER TABLE lessons ADD COLUMN type lesson_type NOT NULL DEFAULT 'group';


CREATE TABLE IF NOT EXISTS lesson_view_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    view_duration_seconds INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);