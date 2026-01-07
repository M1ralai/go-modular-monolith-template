package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	authHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/http"
	authService "github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/service"
	courseHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/course/http"
	courseRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/course/repository"
	courseService "github.com/M1ralai/go-modular-monolith-template/internal/modules/course/service"
	eventHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/event/http"
	eventRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/event/repository"
	eventService "github.com/M1ralai/go-modular-monolith-template/internal/modules/event/service"
	financeHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/http"
	financeRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/repository"
	financeService "github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/service"
	goalHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/http"
	goalRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/repository"
	goalService "github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/service"
	habitHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/http"
	habitRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/repository"
	habitService "github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/service"
	healthHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/health/http"
	journalHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/http"
	journalRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/repository"
	journalService "github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/service"
	lifeareaHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/http"
	lifeareaRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/repository"
	lifeareaService "github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/service"
	noteHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/note/http"
	noteRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/note/repository"
	noteService "github.com/M1ralai/go-modular-monolith-template/internal/modules/note/service"
	peopleHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/people/http"
	peopleRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/people/repository"
	peopleService "github.com/M1ralai/go-modular-monolith-template/internal/modules/people/service"
	taskHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/task/http"
	taskRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/task/repository"
	taskService "github.com/M1ralai/go-modular-monolith-template/internal/modules/task/service"
	userHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/http"
	userRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/repository"
	userService "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/service"

	calendarHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/http"
	calendarRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/repository"
	calendarService "github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/service"

	scheduleHttp "github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/http"
	scheduleRepo "github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/repository"
	scheduleService "github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/service"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/websocket"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"

	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	httpServer *http.Server
	db         *sqlx.DB
	logger     *logger.ZapLogger
}

