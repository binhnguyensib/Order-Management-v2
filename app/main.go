package app

import (
	"intern-project-v2/config"
	_ "intern-project-v2/docs"
	"intern-project-v2/handler"
	"intern-project-v2/middleware"
	"intern-project-v2/repository/mongodb"
	"intern-project-v2/usecase"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Order Management API
// @version 2.0
// @description This is a sample server for managing orders, customers, products, and carts.
// @host order-management-v2.vercel.app
// @BasePath /api
func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func Run() {

	db, err := config.ConnectDB()
	{
		if err != nil {
			panic("Failed to connect to database: " + err.Error())
		}
	}
	customerRepo := mongodb.NewCustomerRepository(db.DB)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	productRepo := mongodb.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	orderRepo := mongodb.NewOrderRepository(db.DB)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	cartRepo := mongodb.NewCartRepository(db.DB)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productRepo, customerRepo)
	cartHandler := handler.NewCartHandler(cartUsecase)

	authRepo := mongodb.NewAuthRepository(db.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(middleware.SetupCORS())
	router.Use(middleware.RateLimit(3))
	router.Use(middleware.RequestLogging())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	config.InitCache() // Initialize cache store

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		customers := api.Group("/customers")
		{
			customers.GET("/", customerHandler.GetAll)
			customers.GET("/:id", customerHandler.GetByID)
			customers.POST("/", customerHandler.Create)
			customers.PUT("/:id", customerHandler.Update)
			customers.DELETE("/:id", customerHandler.Delete)
		}
		products := api.Group("/products")

		{
			products.GET("/", middleware.CacheMiddleware(time.Minute*15, productHandler.GetAll))
			products.GET("/:id", middleware.CacheMiddleware(time.Minute*15, productHandler.GetByID))
			products.POST("/", productHandler.Create)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}
		orders := api.Group("/orders")
		{
			orders.GET("/", orderHandler.GetAll)
			orders.GET("/:id", orderHandler.GetByID)
			orders.POST("/", orderHandler.Create)
			orders.PUT("/:id", orderHandler.Update)
			orders.DELETE("/:id", orderHandler.Delete)
		}
		carts := customers.Group("/:id")
		{
			//carts.POST("/cart/item", cartHandler.AddToCart)
			//carts.GET("/cart", cartHandler.GetCartByCustomerId)
			carts.PUT("/cart/item", cartHandler.UpdateCartItem)
			carts.DELETE("/cart/item/:product_id", cartHandler.RemoveCartItem)
			carts.DELETE("/cart", cartHandler.ClearCart)
		}
	}
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/customers/:id/cart", cartHandler.GetCartByCustomerId)
		protected.POST("/customers/:id/cart/item", cartHandler.AddToCart)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
