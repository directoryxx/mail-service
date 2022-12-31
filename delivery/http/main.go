package http

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func RunAPI() {
	app := echo.New()

	log.Println("[INFO] App Mode : API")

	app.Validator = &CustomValidator{validator: validator.New()}

	log.Println("[INFO] Starting mail Service on port", os.Getenv("APPLICATION_PORT"))

	// log.Println("[INFO] Loading Database")
	// dbMongo, err := infrastructure.ConnectMongo()

	// if err != nil {
	// 	log.Fatalf("Could not initialize Mongo connection using client %s", err)
	// }

	// defer dbMongo.Disconnect(context.TODO())

	// log.Println("[INFO] Loading Redis")
	// redisConnect := infrastructure.OpenRedis()

	// defer redisConnect.Close()

	// log.Println("[INFO] Loading Repository")
	// logRepo := repository.NewLogRepository(dbMongo)

	// log.Println("[INFO] Loading Usecase")
	// logUsecase := usecase.NewLogUseCase(logRepo)

	// log.Println("[INFO] Loading Controller")
	// logController := controller.NewLogController(logUsecase)

	// log.Println("[INFO] Loading Middleware")
	// SetMiddleware(app)

	// log.Println("[INFO] Loading Routes")
	// api.Routes(app, logController)

	// log.Println("[INFO] Loading JWT Middleware")
	// SetPrivateMiddleware(app)

	// log.Println("[INFO] Loading Protected Endpoint")
	// api.PrivateRoutes(app, userController)

	log.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("APPLICATION_PORT"))))
}

func SetMiddleware(r *echo.Echo) {
	// Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	// Cors
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "http://localhost"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
}
