package main

import (
	"messenging_test/appi18n"
	appContext "messenging_test/components/app_context"
	appConfig "messenging_test/config/app"
	storage "messenging_test/config/postgres"
	"messenging_test/config/pubsub"
	"messenging_test/config/redis"
	_ "messenging_test/docs"
	"messenging_test/logger"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// @title Insider Assessment Project API
// @version 1.0
// @description API for an automatic message sending system: fetch unsent messages from DB, send periodically, manage sent status, provide start/stop auto sender API, retrieve sent messages list, bonus: cache messageId to Redis.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file - continuing with system environment variables")
	}

	appi18n.Init()

	logger.InitializeLogger()

	//Pub sub
	pubSub := pubsub.NewPubSub()
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		log.Warn().Err(err).Msg("Redis connection failed - continuing without Redis functionality")
	} else {
		pubsub.ListenEvent(redisClient.GetClient(), pubSub)
	}

	// Initialize OKX service
	log.Info().Msg("J City EG service initialized successfully")

	db, err := storage.NewConnection()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load the database")
		return
	}

	//Migrate data
	storage.AutoMigrate(db)

	// Initialize Gin router
	app := gin.Default()

	// Initialize AppContext with DB and app
	appContext := appContext.NewAppContext(db, redisClient.GetClient(), pubSub, app)

	// Initialize application config
	appConfig.InitializeApp(appContext)

	// Start the application on port 80
	if err := app.Run(":80"); err != nil {
		log.Fatal().Err(err).Msg("failed to start the application")
		return
	}
}
