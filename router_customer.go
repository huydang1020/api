package main

func (r *Router) mappingRouterCustomer() {
	// r.route.POST("/api/customer/sign-in", r.handleSignInCustomer)
	// r.route.POST("/api/customer/sign-up", r.handleSignUpCustomer)
	r1 := r.route.Group("/api/customer", authMiddleware(r))

	// user
	r1.POST("/sign-up", r.handleSignUpCustomer)
	r1.POST("/sign-in", r.handleSignInCustomer)
	r1.GET("/me", r.handleGetMe)
	r1.POST("/sign-out", r.handleSignOutCustomer)

	// home
	r.route.GET("/api/home", r.handleHome)

	// productType
	r1.GET("/product", r.handleListProduct)
	r1.GET("/product/:id", r.handleGetProduct)

}
