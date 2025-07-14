package main

func (r *Router) mappingRouterCustomer() {

	r.route.GET("/api/customer/home", r.handleHome)
	r.route.POST("/api/customer/sign-in", r.handleSignInCustomer)
	r.route.POST("/api/customer/sign-in-after-verify-otp", r.handleSignInCustomerAfterVerifyOtp)
	r.route.POST("/api/customer/sign-up", r.handleSignUpCustomer)
	r.route.POST("/api/customer/verify-otp", r.handleVerifyOtp)
	r.route.POST("/api/customer/send-otp", r.handleSendOtp)
	r.route.POST("/api/customer/upload-image", r.handleUploadImage)
	r.route.POST("/api/customer/send-reset-password-otp", r.handleSendResetPasswordOtp)
	r.route.POST("/api/customer/reset-password", r.handleResetPassword)

	// plan
	r.route.GET("/api/customer/plan", r.handleListPlan)
	r.route.GET("/api/customer/plan/:id", r.handleGetPlan)

	// category
	r.route.GET("/api/customer/category", r.handleListCategoryCustomer)

	// productType
	r.route.GET("/api/customer/product-type", r.handleListProductTypeCustomer)
	r.route.GET("/api/customer/product-type/:id", r.handleGetProductTypeCustomer)
	r.route.POST("/api/customer/product-type/update-views", r.handleUpdateProductTypeCustomer)

	// store
	r.route.GET("/api/customer/store/:slug", r.handleGetStoreCustomer)

	r1 := r.route.Group("/api/customer", authMiddleware(r))

	// user
	r1.GET("/me", r.handleGetMe)
	r1.POST("/sign-out", r.handleSignOutCustomer)
	r1.POST("/register-seller", r.handleRegisterSellerAccount)
	r1.POST("/update-profile", r.handleUpdateProfile)
	r1.POST("/change-password", r.handleChangePassword)

	// user address
	r1.GET("/user-address", r.handleListUserAddress)
	r1.POST("/user-address", r.handleCreateUserAddress)
	r1.PUT("/user-address/:id", r.handleUpdateUserAddress)
	r1.DELETE("/user-address/:id", r.handleDeleteUserAddress)

	// cart
	r1.POST("/cart", r.handleUpsertCart)
	r1.GET("/cart", r.handleListCart)
	r1.DELETE("/cart", r.handleDeleteAllItemCart)
	r1.DELETE("/cart-item", r.handleDeleteItemCart)

	// point exchange
	r1.GET("/point-exchange", r.handleListPointExchange)

	// order
	r1.POST("/order", r.handleCreateOrder)
	r1.POST("/order/vnpay", r.handleCreateOrderVNpay)
	r1.GET("/order", r.handleListOrder)
	r1.GET("/order/:id", r.handleGetOrder)
	r1.POST("/order/:id/cancel", r.handleCancelOrder)

	// partner registration
	// r1.POST("/partner-registration", r.handleCreatePartnerRegistration)

	// order plan
	r1.POST("/order-plan", r.handleCreateOrderPlan)
	r1.GET("/order-plan", r.handleListOrderPlan)
	r1.GET("/order-plan/:id", r.handleGetOrderPlan)
	r1.POST("order-plan/vnpay", r.handleCreateOrderPlanVNPay)

	// review
	r.route.GET("api/customer/reviews", r.handleListReviews)
	r1.GET("/reviews-of-customer", r.handleListReviewsOfCustomer)
	r1.POST("/reviews", r.handleCreateReviews)
	r1.GET("/reviews/:id", r.handleGetReviews)

	// voucher
	r.route.GET("/api/customer/voucher", r.handleGetListVoucher)
	r.route.GET("/api/customer/voucher/:id", r.handleGetVoucher)
	// danh sách voucher free đã lấy của user
	r1.GET("/voucher/user-voucher/free", r.handleListUserVoucherFree)
	r1.GET("/voucher/user-voucher/:id", r.handleUserVoucher)
	r1.GET("/voucher/user-voucher", r.handleListUserVoucher)
	r1.POST("/voucher/verify-code", r.handleVerifyCodeVoucher)
	r1.POST("/voucher/buy", r.handleBuyVoucher)

	// favorites
	r1.POST("/favorite-product", r.handleAddFavorites)
	r1.GET("/favorite-product", r.handleListFavorites)
	r1.DELETE("/favorite-product", r.handleDeleteAllFavorites)
	r1.DELETE("/favorite-product/:id", r.handleDeleteOneFavorite)
}
