package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/handlers"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler,
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

	// Pricing routes
	app.Post("/api/v1/pricing/calculate", pricingHandler.CalculatePrice)
	app.Get("/api/v1/pricing/calculate", pricingHandler.CalculatePriceByParams)
	app.Post("/api/v1/pricing/shipping", pricingHandler.CalculateShippingCost)
	app.Post("/api/v1/pricing/order-summary", pricingHandler.CalculateOrderSummary)
	app.Post("/api/v1/pricing/coupons/validate", pricingHandler.ValidateCoupon)
	app.Get("/api/v1/pricing/products/:product_id", pricingHandler.GetPricingInfo)

	// Catalog routes
	app.Post("/api/v1/products", productHandler.CreateProduct)
	app.Get("/api/v1/products", productHandler.GetProducts)
	app.Get("/api/v1/products/:id", productHandler.GetProductByID)
	app.Put("/api/v1/products/:id", productHandler.UpdateProduct)
	app.Delete("/api/v1/products/:id", productHandler.DeleteProduct)
	app.Patch("/api/v1/products/:id/stock", productHandler.UpdateProductStock)
	app.Post("/api/v1/products/:id/media", productHandler.UploadProductMedia)
	app.Get("/api/v1/products/:id/media", productHandler.GetProductMedia)
	app.Get("/api/v1/products/slug/:slug", productHandler.GetProductBySlug)
	app.Get("/api/v1/products/search", productHandler.SearchProducts)
	app.Get("/api/v1/products/category/:category", productHandler.GetProductsByCategory)
	app.Get("/api/v1/products/featured", productHandler.GetFeaturedProducts)
	app.Get("/api/v1/products/bestsellers", productHandler.GetBestSellers)

	// Variant routes
	app.Post("/api/v1/variants", variantHandler.CreateVariant)
	app.Get("/api/v1/variants/:id", variantHandler.GetVariantByID)
	app.Put("/api/v1/variants/:id", variantHandler.UpdateVariant)
	app.Delete("/api/v1/variants/:id", variantHandler.DeleteVariant)
	app.Patch("/api/v1/variants/:id/stock", variantHandler.UpdateVariantStock)
	app.Get("/api/v1/products/:product_id/variants", variantHandler.GetVariantsByProductID)

	// Category routes
	// app.Post("/api/v1/categories", categoryHandler.CreateCategory)
	app.Get("/api/v1/categories", categoryHandler.GetAllCategories)
	// app.Get("/api/v1/categories/:id", categoryHandler.GetCategory)
	// app.Put("/api/v1/categories/:id", categoryHandler.UpdateCategory)
	// app.Delete("/api/v1/categories/:id", categoryHandler.DeleteCategory)

	// Media routes
	// app.Post("/api/v1/media", mediaHandler.CreateMedia)
	// app.Get("/api/v1/media", mediaHandler.GetMedia)
	// app.Get("/api/v1/media/:id", mediaHandler.GetMedia)
	// app.Put("/api/v1/media/:id", mediaHandler.UpdateMedia)
	// app.Delete("/api/v1/media/:id", mediaHandler.DeleteMedia)
	// app.Post("/api/v1/products/:product_id/media", mediaHandler.AddProductMedia)
	// app.Delete("/api/v1/products/:product_id/media/:media_id", mediaHandler.RemoveProductMedia)

	// User routes
	// app.Post("/api/v1/users", handlers.NewUserHandler(handlers.UserService{}).CreateUser)
	// app.Get("/api/v1/users", handlers.NewUserHandler(handlers.UserService{}).GetUsers)
	// app.Get("/api/v1/users/:id", handlers.NewUserHandler(handlers.UserService{}).GetUser)
	// app.Put("/api/v1/users/:id", handlers.NewUserHandler(handlers.UserService{}).UpdateUser)
	// app.Delete("/api/v1/users/:id", handlers.NewUserHandler(handlers.UserService{}).DeleteUser)

	// Cart routes
	// app.Post("/api/v1/carts", handlers.NewCartHandler(handlers.CartService{}).CreateCart)
	// app.Get("/api/v1/carts", handlers.NewCartHandler(handlers.CartService{}).GetCart)
	// app.Get("/api/v1/carts/:id", handlers.NewCartHandler(handlers.CartService{}).GetCart)
	// app.Put("/api/v1/carts/:id", handlers.NewCartHandler(handlers.CartService{}).UpdateCart)
	// app.Delete("/api/v1/carts/:id", handlers.NewCartHandler(handlers.CartService{}).DeleteCart)
	// app.Post("/api/v1/carts/:id/items", handlers.NewCartHandler(handlers.CartService{}).AddCartItem)
	// app.Delete("/api/v1/carts/:id/items/:item_id", handlers.NewCartHandler(handlers.CartService{}).RemoveCartItem)

	// Wishlist routes
	// app.Post("/api/v1/wishlists", handlers.NewWishlistHandler(handlers.WishlistService{}).CreateWishlist)
	// app.Get("/api/v1/wishlists", handlers.NewWishlistHandler(handlers.WishlistService{}).GetWishlists)
	// app.Get("/api/v1/wishlists/:id", handlers.NewWishlistHandler(handlers.WishlistService{}).GetWishlist)
	// app.Put("/api/v1/wishlists/:id", handlers.NewWishlistHandler(handlers.WishlistService{}).UpdateWishlist)
	// app.Delete("/api/v1/wishlists/:id", handlers.NewWishlistHandler(handlers.WishlistService{}).DeleteWishlist)
	// app.Post("/api/v1/wishlists/:id/items", handlers.NewWishlistHandler(handlers.WishlistService{}).AddWishlistItem)
	// app.Delete("/api/v1/wishlists/:id/items/:item_id", handlers.NewWishlistHandler(handlers.WishlistService{}).RemoveWishlistItem)

	// Flash sale routes
	// app.Post("/api/v1/flash-sales", handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).CreateFlashSale)
	// app.Get("/api/v1/flash-sales", handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).GetFlashSales)
	// app.Get("/api/v1/flash-sales/:id", handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).GetFlashSale)
	// app.Put("/api/v1/flash-sales/:id", handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).UpdateFlashSale)
	// app.Delete("/api/v1/flash-sales/:id", handlers.NewFlashSaleHandler(handlers.FlashSaleService{}).DeleteFlashSale)

	// Checkout routes
	app.Post("/api/v1/checkout", checkoutHandler.Checkout)
	// app.Get("/api/v1/checkout/:id", checkoutHandler.GetCheckout)
	// app.Put("/api/v1/checkout/:id/confirm", checkoutHandler.ConfirmOrder)
	// app.Put("/api/v1/checkout/:id/cancel", checkoutHandler.CancelOrder)
	// app.Put("/api/v1/checkout/:id/ship", checkoutHandler.ShipOrder)
	// app.Put("/api/v1/checkout/:id/deliver", checkoutHandler.DeliverOrder)
	// app.Put("/api/v1/checkout/:id/refund", checkoutHandler.RefundOrder)

	// Order routes
	// Order routes
	app.Get("/api/v1/orders", authMiddleware, orderHandler.GetOrders)
	app.Get("/api/v1/orders/track", orderHandler.TrackOrder)
	app.Get("/api/v1/orders/:id", authMiddleware, orderHandler.GetOrder)

	// Review routes
	// app.Post("/api/v1/reviews", handlers.NewReviewHandler(handlers.ReviewService{}).CreateReview)
	// app.Get("/api/v1/reviews", handlers.NewReviewHandler(handlers.ReviewService{}).GetReviews)
	// app.Get("/api/v1/reviews/:id", handlers.NewReviewHandler(handlers.ReviewService{}).GetReview)
	// app.Put("/api/v1/reviews/:id", handlers.NewReviewHandler(handlers.ReviewService{}).UpdateReview)
	// app.Delete("/api/v1/reviews/:id", handlers.NewReviewHandler(handlers.ReviewService{}).DeleteReview)

	// Pricing routes
	app.Post("/api/v1/pricing/calculate", pricingHandler.CalculatePrice)
	app.Get("/api/v1/pricing/calculate", pricingHandler.CalculatePriceByParams)
	app.Post("/api/v1/pricing/shipping", pricingHandler.CalculateShippingCost)
	app.Post("/api/v1/pricing/order-summary", pricingHandler.CalculateOrderSummary)
	app.Post("/api/v1/pricing/coupons/validate", pricingHandler.ValidateCoupon)
	app.Get("/api/v1/pricing/products/:product_id", pricingHandler.GetPricingInfo)

	// Komerce routes
	// app.Post("/api/v1/komerce/calculate", komerceHandler.CalculateShippingCost)
	// app.Post("/api/v1/komerce/orders", komerceHandler.CreateOrder)
	// app.Get("/api/v1/komerce/orders/:order_no", komerceHandler.GetOrderDetail)
	// app.Put("/api/v1/komerce/orders/cancel", komerceHandler.CancelOrder)
	// app.Post("/api/v1/komerce/pickup", komerceHandler.RequestPickup)
	// app.Post("/api/v1/komerce/orders/print-label", komerceHandler.PrintLabel)
	// app.Get("/api/v1/komerce/orders/track", komerceHandler.TrackOrder)

	// WhatsApp Notification routes
	app.Post("/api/v1/whatsapp/send", whatsappHandler.SendWhatsAppMessage)
	app.Get("/api/v1/whatsapp/order-created/:order_id", whatsappHandler.SendOrderCreatedNotification)
	app.Get("/api/v1/whatsapp/payment-success/:order_id", whatsappHandler.SendPaymentSuccessNotification)
	app.Post("/api/v1/whatsapp/webhook", whatsappHandler.ProcessWhatsAppWebhook)
	app.Get("/api/v1/whatsapp/status", whatsappHandler.GetWhatsAppStatus)
	app.Post("/api/v1/whatsapp/test", whatsappHandler.SendTestWhatsAppMessage)
	app.Get("/api/v1/whatsapp/webhook-url", whatsappHandler.GetWhatsAppWebhookURL)

	// Shipping routes (Komerce API)
	app.Get("/api/v1/shipping/destination/search", komerceHandler.SearchDestination)
	app.Get("/api/v1/shipping/calculate", komerceHandler.CalculateShippingCost)

	// Documentation routes
	app.Get("/swagger/*", swaggerHandler.ServeSwagger)

	// Health check
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "up",
			"message": "Server is healthy",
		})
	})
}
