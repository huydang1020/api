package main

import "github.com/gin-gonic/gin"

func (r *Router) mappingRouterAdmin() {
	ro := r.route.Group("/admin")
	ro.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "okeeeee"})
	})
}
