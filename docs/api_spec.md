# Спецификация API — Cap Education LMS

## Базовый URL

```
Production: https://platform.capedu.kz
```

## Аутентификация

### Формат токена (plaintext)

```
userID:role
```

Пример: `550e8400-e29b-41d4-a716-446655440000:admin`

### Как передавать

| Метод | Header / Cookie |
|---|---|
| **Cookie** (автоматически после входа) | `auth_token=<token>` |
| **Header** (вручную) | `Authorization: Bearer <token>` |
| **Ответ при входе** | JSON `{"token": "<token>"}` (фронтенд сохраняет) |

### Вход в систему

```http
POST /auth/login
Content-Type: application/json

{
  "email": "admin@capedu.kz",
  "password": "demo123456"
}
```

Ответ:

```json
{
  "token": "550e8400-...:admin",
  "user": {
    "id": "550e8400-...",
    "first_name": "Admin",
    "last_name": "User",
    "email": "admin@capedu.kz",
    "role": "admin"
  }
}
```

### Роли

| Роль | Уровень доступа |
|---|---|
| `admin` | Полный доступ — все эндпоинты |
| `teacher` | Персонал — проверка ДЗ, посещаемость, свой дашборд |
| `curator` | Персонал — группы, заморозки/доступ, **не может** проверять ДЗ |
| `moderator` | Персонал — та же группа что admin/teacher/curator |
| `student` | Ограниченный — свои курсы, уроки, задания |
| `parent` | Только данные — нет API эндпоинтов |
| `superadmin` | Не используется |

### Группы middleware

| Группа | Требования | Маршруты |
|---|---|---|
| **Публичные** | Нет | `/auth/*`, `/swagger/*` |
| **Группа 1** | Auth + Роль ∈ {admin, teacher, moderator, curator} | `/admin/*`, `/api/admin/*`, `/staff/*`, `/api/groups/*`, `/api/attendance/*`, `/api/freeze-requests/*`, `/api/comments`, `/api/notifications`, `/api/access-requests`, `/api/statistics/*`, `/api/reports/*` |
| **Группа 2** | Auth (любая роль) | `/dashboard/home`, `/my-courses`, `/courses/{id}`, `/lessons/{id}`, `/profile`, `/schedule/*`, `/chat/*`, `/teachers/*`, `/api/notifications`, `/api/banner/active` |
| **Системный** | Header `X-System-Secret` | `/system/reset-password` |

---

## Admin

Все эндпоинты требуют роль `admin` (или персонал, где указано).

### Дашборд

#### Получить статистику дашборда

```http
GET /admin/dashboard/stats
Authorization: Bearer <token>
```

Ответ:

```json
{
  "total_students": 150,
  "total_teachers": 12,
  "total_curators": 5,
  "active_courses": 8,
  "revenue": 1250000,
  "performance_zones": {
    "green": 80,
    "yellow": 45,
    "red": 25
  },
  "lesson_activity": [
    { "date": "2026-06-01", "lessons": 15 },
    { "date": "2026-06-02", "lessons": 12 }
  ]
}
```

---

### Курсы

#### Список всех курсов

```http
GET /admin/courses
Authorization: Bearer <token>
```

#### Создать курс

```http
POST /admin/courses
Authorization: Bearer <token>
Content-Type: multipart/form-data

title=Python Basics
description=Intro course
image=@cover.jpg
```

#### Получить структуру курса

```http
GET /admin/courses/{courseId}/structure
Authorization: Bearer <token>
```

Ответ:

```json
{
  "course": { "id": "...", "title": "Python Basics" },
  "modules": [
    {
      "id": "...",
      "title": "Module 1",
      "order_num": 1,
      "lessons": [
        { "id": "...", "title": "Variables", "order_num": 1, "duration_min": 45 },
        { "id": "...", "title": "Functions", "order_num": 2, "duration_min": 60 }
      ],
      "tests": [],
      "projects": []
    }
  ]
}
```

#### Обновить настройки курса

```http
PUT /admin/courses/{courseId}/settings
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Python Advanced",
  "status": "active",
  "has_homework": true,
  "has_test": true,
  "has_project": false,
  "has_discord": true,
  "disable_copy_paste": false,
  "cover_image_url": "https://...",
  "teacher_ids": ["uuid-1", "uuid-2"]
}
```

