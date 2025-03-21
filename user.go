package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	userpb "github.com/huyshop/header/user"
)

type Response struct {
	Code    int         `json:"code"` // 0: success, -1: error
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Router) handleGetListUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.UserRequest{}
	utils.BindQuery(req, ctx)
	users, err := r.userSer.ListUsers(c, req)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{
		Code:    0,
		Message: "success",
		Data:    users,
	})
}

func (r *Router) handleGetUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: id})
	if err != nil {
		ctx.JSON(500, Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{
		Code:    0,
		Message: "success",
		Data:    user,
	})
}

func (r *Router) handleCreateUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	_, err := r.userSer.CreateUser(c, req)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{
		Code:    0,
		Message: "success",
	})
}

func (r *Router) handleUpdateUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	_, err := r.userSer.UpdateUser(c, req)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{
		Code:    0,
		Message: "success",
	})
}

func (r *Router) handleDeleteUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	_, err := r.userSer.DeleteUser(c, &userpb.User{Id: id})
	if err != nil {
		ctx.JSON(500, Response{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{
		Code:    0,
		Message: "success",
	})
}
