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
)

// @title Cap Education LMS - Auth API
// @version 1.0
// @description API для LMS платформы Cap Education.
// @host localhost:8080
// @BasePath /
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
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

	authRepoImpl := repository.NewAuthRepository(db)
	authUsecase := authUseCase.NewAuthUsecase(authRepoImpl)
	authHandler := authHttp.NewAuthHandler(authUsecase)

	dashboardRepository := dashboardRepo.NewDashboardRepository(db)
	dashboardUsecase := dashboardUseCase.NewDashboardUseCase(dashboardRepository)
	dashboardHandler := dashboardHttp.NewDashboardHandler(dashboardUsecase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {

		r.Use(authMiddleware.AuthMiddleware)

		r.Get("/dashboard/home", dashboardHandler.GetUserHome)

	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))
	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))).ServeHTTP(w, r)
	})

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
