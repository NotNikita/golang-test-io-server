package app

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"test-server/internal/app/handlers"
	config "test-server/internal/config"
	"test-server/internal/domain/task/repository"
	"test-server/internal/domain/task/service"
	middleware "test-server/internal/middleware"
)

type App struct {
	config *config.Config
	server *fiber.App
}

func NewApp(configPath string) (*App, error) {
	c, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{config: c}
	httpServer := app.BootstrapHandlers()
	app.server = httpServer

	return app, nil
}

func (a *App) BootstrapHandlers() *fiber.App {
	fiberApp := fiber.New()
	middleware.CorsMiddleware(fiberApp)
	middleware.LoggerMiddleware(fiberApp)

	tasksRepo := repository.NewTasksRepository()
	tasksService := service.NewTasksService(a.config.Service.Interval, tasksRepo)
	handler := handlers.NewHandler(tasksService)

	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	})
	fiberApp.Post("api/tasks", handler.PostRegisterTask)
	fiberApp.Get("api/tasks/:id", handler.GetTaskInfo)
	fiberApp.Delete("api/tasks/:id", handler.DeleteTask)
	return fiberApp
}

func (a *App) ListenAndServe() error {
	// Setup graceful shutdown
	shutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	serverAddress := fmt.Sprintf("%s:%s", a.config.Service.Host, strconv.Itoa(a.config.Service.Port))

	go func() {
		fmt.Println("Golang test IO server started")
		if servErr := a.server.Listen(serverAddress); servErr != nil {
			log.Fatalf("app.ListenAndServe: failed to start server: %v", servErr)
		}
	}()

	<-shutdownCtx.Done()
	fmt.Println("Shutdown signal received, starting graceful shutdown...")

	// Create a timeout context for shutdown
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeoutCancel()

	// Shutdown HTTP server
	if err := a.server.ShutdownWithContext(timeoutCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	fmt.Println("Graceful shutdown completed")
	return nil
}
