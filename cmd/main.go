package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gofiber/fiber/v2"
    "github.com/inlovewithgo/transit-backend/main/middlewares"
    "github.com/inlovewithgo/transit-backend/main/routes"
    "github.com/inlovewithgo/transit-backend/main/config"
)

func init() {
    config.InitDatabase()
}

func main() {
    fmt.Println("Transit Backend Service is starting...")

    app := fiber.New(fiber.Config{
        AppName:      "Transit Backend API",
        ServerHeader: "Transit-Backend",
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            
            log.Printf("Error: %v", err)
            
            return c.Status(code).JSON(fiber.Map{
                "error":   true,
                "message": err.Error(),
                "code":    code,
            })
        },
    })

    middleware.SetupMiddleware(app)
    routes.SetupRoutes(app)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        fmt.Println("\nGracefully shutting down Transit Backend...")
        config.ShutdownDatabase()
        if err := app.Shutdown(); err != nil {
            log.Printf("Error during shutdown: %v", err)
        }
        
        fmt.Println("Transit Backend shutdown complete")
        os.Exit(0)
    }()

    port := os.Getenv("PORT")
    if port == "" {
        port = "3030"
    }

    fmt.Printf("Transit Backend Service started successfully on port %s\n", port)
    fmt.Printf("Health check available at: http://localhost:%s/health\n", port)
    fmt.Printf("API documentation available at: http://localhost:%s/api/v1/health/detailed\n", port)

    if err := app.Listen(":" + port); err != nil {
        log.Fatalf("Failed to start Transit Backend Service: %v", err)
    }
}