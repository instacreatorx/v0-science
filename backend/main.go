package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kubektl/v0-blog-backend/config"
	"github.com/kubektl/v0-blog-backend/database"
	"github.com/kubektl/v0-blog-backend/handlers"
	"github.com/kubektl/v0-blog-backend/middleware"
	"github.com/kubektl/v0-blog-backend/models"
	"github.com/kubektl/v0-blog-backend/repositories"
	"github.com/kubektl/v0-blog-backend/services"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	if err := database.Connect(cfg.DatabaseURL, cfg.GinMode); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Repositories
	userRepo := repositories.NewUserRepository(database.DB)
	articleRepo := repositories.NewArticleRepository(database.DB)
	authRepo := repositories.NewAuthRepository(database.DB)
	socialRepo := repositories.NewSocialRepository(database.DB)
	teamRepo := repositories.NewTeamRepository(database.DB)

	// Services
	authService := services.NewAuthService(cfg, userRepo, authRepo)
	userService := services.NewUserService(userRepo)
	articleService := services.NewArticleService(articleRepo, socialRepo)
	socialService := services.NewSocialService(socialRepo, articleRepo)
	teamService := services.NewTeamService(teamRepo, userRepo, articleRepo)

	h := handlers.NewHandlers(authService, userService, articleService, socialService, teamService)
	authMW := middleware.NewAuthMiddleware(authService)

	r := gin.Default()
	// r.Use(cors.New(cors.Config{
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return cfg.IsAllowedOrigin(origin)
	// 	},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/send-otp", h.Auth.SendOTP)
			auth.POST("/verify-otp", h.Auth.VerifyOTP)
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.Refresh)
			auth.POST("/logout", authMW.AuthRequired(), h.Auth.Logout)
		}

		users := api.Group("/users")
		{
			users.GET("/me", authMW.AuthRequired(), h.Auth.GetCurrentUser)
			users.PUT("/me", authMW.AuthRequired(), h.Auth.UpdateUser)
			users.GET("/:id", authMW.OptionalAuth(), h.User.GetUser)
			users.GET("/:id/articles", authMW.OptionalAuth(), h.User.GetUserArticles)
			users.POST("/:id/follow", authMW.AuthRequired(), h.User.FollowUser)
		}

		articles := api.Group("/articles")
		{
			articles.GET("", authMW.OptionalAuth(), h.Article.GetArticles)
			articles.GET("/trending", authMW.OptionalAuth(), h.Article.GetTrendingArticles)
			articles.GET("/search", authMW.OptionalAuth(), h.Article.SearchArticles)
			articles.GET("/me", authMW.AuthRequired(), h.Article.GetMyArticles)
			articles.GET("/slug/:slug", authMW.OptionalAuth(), h.Article.GetArticleBySlug)
			articles.GET("/:id", authMW.OptionalAuth(), h.Article.GetArticle)
			articles.POST("", authMW.AuthRequired(), h.Article.CreateArticle)
			articles.PUT("/:id", authMW.AuthRequired(), h.Article.UpdateArticle)
			articles.DELETE("/:id", authMW.AuthRequired(), h.Article.DeleteArticle)
			articles.POST("/:id/publish", authMW.AuthRequired(), h.Article.PublishArticle)
			articles.POST("/:id/unpublish", authMW.AuthRequired(), h.Article.UnpublishArticle)
			articles.POST("/:id/archive", authMW.AuthRequired(), h.Article.ArchiveArticle)
			articles.POST("/:id/like", authMW.AuthRequired(), h.Social.LikeArticle)
			articles.POST("/:id/bookmark", authMW.AuthRequired(), h.Social.BookmarkArticle)
			articles.GET("/:id/comments", h.Social.GetComments)
			articles.POST("/:id/comments", authMW.AuthRequired(), h.Social.CreateComment)
			articles.DELETE("/:id/comments/:commentId", authMW.AuthRequired(), h.Social.DeleteComment)
		}

		api.GET("/bookmarks", authMW.AuthRequired(), h.Social.GetBookmarks)
		api.GET("/feed", authMW.AuthRequired(), h.Article.GetFeed)

		teams := api.Group("/teams")
		{
			teams.POST("", authMW.AuthRequired(), h.Team.CreateTeam)
			teams.GET("/mine", authMW.AuthRequired(), h.Team.ListMyTeams)
			teams.GET("/slug/:slug", h.Team.GetTeamBySlug)
			teams.PUT("/:id", authMW.AuthRequired(), h.Team.UpdateTeam)
			teams.GET("/:id/members", h.Team.ListMembers)
			teams.POST("/:id/members", authMW.AuthRequired(), h.Team.AddMember)
			teams.DELETE("/:id/members/:userId", authMW.AuthRequired(), h.Team.RemoveMember)
			teams.POST("/:id/verify-request", authMW.AuthRequired(), h.Team.SubmitVerification)
		}

		admin := api.Group("/admin", authMW.AuthRequired(), middleware.RequireSuperAdmin(func(id int64) (*models.User, error) {
			return userService.GetCurrent(id)
		}))
		{
			admin.GET("/verification-requests", h.Admin.ListVerificationRequests)
			admin.POST("/verification-requests/:id/approve", h.Admin.ApproveVerification)
			admin.POST("/verification-requests/:id/reject", h.Admin.RejectVerification)
		}
	}

	log.Printf("Server starting on port %s (CORS origins: %v)", cfg.Port, cfg.CORSOrigins)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