#### Список студентов курса

```http
GET /admin/courses/{courseId}/students
Authorization: Bearer <token>
```

#### Статистика курса

```http
GET /admin/courses/{courseId}/stats
Authorization: Bearer <token>
```

#### Массовое создание курса (с модулями и уроками)

```http
POST /admin/courses/bulk
Authorization: Bearer <token>
Content-Type: application/json

{
  "course": { "title": "Python", "description": "..." },
  "modules": [
    {
      "title": "Module 1",
      "lessons": [
        { "title": "Intro", "duration_min": 30, "content_text": "# Hello" }
      ]
    }
  ]
}
```

#### Отчислить студента

```http
DELETE /admin/courses/{courseId}/enroll/{userId}
Authorization: Bearer <token>
```

---

### Модули

#### Создать модуль

```http
POST /admin/modules
Authorization: Bearer <token>
Content-Type: application/json

{
  "course_id": "uuid",
  "title": "Module 1",
  "order_num": 1
}
```

#### Удалить модуль

```http
DELETE /admin/modules/{moduleId}
Authorization: Bearer <token>
```

#### Массовое создание модулей

```http
POST /admin/modules/bulk
Authorization: Bearer <token>
Content-Type: application/json

{
  "course_id": "uuid",
  "modules": [
    { "title": "Module 1", "order_num": 1 },
    { "title": "Module 2", "order_num": 2 }
  ]
}
```

---

### Уроки

#### Создать урок (JSON)

```http
POST /admin/lessons
Authorization: Bearer <token>
Content-Type: application/json

{
  "module_id": "uuid",
  "title": "Variables",
  "order_num": 1,
  "duration_min": 45,
  "content_text": "# Variables\n\nIn Python...",
  "video_url": "https://...",
  "presentation_url": "https://..."
}
```

#### Создать урок (multipart с видео/презентацией)

```http
POST /admin/lessons
Authorization: Bearer <token>
Content-Type: multipart/form-data

module_id=uuid
title=Variables
order_num=1
duration_min=45
video=@lesson.mp4
presentation=@slides.pdf
```

#### Получить урок (для редактора)

```http
GET /admin/lessons/{lessonId}
Authorization: Bearer <token>
```

#### Обновить урок

```http
PUT /admin/lessons/{lessonId}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "content_text": "# New content",
  "video_url": "https://..."
}
```

#### Удалить урок

```http
DELETE /admin/lessons/{lessonId}
Authorization: Bearer <token>
```

#### Отменить урок

```http
POST /admin/lessons/{lessonId}/cancel
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "Teacher is sick"
}
```

#### Замена учителя

```http
POST /admin/lessons/{lessonId}/substitute
Authorization: Bearer <token>
Content-Type: application/json

{
  "teacher_id": "uuid"
}
```

#### Массовое создание уроков

```http
POST /admin/lessons/bulk
Authorization: Bearer <token>
Content-Type: application/json

{
  "module_id": "uuid",
  "lessons": [
    { "title": "Lesson 1", "order_num": 1, "duration_min": 30 },
    { "title": "Lesson 2", "order_num": 2, "duration_min": 45 }
  ]
}
```

---

### Тесты

#### Создать тест

```http
POST /admin/tests
Authorization: Bearer <token>
Content-Type: application/json

{
  "lesson_id": "uuid",
  "title": "Quiz 1",
  "questions": [
    {
      "question_text": "What is 2+2?",
      "options": ["3", "4", "5"],
      "correct_answer": 1
    }
  ]
}
```

#### Получить тест

```http
GET /admin/tests/{testId}
Authorization: Bearer <token>
```

#### Удалить тест

```http
DELETE /admin/tests/{testId}
Authorization: Bearer <token>
```

---

### Проекты

#### Создать проект

```http
POST /admin/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "lesson_id": "uuid",
  "title": "Final Project",
  "description": "Build a calculator"
}
```

#### Получить проект

```http
GET /admin/projects/{projectId}
Authorization: Bearer <token>
```

#### Удалить проект

```http
DELETE /admin/projects/{projectId}
Authorization: Bearer <token>
```

---

### Загрузка медиа

```http
POST /admin/media/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file=@image.jpg
```

Ответ:

```json
{
  "url": "https://storage.capedu.kz/uploads/uuid-image.jpg"
}
```

