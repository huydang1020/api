package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	userpb "github.com/huyshop/header/user"
)

func (r *Router) handleCreateOrderPlan(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.OrderPlan{}
	ctx.ShouldBindJSON(&req)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	uid := claims.UserId
	if uid == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	if req.Action == "" {
		req.Action = userpb.OrderPlan_create.String()
	}
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: uid})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if user == nil {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user))
		return
	}
	req.UserId = uid
	req.IpAddress = ctx.ClientIP()
	orderPlan, err := r.userSer.CreateOrderPlan(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: orderPlan})
}

func (r *Router) handleListOrderPlan(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.OrderPlanRequest{}
	utils.BindQuery(req, ctx)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	uid := claims.UserId
	if uid == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: uid})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if user == nil {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user))
		return
	}
	resp, err := r.userSer.ListOrderPlan(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetOrderPlan(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	uid := claims.UserId
	if uid == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	resp, err := r.userSer.GetOrderPlan(c, &userpb.OrderPlan{Id: id, UserId: uid})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleListOrderPlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.OrderPlanRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "order_plan", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp, err := r.userSer.ListOrderPlan(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetOrderPlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "order_plan", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp, err := r.userSer.GetOrderPlan(c, &userpb.OrderPlan{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleCreateOrderPlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.OrderPlan{}
	ctx.ShouldBindJSON(req)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	req.UserId = claims.UserId
	req.Action = userpb.OrderPlan_renew.String()
	orderPlan, err := r.userSer.CreateOrderPlan(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: orderPlan})
}

func (r *Router) handleCreateOrderPlanVNPay(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.OrderPlan{}
	ctx.ShouldBindJSON(req)
	_, err := r.userSer.CreateOrderPlanVNpay(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}
