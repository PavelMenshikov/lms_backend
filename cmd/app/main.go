package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"

	authHttp "lms_backend/internal/auth/delivery/http"
	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/auth/repository"
	authUseCase "lms_backend/internal/auth/usecase"

	dashboardHttp "lms_backend/internal/dashboard/delivery/http"
	dashboardRepo "lms_backend/internal/dashboard/repository"
	dashboardUseCase "lms_backend/internal/dashboard/usecase"

	contentAdminHttp "lms_backend/internal/content_admin/delivery/http"
	contentAdminRepo "lms_backend/internal/content_admin/repository"
	contentAdminUseCase "lms_backend/internal/content_admin/usecase"

	learningHttp "lms_backend/internal/learning/delivery/http"
	learningRepo "lms_backend/internal/learning/repository"
	learningUseCase "lms_backend/internal/learning/usecase"

	teacherDashboardHttp "lms_backend/internal/teacher_dashboard/delivery/http"
	teacherDashboardRepo "lms_backend/internal/teacher_dashboard/repository"
	teacherDashboardUseCase "lms_backend/internal/teacher_dashboard/usecase"

	reviewHttp "lms_backend/internal/review/delivery/http"
	reviewRepo "lms_backend/internal/review/repository"
	reviewUseCase "lms_backend/internal/review/usecase"

	profileHttp "lms_backend/internal/profile/delivery/http"
	profileRepo "lms_backend/internal/profile/repository"
	profileUseCase "lms_backend/internal/profile/usecase"

	scheduleHttp "lms_backend/internal/schedule/delivery/http"
	scheduleRepo "lms_backend/internal/schedule/repository"
	scheduleUseCase "lms_backend/internal/schedule/usecase"

	chatHttp "lms_backend/internal/chat/delivery/http"
	chatRepo "lms_backend/internal/chat/repository"
	chatUseCase "lms_backend/internal/chat/usecase"

	attendanceHttp "lms_backend/internal/attendance/delivery/http"
	attendanceRepo "lms_backend/internal/attendance/repository"
	attendanceUseCase "lms_backend/internal/attendance/usecase"

	freezeHttp "lms_backend/internal/freeze/delivery/http"
	freezeRepo "lms_backend/internal/freeze/repository"
	freezeUseCase "lms_backend/internal/freeze/usecase"

	commentHttp "lms_backend/internal/comment/delivery/http"
	commentRepo "lms_backend/internal/comment/repository"
	commentUseCase "lms_backend/internal/comment/usecase"

	notificationHttp "lms_backend/internal/notification/delivery/http"
	notificationRepo "lms_backend/internal/notification/repository"
	notificationUseCase "lms_backend/internal/notification/usecase"

	accessHttp "lms_backend/internal/access/delivery/http"
	accessRepo "lms_backend/internal/access/repository"
	accessUseCase "lms_backend/internal/access/usecase"

	bannerHttp "lms_backend/internal/banner/delivery/http"
	bannerRepo "lms_backend/internal/banner/repository"
	bannerUseCase "lms_backend/internal/banner/usecase"

	auditRepo "lms_backend/internal/audit/repository"
	auditUseCase "lms_backend/internal/audit/usecase"

	statisticsHttp "lms_backend/internal/statistics/delivery/http"
	statisticsRepo "lms_backend/internal/statistics/repository"
	statisticsUseCase "lms_backend/internal/statistics/usecase"

	groupsHttp "lms_backend/internal/groups/delivery/http"
	groupsRepo "lms_backend/internal/groups/repository"
	groupsUseCase "lms_backend/internal/groups/usecase"

	reportsService "lms_backend/internal/reports"
	reportsHttp "lms_backend/internal/reports/delivery/http"

	"lms_backend/internal/domain"
	"lms_backend/internal/httperror"
	dbPkg "lms_backend/pkg/database"
	"lms_backend/pkg/logger"
	storageService "lms_backend/pkg/storage"
)

