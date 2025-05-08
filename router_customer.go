package main

func (r *Router) mappingRouterCustomer() {

	r.route.GET("/api/home", r.handleHome)
	r.route.POST("/api/customer/sign-in", r.handleSignInCustomer)
	r.route.POST("/api/customer/sign-up", r.handleSignUpCustomer)
	r.route.POST("/api/customer/verify-email", r.handleVerifyEmail)
	r.route.POST("/api/customer/send-code", r.handleSendCode)

	// productType
	r.route.GET("/api/product", r.handleListProduct)
	r.route.GET("/api/product/:id", r.handleGetProduct)

	r1 := r.route.Group("/api/customer", authMiddleware(r))

	// user
	r1.GET("/me", r.handleGetMe)
	r1.POST("/sign-out", r.handleSignOutCustomer)

}
