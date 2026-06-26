package router

import (
	"BE/config"
	"BE/database"
	"BE/internal/domain"
	"BE/internal/handler"
	"BE/internal/hub"
	"BE/internal/repository"
	"BE/internal/service"
	"BE/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	wsfiber "github.com/gofiber/websocket/v2"
)

func Setup(app *fiber.App) {
	app.Use(middleware.CORS())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// Init upload service DULU sebelum service lain
	uploadService, err := service.NewUploadService(
		config.App.CloudinaryCloudName,
		config.App.CloudinaryAPIKey,
		config.App.CloudinaryAPISecret,
	)
	if err != nil {
		log.Fatalf("Failed to init Cloudinary: %v", err)
	}
	uploadHandler := handler.NewUploadHandler(uploadService)

	// Init repositories
	userRepo        := repository.NewUserRepository(database.DB)
	productRepo     := repository.NewProductRepository(database.DB)
	portfolioRepo   := repository.NewPortfolioRepository(database.DB)
	articleRepo     := repository.NewArticleRepository(database.DB)
	testimonialRepo := repository.NewTestimonialRepository(database.DB)
	chatRepo        := repository.NewChatRepository(database.DB)
	settingRepo     := repository.NewSettingRepository(database.DB)
	orderRepo       := repository.NewOrderRepository(database.DB)

	// Init services
	authService        := service.NewAuthService(userRepo)
	userService        := service.NewUserService(userRepo)
	productService     := service.NewProductService(productRepo, uploadService)
	portfolioService   := service.NewPortfolioService(portfolioRepo, uploadService)
	articleService     := service.NewArticleService(articleRepo, uploadService)
	testimonialService := service.NewTestimonialService(testimonialRepo)
	chatService        := service.NewChatService(chatRepo)
	settingService     := service.NewSettingService(settingRepo, productRepo) // Hapus orderRepo
	orderService       := service.NewOrderService(orderRepo)                  // Tambahkan ini

	// Init handlers
	authHandler        := handler.NewAuthHandler(authService)
	userHandler        := handler.NewUserHandler(userService)
	productHandler     := handler.NewProductHandler(productService)
	portfolioHandler   := handler.NewPortfolioHandler(portfolioService)
	articleHandler     := handler.NewArticleHandler(articleService)
	testimonialHandler := handler.NewTestimonialHandler(testimonialService)
	chatHandler        := handler.NewChatHandler(chatService)
	settingHandler     := handler.NewSettingHandler(settingService)
	orderHandler       := handler.NewOrderHandler(orderService) // Tambahkan ini

	// WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if wsfiber.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/admin", wsfiber.New(func(c *wsfiber.Conn) {
		client := &hub.Client{Conn: c, Role: "admin"}
		hub.H.Register(client)
		defer hub.H.Unregister(client)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	app.Get("/ws/public", wsfiber.New(func(c *wsfiber.Conn) {
		client := &hub.Client{Conn: c, Role: "visitor"}
		hub.H.Register(client)
		defer hub.H.Unregister(client)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	app.Get("/ws/chat/:session_id", wsfiber.New(func(c *wsfiber.Conn) {
		sessionID := c.Params("session_id")
		client := &hub.Client{Conn: c, Role: "visitor", SessionID: sessionID}
		hub.H.Register(client)
		defer hub.H.Unregister(client)

		chatRepoWS := repository.NewChatRepository(database.DB)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			message := &domain.ChatMessage{
				ID:         uuid.New(),
				SessionID:  uuid.MustParse(sessionID),
				SenderRole: "visitor",
				Message:    string(msg),
			}
			chatRepoWS.CreateMessage(message)

			hub.H.BroadcastToAdminSession(sessionID, fiber.Map{
				"type":        "chat:message",
				"session_id":  sessionID,
				"message":     string(msg),
				"sender_role": "visitor",
				"created_at":  message.CreatedAt,
			})

			hub.H.BroadcastToAdmins(fiber.Map{
				"type":        "chat:new_message",
				"session_id":  sessionID,
				"sender_name": "Visitor",
			})
		}
	}))

	app.Get("/ws/chat/admin/:session_id", wsfiber.New(func(c *wsfiber.Conn) {
		sessionID := c.Params("session_id")
		client := &hub.Client{Conn: c, Role: "admin", SessionID: sessionID}
		hub.H.Register(client)
		defer hub.H.Unregister(client)

		chatRepoWS := repository.NewChatRepository(database.DB)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			message := &domain.ChatMessage{
				ID:         uuid.New(),
				SessionID:  uuid.MustParse(sessionID),
				SenderRole: "admin",
				Message:    string(msg),
			}
			chatRepoWS.CreateMessage(message)

			hub.H.BroadcastToSession(sessionID, fiber.Map{
				"type":        "chat:message",
				"session_id":  sessionID,
				"message":     string(msg),
				"sender_role": "admin",
				"created_at":  message.CreatedAt,
			})
		}
	}))

	// API
	api := app.Group("/api")

	// Auth
	api.Post("/auth/login", authHandler.Login)

	// Public
	api.Get("/products",            productHandler.GetAll)
	api.Get("/products/categories", productHandler.GetCategories)
	api.Get("/products/:id",        productHandler.GetByID)

	api.Get("/portfolio",     portfolioHandler.GetAll)
	api.Get("/portfolio/:id", portfolioHandler.GetByID)

	api.Get("/articles",            articleHandler.GetPublished)
	api.Get("/articles/categories", articleHandler.GetCategories)
	api.Get("/articles/:slug",      articleHandler.GetBySlug)

	api.Get("/testimonials",  testimonialHandler.GetVisible)
	api.Post("/testimonials", testimonialHandler.SubmitPublic)

	api.Get("/settings",              settingHandler.GetPublicSettings)
	api.Post("/orders/whatsapp",      orderHandler.OrderViaWhatsapp) 
	api.Post("/chat/session",         chatHandler.CreateSession)
	api.Get("/whatsapp",              settingHandler.GetPublicWhatsapp)
	api.Get("/chat/history/:session_id", chatHandler.GetSessionHistory)

	// Admin
	admin := api.Group("/admin", middleware.Protected())

	admin.Get("/account",         userHandler.GetProfile)
	admin.Put("/account",         userHandler.UpdateAccount)
	admin.Get("/dashboard/stats", chatHandler.GetDashboardStats)

	admin.Get("/orders",          orderHandler.GetOrders)
	admin.Get("/orders/:id",      orderHandler.GetOrderByID)
	admin.Put("/orders/:id/status", orderHandler.UpdateOrderStatus)
admin.Delete("/orders/:id",   orderHandler.DeleteOrder)

	admin.Post("/products",                  productHandler.Create)
	admin.Put("/products/:id",               productHandler.Update)
	admin.Put("/products/:id/stock",         productHandler.UpdateStock)
	admin.Delete("/products/:id",            productHandler.Delete)
	admin.Post("/products/categories",       productHandler.CreateCategory)
	admin.Delete("/products/categories/:id", productHandler.DeleteCategory)

	admin.Post("/portfolio",       portfolioHandler.Create)
	admin.Put("/portfolio/:id",    portfolioHandler.Update)
	admin.Delete("/portfolio/:id", portfolioHandler.Delete)

	admin.Get("/articles",                   articleHandler.GetAll)
	admin.Post("/articles",                  articleHandler.Create)
	admin.Put("/articles/:id",               articleHandler.Update)
	admin.Delete("/articles/:id",            articleHandler.Delete)
	admin.Post("/articles/categories",       articleHandler.CreateCategory)
	admin.Delete("/articles/categories/:id", articleHandler.DeleteCategory)

	admin.Get("/testimonials",        testimonialHandler.GetAll)
	admin.Post("/testimonials",       testimonialHandler.Create)
	admin.Put("/testimonials/:id",    testimonialHandler.Update)
	admin.Delete("/testimonials/:id", testimonialHandler.Delete)

	admin.Get("/chat/sessions",     chatHandler.GetAllSessions)
	admin.Get("/chat/sessions/:id", chatHandler.GetSessionMessages)

	admin.Put("/settings",        settingHandler.UpdateSettings)
	admin.Get("/whatsapp-config", settingHandler.GetWhatsappConfig)
	admin.Put("/whatsapp-config", settingHandler.UpdateWhatsappConfig)

	admin.Post("/upload/:folder", uploadHandler.Upload)
	admin.Delete("/upload",       uploadHandler.Delete)
}