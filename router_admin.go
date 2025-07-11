package main

func (r *Router) mappingRouterAdmin() {

	r.route.POST("/api/admin/user/sign-in", r.handleSignInAdmin)
	r1 := r.route.Group("/api/admin", authMiddleware(r))

	r1.GET("/user/page", r.handleListUserPage)
	r1.POST("/upload-image", r.handleUploadImage)

	// user
	r1.POST("/user/sign-out", r.handleSignOutAdmin)
	r1.GET("/me", r.handleGetMe)
	r1.GET("/user", r.handleGetListUser)
	r1.POST("/user", r.handleCreateUser)
	r1.GET("/user/:id", r.handleGetUser)
	r1.PUT("/user/:id", r.handleUpdateUser)
	r1.DELETE("/user/:id", r.handleDeleteUser)
	r1.POST("/user/create-point-transaction", r.handleCreatePointTransaction)
	r1.GET("/user/point-exchange", r.handleListPointExchangeAdmin)

	// role
	r1.GET("/role", r.handleListRole)
	r1.POST("/role", r.handleCreateRole)
	r1.GET("/role/:id", r.handleGetRole)
	r1.PUT("/role/:id", r.handleUpdateRole)
	r1.DELETE("/role/:id", r.handleDeleteRole)

	// page
	r1.GET("/page", r.handleListPage)
	r1.POST("/page", r.handleCreatePage)
	r1.GET("/page/:id", r.handleGetPage)
	r1.PUT("/page/:id", r.handleUpdatePage)
	r1.DELETE("/page/:id", r.handleDeletePage)

	// store
	r1.GET("/store", r.handleListStore)
	r1.POST("/store", r.handleCreateStore)
	r1.GET("/store/:id", r.handleGetStore)
	r1.PUT("/store/:id", r.handleUpdateStore)
	r1.DELETE("/store/:id", r.handleDeleteStore)

	// partner
	r1.GET("/partner", r.handleListPartner)
	r1.POST("/partner", r.handleCreatePartner)
	r1.GET("/partner/:id", r.handleGetPartner)
	r1.PUT("/partner/:id", r.handleUpdatePartner)
	r1.DELETE("/partner/:id", r.handleDeletePartner)

	// product
	r1.GET("/product", r.handleListProduct)
	r1.GET("/product/:id", r.handleGetProduct)

	// product type
	r1.GET("/product-type", r.handleListProductType)
	r1.POST("/product-type", r.handleCreateProductType)
	r1.GET("/product-type/:id", r.handleGetProductType)
	r1.PUT("/product-type/:id", r.handleUpdateProductTypeAdmin)
	r1.DELETE("/product-type/:id", r.handleDeleteProductType)
	r1.PUT("product-type/state/:id", r.handleUpdateStateProductType)

	// product category
	r1.GET("/category", r.handleListCategory)
	r1.POST("/category", r.handleCreateCategory)
	r1.GET("/category/:id", r.handleGetCategory)
	r1.PUT("/category/:id", r.handleUpdateCategory)
	r1.DELETE("/category/:id", r.handleDeleteCategory)

	// banner
	r1.GET("/banner", r.handleListBanner)
	r1.POST("/banner", r.handleCreateBanner)
	r1.GET("/banner/:id", r.handleGetBanner)
	r1.PUT("/banner/:id", r.handleUpdateBanner)
	r1.DELETE("/banner/:id", r.handleDeleteBanner)

	// order
	r1.PUT("/order/:id", r.handleUpdateStateOrderAdmin)
	r1.GET("/order", r.handleListOrderAdmin)
	r1.GET("/order/:id", r.handleGetOrderAdmin)

	// plan
	r1.GET("/plan", r.handleListPlanAdmin)
	r1.POST("/plan", r.handleCreatePlanAdmin)
	r1.GET("/plan/:id", r.handleGetPlanAdmin)
	r1.PUT("/plan/:id", r.handleUpdatePlanAdmin)
	r1.DELETE("/plan/:id", r.handleDeletePlanAdmin)

	// order plan
	r1.GET("/order-plan", r.handleListOrderPlanAdmin)
	r1.POST("/order-plan", r.handleCreateOrderPlanAdmin)
	r1.GET("/order-plan/:id", r.handleGetOrderPlanAdmin)
	r1.POST("order-plan/vnpay", r.handleCreateOrderPlanVNPay)

	// review
	r1.GET("/reviews", r.handleListReviewsByAdmin)
	r1.POST("/reviews/reply", r.handleReplyReviewsByAdmin)
	r1.GET("/review/:id", r.handleGetReviewsByAdmin)
	r1.PUT("/reviews/:id", r.handleUpdateReviews)
	r1.DELETE("/reviews/:id", r.handleDeleteReviews)

	// voucher
	r1.GET("/voucher", r.handleGetListVoucherAdmin)
	r1.POST("/voucher", r.handleCreateVoucher)
	r1.GET("/voucher/:id", r.handleGetOneVoucher)
	r1.PUT("/voucher/:id", r.handleUpdateVoucher)
	r1.DELETE("/voucher/:id", r.handleDeleteVoucher)
	r1.GET("/voucher/user-voucher", r.handleListUserVoucherAdmin)

	// report
	r1.GET("/report/overview", r.handleGetReportOverview)
	r1.GET("/report/revenue", r.handleGetReportRevenue)
	r1.GET("/report/store/revenue", r.handleGetReportRevenueByStore)
	r1.GET("/report/user", r.handleGetReportUser)
	r1.GET("/report/product", r.handleGetReportProduct)
	r1.GET("/report/top-products", r.handleGetReportTopProducts)
}
