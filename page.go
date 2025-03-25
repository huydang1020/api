package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	ppb "github.com/huyshop/header/permission"
)

func (r *Router) handleListPage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.PageRequest{}
	utils.BindQuery(req, ctx)
	pages, err := r.permSer.ListPages(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: pages})
}

func (r *Router) handleGetPage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	page, err := r.permSer.GetPage(c, &ppb.PageRequest{Id: id})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: page})
}

func (r *Router) handleCreatePage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.Page{}
	ctx.ShouldBindJSON(req)
	_, err := r.permSer.CreatePage(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdatePage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ppb.Page{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	_, err := r.permSer.UpdatePage(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeletePage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	_, err := r.permSer.DeletePage(c, &ppb.Page{Id: id})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}
