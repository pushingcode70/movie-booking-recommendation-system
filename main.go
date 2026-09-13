package main

import (
	"log"

	"movie-booking/config"
	"movie-booking/database"
	"movie-booking/handlers"
	"movie-booking/middleware"
	"movie-booking/models"
	"movie-booking/repositories"
	"movie-booking/routes"
	"movie-booking/services"

	"github.com/joho/godotenv"

	"time"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// load all environment variables into appconfig
	config.LoadConfig()

	// connect to postgresql
	database.ConnectDB()

	// auto migrate database tables
	err = database.DB.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Theatre{},
		&models.Screen{},
		&models.Seat{},
		&models.Show{},
		&models.Booking{},
		&models.BookingSeat{},
		&models.Payment{},
		&models.Wishlist{},
		&models.WatchedMovie{},
		&models.Genre{},
		&models.MovieEmbedding{},
	)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database Migrated Successfully!")

	// movie
	movieRepo := repositories.NewMovieRepository(database.DB)
	tmdbService := services.NewTMDBService(movieRepo)
	movieEmbeddingService := services.NewMovieEmbeddingService(tmdbService)
	movieEmbeddingRepo := repositories.NewMovieEmbeddingRepository(database.DB)
	// sync tmdb genres
	if err := tmdbService.SyncGenres(); err != nil {
		log.Printf("Genre synchronization failed: %v", err)
	}

	movieService := services.NewMovieService(movieRepo, tmdbService)
	movieHandler := handlers.NewMovieHandler(movieService)

	// genre
	genreRepo := repositories.NewGenreRepository(database.DB)
	genreService := services.NewGenreService(genreRepo)
	genreHandler := handlers.NewGenreHandler(genreService)

	// user genre
	userGenreRepo := repositories.NewUserGenreRepository(database.DB)
	userGenreService := services.NewUserGenreService(userGenreRepo)
	userGenreHandler := handlers.NewUserGenreHandler(userGenreService)

	// user
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// email
	emailService := services.NewEmailService()

	// auth
	authService := services.NewAuthService(userRepo, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	// theatre
	theatreRepo := repositories.NewTheatreRepository(database.DB)
	theatreService := services.NewTheatreService(theatreRepo)
	theatreHandler := handlers.NewTheatreHandler(theatreService)

	// screen
	screenRepo := repositories.NewScreenRepository(database.DB)
	screenService := services.NewScreenService(screenRepo)
	screenHandler := handlers.NewScreenHandler(screenService)

	// show repository
	showRepo := repositories.NewShowRepository(database.DB)

	// seat
	seatRepo := repositories.NewSeatRepository(database.DB)
	seatService := services.NewSeatService(
		database.DB,
		seatRepo,
		screenRepo,
		showRepo,
	)
	seatHandler := handlers.NewSeatHandler(seatService)

	// wishlist repository
	wishlistRepo := repositories.NewWishlistRepository(database.DB)

	// watched repository
	watchedRepo := repositories.NewWatchedMovieRepository(database.DB)

	// recommendation
	recommendationService := services.NewRecommendationService(userGenreRepo, wishlistRepo, watchedRepo, movieRepo,
		movieEmbeddingService, movieEmbeddingRepo)
	tmdbHandler := handlers.NewTMDBHandler(tmdbService, recommendationService)

	recommendationHandler := handlers.NewRecommendationHandler(recommendationService)
	embeddingHandler := handlers.NewEmbeddingHandler(recommendationService)

	// wishlist
	wishlistService := services.NewWishlistService(
		wishlistRepo,
		watchedRepo,
		tmdbService,
	)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)

	// watched
	watchedService := services.NewWatchedMovieService(
		watchedRepo,
		wishlistRepo,
		tmdbService,
	)
	watchedHandler := handlers.NewWatchedMovieHandler(watchedService)
	// show
	showService := services.NewShowService(showRepo, movieRepo, screenRepo)
	showHandler := handlers.NewShowHandler(showService)

	// payment
	paymentRepo := repositories.NewPaymentRepository(database.DB)

	// booking
	bookingRepo := repositories.NewBookingRepository(database.DB)
	bookingSeatRepo := repositories.NewBookingSeatRepository(database.DB)

	// razorpay
	razorpayService := services.NewRazorpayService()

	adminRepo := repositories.NewAdminRepository(database.DB)
	adminService := services.NewAdminService(adminRepo)
	adminHandler := handlers.NewAdminHandler(adminService)

	bookingService := services.NewBookingService(
		database.DB,
		bookingRepo,
		bookingSeatRepo,
		paymentRepo,
		showRepo,
		seatRepo,
		razorpayService,
	)

	paymentService := services.NewPaymentService(
		database.DB,
		paymentRepo,
		bookingRepo,
		userRepo,
		showRepo,
		movieRepo,
		screenRepo,
		theatreRepo,
		bookingSeatRepo,
		seatRepo,
		razorpayService,
		emailService,
	)

	bookingHandler := handlers.NewBookingHandler(bookingService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// create gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5174",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// register routes
	routes.RegisterRoutes(router, movieHandler)
	routes.RegisterAuthRoutes(router, authHandler, userHandler)
	routes.RegisterUserRoutes(router, userHandler)
	routes.RegisterTheatreRoutes(router, theatreHandler)
	routes.RegisterScreenRoutes(router, screenHandler)
	routes.RegisterSeatRoutes(router, seatHandler)
	routes.RegisterShowRoutes(router, showHandler)
	routes.RegisterBookingRoutes(router, bookingHandler)
	routes.RegisterPaymentRoutes(router, paymentHandler)
	routes.RegisterTMDBRoutes(router, tmdbHandler, middleware.AuthMiddleware(), middleware.AdminMiddleware())
	routes.RegisterAdminRoutes(router, adminHandler)
	routes.RegisterWishlistRoutes(router, wishlistHandler)
	routes.RegisterWatchedMovieRoutes(router, watchedHandler)
	routes.RegisterGenreRoutes(router, genreHandler)
	routes.RegisterUserGenreRoutes(router, userGenreHandler)
	routes.RegisterRecommendationRoutes(router, recommendationHandler)
	routes.RegisterEmbeddingRoutes(router, embeddingHandler)

	// start server
	if err := router.Run(":8000"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
