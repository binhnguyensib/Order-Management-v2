package handler

import (
	"net/http"
	"sync"

	"intern-project-v2/config"
	appHandler "intern-project-v2/handler"
	"intern-project-v2/middleware"
	"intern-project-v2/repository/mongodb"
	"intern-project-v2/usecase"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	router *gin.Engine
	once   sync.Once
)

// Main serverless handler
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initializeApp()
	})

	// Handle request
	router.ServeHTTP(w, r)
}

func initializeApp() {
	// Set production mode cho Vercel
	gin.SetMode(gin.ReleaseMode)

	// Initialize database connection
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Initialize cache
	config.InitCache()

	// Setup dependencies
	deps := setupDependencies(db)

	// Setup router
	router = setupRouter(deps)
}

type Dependencies struct {
	CustomerHandler interface {
		GetAll(c *gin.Context)
		GetByID(c *gin.Context)
		Create(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
	}
	ProductHandler interface {
		GetAll(c *gin.Context)
		GetByID(c *gin.Context)
		Create(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
	}
	OrderHandler interface {
		GetAll(c *gin.Context)
		GetByID(c *gin.Context)
		Create(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
	}
	CartHandler interface {
		AddToCart(c *gin.Context)
		GetCartByCustomerId(c *gin.Context)
		UpdateCartItem(c *gin.Context)
		RemoveCartItem(c *gin.Context)
		ClearCart(c *gin.Context)
	}
	AuthHandler interface {
		Register(c *gin.Context)
		Login(c *gin.Context)
	}
}

func setupDependencies(db *config.Database) *Dependencies {
	// Customer dependencies
	customerRepo := mongodb.NewCustomerRepository(db.DB)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	customerHandler := appHandler.NewCustomerHandler(customerUsecase)

	// Product dependencies
	productRepo := mongodb.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := appHandler.NewProductHandler(productUsecase)

	// Order dependencies
	orderRepo := mongodb.NewOrderRepository(db.DB)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := appHandler.NewOrderHandler(orderUsecase)

	// Cart dependencies
	cartRepo := mongodb.NewCartRepository(db.DB)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productRepo, customerRepo)
	cartHandler := appHandler.NewCartHandler(cartUsecase)

	// Auth dependencies
	authRepo := mongodb.NewAuthRepository(db.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := appHandler.NewAuthHandler(authUsecase)

	return &Dependencies{
		CustomerHandler: customerHandler,
		ProductHandler:  productHandler,
		OrderHandler:    orderHandler,
		CartHandler:     cartHandler,
		AuthHandler:     authHandler,
	}
}

func setupRouter(deps *Dependencies) *gin.Engine {
	router := gin.New()

	// ✅ Simplified middleware cho serverless
	router.Use(gin.Recovery())
	router.Use(middleware.SetupCORS())
	router.Use(middleware.RateLimit(3))
	// - middleware.RequestLogging()
	// - gin.Logger()

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running on Vercel",
			"version": "2.0",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running on Vercel",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup API routes
	setupAPIRoutes(router, deps)

	return router
}

func setupAPIRoutes(router *gin.Engine, deps *Dependencies) {
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", deps.AuthHandler.Register)
			auth.POST("/login", deps.AuthHandler.Login)
		}

		// Customer routes
		customers := api.Group("/customers")
		{
			customers.GET("/", deps.CustomerHandler.GetAll)
			customers.GET("/:id", deps.CustomerHandler.GetByID)
			customers.POST("/", deps.CustomerHandler.Create)
			customers.PUT("/:id", deps.CustomerHandler.Update)
			customers.DELETE("/:id", deps.CustomerHandler.Delete)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("/", middleware.CacheMiddleware(15*60, deps.ProductHandler.GetAll))
			products.GET("/:id", middleware.CacheMiddleware(15*60, deps.ProductHandler.GetByID))
			products.POST("/", deps.ProductHandler.Create)
			products.PUT("/:id", deps.ProductHandler.Update)
			products.DELETE("/:id", deps.ProductHandler.Delete)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.GET("/", deps.OrderHandler.GetAll)
			orders.GET("/:id", deps.OrderHandler.GetByID)
			orders.POST("/", deps.OrderHandler.Create)
			orders.PUT("/:id", deps.OrderHandler.Update)
			orders.DELETE("/:id", deps.OrderHandler.Delete)
		}

		// Cart routes (nested under customers)
		carts := customers.Group("/:id")
		{
			carts.PUT("/cart/item", deps.CartHandler.UpdateCartItem)
			carts.DELETE("/cart/item/:product_id", deps.CartHandler.RemoveCartItem)
			carts.DELETE("/cart", deps.CartHandler.ClearCart)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/customers/:id/cart", deps.CartHandler.GetCartByCustomerId)
			protected.POST("/customers/:id/cart/item", deps.CartHandler.AddToCart)
		}
	}
}
