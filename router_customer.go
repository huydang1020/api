package main

func (r *Router) mappingRouterCustomer() {

	r.route.GET("/api/customer/home", r.handleHome)
	r.route.POST("/api/customer/sign-in", r.handleSignInCustomer)
	r.route.POST("/api/customer/sign-up", r.handleSignUpCustomer)
	r.route.POST("/api/customer/verify-email", r.handleVerifyEmail)
	r.route.POST("/api/customer/send-otp", r.handleSendOtp)

	// productType
	r.route.GET("/api/customer/product-type", r.handleListProductTypeCustomer)
	r.route.GET("/api/customer/product-type/:id", r.handleGetProductTypeCustomer)

	r1 := r.route.Group("/api/customer", authMiddleware(r))

	// user
	r1.GET("/me", r.handleGetMe)
	r1.POST("/sign-out", r.handleSignOutCustomer)

	// cart
	r1.POST("/cart", r.handleUpsertCart)
	r1.GET("/cart", r.handleListCart)
	r1.DELETE("/cart", r.handleDeleteAllItemCart)
	r1.DELETE("/cart-item", r.handleDeleteItemCart)

	// order customer
	r1.POST("/order", r.handleCreateOrder)
	r1.POST("/order/vnpay", r.handleCreateOrderVNpay)
	r1.GET("/order", r.handleListOrder)
	r1.GET("/order/:id", r.handleGetOrder)
	r1.PUT("/order/:id", r.handleCancelOrder)
}
