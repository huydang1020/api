package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
	"github.com/huyshop/header/user"
)

func (r *Router) handleListReviews(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ReviewsRequest{}
	utils.BindQuery(req, ctx)
	resp, err := r.productSer.ListReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, rv := range resp.GetReviews() {
		user, err := r.userSer.GetUser(c, &user.UserRequest{Id: rv.GetUserId()})
		if err != nil {
			log.Println("get user err:", err, rv.GetUserId())
			continue
		}
		rv.User = user
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleListReviewsOfCustomer(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ReviewsRequest{}
	utils.BindQuery(req, ctx)
	req.UserId = claims.UserId
	resp, err := r.productSer.ListReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleCreateReviews(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Reviews{}
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
	_, err := r.productSer.CreateReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleGetReviews(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.ReviewsRequest{Id: id, UserId: claims.UserId}
	resp, err := r.productSer.GetReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleUpdateReviews(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.Reviews{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "reviews", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	resp, err := r.productSer.UpdateReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

// handle review by admin
func (r *Router) handleListReviewsByAdmin(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ReviewsRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "reviews", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	listOrder, err := r.productSer.ListOrder(c, &ptpb.OrderRequest{PartnerId: claims.PartnerId})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if listOrder == nil {
		utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
	}
	var ordIds []string
	for _, ord := range listOrder.GetOrders() {
		ordIds = append(ordIds, ord.Id)
	}
	req.OrderIds = ordIds
	resp, err := r.productSer.ListReviews(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, rv := range resp.GetReviews() {
		user, err := r.userSer.GetUser(c, &user.UserRequest{Id: rv.GetUserId()})
		if err != nil {
			log.Println("get user err:", err, rv.GetUserId())
			continue
		}
		rv.User = user
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetReviewsByAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "reviews", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req := &ptpb.ReviewsRequest{Id: id}
	resp, err := r.productSer.GetReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleReplyReviewsByAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.Reviews{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "reviews", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	resp, err := r.productSer.ReplyReviews(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleDeleteReviews(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "reviews", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req := &ptpb.Reviews{Id: id}
	if _, err := r.productSer.DeleteReviews(c, req); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}
