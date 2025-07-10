package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	upb "github.com/huyshop/header/user"
	vpb "github.com/huyshop/header/voucher"
)

// voucher
func (r *Router) handleCreateVoucher(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.Voucher{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "voucher", "c"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.RemainingQuantity = req.TotalQuantity
	req.PartnerId = claims.PartnerId
	vou, err := r.voucherSer.CreateVoucher(c, req)
	if err != nil {
		log.Println("insert voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleGetListVoucherAdmin(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.VoucherRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "voucher", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType != upb.Partner_admin.String() {
		req.PartnerId = claims.PartnerId
	}
	vou, err := r.voucherSer.ListVouchers(c, req)
	if err != nil {
		log.Println("list voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, vou := range vou.Vouchers {
		partner, err := r.userSer.GetPartner(c, &upb.PartnerRequest{Id: vou.PartnerId})
		if err != nil {
			log.Println("get partner err:", err)
			continue
		}
		vou.Partner = partner
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleGetListVoucher(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.VoucherRequest{}
	utils.BindQuery(req, ctx)
	req.State = vpb.Voucher_active.String()
	vou, err := r.voucherSer.ListVouchersCustomer(c, req)
	if err != nil {
		log.Println("list voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, data := range vou.Vouchers {
		partner, err := r.userSer.GetPartner(c, &upb.PartnerRequest{Id: data.PartnerId})
		if err != nil {
			log.Println("get partner err:", err)
			continue
		}
		data.Partner = partner
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleGetVoucher(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	vou, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: id})
	if err != nil {
		log.Println("get voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	partner, err := r.userSer.GetPartner(c, &upb.PartnerRequest{Id: vou.PartnerId})
	if err != nil {
		log.Println("get partner err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	vou.Partner = partner
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleListUserVoucherFree(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &vpb.UserVoucherRequest{}
	utils.BindQuery(req, ctx)
	req.Type = vpb.Voucher_free.String()
	req.UserId = claims.UserId
	vou, err := r.voucherSer.ListUserVouchers(c, req)
	if err != nil {
		log.Println("list user voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleGetOneVoucher(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "voucher", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	vou, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: id})
	if err != nil {
		log.Println("get one voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleUpdateVoucher(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "voucher", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req := &vpb.Voucher{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	req.UpdatedAt = time.Now().Unix()
	vou, err := r.voucherSer.UpdateVoucher(c, req)
	if err != nil {
		log.Println("update voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: vou})
}

func (r *Router) handleDeleteVoucher(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "voucher", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if _, err := r.voucherSer.DeleteVoucher(c, &vpb.Voucher{Id: id}); err != nil {
		log.Println("delete vou err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUserVoucher(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.UserVoucher{}
	utils.BindQuery(req, ctx)
	id := ctx.Param("id")
	uid := claims.UserId
	uv, err := r.voucherSer.GetUserVoucher(c, &vpb.UserVoucher{Id: req.VoucherId, UserId: uid, CodeId: id})
	if err != nil {
		log.Println("user_voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	vou, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: req.VoucherId})
	if err != nil {
		log.Println("voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	uv.Voucher = vou
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: uv})
}

func (r *Router) handleListUserVoucher(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.UserVoucherRequest{}
	utils.BindQuery(req, ctx)
	req.State = ctx.Query("state")
	req.UserId = claims.UserId
	list, err := r.voucherSer.ListUserVouchers(c, req)
	if err != nil {
		log.Println("user_voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, uv := range list.UserVouchers {
		vou, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: uv.VoucherId})
		if err != nil {
			log.Println("voucher err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		uv.Voucher = vou
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: list})
}

func (r *Router) handleVerifyCodeVoucher(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.VerifyCodeRequest{}
	ctx.ShouldBindJSON(req)
	_, err := r.voucherSer.VerifyCode(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListUserVoucherAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.UserVoucher{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "user_voucher", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	list, err := r.voucherSer.ListUserVouchers(c, &vpb.UserVoucherRequest{})
	if err != nil {
		log.Println("user_voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, uv := range list.UserVouchers {
		vou, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: uv.VoucherId})
		if err != nil {
			log.Println("voucher err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		uv.Voucher = vou
		user, err := r.userSer.GetUser(c, &upb.UserRequest{Id: uv.UserId})
		if err != nil {
			log.Println("user err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		uv.User = user
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: list})
}

func (r *Router) handleBuyVoucher(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &vpb.UserVoucher{}
	ctx.ShouldBindJSON(req)
	if req.VoucherId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_invalid_voucher_id))
		return
	}
	voucher, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: req.VoucherId})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if voucher.State != vpb.Voucher_active.String() {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_voucher_not_active))
		return
	}
	if voucher.TotalQuantity <= 0 {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_total_quantity_empty))
		return
	}
	if voucher.RemainingQuantity <= 0 {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_remaining_quantity_empty))
		return
	}
	if voucher.StartAt > time.Now().Unix() {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_start_at_in_the_future))
		return
	}
	if voucher.EndAt < time.Now().Unix() {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_end_at_in_the_past))
		return
	}
	req.UserId = claims.UserId
	if voucher.PointExchange > 0 {
		userpoint, err := r.userSer.GetUserPoint(c, &upb.UserPointRequest{UserId: claims.UserId})
		if err != nil {
			log.Println("user err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		if userpoint.Points < int64(voucher.PointExchange) {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_point_not_enough))
			return
		}
		if _, err := r.userSer.CreatePointExchange(c, &upb.PointExchange{
			ReceiverId:  claims.UserId,
			Points:      -int64(voucher.PointExchange),
			Description: fmt.Sprintf("Trừ %d điểm để đổi ưu đãi: %v", voucher.PointExchange, voucher.Name),
		}); err != nil {
			log.Println("point_exchange err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
	}
	uv, err := r.voucherSer.CreateUserVoucher(c, req)
	if err != nil {
		log.Println("user_voucher err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: uv})
}