---

### Пользователи

#### Список пользователей (с фильтром)

```http
GET /admin/users?role=student&limit=50&offset=0
Authorization: Bearer <token>
```

#### Получить карточку пользователя

```http
GET /admin/users/{userId}
Authorization: Bearer <token>
```

Ответ включает: профиль, записанные курсы, родителей (если студент).

#### Создать пользователя

```http
POST /admin/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "secret123",
  "role": "student",
  "phone": "+77001112233",
  "city": "Almaty",
  "group_id": "uuid",
  "course_ids": ["uuid-1", "uuid-2"],
  "parents": [
    { "first_name": "Jane", "last_name": "Doe", "email": "jane@example.com", "phone": "+77001112244" }
  ]
}
```

#### Обновить пользователя

```http
PUT /admin/users/{userId}
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "John Updated",
  "phone": "+77009998877"
}
```

#### Удалить пользователя

```http
DELETE /admin/users/{userId}
Authorization: Bearer <token>
```

#### Зачислить студента

```http
POST /admin/enroll
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": "uuid",
  "course_id": "uuid",
  "group_id": "uuid"
}
```

#### Все пользователи (таблица)

```http
GET /admin/users/all
Authorization: Bearer <token>
```

#### Детальный список студентов

```http
GET /admin/students/detailed
Authorization: Bearer <token>
```

Ответ включает: курс, группа, куратор, учитель, успеваемость.

#### Детальный список учителей

```http
GET /admin/teachers/detailed
Authorization: Bearer <token>
```

#### Детальный список кураторов

```http
GET /admin/curators/detailed
Authorization: Bearer <token>
```

#### Детальный список модераторов

```http
GET /admin/moderators/detailed
Authorization: Bearer <token>
```

---

### Потоки и Группы

#### Создать поток

```http
POST /admin/streams
Authorization: Bearer <token>
Content-Type: application/json

{
  "course_id": "uuid",
  "title": "Spring 2026"
}
```

#### Список потоков

```http
GET /admin/streams?course_id=uuid
Authorization: Bearer <token>
```

#### Создать группу

```http
POST /admin/groups
Authorization: Bearer <token>
Content-Type: application/json

{
  "stream_id": "uuid",
  "title": "Group A",
  "teacher_id": "uuid"
}
```

#### Список групп

```http
GET /admin/groups?stream_id=uuid
Authorization: Bearer <token>
```

---

### Баннеры

#### Создать баннер

```http
POST /api/admin/banner
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "New Course",
  "content": "Python course starts soon",
  "type": "INFO",
  "is_active": true,
  "priority": 1,
  "start_date": "2026-06-20T00:00:00Z",
  "end_date": "2026-07-20T00:00:00Z",
  "target_roles": ["student", "teacher"]
}
```

#### Обновить баннер

```http
PATCH /api/admin/banner/{bannerId}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Banner",
  "is_active": false
}
```

#### Удалить баннер

```http
DELETE /api/admin/banner/{bannerId}
Authorization: Bearer <token>
```

---

### Группы (CRUD)

#### Обновить группу

```http
PATCH /api/groups/{groupId}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Group B",
  "teacher_id": "uuid"
}
```

#### Добавить студента в группу

```http
POST /api/groups/{groupId}/students
Authorization: Bearer <token>
Content-Type: application/json

{
  "student_id": "uuid"
}
```

#### Удалить студента из группы

```http
DELETE /api/groups/{groupId}/students/{studentId}
Authorization: Bearer <token>
```

#### Сменить группу студента

```http
PATCH /api/students/{studentId}/group
Authorization: Bearer <token>
Content-Type: application/json

{
  "group_id": "uuid"
}
```

#### Сменить группу учителя

```http
PATCH /api/teachers/{teacherId}/group
Authorization: Bearer <token>
Content-Type: application/json

{
  "group_id": "uuid"
}
```

---

### Посещаемость (Admin/Teacher/Curator)

#### Календарь посещаемости студента

```http
GET /api/attendance/students/{studentId}/calendar?start=2026-06-01&end=2026-06-30
Authorization: Bearer <token>
```

#### Отметить посещаемость

```http
PATCH /api/attendance/lessons/{lessonId}
Authorization: Bearer <token>
Content-Type: application/json

{
  "student_id": "uuid",
  "status": "ATTENDED",
  "reason": null,
  "comment": null
}
```

