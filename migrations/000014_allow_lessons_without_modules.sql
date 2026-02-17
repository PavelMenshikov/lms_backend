-- +goose Up
ALTER TABLE lessons ADD COLUMN course_id UUID;

UPDATE lessons 
SET course_id = (SELECT course_id FROM modules WHERE modules.id = lessons.module_id)
WHERE module_id IS NOT NULL;
ALTER TABLE lessons ALTER COLUMN course_id SET NOT NULL;

ALTER TABLE lessons 
ADD CONSTRAINT fk_lessons_course 
FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;

ALTER TABLE lessons ALTER COLUMN module_id DROP NOT NULL;