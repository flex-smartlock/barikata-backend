package main

import (
	"os"
	"strconv"

	cr "github.com/flex-smartlock/barikata-backend/routes/control"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	isRunningSecure, err := strconv.ParseBool(os.Getenv("SERVER_RUNNING_SECURE"))
	if err != nil {
		isRunningSecure = false
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	// e.Use(middleware.CORS()) // ! 最終的にCORSは有効にする！
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 1,
	}))
	if isRunningSecure {
		e.Use(middleware.Secure())
		e.Use(middleware.HTTPSRedirect())
	}
	e.Use(middleware.RemoveTrailingSlash())

	v1 := e.Group("/v1")
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	v1.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	control := e.Group("/control")
	control.GET("/test", cr.Test)

	serverPort := os.Getenv("SERVER_PORT")
	if isRunningSecure {
		serverCertFile := os.Getenv("SERVER_CERT_FILE")
		serverKeyFile := os.Getenv("SERVER_KEY_FILE")
		e.Logger.Fatal(e.StartTLS(serverPort, serverCertFile, serverKeyFile))
	} else {
		e.Logger.Fatal(e.Start(serverPort))
	}
}
