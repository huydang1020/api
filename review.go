package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
)

// handle review by user
func (r *Router) handleListRewivewByUser(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ReviewRequest{}
	ctx.ShouldBindJSON(req)
	req.UserId = claims.UserId
	resp, err := r.productSer.ListReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleCreateReview(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Review{}
	ctx.ShouldBindJSON(req)
	req.UserId = claims.UserId
	if req.ProductId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_product_id))
		return
	}
	if req.OrderId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_order_id))
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_invalid_rating))
		return
	}
	resp, err := r.productSer.CreateReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetReview(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.ReviewRequest{Id: id, UserId: claims.UserId}
	resp, err := r.productSer.GetReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleUpdateReview(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.Review{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "review", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	req.UserId = claims.UserId
	resp, err := r.productSer.UpdateReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

// handle review by admin
func (r *Router) handleListReviewByAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ReviewRequest{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "review", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp, err := r.productSer.ListReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetReviewByAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "review", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req := &ptpb.ReviewRequest{Id: id}
	resp, err := r.productSer.GetReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleReplyReviewByAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.Review{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "review", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	resp, err := r.productSer.ReplyReview(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleDeleteReview(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "review", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req := &ptpb.Review{Id: id}
	if _, err := r.productSer.DeleteReview(c, req); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}
