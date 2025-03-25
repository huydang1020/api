package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	ppb "github.com/huyshop/header/permission"
)

func (r *Router) handleListRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.RoleRequest{}
	utils.BindQuery(req, ctx)
	roles, err := r.permSer.ListRoles(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: roles})
}

func (r *Router) handleGetRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	role, err := r.permSer.GetRole(c, &ppb.RoleRequest{Id: id})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: role})
}

func (r *Router) handleCreateRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.Role{}
	ctx.ShouldBindJSON(req)
	_, err := r.permSer.CreateRole(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ppb.Role{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	_, err := r.permSer.UpdateRole(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	_, err := r.permSer.DeleteRole(c, &ppb.Role{Id: id})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}