Статусы: `ATTENDED`, `ABSENT_EXCUSED`, `ABSENT_UNEXCUSED`, `FREEZE`

#### Статистика посещаемости студента

```http
GET /api/attendance/students/{studentId}/stats
Authorization: Bearer <token>
```

Ответ:

```json
{
  "attended": 12,
  "absent_excused": 2,
  "absent_unexcused": 1,
  "freeze": 0
}
```

#### Записи посещаемости урока

```http
GET /api/attendance/lessons/{lessonId}
Authorization: Bearer <token>
```

---

### Запросы на заморозку

#### Создать запрос на заморозку

```http
POST /api/freeze-requests
Authorization: Bearer <token>
Content-Type: application/json

{
  "student_id": "uuid",
  "start_date": "2026-07-01T00:00:00Z",
  "end_date": "2026-07-14T00:00:00Z",
  "reason": "Family vacation"
}
```

#### Получить ожидающие запросы

```http
GET /api/freeze-requests
Authorization: Bearer <token>
```

#### Одобрить заморозку

```http
PATCH /api/freeze-requests/{requestId}/approve
Authorization: Bearer <token>
Content-Type: application/json

{
  "review_comment": "Approved"
}
```

#### Отклонить заморозку

```http
PATCH /api/freeze-requests/{requestId}/reject
Authorization: Bearer <token>
Content-Type: application/json

{
  "review_comment": "Not eligible"
}
```

#### Статус заморозки студента

```http
GET /api/students/{studentId}/freeze-status
Authorization: Bearer <token>
```

---

### Запросы доступа

#### Создать запрос доступа

```http
POST /api/access-requests
Authorization: Bearer <token>
Content-Type: application/json

{
  "resource_type": "course",
  "resource_id": "uuid",
  "reason": "Need access to materials"
}
```

#### Получить ожидающие запросы

```http
GET /api/access-requests
Authorization: Bearer <token>
```

#### Одобрить доступ

```http
PATCH /api/access-requests/{requestId}/approve
Authorization: Bearer <token>
Content-Type: application/json

{
  "review_comment": "Approved"
}
```

#### Отклонить доступ

```http
PATCH /api/access-requests/{requestId}/reject
Authorization: Bearer <token>
Content-Type: application/json

{
  "review_comment": "Denied"
}
```

---

### Комментарии

#### Создать комментарий

```http
POST /api/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "student_id": "uuid",
  "lesson_id": null,
  "recipient_id": "uuid",
  "content": "Great progress this week!",
  "parent_comment_id": null
}
```

#### Получить комментарии студента

```http
GET /api/comments?student_id=uuid
Authorization: Bearer <token>
```

#### Отметить прочитанным

```http
PATCH /api/comments/{commentId}/read
Authorization: Bearer <token>
```

---

### Уведомления (управление)

#### Создать уведомление

```http
POST /api/notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "recipient_id": "uuid",
  "title": "Reminder",
  "content": "You have homework due tomorrow",
  "type": "INFO",
  "sender_id": "uuid",
  "link_url": "https://..."
}
```

---

### Статистика

#### Получить статистику студента

```http
GET /api/statistics/students/{studentId}
Authorization: Bearer <token>
```

Ответ включает: прогресс по курсу, % посещаемости, % выполнения ДЗ, средняя оценка.

#### Обновить статистику студента

```http
POST /api/statistics/students/{studentId}/refresh
Authorization: Bearer <token>
```

---

### Отчёты

#### Скачать Excel-отчёт по урокам

```http
GET /api/reports/lessons.xlsx
Authorization: Bearer <token>
```

---

### Дашборд куратора

```http
GET /admin/curator/dashboard
Authorization: Bearer <token>
```

Ответ:

```json
{
  "groups": [
    {
      "group_name": "Group A",
      "student_count": 15,
      "attendance_percent": 85.5,
      "homework_completion_percent": 72.3
    }
  ],
  "performance_zones": {
    "green": 10,
    "yellow": 3,
    "red": 2
  }
}
```

---

## Teacher

Все эндпоинты требуют роль `teacher`.

### Дашборд учителя

```http
GET /teacher/profile
Authorization: Bearer <token>
```

Ответ:

