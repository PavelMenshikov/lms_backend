-- +goose Up
-- +goose StatementBegin

CREATE INDEX IF NOT EXISTS idx_lessons_course_id_order ON lessons(course_id, order_num);
CREATE INDEX IF NOT EXISTS idx_lessons_module_id ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_tests_lesson_id ON tests(lesson_id);
CREATE INDEX IF NOT EXISTS idx_projects_lesson_id ON projects(lesson_id);
CREATE INDEX IF NOT EXISTS idx_users_role_city ON users(role, city);
CREATE INDEX IF NOT EXISTS idx_teacher_reviews_teacher_id ON teacher_reviews(teacher_id);
CREATE INDEX IF NOT EXISTS idx_user_courses_user_course ON user_courses(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_user_assignments_submission_lookup ON user_assignments_submission(user_id, assignment_id);

-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS idx_lessons_course_id_order;
DROP INDEX IF EXISTS idx_lessons_module_id;
DROP INDEX IF EXISTS idx_tests_lesson_id;
DROP INDEX IF EXISTS idx_projects_lesson_id;
DROP INDEX IF EXISTS idx_users_role_city;
DROP INDEX IF EXISTS idx_teacher_reviews_teacher_id;
DROP INDEX IF EXISTS idx_user_courses_user_course;
DROP INDEX IF EXISTS idx_user_assignments_submission_lookup;