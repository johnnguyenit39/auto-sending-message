package app

import (
	appContext "messenging_test/components/app_context"
	"messenging_test/middlewares"

	// ginOtp "messenging_test/modules/otp/transport/gin"
	// ginPermission "messenging_test/modules/permission/transport/gin"
	// ginSubscription "messenging_test/modules/subscription/transport/gin"
	autoSender "messenging_test/modules/message/auto_sender"
	messageGinHandler "messenging_test/modules/message/transport/gin"
	"os"
	"time"

	"context"
	"encoding/json"

	httprequest "messenging_test/http_request"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeApp(appContext appContext.AppContext) {
	router := appContext.GetGinApp()

	// Add swagger
	swaggerGroup := router.Group("/docs")
	// Add basic auth middleware for Swagger UI
	swaggerGroup.Use(BasicAuthMiddleware())

	// Add swagger
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Use CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Apply middleware
	router.Use(middlewares.PanicRecoveryMiddleware())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.LoggerMiddleware())
	router.Use(middlewares.LanguageMiddleware())

	v1 := router.Group("/api/v1")

	// Khởi tạo http_request service
	httpRepo := httprequest.NewHttpRepository()
	httpSvc := httprequest.NewHttpService(httpRepo)

	// Message & Auto-sender endpoints
	{
		webhookURL := os.Getenv("WEBHOOK_URL")
		webhookKey := os.Getenv("WEBHOOK_KEY")
		sender := autoSender.NewAutoSender(appContext.GetMainDBConnection(), appContext.GetClient(), webhookURL, webhookKey, 2*time.Minute, appContext.GetPubSub(), httpSvc)

		// Subscribe to message.sent event and save messageId to Redis
		if appContext.GetClient() != nil && appContext.GetPubSub() != nil {
			ch := appContext.GetPubSub().Subscribe("message.sent")
			go func() {
				for msg := range ch {
					var event autoSender.MessageSentEvent
					if err := json.Unmarshal([]byte(msg), &event); err == nil {
						key := "sent_message:" + event.ID
						val, _ := json.Marshal(map[string]interface{}{"messageId": event.MessageID, "sent_at": event.SentAt.Format(time.RFC3339)})
						appContext.GetClient().Set(context.Background(), key, val, 0)
					}
				}
			}()
		}

		v1.GET("/messages/sent", messageGinHandler.GetSentMessagesHandler(appContext.GetMainDBConnection()))
		v1.POST("/auto-sender/start", messageGinHandler.StartAutoSenderHandler(func() { sender.Start() }))
		v1.POST("/auto-sender/stop", messageGinHandler.StopAutoSenderHandler(func() { sender.Stop() }))
	}

}