```json
{
  "profile": {
    "id": "...",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@school.com",
    "avatar_url": "https://..."
  },
  "assigned_courses": [
    { "id": "...", "title": "Python Basics", "progress_percent": 65 }
  ],
  "my_reviews": [
    { "id": "...", "student_name": "Alice", "rating": 5, "comment": "Great teacher!", "created_at": "..." }
  ],
  "substitutions": [
    { "id": "...", "title": "Math", "scheduled_at": "..." }
  ],
  "cancelled_lessons": [],
  "upcoming_lessons": []
}
```

### Месячный отчёт

```http
GET /teacher/monthly-report?year=2026&month=6
Authorization: Bearer <token>
```

Ответ:

```json
{
  "teacher_id": "uuid",
  "year": 2026,
  "month": 6,
  "total_lessons": 42,
  "substitutions_count": 3,
  "replaced_count": 1,
  "avg_rating": 4.7,
  "total_students": 25,
  "attendance_avg": 88.5
}
```

### Сохранить расписание учителя

```http
PUT /profile/teacher/schedule
Authorization: Bearer <token>
Content-Type: application/json

{
  "schedule": {
    "monday": ["10:00-11:00", "14:00-15:00"],
    "wednesday": ["10:00-11:00"]
  }
}
```

### Проверка ДЗ

#### Получить ожидающие работы

```http
GET /staff/submissions
Authorization: Bearer <token>
```

Ответ:

```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "student_name": "Alice Johnson",
    "course_title": "Python Basics",
    "module_order": 1,
    "lesson_order": 2,
    "lesson_title": "Functions",
    "submission_text": "def add(a, b): ...",
    "submission_files": ["hw.py"],
    "status": "pending",
    "grade": 0,
    "teacher_comment": "",
    "submitted_at": "2026-06-17T15:30:00Z"
  }
]
```

#### Оценить работу

```http
POST /staff/submissions/{submissionId}/evaluate
Authorization: Bearer <token>
Content-Type: application/json

{
  "grade": 80,
  "comment": "Good work, but missing error handling",
  "status": "accepted"
}
```

Допустимые оценки для `accepted`: `20`, `40`, `60`, `80`, `100`. Недопустимая оценка → `100`.  
Если `status` = `rejected`, оценка принудительно `0`.

> **Важно:** Куратор **не может** проверять — endpoint возвращает 403.

---

### Посещаемость (Teacher)

Те же эндпоинты что в разделе Admin: `GET/PATCH /api/attendance/*`.

### Комментарии (Teacher)

Те же что в Admin: `POST /api/comments`, `GET /api/comments`, `PATCH /api/comments/{id}/read`.

### Уведомления (Teacher)

Те же что в Admin: `POST /api/notifications`, `GET /api/notifications`, `PATCH /api/notifications/{id}/read`.

### Запросы заморозки/доступа (Teacher)

Как в Admin: создать, просмотреть ожидающие.
**Не может** одобрять/отклонять.

### Статистика (Teacher)

Как в Admin: `GET /api/statistics/students/{id}`, `POST /api/statistics/students/{id}/refresh`.

---

## Student

Все эндпоинты требуют авторизации (любая роль — Группа 2).

### Мои курсы

```http
GET /my-courses
Authorization: Bearer <token>
```

Ответ:

```json
[
  {
    "id": "uuid",
    "title": "Python Basics",
    "description": "Intro to Python",
    "image_url": "https://...",
    "progress_percent": 45,
    "is_main": true
  }
]
```

### Содержимое курса

```http
GET /courses/{courseId}
Authorization: Bearer <token>
```

Ответ:

```json
{
  "course": { "id": "...", "title": "Python Basics" },
  "modules": [
    {
      "id": "...",
      "title": "Module 1",
      "description": "...",
      "order_num": 1,
      "lessons": [
        {
          "id": "...",
          "title": "Variables",
          "order_num": 1,
          "duration_min": 45,
          "is_completed": true,
          "is_locked": false,
          "tests": [],
          "projects": []
        }
      ]
    }
  ],
  "root_lessons": [],
  "root_tests": [],
  "root_projects": []
}
```

### Детали урока

```http
GET /lessons/{lessonId}
Authorization: Bearer <token>
```

Ответ:

