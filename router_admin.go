package main

func (r *Router) mappingRouterAdmin() {

	r.route.POST("/api/admin/user/sign-in", r.handleSignInAdmin)
	r1 := r.route.Group("/api/admin", authMiddleware(r))

	r1.GET("/user/page", r.handleListUserPage)
	r1.POST("/upload-image", r.handleUploadImage)

	// user
	r1.GET("/me", r.handleGetMe)
	r1.GET("/user", r.handleGetListUser)
	r1.POST("/user", r.handleCreateUser)
	r1.GET("/user/:id", r.handleGetUser)
	r1.PUT("/user/:id", r.handleUpdateUser)
	r1.DELETE("/user/:id", r.handleDeleteUser)

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
}
