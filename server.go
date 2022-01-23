package main

import "os"
import "log"
import "fmt"
import "time"
import "syscall"
import "os/signal"
import "github.com/gofiber/fiber/v2"
import "github.com/gofiber/fiber/v2/middleware/etag"

const idleTimeout = 5 * time.Second

func main() {
    appConfig := fiber.Config{
        Prefork:              false,
        ServerHeader:         "Project",
        StrictRouting:        true,
        CaseSensitive:        true,
        Immutable:            true,
        UnescapePath:         true,
        BodyLimit:            2 * 1024 * 1024,
        Concurrency:          256 * 1024,
        Views:                nil,
        ReadBufferSize:       2048,
        WriteBufferSize:      2048,
        CompressedFileSuffix: ".project.gz",
        AppName:              "Project",
        IdleTimeout:          idleTimeout,
    }

    app := fiber.New(appConfig)

    eTagConfig := etag.Config{
        Weak: true,
    }

    eTag := etag.New(eTagConfig)

    app.Use(eTag)

    app.Get("/", helloWorld)

    go func() {
        if err := app.Listen(":3000"); err != nil {
            log.Panic(err)
        }
    }()

    c := make(chan os.Signal, 1)

    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    _ = <-c

    fmt.Println("\nShutting down a project...")

    _ = app.Shutdown()

    fmt.Println("Project was successfully shutdown.")
}

func helloWorld(c *fiber.Ctx) error {
    response := fiber.Map{
        "message": "Hello World!",
    }

    return c.JSON(response)
}