```json
{
  "lesson": {
    "id": "...",
    "title": "Variables",
    "content_text": "# Variables\n\nIn Python...",
    "video_url": "https://...",
    "video_duration": 900,
    "presentation_url": "https://...",
    "order_num": 1,
    "duration_min": 45
  },
  "materials": [
    { "id": "...", "title": "Cheatsheet", "s3_path": "...", "public_url": "https://..." }
  ],
  "previous_lesson_id": "",
  "next_lesson_id": "...",
  "is_completed": false,
  "attendance_status": "",
  "recording_url": "",
  "assignment_status": "pending",
  "teacher_comment": "",
  "grade": 0
}
```

### Отправить задание

```http
POST /lessons/{lessonId}/assignment
Authorization: Bearer <token>
Content-Type: multipart/form-data

text_answer=def add(a,b): return a+b
file=@homework.py
```

### Самостоятельная отметка посещаемости

```http
POST /lessons/{lessonId}/attendance
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "visited"
}
```

Статусы: `visited`, `missing_valid`, `missing_invalid`, `frozen`, `trial`

### Отзывы об учителях

#### Список учителей

```http
GET /teachers
Authorization: Bearer <token>
```

#### Детали учителя

```http
GET /teachers/{teacherId}
Authorization: Bearer <token>
```

#### Оставить отзыв

```http
POST /teachers/{teacherId}/reviews
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 5,
  "comment": "Excellent teacher!"
}
```

Оценка от 1 до 5.

### Просмотр теста

```http
GET /tests/{testId}
Authorization: Bearer <token>
```

### Просмотр проекта

```http
GET /projects/{projectId}
Authorization: Bearer <token>
```

---

## Общие (все авторизованные пользователи)

### Дашборд

```http
GET /dashboard/home
Authorization: Bearer <token>
```

Ответ (зависит от роли):

```json
{
  "last_lesson": { "id": "...", "title": "Variables", "scheduled_at": "..." },
  "active_courses": [
    { "id": "...", "title": "Python", "progress_percent": 45, "attendance_percent": 80 }
  ],
  "attendance_stats": {
    "total": 20,
    "attended": 16,
    "percent": 80
  },
  "assignment_stats": {
    "total": 10,
    "completed": 7,
    "percent": 70
  },
  "upcoming_lessons": [
    { "id": "...", "title": "Functions", "scheduled_at": "2026-06-20T10:00:00Z" }
  ]
}
```

### Профиль

#### Просмотр профиля

```http
GET /profile
Authorization: Bearer <token>
```

#### Обновление профиля

```http
PUT /profile
Authorization: Bearer <token>
Content-Type: multipart/form-data

first_name=John
last_name=Doe
phone=+77001112233
avatar=@avatar.jpg
```

### Расписание

#### На неделю

```http
GET /schedule/weekly
Authorization: Bearer <token>
```

#### На месяц

```http
GET /schedule/monthly?year=2026&month=6
Authorization: Bearer <token>
```

### Чат

#### WebSocket

```http
GET /chat/ws
Authorization: Bearer <token>
Connection: Upgrade
Upgrade: websocket
```

#### История чата

```http
GET /chat/history?limit=50&offset=0
Authorization: Bearer <token>
```

### Уведомления (свои)

#### Получить уведомления

```http
GET /api/notifications
Authorization: Bearer <token>
```

#### Отметить прочитанным

```http
PATCH /api/notifications/{notificationId}/read
Authorization: Bearer <token>
```

### Активные баннеры

```http
GET /api/banner/active
Authorization: Bearer <token>
```

Опциональный параметр: `?role=student` для фильтрации по роли.

---

## Системное (Backdoor — только для демо)

### Сброс пароля

```http
POST /system/reset-password
X-System-Secret: my_ultra_secret_backdoor_key_2026
Content-Type: application/json

{
  "user_id": "uuid",
  "new_password": "newpassword123"
}
```

> **⚠️ После демо:** удалить эндпоинт, сменить все секреты.

---

## Формат ошибок

Все ошибки:

```json
{
  "error": "человекочитаемое сообщение"
}
```

HTTP статусы:
| Код | Значение |
|---|---|
| 400 | Неверный запрос (невалидный JSON, пропущены поля) |
| 401 | Не авторизован (нет/невалидный токен) |
| 403 | Запрещено (неподходящая роль) |
| 404 | Ресурс не найден |
| 409 | Конфликт (дубликат) |
| 500 | Внутренняя ошибка сервера (без стектрейса) |

