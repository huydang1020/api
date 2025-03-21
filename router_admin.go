package main

func (r *Router) mappingRouterAdmin() {
	r1 := r.route.Group("/admin")

	r1.GET("/user", r.handleGetListUser)
	r1.POST("/user", r.handleCreateUser)
	r1.GET("/user/:id", r.handleGetUser)
	r1.PUT("/user/:id", r.handleUpdateUser)
	r1.DELETE("/user/:id", r.handleDeleteUser)

}
