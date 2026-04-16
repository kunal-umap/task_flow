package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskflow/internal/config"
	"taskflow/internal/db"
	"taskflow/internal/handler"
	"taskflow/internal/middleware"
	"taskflow/internal/repository"
	service "taskflow/internal/services"
	"taskflow/seed"
	"time"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	// "github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Critical error could not connect to database: %v", err)
	}
	defer pool.Close()
	runMigrations(cfg.DBUrl)
	seed.Run(&db.Database{Pool: pool})

	userRepo := repository.NewUserRepository(pool)
	authServices := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authServices)
	projectRepo := repository.NewProjectRepository(pool)
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	taskRepo := repository.NewTaskRepository(pool)
	taskService := service.NewTaskService(taskRepo, projectRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/projects", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/", projectHandler.GetProjects)
		r.Post("/", projectHandler.CreateProject)

		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", projectHandler.GetProjectByID)
			r.Patch("/", projectHandler.UpdateProject)
			r.Delete("/", projectHandler.DeleteProject)

			// nested tasks
			r.Get("/tasks", taskHandler.GetTasks)
			r.Post("/tasks", taskHandler.CreateTask)
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/{taskID}", taskHandler.GetTaskByID)
		r.Patch("/{taskID}", taskHandler.UpdateTask)
		r.Delete("/{taskID}", taskHandler.DeleteTask)
	})

	//  Router setup
	// mux := http.NewServeMux()

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Task flow Api is running"))
	// })

	// mux.HandleFunc("/auth/register", authHandler.Register)
	// mux.HandleFunc("/auth/login", authHandler.Login)

	// // For testing protected route
	// mux.Handle("/protected", middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Protected route accessed"))
	// })))
	// // Configuration of server
	// mux.Handle("/projects", authMiddleware(http.HandlerFunc(projectHandler.GetProjects)))
	// mux.Handle("/projects/create", authMiddleware(http.HandlerFunc(projectHandler.CreateProject)))
	// mux.Handle("/projects/get", authMiddleware(http.HandlerFunc(projectHandler.GetProjectByID)))
	// mux.Handle("/projects/update", authMiddleware(http.HandlerFunc(projectHandler.UpdateProject)))
	// mux.Handle("/projects/delete", authMiddleware(http.HandlerFunc(projectHandler.DeleteProject)))
	// mux.Handle("/tasks/get", authMiddleware(http.HandlerFunc(taskHandler.GetTaskByID)))
	// mux.Handle("/tasks/create", authMiddleware(http.HandlerFunc(taskHandler.CreateTask)))
	// mux.Handle("/tasks/list", authMiddleware(http.HandlerFunc(taskHandler.GetTasks)))
	// mux.Handle("/tasks/update", authMiddleware(http.HandlerFunc(taskHandler.UpdateTask)))
	// mux.Handle("/tasks/delete", authMiddleware(http.HandlerFunc(taskHandler.DeleteTask)))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//  Shutdown logic

	go func() {
		log.Printf("🚀 Server starting on port %s ", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	// intrupt for shutdown

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exit safely")
}

func runMigrations(dbURL string) {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("migration init failed:", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal("migration failed:", err)
	}

	log.Println("migrations applied")
}