---

## Матрица фич

| Фича | Admin | Teacher | Curator | Moderator | Student |
|---|---|---|---|---|---|
| Статистика дашборда админа | ✅ | ✅ | ✅ | ✅ | ❌ |
| Дашборд куратора | ✅ | ❌ | ✅ | ✅ | ❌ |
| Дашборд учителя | ❌ | ✅ | ❌ | ❌ | ❌ |
| Домашний дашборд | ✅ | ✅ | ✅ | ✅ | ✅ |
| Месячный отчёт | ❌ | ✅ | ❌ | ❌ | ❌ |
| CRUD курсов | ✅ | ❌ | ❌ | ❌ | ❌ |
| Просмотр структуры курса | ✅ | ✅ | ✅ | ✅ | ❌ |
| CRUD модулей | ✅ | ❌ | ❌ | ❌ | ❌ |
| CRUD уроков | ✅ | ❌ | ❌ | ❌ | ❌ |
| Детали урока | ✅ | ✅ | ✅ | ✅ | ✅ |
| CRUD тестов/проектов | ✅ | ❌ | ❌ | ❌ | ❌ |
| Отмена/замена урока | ✅ | ❌ | ❌ | ❌ | ❌ |
| Загрузка медиа | ✅ | ❌ | ❌ | ❌ | ❌ |
| CRUD пользователей | ✅ | ❌ | ❌ | ❌ | ❌ |
| Детальные списки | ✅ | ❌ | ❌ | ❌ | ❌ |
| Зачисление/отчисление | ✅ | ❌ | ❌ | ❌ | ❌ |
| CRUD потоков/групп | ✅ | ❌ | ❌ | ❌ | ❌ |
| Управление группами | ✅ | ❌ | ❌ | ✅ | ❌ |
| Мои курсы | ❌ | ❌ | ❌ | ❌ | ✅ |
| Контент курса (студент) | ❌ | ❌ | ❌ | ❌ | ✅ |
| Отправить ДЗ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Отметить посещаемость (сам) | ❌ | ❌ | ❌ | ❌ | ✅ |
| Просмотр ожидающих ДЗ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Оценка ДЗ | ✅ | ✅ | ❌ (403) | ✅ | ❌ |
| Посещаемость (управление) | ✅ | ✅ | ✅ | ✅ | ❌ |
| Запросы заморозки | ✅ | ✅ | ✅ | ✅ | ❌ |
| Одобрение заморозки | ✅ | ❌ | ✅ | ✅ | ❌ |
| Запросы доступа | ✅ | ✅ | ✅ | ✅ | ❌ |
| Одобрение доступа | ✅ | ❌ | ✅ | ✅ | ❌ |
| Комментарии | ✅ | ✅ | ✅ | ✅ | ❌ |
| Создание уведомлений | ✅ | ✅ | ✅ | ✅ | ❌ |
| Свои уведомления | ✅ | ✅ | ✅ | ✅ | ✅ |
| Статистика | ✅ | ✅ | ✅ | ✅ | ❌ |
| Excel отчёты | ✅ | ✅ | ✅ | ✅ | ❌ |
| CRUD баннеров | ✅ | ❌ | ❌ | ❌ | ❌ |
| Просмотр баннеров | ✅ | ✅ | ✅ | ✅ | ✅ |
| Профиль | ✅ | ✅ | ✅ | ✅ | ✅ |
| Расписание учителя | ❌ | ✅ | ❌ | ❌ | ❌ |
| Расписание (нед/мес) | ✅ | ✅ | ✅ | ✅ | ✅ |
| Чат (WebSocket) | ✅ | ✅ | ✅ | ✅ | ✅ |
| Список/детали учителей | ✅ | ✅ | ✅ | ✅ | ✅ |
| Отзывы учителям | ❌ | ❌ | ❌ | ❌ | ✅ |

---

## Известные проблемы

- Куратор не может оценивать ДЗ (явный 403)
- `not_started` enum: исправлено `COALESCE(status::text, 'not_started')`
- NULL scan: все nullable колонки обёрнуты COALESCE (исправлено)
- Токен аутентификации — plaintext `userID:role` (без JWT, после демо)