// @title Cap Education LMS - API
// @version 1.0
// @description API для LMS платформы Cap Education.
// @host platform.capedu.kz
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file, using system envs")
	}

	logger.Init(os.Getenv("APP_ENV"))

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := "host=" + host + " port=" + port + " user=" + user +
		" password=" + password + " dbname=" + dbname + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("failed to open db", logger.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	slog.Info("connected to postgresql")

	if err := dbPkg.RunMigrations(db); err != nil {
		slog.Error("Critical Error running migrations", logger.Err(err))
		os.Exit(1)
	}

	s3Client, err := storageService.NewS3Client(
		os.Getenv("S3_ENDPOINT_URL"),
		os.Getenv("S3_REGION"),
		os.Getenv("S3_BUCKET_NAME"),
		os.Getenv("S3_ACCESS_KEY_ID"),
		os.Getenv("S3_SECRET_ACCESS_KEY"),
	)
	if err != nil {
		slog.Error("Failed to initialize S3 Storage", logger.Err(err))
		os.Exit(1)
	}

	authRepoImpl := repository.NewAuthRepository(db)
	authUsecase := authUseCase.NewAuthUsecase(authRepoImpl)
	authHandler := authHttp.NewAuthHandler(authUsecase)

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	dashboardRepository := dashboardRepo.NewCachedDashboardRepo(
		dashboardRepo.NewDashboardRepository(db), rdb,
	)
	dashboardUsecase := dashboardUseCase.NewDashboardUseCase(dashboardRepository)
	dashboardHandler := dashboardHttp.NewDashboardHandler(dashboardUsecase)

	adminRepo := contentAdminRepo.NewContentAdminRepository(db)
	adminUsecase := contentAdminUseCase.NewContentAdminUseCase(adminRepo, s3Client)
	adminHandler := contentAdminHttp.NewContentAdminHandler(adminUsecase)

	learningRepoImpl := learningRepo.NewLearningRepository(db)
	learningUC := learningUseCase.NewLearningUseCase(learningRepoImpl, s3Client)
	learningHandler := learningHttp.NewLearningHandler(learningUC)

	teacherDashboardRepoImpl := teacherDashboardRepo.NewTeacherDashboardRepository(db)
	teacherDashboardUC := teacherDashboardUseCase.NewTeacherDashboardUseCase(teacherDashboardRepoImpl)
	teacherDashboardHandler := teacherDashboardHttp.NewTeacherDashboardHandler(teacherDashboardUC)

	reviewRepoImpl := reviewRepo.NewReviewRepository(db)
	reviewUC := reviewUseCase.NewReviewUseCase(reviewRepoImpl)
	reviewHandler := reviewHttp.NewReviewHandler(reviewUC)

	profileRepoImpl := profileRepo.NewProfileRepository(db)
	profileUC := profileUseCase.NewProfileUseCase(profileRepoImpl, s3Client)
	profileHandler := profileHttp.NewProfileHandler(profileUC)

	scheduleRepoImpl := scheduleRepo.NewScheduleRepository(db)
	scheduleUC := scheduleUseCase.NewScheduleUseCase(scheduleRepoImpl)
	scheduleHandler := scheduleHttp.NewScheduleHandler(scheduleUC)

	chatRepoImpl := chatRepo.NewChatRepository(db)
	chatUC := chatUseCase.NewChatUseCase(chatRepoImpl)
	chatHandler := chatHttp.NewChatHandler(chatUC)

	attendanceRepoImpl := attendanceRepo.NewAttendanceRepository(db)
	attendanceUC := attendanceUseCase.NewAttendanceUseCase(attendanceRepoImpl)
	attendanceHandler := attendanceHttp.NewAttendanceHandler(attendanceUC)

	freezeRepoImpl := freezeRepo.NewFreezeRepository(db)
	freezeUC := freezeUseCase.NewFreezeUseCase(freezeRepoImpl)
	freezeHandler := freezeHttp.NewFreezeHandler(freezeUC)

	commentRepoImpl := commentRepo.NewCommentRepository(db)
	commentUC := commentUseCase.NewCommentUseCase(commentRepoImpl)
	commentHandler := commentHttp.NewCommentHandler(commentUC)

	notificationRepoImpl := notificationRepo.NewNotificationRepository(db)
	notificationUC := notificationUseCase.NewNotificationUseCase(notificationRepoImpl)
	notificationHandler := notificationHttp.NewNotificationHandler(notificationUC)

	accessRepoImpl := accessRepo.NewAccessRepository(db)
	accessUC := accessUseCase.NewAccessUseCase(accessRepoImpl)
	accessHandler := accessHttp.NewAccessHandler(accessUC)

	bannerRepoImpl := bannerRepo.NewBannerRepository(db)
	bannerUC := bannerUseCase.NewBannerUseCase(bannerRepoImpl)
	bannerHandler := bannerHttp.NewBannerHandler(bannerUC)

	auditRepoImpl := auditRepo.NewAuditRepository(db)
	_ = auditUseCase.NewAuditUseCase(auditRepoImpl)

	statisticsRepoImpl := statisticsRepo.NewStatisticsRepository(db)
	statisticsUC := statisticsUseCase.NewStatisticsUseCase(statisticsRepoImpl)
	statisticsHandler := statisticsHttp.NewStatisticsHandler(statisticsUC)

	groupsRepoImpl := groupsRepo.NewGroupRepository(db)
	groupsUC := groupsUseCase.NewGroupUseCase(groupsRepoImpl)
	groupsHandler := groupsHttp.NewGroupHandler(groupsUC)

	reportsServiceImpl := reportsService.NewReportsService(db)
	reportsHandler := reportsHttp.NewReportsHandler(reportsServiceImpl)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://localhost:3000",
			"http://test.xcx.cx",
			"https://test.xcx.cx",
			"http://capapi.grtsq.ru",
			"https://capapi.grtsq.ru",
			"https://cap.grtsq.ru",
			"https://cap-education.vercel.app",
			"http://platform.capedu.kz",
			"https://platform.capedu.kz",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/logout", authHandler.Logout)
	r.Post("/auth/forgot-password", authHandler.ForgotPassword)
	r.Post("/auth/reset-password", authHandler.ResetPassword)

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/logout", authHandler.Logout)
	r.Post("/api/auth/forgot-password", authHandler.ForgotPassword)
	r.Post("/api/auth/reset-password", authHandler.ResetPassword)

	r.Post("/system/reset-password", func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SYSTEM_SECRET")
		if secret == "" || r.Header.Get("X-System-Secret") != secret {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
		if err != nil {
			http.Error(w, "Hash Error", http.StatusInternalServerError)
			return
		}

		res, err := db.Exec("UPDATE users SET password_hash=$1 WHERE email=$2", string(hash), body.Email)
		if err != nil {
			httperror.Internal(w, err)
			return
		}

		affected, _ := res.RowsAffected()
		if affected == 0 {
			_, err = db.Exec(`INSERT INTO users (id, first_name, last_name, email, password_hash, role)
				VALUES (gen_random_uuid(), 'Root', 'Admin', $1, $2, 'admin')`, body.Email, string(hash))
			if err != nil {
				httperror.Internal(w, err)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("User Created and Password Set"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password Updated Successfully"))
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware, authMiddleware.RoleRequiredMiddleware(domain.RoleAdmin, domain.RoleTeacher, domain.RoleModerator, domain.RoleCurator))

		r.Get("/admin/dashboard/stats", dashboardHandler.GetAdminDashboard)
		r.Get("/admin/curator/dashboard", dashboardHandler.GetCuratorDashboard)
		r.Get("/admin/courses", adminHandler.GetAllCourses)
		r.Post("/admin/courses", adminHandler.CreateCourse)
		r.Put("/admin/courses/{id}/settings", adminHandler.UpdateCourseSettings)
		r.Get("/admin/courses/{id}/structure", adminHandler.GetCourseStructure)
		r.Get("/admin/courses/{id}/students", adminHandler.GetCourseStudents)
		r.Get("/admin/courses/{id}/stats", adminHandler.GetCourseStats)
		r.Post("/admin/modules", adminHandler.CreateModule)
		r.Delete("/admin/modules/{id}", adminHandler.DeleteModule)
		r.Get("/admin/tests/{id}", adminHandler.GetTest)
		r.Get("/admin/projects/{id}", adminHandler.GetProject)
		r.Post("/admin/lessons", adminHandler.CreateLesson)
		r.Get("/admin/lessons/{id}", adminHandler.GetLesson)
		r.Put("/admin/lessons/{id}", adminHandler.UpdateLesson)
		r.Delete("/admin/lessons/{id}", adminHandler.DeleteLesson)
		r.Post("/admin/lessons/{id}/cancel", adminHandler.CancelLesson)
		r.Post("/admin/lessons/{id}/substitute", adminHandler.SubstituteTeacher)
		r.Post("/admin/modules/bulk", adminHandler.CreateModulesBulk)
		r.Post("/admin/lessons/bulk", adminHandler.CreateLessonsBulk)
		r.Post("/admin/tests", adminHandler.CreateTest)
		r.Delete("/admin/tests/{id}", adminHandler.DeleteTest)
		r.Post("/admin/projects", adminHandler.CreateProject)
		r.Delete("/admin/projects/{id}", adminHandler.DeleteProject)
		r.Post("/admin/media/upload", adminHandler.UploadMedia)
		r.Get("/admin/users", adminHandler.GetUsersList)
		r.Get("/admin/users/{id}", adminHandler.GetUserInfo)
		r.Post("/admin/users", adminHandler.CreateUser)
		r.Put("/admin/users/{id}", adminHandler.UpdateUser)
		r.Delete("/admin/users/{id}", adminHandler.DeleteUser)
		r.Post("/admin/enroll", adminHandler.EnrollUser)
		r.Delete("/admin/courses/{id}/enroll/{user_id}", adminHandler.UnenrollStudent)
		r.Get("/admin/users/all", adminHandler.GetAllUsersTable)
		r.Get("/admin/students/detailed", adminHandler.GetDetailedStudents)
		r.Get("/admin/teachers/detailed", adminHandler.GetDetailedTeachers)
		r.Get("/admin/teachers/{id}", adminHandler.GetUserInfo)
		r.Get("/admin/curators/detailed", adminHandler.GetDetailedCurators)
		r.Get("/admin/moderators/detailed", adminHandler.GetDetailedModerators)
		r.Post("/admin/streams", adminHandler.CreateStream)
		r.Get("/admin/streams", adminHandler.GetStreams)
		r.Post("/admin/groups", adminHandler.CreateGroup)
		r.Get("/admin/groups", adminHandler.GetGroups)
		r.Post("/admin/courses/bulk", adminHandler.CreateFullCourse)

		r.Get("/api/admin/dashboard/stats", dashboardHandler.GetAdminDashboard)
		r.Get("/api/admin/curator/dashboard", dashboardHandler.GetCuratorDashboard)
		r.Get("/api/admin/courses", adminHandler.GetAllCourses)
		r.Post("/api/admin/courses", adminHandler.CreateCourse)
		r.Put("/api/admin/courses/{id}/settings", adminHandler.UpdateCourseSettings)
		r.Get("/api/admin/courses/{id}/structure", adminHandler.GetCourseStructure)
		r.Get("/api/admin/courses/{id}/students", adminHandler.GetCourseStudents)
		r.Get("/api/admin/courses/{id}/stats", adminHandler.GetCourseStats)
		r.Post("/api/admin/modules", adminHandler.CreateModule)
		r.Delete("/api/admin/modules/{id}", adminHandler.DeleteModule)
		r.Get("/api/admin/tests/{id}", adminHandler.GetTest)
		r.Get("/api/admin/projects/{id}", adminHandler.GetProject)
		r.Post("/api/admin/lessons", adminHandler.CreateLesson)
		r.Get("/api/admin/lessons/{id}", adminHandler.GetLesson)
		r.Put("/api/admin/lessons/{id}", adminHandler.UpdateLesson)
		r.Delete("/api/admin/lessons/{id}", adminHandler.DeleteLesson)
		r.Post("/api/admin/lessons/{id}/cancel", adminHandler.CancelLesson)
		r.Post("/api/admin/lessons/{id}/substitute", adminHandler.SubstituteTeacher)
		r.Post("/api/admin/modules/bulk", adminHandler.CreateModulesBulk)
		r.Post("/api/admin/lessons/bulk", adminHandler.CreateLessonsBulk)
		r.Post("/api/admin/tests", adminHandler.CreateTest)
		r.Delete("/api/admin/tests/{id}", adminHandler.DeleteTest)
		r.Post("/api/admin/projects", adminHandler.CreateProject)
		r.Delete("/api/admin/projects/{id}", adminHandler.DeleteProject)
		r.Post("/api/admin/media/upload", adminHandler.UploadMedia)
		r.Get("/api/admin/users", adminHandler.GetUsersList)
		r.Get("/api/admin/users/{id}", adminHandler.GetUserInfo)
		r.Post("/api/admin/users", adminHandler.CreateUser)
		r.Put("/api/admin/users/{id}", adminHandler.UpdateUser)
		r.Delete("/api/admin/users/{id}", adminHandler.DeleteUser)
		r.Post("/api/admin/enroll", adminHandler.EnrollUser)
		r.Delete("/api/admin/courses/{id}/enroll/{user_id}", adminHandler.UnenrollStudent)
		r.Get("/api/admin/users/all", adminHandler.GetAllUsersTable)
		r.Get("/api/admin/students/detailed", adminHandler.GetDetailedStudents)
		r.Get("/api/admin/teachers/detailed", adminHandler.GetDetailedTeachers)
		r.Get("/api/admin/teachers/{id}", adminHandler.GetUserInfo)
		r.Get("/api/admin/curators/detailed", adminHandler.GetDetailedCurators)
		r.Get("/api/admin/moderators/detailed", adminHandler.GetDetailedModerators)
		r.Post("/api/admin/streams", adminHandler.CreateStream)
		r.Get("/api/admin/streams", adminHandler.GetStreams)
		r.Post("/api/admin/groups", adminHandler.CreateGroup)
		r.Get("/api/admin/groups", adminHandler.GetGroups)
		r.Post("/api/admin/courses/bulk", adminHandler.CreateFullCourse)

		r.Patch("/api/groups/{groupId}", groupsHandler.UpdateGroup)
		r.Post("/api/groups/{groupId}/students", groupsHandler.AddStudentToGroup)
		r.Delete("/api/groups/{groupId}/students/{studentId}", groupsHandler.RemoveStudentFromGroup)
		r.Patch("/api/students/{studentId}/group", groupsHandler.ChangeStudentGroup)
		r.Patch("/api/teachers/{teacherId}/group", groupsHandler.ChangeTeacherGroup)

		r.Get("/api/attendance/students/{studentId}/calendar", attendanceHandler.GetStudentCalendar)
		r.Patch("/api/attendance/lessons/{lessonId}", attendanceHandler.MarkLessonAttendance)
		r.Get("/api/attendance/students/{studentId}/stats", attendanceHandler.GetStudentStats)
		r.Get("/api/attendance/lessons/{lessonId}", attendanceHandler.GetLessonAttendance)

		r.Post("/api/freeze-requests", freezeHandler.CreateFreezeRequest)
		r.Get("/api/freeze-requests", freezeHandler.GetPendingRequests)
		r.Patch("/api/freeze-requests/{requestId}/approve", freezeHandler.ApproveRequest)
		r.Patch("/api/freeze-requests/{requestId}/reject", freezeHandler.RejectRequest)
		r.Get("/api/students/{studentId}/freeze-status", freezeHandler.GetStudentFreezeStatus)

		r.Post("/api/comments", commentHandler.CreateComment)
		r.Get("/api/comments", commentHandler.GetComments)
		r.Patch("/api/comments/{commentId}/read", commentHandler.MarkCommentAsRead)

		r.Post("/api/notifications", notificationHandler.CreateNotification)

		r.Get("/api/access-requests", accessHandler.GetPendingRequests)
		r.Patch("/api/access-requests/{requestId}/approve", accessHandler.ApproveRequest)
		r.Patch("/api/access-requests/{requestId}/reject", accessHandler.RejectRequest)

		r.Get("/api/statistics/students/{studentId}", statisticsHandler.GetStudentStatistics)
		r.Post("/api/statistics/students/{studentId}/refresh", statisticsHandler.RefreshStudentStatistics)

		r.Get("/api/reports/lessons.xlsx", reportsHandler.DownloadLessonsReport)

		r.Post("/api/admin/banner", bannerHandler.CreateBanner)
		r.Patch("/api/admin/banner/{bannerId}", bannerHandler.UpdateBanner)
		r.Delete("/api/admin/banner/{bannerId}", bannerHandler.DeleteBanner)

		r.Get("/staff/submissions", reviewHandler.GetPendingSubmissions)
		r.Post("/staff/submissions/{id}/evaluate", reviewHandler.EvaluateSubmission)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware)

		r.Get("/teachers", learningHandler.GetTeachers)
		r.Get("/teachers/{id}", learningHandler.GetTeacherDetails)
		r.Post("/teachers/{id}/reviews", learningHandler.AddReview)
		r.Get("/teacher/profile", learningHandler.GetTeacherDashboard)
		r.Get("/teacher/monthly-report", teacherDashboardHandler.GetTeacherMonthlyReport)
		r.Get("/dashboard/home", dashboardHandler.GetUserHome)
		r.Get("/my-courses", learningHandler.GetMyCourses)
		r.Get("/courses/{id}", learningHandler.GetCourseContent)
		r.Get("/lessons/{id}", learningHandler.GetLessonDetail)
		r.Post("/lessons/{id}/assignment", learningHandler.SubmitAssignment)
		r.Post("/lessons/{id}/attendance", learningHandler.SetLessonAttendance)
		r.Post("/admin/courses/bulk", adminHandler.CreateFullCourse)
		r.Get("/tests/{id}", learningHandler.GetTest)
		r.Post("/tests/{id}/submit", learningHandler.SubmitTest)
		r.Get("/projects/{id}", learningHandler.GetProject)
		r.Post("/projects/{id}/submission", learningHandler.SubmitProject)
		r.Get("/profile", profileHandler.GetProfile)
		r.Put("/profile", profileHandler.UpdateProfile)
		r.Put("/profile/teacher/schedule", profileHandler.UpdateTeacherSchedule)
		r.Get("/schedule/weekly", scheduleHandler.GetWeeklySchedule)
		r.Get("/schedule/monthly", scheduleHandler.GetMonthlySchedule)
		r.Get("/chat/ws", chatHandler.ConnectToChat)
		r.Get("/chat/history", chatHandler.GetChatHistory)

		r.Get("/teacher/certificates", learningHandler.GetTeacherCertificates)

		r.Get("/api/teachers", learningHandler.GetTeachers)
		r.Get("/api/teachers/{id}", learningHandler.GetTeacherDetails)
		r.Post("/api/teachers/{id}/reviews", learningHandler.AddReview)
		r.Get("/api/teacher/profile", learningHandler.GetTeacherDashboard)
		r.Get("/api/teacher/certificates", learningHandler.GetTeacherCertificates)
		r.Get("/api/teacher/monthly-report", teacherDashboardHandler.GetTeacherMonthlyReport)
		r.Get("/api/dashboard/home", dashboardHandler.GetUserHome)
		r.Get("/api/my-courses", learningHandler.GetMyCourses)
		r.Get("/api/courses/{id}", learningHandler.GetCourseContent)
		r.Get("/api/lessons/{id}", learningHandler.GetLessonDetail)
		r.Post("/api/lessons/{id}/assignment", learningHandler.SubmitAssignment)
		r.Post("/api/lessons/{id}/attendance", learningHandler.SetLessonAttendance)
		r.Get("/api/tests/{id}", learningHandler.GetTest)
		r.Post("/api/tests/{id}/submit", learningHandler.SubmitTest)
		r.Get("/api/projects/{id}", learningHandler.GetProject)
		r.Post("/api/projects/{id}/submission", learningHandler.SubmitProject)
		r.Get("/api/profile", profileHandler.GetProfile)
		r.Put("/api/profile", profileHandler.UpdateProfile)
		r.Put("/api/profile/teacher/schedule", profileHandler.UpdateTeacherSchedule)
		r.Get("/api/schedule/weekly", scheduleHandler.GetWeeklySchedule)
		r.Get("/api/schedule/monthly", scheduleHandler.GetMonthlySchedule)
		r.Get("/api/chat/ws", chatHandler.ConnectToChat)
		r.Get("/api/chat/history", chatHandler.GetChatHistory)

		r.Get("/api/notifications", notificationHandler.GetNotifications)
		r.Patch("/api/notifications/{notificationId}/read", notificationHandler.MarkNotificationAsRead)

		r.Post("/api/access-requests", accessHandler.CreateAccessRequest)

		r.Get("/api/banner/active", bannerHandler.GetActiveBanners)

		r.Get("/api/students/{studentId}/freeze-status", freezeHandler.GetStudentFreezeStatus)
		r.Get("/api/statistics/students/{studentId}", statisticsHandler.GetStudentStatistics)
		r.Get("/api/courses", learningHandler.GetAllCourses)

		r.Get("/students/{studentId}/freeze-status", freezeHandler.GetStudentFreezeStatus)
		r.Get("/statistics/students/{studentId}", statisticsHandler.GetStudentStatistics)
		r.Get("/courses", learningHandler.GetAllCourses)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))
	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))).ServeHTTP(w, r)
	})

	listenAddr := ":" + apiPort
	slog.Info("server started", slog.String("port", apiPort))

	if err := http.ListenAndServe(listenAddr, r); err != nil {
		slog.Error("server failed", logger.Err(err))
	}
}
