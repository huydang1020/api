package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
	upb "github.com/huyshop/header/user"
)

func (r *Router) handleUpsertCart(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := []*ptpb.OrderDetail{}
	ctx.ShouldBindJSON(&req)
	user_id := claims.UserId
	_, err := r.productSer.AddToCart(c, &ptpb.Cart{Item: req, UserId: user_id})
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListCart(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Cart{}
	ctx.ShouldBindQuery(&req)
	req.UserId = claims.UserId
	cart, err := r.productSer.ListCart(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: cart})
}

func (r *Router) handleDeleteAllItemCart(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	user_id := claims.UserId
	defer cancel()
	if user_id == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	_, err := r.productSer.DeleteCart(c, &ptpb.Cart{UserId: user_id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteItemCart(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	user_id := claims.UserId
	req := []*ptpb.OrderDetail{}
	ctx.ShouldBindJSON(&req)
	if user_id == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	_, err := r.productSer.DeleteCartItem(c, &ptpb.Cart{Item: req, UserId: user_id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleCreateOrder(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Order{}
	ctx.ShouldBindJSON(&req)
	user_id := claims.UserId
	if user_id == "" {
		log.Println("user_id", user_id)
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	req.UserId = user_id
	user, err := r.userSer.GetUser(c, &upb.UserRequest{Id: user_id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if user == nil {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user))
		return
	}
	if user.State != upb.User_active.String() {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_access_denied))
		return
	}
	req.IpAddress = ctx.ClientIP()
	order, err := r.productSer.CreateOrder(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if order == nil {
		utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: order})
}

func (r *Router) handleCreateOrderVNpay(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Order{}
	ctx.ShouldBindJSON(&req)
	req.UserId = claims.UserId
	_, err := r.productSer.CreateOrderVNpay(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListOrder(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.OrderRequest{}
	ctx.ShouldBindQuery(&req)
	orders, err := r.productSer.ListOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: orders})
}

func (r *Router) handleGetOrder(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.OrderRequest{}
	ctx.ShouldBindQuery(&req)
	order, err := r.productSer.GetOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: order})
}

func (r *Router) handleCancelOrder(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Order{}
	ctx.ShouldBindJSON(&req)
	order, err := r.productSer.CancelOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: order})
}
