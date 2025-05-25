package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	config "github.com/magnete-library/config"
	handlers "github.com/magnete-library/handlers"
)

// CustomValidator wraps validator.Validate for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file: %v", err)
	}

	// Get allowed client URL for CORS from env
	client := os.Getenv("CLIENT")

	// Create new Echo instance
	e := echo.New()

	// Setup validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Use CORS middleware with config to allow your client URL
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{client},  // your frontend origin, e.g. "http://localhost:3000"
    AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
    AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
    AllowCredentials: true,  // Important: allow credentials
}))


	// Connect to MongoDB
	config.ConnectDB()

	// Routes
	e.GET("/", handlers.HomePage)
	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	
	e.POST("/students/:month", handlers.CreateStudent)
	e.GET("/students/:month", handlers.GetStudents)
	e.GET("/students/:month/search", handlers.SearchStudents)
	e.PUT("/students/:month/:id", handlers.UpdateStudent)
	e.PATCH("/students/:month/:id/payment", handlers.UpdatePayment)
	e.PATCH("/students/:month/:id/status", handlers.ToggleActiveStatus)
	e.PATCH("/students/:month/:id/seat", handlers.UpdateSeatNumber)
	e.POST("/students/migrate", handlers.MigrateMonth)
	e.GET("/collections", handlers.ListCollections)

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
