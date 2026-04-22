package main

import (
	authHandler "cinema-booking-api/internal/auth/delivery/http"
	authRepo "cinema-booking-api/internal/auth/repository"
	authUsecase "cinema-booking-api/internal/auth/usecase"
	"cinema-booking-api/internal/database"
	movieHandler "cinema-booking-api/internal/movie/delivery/http"
	movieDomain "cinema-booking-api/internal/movie/domain"
	movieRepo "cinema-booking-api/internal/movie/repository"
	movieUsecase "cinema-booking-api/internal/movie/usecase"
	userHandler "cinema-booking-api/internal/user/delivery/http"
	userDomain "cinema-booking-api/internal/user/domain"
	userRepo "cinema-booking-api/internal/user/repository"
	userUsecase "cinema-booking-api/internal/user/usecase"
	"cinema-booking-api/pkg/jwt"
	"cinema-booking-api/pkg/logger"
	"cinema-booking-api/pkg/middleware"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Logger
	logger.InitLogger()

	// Initialize Database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Run Database Migration
	if err := database.Migrate(db, &movieDomain.Movie{}, &userDomain.User{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Initialize JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}
	jwtInstance := jwt.NewJWT(
		jwtSecret,
		time.Hour,      // Access token expiry: 1 hour
		time.Hour*24*7, // Refresh token expiry: 7 days
	)

	// Initialize Auth Module
	authRepository := authRepo.NewAuthRepository(db)
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		googleRedirectURL = "http://localhost:5050/api/v1/auth/google/callback"
	}
	authUC := authUsecase.NewAuthUsecase(
		authRepository,
		jwtInstance,
		googleClientID,
		googleSecret,
		googleRedirectURL,
	)
	authHD := authHandler.NewAuthHandler(authUC)

	// Initialize Movie Module
	movieRepository := movieRepo.NewMovieRepository(db)
	movieUC := movieUsecase.NewMovieUsecase(movieRepository)
	movieHD := movieHandler.NewMovieHandler(movieUC)

	// Initialize User Module
	userRepository := userRepo.NewUserRepository(db)
	userUC := userUsecase.NewUserUsecase(userRepository)
	userHD := userHandler.NewUserHandler(userUC)

	// Initialize Auth Middleware
	authMiddleware := middleware.NewAuthMiddleware(func(token string) (*userDomain.User, error) {
		return authUC.ValidateToken(token)
	})

	// Initialize RBAC
	rbac := middleware.NewRBAC()
	// Register permissions
	rbac.RegisterPermission("/api/v1/movies", middleware.RoleAdmin)
	rbac.RegisterPermission("/api/v1/movies/:id", middleware.RoleAdmin)
	rbac.RegisterPermission("/api/v1/users", middleware.RoleAdmin)
	rbac.RegisterPermission("/api/v1/users/:username", middleware.RoleAdmin, middleware.RoleUser)

	r := gin.New() // Use New() to avoid default middleware
	r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":   "Hello, Cinema Booking API",
				"db_status": "connected",
			})
		})

		// Auth Routes (public)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHD.Register)
			authRoutes.POST("/login", authHD.Login)
			authRoutes.GET("/google", authHD.GoogleLogin)
			authRoutes.GET("/google/callback", authHD.GoogleCallback)
			authRoutes.POST("/refresh", authHD.RefreshToken)
		}

		movieRoutes := api.Group("/movies")
		{
			movieRoutes.GET("/", movieHD.GetAllMovies)
			movieRoutes.GET("/:id", movieHD.GetMovieByID)
		}

		// Protected Routes (require auth)
		protectedRoutes := api.Group("")
		protectedRoutes.Use(authMiddleware.Authenticate())
		{
			// Movie Routes - admin only
			movieRoutes := protectedRoutes.Group("/movies")
			movieRoutes.Use(rbac.Guard())
			{
				movieRoutes.POST("/", movieHD.CreateMovie)
				movieRoutes.PUT("/:id", movieHD.UpdateMovie)
				movieRoutes.DELETE("/:id", movieHD.DeleteMovie)
			}

			// User Routes
			userRoutes := protectedRoutes.Group("/users")
			{
				userRoutes.GET("/", userHD.GetAllUsers)
				userRoutes.GET("/:username", userHD.GetUserByUsername)
			}

			// Example: Get current user
			protectedRoutes.GET("/me", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(200, gin.H{"user": user})
			})
		}

	}
	r.Run(":5050")
}
