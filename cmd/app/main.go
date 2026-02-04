package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

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

	reviewHttp "lms_backend/internal/review/delivery/http"
	reviewRepo "lms_backend/internal/review/repository"
	reviewUseCase "lms_backend/internal/review/usecase"

	profileHttp "lms_backend/internal/profile/delivery/http"
	profileRepo "lms_backend/internal/profile/repository"
	profileUseCase "lms_backend/internal/profile/usecase"

	scheduleHttp "lms_backend/internal/schedule/delivery/http"
	scheduleRepo "lms_backend/internal/schedule/repository"
	scheduleUseCase "lms_backend/internal/schedule/usecase"

	"lms_backend/internal/domain"
	storageService "lms_backend/pkg/storage"
)

// @title Cap Education LMS - API
// @version 1.0
// @description API для LMS платформы Cap Education.
// @host localhost:8000
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (using system envs)")
	}

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
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Успешное подключение к PostgreSQL")

	s3Client, err := storageService.NewS3Client(
		os.Getenv("S3_ENDPOINT_URL"),
		os.Getenv("S3_REGION"),
		os.Getenv("S3_BUCKET_NAME"),
		os.Getenv("S3_ACCESS_KEY_ID"),
		os.Getenv("S3_SECRET_ACCESS_KEY"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize S3 Storage: %v", err)
	}
	log.Println("Успешная инициализация S3 Storage (Mail.ru CS).")

	authRepoImpl := repository.NewAuthRepository(db)
	authUsecase := authUseCase.NewAuthUsecase(authRepoImpl)
	authHandler := authHttp.NewAuthHandler(authUsecase)

	dashboardRepository := dashboardRepo.NewDashboardRepository(db)
	dashboardUsecase := dashboardUseCase.NewDashboardUseCase(dashboardRepository)
	dashboardHandler := dashboardHttp.NewDashboardHandler(dashboardUsecase)

	adminRepo := contentAdminRepo.NewContentAdminRepository(db)
	adminUsecase := contentAdminUseCase.NewContentAdminUseCase(adminRepo, s3Client)
	adminHandler := contentAdminHttp.NewContentAdminHandler(adminUsecase)

	learningRepoImpl := learningRepo.NewLearningRepository(db)
	learningUC := learningUseCase.NewLearningUseCase(learningRepoImpl, s3Client)
	learningHandler := learningHttp.NewLearningHandler(learningUC)

	reviewRepoImpl := reviewRepo.NewReviewRepository(db)
	reviewUC := reviewUseCase.NewReviewUseCase(reviewRepoImpl)
	reviewHandler := reviewHttp.NewReviewHandler(reviewUC)

	profileRepoImpl := profileRepo.NewProfileRepository(db)
	profileUC := profileUseCase.NewProfileUseCase(profileRepoImpl, s3Client)
	profileHandler := profileHttp.NewProfileHandler(profileUC)

	scheduleRepoImpl := scheduleRepo.NewScheduleRepository(db)
	scheduleUC := scheduleUseCase.NewScheduleUseCase(scheduleRepoImpl)
	scheduleHandler := scheduleHttp.NewScheduleHandler(scheduleUC)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware, authMiddleware.RoleRequiredMiddleware(domain.RoleAdmin, domain.RoleTeacher, domain.RoleModerator, domain.RoleCurator))

		r.Get("/admin/dashboard/stats", dashboardHandler.GetAdminDashboard)
		r.Get("/admin/courses", adminHandler.GetAllCourses)
		r.Post("/admin/courses", adminHandler.CreateCourse)
		r.Put("/admin/courses/{id}/settings", adminHandler.UpdateCourseSettings)
		r.Get("/admin/courses/{id}/structure", adminHandler.GetCourseStructure)
		r.Get("/admin/courses/{id}/students", adminHandler.GetCourseStudents)
		r.Get("/admin/courses/{id}/stats", adminHandler.GetCourseStats)

		r.Post("/admin/modules", adminHandler.CreateModule)
		r.Post("/admin/lessons", adminHandler.CreateLesson)
		r.Post("/admin/tests", adminHandler.CreateTest)
		r.Post("/admin/projects", adminHandler.CreateProject)

		r.Post("/admin/media/upload", adminHandler.UploadMedia)

		r.Get("/admin/users", adminHandler.GetUsersList)
		r.Post("/admin/users", adminHandler.CreateUser)
		r.Put("/admin/users/{id}", adminHandler.UpdateUser)
		r.Delete("/admin/users/{id}", adminHandler.DeleteUser)
		r.Post("/admin/enroll", adminHandler.EnrollUser)

		r.Get("/staff/submissions", reviewHandler.GetPendingSubmissions)
		r.Post("/staff/submissions/{id}/evaluate", reviewHandler.EvaluateSubmission)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware)
		r.Get("/dashboard/home", dashboardHandler.GetUserHome)
		r.Get("/my-courses", learningHandler.GetMyCourses)
		r.Get("/courses/{id}", learningHandler.GetCourseContent)
		r.Get("/lessons/{id}", learningHandler.GetLessonDetail)
		r.Post("/lessons/{id}/assignment", learningHandler.SubmitAssignment)
		r.Post("/lessons/{id}/complete", learningHandler.CompleteLesson)
		r.Get("/profile", profileHandler.GetProfile)
		r.Put("/profile", profileHandler.UpdateProfile)
		r.Get("/schedule/weekly", scheduleHandler.GetWeeklySchedule)
		r.Get("/schedule/monthly", scheduleHandler.GetMonthlySchedule)
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/docs/swagger.json")))
	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))).ServeHTTP(w, r)
	})

	listenAddr := ":" + apiPort
	log.Printf("Server started on port: %s", apiPort)
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
