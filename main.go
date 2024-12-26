package main

import (
	"context"
	"log"

	"github.com/fleimkeipa/maker-checker/controller"
	_ "github.com/fleimkeipa/maker-checker/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/maker-checker/pkg"
	"github.com/fleimkeipa/maker-checker/repositories"
	"github.com/fleimkeipa/maker-checker/uc"
	"github.com/fleimkeipa/maker-checker/util"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type \"Bearer \" and then your API Token
func main() {
	// Start the application
	serveApplication()
}

func serveApplication() {
	// init config
	// loadConfig()

	// Create a new Echo instance
	e := echo.New()

	// Configure Echo settings
	configureEcho(e)

	// Configure CORS middleware
	configureCORS(e)

	// Configure the logger
	sugar := configureLogger(e)
	defer sugar.Sync() // Clean up logger at the end

	mongoClient := initMongo()
	defer mongoClient.Client().Disconnect(context.TODO())

	// Initialize the user use case
	userMongoRepo := repositories.NewUserMongoRepo(mongoClient)
	userUC := uc.NewUserUC(userMongoRepo)
	userController := controller.NewUserHandlers(userUC)

	authHandlers := controller.NewAuthHandlers(userUC)

	messageMongoRepo := repositories.NewMsgMongoRepo(mongoClient)
	messageUC := uc.NewMessageUC(messageMongoRepo)
	messageController := controller.NewMessageHandlers(messageUC)

	// Define authentication routes and handlers
	authRoutes := e.Group("/auth")
	authRoutes.POST("/login", authHandlers.Login)
	authRoutes.POST("/register", authHandlers.Register)

	// Define user routes
	userRoutes := e.Group("")
	userRoutes.Use(util.JWTAuthUser)

	// Define user routes
	usersRoutes := userRoutes.Group("/users")
	usersRoutes.GET("/:id", userController.GetByID)
	usersRoutes.POST("", userController.Create)
	usersRoutes.PATCH("/:id", userController.UpdateUser)
	usersRoutes.DELETE("/:id", userController.DeleteUser)

	// Define message routes
	messageRoutes := userRoutes.Group("/messages")
	messageRoutes.GET("/:id", messageController.GetByID)
	messageRoutes.POST("", messageController.Create)
	messageRoutes.PATCH("/:id", messageController.Update)
	messageRoutes.GET("", messageController.List)

	e.Logger.Fatal(e.Start(":8080"))
}

// func loadConfig() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}
// }

// Configures the Echo instance
func configureEcho(e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true

	// Add Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add Recover middleware
	e.Use(middleware.Recover())
}

// Configures CORS settings
func configureCORS(e *echo.Echo) {
	corsConfig := middleware.CORSWithConfig(middleware.CORSConfig{
		UnsafeWildcardOriginWithAllowCredentials: true,
		AllowCredentials:                         true,
		AllowOrigins:                             []string{"*"},
		AllowMethods:                             []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
		AllowHeaders:                             []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})

	e.Use(corsConfig)
}

// Configures the logger and adds it as middleware
func configureLogger(e *echo.Echo) *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(pkg.ZapLogger(logger))

	sugar := logger.Sugar()
	loggerHandler := controller.NewLogger(sugar)
	e.Use(loggerHandler.LoggerMiddleware)

	return sugar
}

func initMongo() *mongo.Database {
	mongo, err := pkg.MongoConnect()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	log.Println("MongoDB client initialized successfully")

	return mongo
}
