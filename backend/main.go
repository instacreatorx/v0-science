package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/handlers"
	"github.com/kubektl/v0-blog-backend/middleware"
)

func main() {
	// Load .env file if exists
	godotenv.Load()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Setup Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://*.vercel.app", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	r.Use(cors.Default()) // All origins allowed by default

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/send-otp", handlers.SendOTP)
			auth.POST("/verify-otp", handlers.VerifyOTP)
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/me", middleware.AuthRequired(), handlers.GetCurrentUser)
			users.PUT("/me", middleware.AuthRequired(), handlers.UpdateUser)
			users.GET("/:id", handlers.GetUser)
			users.GET("/:id/articles", handlers.GetUserArticles)
			users.POST("/:id/follow", middleware.AuthRequired(), handlers.FollowUser)
		}

		// Article routes
		articles := api.Group("/articles")
		{
			articles.GET("", handlers.GetArticles)
			articles.GET("/trending", handlers.GetTrendingArticles)
			articles.GET("/search", handlers.SearchArticles)
			articles.GET("/:id", handlers.GetArticle)
			articles.POST("", middleware.AuthRequired(), handlers.CreateArticle)
			articles.PUT("/:id", middleware.AuthRequired(), handlers.UpdateArticle)
			articles.DELETE("/:id", middleware.AuthRequired(), handlers.DeleteArticle)

			// Article interactions
			articles.POST("/:id/like", middleware.AuthRequired(), handlers.LikeArticle)
			articles.POST("/:id/bookmark", middleware.AuthRequired(), handlers.BookmarkArticle)

			// Comments
			articles.GET("/:id/comments", handlers.GetComments)
			articles.POST("/:id/comments", middleware.AuthRequired(), handlers.CreateComment)
			articles.DELETE("/:id/comments/:commentId", middleware.AuthRequired(), handlers.DeleteComment)
		}

		// Bookmarks
		api.GET("/bookmarks", middleware.AuthRequired(), handlers.GetBookmarks)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
