package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/handlers"
	"github.com/karima-store/internal/middleware"
)

// RegisterRoutes registers all application routes with proper authentication
func RegisterRoutes(app *fiber.App,
	auth middleware.KratosMiddleware,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	variantHandler *handlers.VariantHandler,
	categoryHandler *handlers.CategoryHandler,
	pricingHandler *handlers.PricingHandler,
	mediaHandler *handlers.MediaHandler,
	checkoutHandler *handlers.CheckoutHandler,
	komerceHandler *handlers.KomerceHandler,
	orderHandler *handlers.OrderHandler,
	whatsappHandler *handlers.WhatsAppHandler,
	swaggerHandler *handlers.SwaggerHandler) {

	// ===================================================================
	// CSRF PROTECTION MIDDLEWARE
	// ===================================================================

	// Apply CSRF middleware globally with excluded paths
	csrfConfig := middleware.CSRFConfig{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieSecure:   true,
		CookieHTTPOnly: false,
		CookieSameSite: "Strict",
		Expiration:     24 * 60 * 60, // 24 hours
		ContextKey:     "token",
		Next: func(c *fiber.Ctx) bool {
			// Skip CSRF for public endpoints and safe methods
			excludedPaths := []string{
				"/api/v1/health",
				"/swagger",
				"/api/v1/pricing",
				"/api/v1/products",
				"/api/v1/variants",
				"/api/v1/categories",
				"/api/v1/shipping",
				"/api/v1/orders/track",
				"/api/v1/whatsapp/webhook",
				"/api/v1/whatsapp/status",
				"/api/v1/whatsapp/webhook-url",
			}

			// Skip safe methods
			if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" || c.Method() == "TRACE" {
				return true
			}

			// Skip excluded paths
			for _, path := range excludedPaths {
				if len(c.Path()) >= len(path) && c.Path()[:len(path)] == path {
					return true
				}
			}

			return false
		},
	}
	app.Use(middleware.CSRF(csrfConfig))

	// ===================================================================
	// PUBLIC ENDPOINTS (No Authentication Required)
	// ===================================================================

	// Health check
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "up",
			"message": "Server is healthy",
		})
	})

	// Prometheus metrics endpoint (production-ready using adaptor pattern)
	app.Get("/metrics", middleware.MetricsHandler())

	// Documentation
	app.Get("/swagger/*", swaggerHandler.ServeSwagger)

	// Pricing routes (Public - Read-only calculations)
	app.Post("/api/v1/pricing/calculate", pricingHandler.CalculatePrice)
	app.Get("/api/v1/pricing/calculate", pricingHandler.CalculatePriceByParams)
	app.Post("/api/v1/pricing/shipping", pricingHandler.CalculateShippingCost)
	app.Post("/api/v1/pricing/order-summary", pricingHandler.CalculateOrderSummary)
	app.Post("/api/v1/pricing/coupons/validate", pricingHandler.ValidateCoupon)
	app.Get("/api/v1/pricing/products/:product_id", pricingHandler.GetPricingInfo)

	// Product browsing (Public - Read-only)
	app.Get("/api/v1/products", productHandler.GetProducts)
	app.Get("/api/v1/products/:id", productHandler.GetProductByID)
	app.Get("/api/v1/products/slug/:slug", productHandler.GetProductBySlug)
	app.Get("/api/v1/products/search", productHandler.SearchProducts)
	app.Get("/api/v1/products/category/:category", productHandler.GetProductsByCategory)
	app.Get("/api/v1/products/featured", productHandler.GetFeaturedProducts)
	app.Get("/api/v1/products/bestsellers", productHandler.GetBestSellers)
	app.Get("/api/v1/products/:id/media", productHandler.GetProductMedia)

	// Variant browsing (Public - Read-only)
	app.Get("/api/v1/variants/:id", variantHandler.GetVariantByID)
	app.Get("/api/v1/products/:product_id/variants", variantHandler.GetVariantsByProductID)

	// Category browsing (Public - Read-only)
	app.Get("/api/v1/categories", categoryHandler.GetAllCategories)

	// Shipping Info (Public - Read-only)
	app.Get("/api/v1/shipping/destination/search", komerceHandler.SearchDestination)
	app.Get("/api/v1/shipping/calculate", komerceHandler.CalculateShippingCost)

	// Order tracking (Public with order number)
	app.Get("/api/v1/orders/track", orderHandler.TrackOrder)

	// WhatsApp webhook (Has its own validation)
	app.Post("/api/v1/whatsapp/webhook", whatsappHandler.ProcessWhatsAppWebhook)
	app.Get("/api/v1/whatsapp/status", whatsappHandler.GetWhatsAppStatus)
	app.Get("/api/v1/whatsapp/webhook-url", whatsappHandler.GetWhatsAppWebhookURL)

	// ===================================================================
	// AUTH ROUTES
	// ===================================================================

	// Auth endpoints (some are public redirects, some protected)
	app.Get("/api/v1/auth/register", authHandler.Register)
	app.Get("/api/v1/auth/login", authHandler.Login)
	app.Get("/api/v1/auth/logout", authHandler.Logout)
	app.Get("/api/v1/auth/me", auth.ValidateToken(), authHandler.Me) // Protected

	// ===================================================================
	// USER MANAGEMENT ENDPOINTS (Admin only)
	// ===================================================================

	// User management (Admin only)
	app.Get("/api/v1/users", auth.ValidateToken(), auth.RequireAdmin(), userHandler.GetUsers)
	app.Get("/api/v1/users/stats", auth.ValidateToken(), auth.RequireAdmin(), userHandler.GetUserStats)
	app.Get("/api/v1/users/me", auth.ValidateToken(), userHandler.GetCurrentUser) // Any authenticated user
	app.Get("/api/v1/users/:id", auth.ValidateToken(), auth.RequireAdmin(), userHandler.GetUser)
	app.Put("/api/v1/users/:id/role", auth.ValidateToken(), auth.RequireAdmin(), userHandler.UpdateUserRole)
	app.Put("/api/v1/users/:id/deactivate", auth.ValidateToken(), auth.RequireAdmin(), userHandler.DeactivateUser)
	app.Put("/api/v1/users/:id/activate", auth.ValidateToken(), auth.RequireAdmin(), userHandler.ActivateUser)

	// ===================================================================
	// AUTHENTICATED USER ENDPOINTS (Requires valid Kratos session)
	// ===================================================================

	// Checkout (Authenticated users)
	app.Post("/api/v1/checkout", auth.ValidateToken(), checkoutHandler.Checkout)

	// Order management (Authenticated users - own orders only)
	app.Get("/api/v1/orders", auth.ValidateToken(), orderHandler.GetOrders)
	app.Get("/api/v1/orders/:id", auth.ValidateToken(), orderHandler.GetOrder)

	// ===================================================================
	// ADMIN ONLY ENDPOINTS (Requires admin role)
	// ===================================================================

	// Product management (Admin only)
	app.Post("/api/v1/products", auth.ValidateToken(), auth.RequireAdmin(), productHandler.CreateProduct)
	app.Put("/api/v1/products/:id", auth.ValidateToken(), auth.RequireAdmin(), productHandler.UpdateProduct)
	app.Delete("/api/v1/products/:id", auth.ValidateToken(), auth.RequireAdmin(), productHandler.DeleteProduct)
	app.Patch("/api/v1/products/:id/stock", auth.ValidateToken(), auth.RequireAdmin(), productHandler.UpdateProductStock)
	app.Post("/api/v1/products/:id/media", auth.ValidateToken(), auth.RequireAdmin(), productHandler.UploadProductMedia)

	// Variant management (Admin only)
	app.Post("/api/v1/variants", auth.ValidateToken(), auth.RequireAdmin(), variantHandler.CreateVariant)
	app.Put("/api/v1/variants/:id", auth.ValidateToken(), auth.RequireAdmin(), variantHandler.UpdateVariant)
	app.Delete("/api/v1/variants/:id", auth.ValidateToken(), auth.RequireAdmin(), variantHandler.DeleteVariant)
	app.Patch("/api/v1/variants/:id/stock", auth.ValidateToken(), auth.RequireAdmin(), variantHandler.UpdateVariantStock)

	// WhatsApp admin operations (Admin only)
	app.Post("/api/v1/whatsapp/send", auth.ValidateToken(), auth.RequireAdmin(), whatsappHandler.SendWhatsAppMessage)
	app.Get("/api/v1/whatsapp/order-created/:order_id", auth.ValidateToken(), auth.RequireAdmin(), whatsappHandler.SendOrderCreatedNotification)
	app.Get("/api/v1/whatsapp/payment-success/:order_id", auth.ValidateToken(), auth.RequireAdmin(), whatsappHandler.SendPaymentSuccessNotification)
	app.Post("/api/v1/whatsapp/test", auth.ValidateToken(), auth.RequireAdmin(), whatsappHandler.SendTestWhatsAppMessage)

	// ===================================================================
	// COMMENTED OUT / FUTURE ENDPOINTS
	// ===================================================================

	// Category management (Admin only - commented out for now)
	// app.Post("/api/v1/categories", auth.ValidateToken(), auth.RequireAdmin(), categoryHandler.CreateCategory)
	// app.Get("/api/v1/categories/:id", categoryHandler.GetCategory)
	// app.Put("/api/v1/categories/:id", auth.ValidateToken(), auth.RequireAdmin(), categoryHandler.UpdateCategory)
	// app.Delete("/api/v1/categories/:id", auth.ValidateToken(), auth.RequireAdmin(), categoryHandler.DeleteCategory)

	// Media management (Admin only - commented out for now)
	// app.Post("/api/v1/media", auth.ValidateToken(), auth.RequireAdmin(), mediaHandler.CreateMedia)
	// app.Get("/api/v1/media", mediaHandler.GetMedia)
	// app.Get("/api/v1/media/:id", mediaHandler.GetMedia)
	// app.Put("/api/v1/media/:id", auth.ValidateToken(), auth.RequireAdmin(), mediaHandler.UpdateMedia)
	// app.Delete("/api/v1/media/:id", auth.ValidateToken(), auth.RequireAdmin(), mediaHandler.DeleteMedia)

	// User management (Admin only - commented out for now)
	// app.Post("/api/v1/users", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewUserHandler(handlers.UserService{}).CreateUser)
	// app.Get("/api/v1/users", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewUserHandler(handlers.UserService{}).GetUsers)
	// app.Get("/api/v1/users/:id", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewUserHandler(handlers.UserService{}).GetUser)
	// app.Put("/api/v1/users/:id", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewUserHandler(handlers.UserService{}).UpdateUser)
	// app.Delete("/api/v1/users/:id", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewUserHandler(handlers.UserService{}).DeleteUser)

	// Cart management (Authenticated users - commented out for now)
	// app.Post("/api/v1/carts", auth.ValidateToken(), handlers.NewCartHandler(handlers.CartService{}).CreateCart)
	// app.Get("/api/v1/carts", auth.ValidateToken(), handlers.NewCartHandler(handlers.CartService{}).GetCart)
	// app.Put("/api/v1/carts/:id", auth.ValidateToken(), handlers.NewCartHandler(handlers.CartService{}).UpdateCart)
	// app.Delete("/api/v1/carts/:id", auth.ValidateToken(), handlers.NewCartHandler(handlers.CartService{}).DeleteCart)
	// app.Post("/api/v1/carts/:id/items", auth.ValidateToken(), handlers.NewCartHandler(handlers.CartService{}).AddCartItem)
	// app.Delete("/api/v1/carts/:id/items/:item_id", auth.ValidateToken(), handlers.NewCartHandler(handlers.CartService{}).RemoveCartItem)

	// Wishlist management (Authenticated users - commented out for now)
	// app.Post("/api/v1/wishlists", auth.ValidateToken(), handlers.NewWishlistHandler(handlers.WishlistService{}).CreateWishlist)
	// app.Get("/api/v1/wishlists", auth.ValidateToken(), handlers.NewWishlistHandler(handlers.WishlistService{}).GetWishlists)
	// app.Put("/api/v1/wishlists/:id", auth.ValidateToken(), handlers.NewWishlistHandler(handlers.WishlistService{}).UpdateWishlist)
	// app.Delete("/api/v1/wishlists/:id", auth.ValidateToken(), handlers.NewWishlistHandler(handlers.WishlistService{}).DeleteWishlist)
	// app.Post("/api/v1/wishlists/:id/items", auth.ValidateToken(), handlers.NewWishlistHandler(handlers.WishlistService{}).AddWishlistItem)
	// app.Delete("/api/v1/wishlists/:id/items/:item_id", auth.ValidateToken(), handlers.NewWishlistHandler(handlers.WishlistService{}).RemoveWishlistItem)

	// Flash sale management (Admin only - commented out for now)
	// app.Post("/api/v1/flash-sales", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).CreateFlashSale)
	// app.Get("/api/v1/flash-sales", handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).GetFlashSales)
	// app.Put("/api/v1/flash-sales/:id", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).UpdateFlashSale)
	// app.Delete("/api/v1/flash-sales/:id", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).DeleteFlashSale)

	// Review management (Authenticated users can create, admin can moderate)
	// app.Post("/api/v1/reviews", auth.ValidateToken(), handlers.NewReviewHandler(handlers.ReviewService{}).CreateReview)
	// app.Get("/api/v1/reviews", handlers.NewReviewHandler(handlers.ReviewService{}).GetReviews)
	// app.Put("/api/v1/reviews/:id", auth.ValidateToken(), handlers.NewReviewHandler(handlers.ReviewService{}).UpdateReview)
	// app.Delete("/api/v1/reviews/:id", auth.ValidateToken(), auth.RequireAdmin(), handlers.NewReviewHandler(handlers.ReviewService{}).DeleteReview)

	// Komerce integration (Admin only - commented out for now)
	// app.Post("/api/v1/komerce/orders", auth.ValidateToken(), auth.RequireAdmin(), komerceHandler.CreateOrder)
	// app.Put("/api/v1/komerce/orders/cancel", auth.ValidateToken(), auth.RequireAdmin(), komerceHandler.CancelOrder)
	// app.Post("/api/v1/komerce/pickup", auth.ValidateToken(), auth.RequireAdmin(), komerceHandler.RequestPickup)
	// app.Post("/api/v1/komerce/orders/print-label", auth.ValidateToken(), auth.RequireAdmin(), komerceHandler.PrintLabel)

	// Order admin operations (Admin only - commented out for now)
	// app.Put("/api/v1/checkout/:id/confirm", auth.ValidateToken(), auth.RequireAdmin(), checkoutHandler.ConfirmOrder)
	// app.Put("/api/v1/checkout/:id/cancel", auth.ValidateToken(), auth.RequireAdmin(), checkoutHandler.CancelOrder)
	// app.Put("/api/v1/checkout/:id/ship", auth.ValidateToken(), auth.RequireAdmin(), checkoutHandler.ShipOrder)
	// app.Put("/api/v1/checkout/:id/deliver", auth.ValidateToken(), auth.RequireAdmin(), checkoutHandler.DeliverOrder)
	// app.Put("/api/v1/checkout/:id/refund", auth.ValidateToken(), auth.RequireAdmin(), checkoutHandler.RefundOrder)
}