func NewServer(db *sqlx.DB, zapLogger *logger.ZapLogger) *Server {
	// WebSocket Hub
	wsHub := websocket.NewHub(zapLogger)
	go wsHub.Run()
	wsHandler := websocket.NewHandler(wsHub, zapLogger)

	// Broadcaster for real-time notifications
	broadcaster := notifService.NewBroadcaster(wsHub, zapLogger)

	// Distributed lock for jobs
	jobLock := jobs.NewDistributedLock(db)

	// Job Pool for async task processing
	jobPool := jobs.NewWorkerPool(5, 100, zapLogger, nil, jobLock)
	jobPool.Start()

	// Health module
	healthHandler := healthHttp.NewHandler()

	// User module
	userRepository := userRepo.NewPostgresRepository(db)
	userSvc := userService.NewUserService(userRepository, zapLogger)
	userHandler := userHttp.NewHandler(userSvc)

	// Auth module
	authSvc := authService.NewAuthService(userRepository, zapLogger)
	authHandler := authHttp.NewHandler(authSvc)

	// LifeArea module
	lifeareaRepository := lifeareaRepo.NewPostgresRepository(db)
	lifeareaSvc := lifeareaService.NewLifeAreaService(lifeareaRepository, zapLogger, broadcaster)
	lifeareaHandler := lifeareaHttp.NewHandler(lifeareaSvc)

	// Course module
	courseRepository := courseRepo.NewPostgresRepository(db)
	courseSvc := courseService.NewCourseService(courseRepository, zapLogger, broadcaster)
	courseHandler := courseHttp.NewHandler(courseSvc)

	// Task module
	taskRepository := taskRepo.NewPostgresRepository(db)
	taskSvc := taskService.NewTaskService(taskRepository, zapLogger, broadcaster)
	taskHandler := taskHttp.NewHandler(taskSvc, jobPool, taskRepository, broadcaster, zapLogger)

	// Note module
	noteRepository := noteRepo.NewPostgresRepository(db)
	noteSvc := noteService.NewNoteService(noteRepository, zapLogger, broadcaster)
	noteHandler := noteHttp.NewHandler(noteSvc)

	// Habit module
	habitRepository := habitRepo.NewPostgresRepository(db)
	habitSvc := habitService.NewHabitService(habitRepository, zapLogger, broadcaster)
	habitHandler := habitHttp.NewHandler(habitSvc, habitRepository, broadcaster, zapLogger)

	// Goal module
	goalRepository := goalRepo.NewPostgresRepository(db)
	goalSvc := goalService.NewGoalService(goalRepository, zapLogger, broadcaster)
	goalHandler := goalHttp.NewHandler(goalSvc)

	// Event module
	eventRepository := eventRepo.NewPostgresRepository(db)
	eventSvc := eventService.NewEventService(eventRepository, zapLogger, broadcaster)
	eventHandler := eventHttp.NewHandler(eventSvc)

	// People module
	peopleRepository := peopleRepo.NewPostgresRepository(db)
	peopleSvc := peopleService.NewPersonService(peopleRepository, zapLogger, broadcaster)
	peopleHandler := peopleHttp.NewHandler(peopleSvc)

	// Journal module
	journalRepository := journalRepo.NewPostgresRepository(db)
	journalSvc := journalService.NewJournalService(journalRepository, zapLogger, broadcaster)
	journalHandler := journalHttp.NewHandler(journalSvc)

	// Finance module
	financeRepository := financeRepo.NewPostgresRepository(db)
	financeSvc := financeService.NewFinanceService(financeRepository, zapLogger, broadcaster)
	financeHandler := financeHttp.NewHandler(financeSvc)

	// Calendar module
	integrationRepository := calendarRepo.NewCalendarIntegrationRepository(db)
	syncQueueRepository := calendarRepo.NewSyncQueueRepository(db)
	calendarSvc := calendarService.NewCalendarService(integrationRepository, syncQueueRepository, zapLogger)
	calendarHandler := calendarHttp.NewHandler(calendarSvc)

	// Schedule module
	blockedSlotRepository := scheduleRepo.NewBlockedTimeSlotRepository(db)
	scheduleSvc := scheduleService.NewScheduleService(blockedSlotRepository, zapLogger, broadcaster)
	scheduleHandler := scheduleHttp.NewHandler(scheduleSvc)

	router := mux.NewRouter()

	// WebSocket route - MUST be registered BEFORE middleware to avoid ResponseWriter wrapping
	// WebSocket requires the original http.ResponseWriter to hijack the connection
	router.HandleFunc("/ws", wsHandler.HandleConnection).Methods("GET")

	// Apply middleware to all routes EXCEPT WebSocket
	router.Use(middleware.RecoveryMiddleware)
	router.Use(zapLogger.Middleware)
	router.Use(middleware.MetricsMiddleware)
	router.Use(middleware.TimeoutMiddleware)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip Content-Type header for WebSocket upgrade
			if r.URL.Path != "/ws" {
				w.Header().Set("Content-Type", "application/json")
			}
			next.ServeHTTP(w, r)
		})
	})

	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Public routes (no auth required)
	authHandler.RegisterRoutes(router)

	// API routes (protected)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Register all module routes
	userHandler.RegisterRoutes(api)
	lifeareaHandler.RegisterRoutes(api)
	courseHandler.RegisterRoutes(api)
	taskHandler.RegisterRoutes(api)
	noteHandler.RegisterRoutes(api)
	habitHandler.RegisterRoutes(api)
	goalHandler.RegisterRoutes(api)
	eventHandler.RegisterRoutes(api)
	peopleHandler.RegisterRoutes(api)
	journalHandler.RegisterRoutes(api)
	financeHandler.RegisterRoutes(api)
	calendarHandler.RegisterRoutes(api)
	scheduleHandler.RegisterRoutes(api)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":8080"
	}
	if len(port) > 0 && port[0] != ':' {
		port = ":" + port
	}

	httpServer := &http.Server{
		Addr:         port,
		Handler:      middleware.CorsMiddleware(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		db:         db,
		logger:     zapLogger,
	}
}

func (s *Server) Start() error {
	errChan := make(chan error, 1)

	go func() {
		log.Printf("✓ Server starting... Port: %s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(100 * time.Millisecond):

		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("✓ Graceful shutdown started...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("database close error: %w", err)
	}

	log.Println("✓ Shutdown completed")
	return nil
}
