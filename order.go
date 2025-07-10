package main

import (
	"errors"
	"log"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
	upb "github.com/huyshop/header/user"
	vpb "github.com/huyshop/header/voucher"
)

func (r *Router) handleUpsertCart(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := []*ptpb.ProductOrdered{}
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
	utils.BindQuery(req, ctx)
	req.UserId = claims.UserId
	log.Println("req: ", req)
	resp := &ptpb.CartDetail{}
	cart, err := r.productSer.ListCart(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if len(cart.Item) == 0 {
		utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
		return
	}
	mapOrd := make(map[string][]*ptpb.ProductOrdered)
	var storeIds []string
	for _, item := range cart.Item {
		pty, err := r.productSer.GetProductType(c, &ptpb.ProductTypeRequest{Id: item.Product.ProductTypeId})
		if err != nil {
			log.Println("err ", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		storeIds = append(storeIds, pty.StoreId)
		mapOrd[pty.StoreId] = append(mapOrd[pty.StoreId], item)
	}
	stores, err := r.userSer.ListStore(c, &upb.StoreRequest{Ids: storeIds})
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, sto := range stores.Stores {
		if sto == nil {
			continue
		}
		if sto.State != upb.Store_active.String() {
			continue
		}
		detail := &ptpb.ProductStore{
			Id:          sto.GetId(),
			Name:        sto.GetName(),
			Products:    mapOrd[sto.GetId()],
			Logo:        sto.GetLogo(),
			Address:     sto.GetAddress(),
			PhoneNumber: sto.GetPhoneNumber(),
			State:       sto.GetState(),
			Description: sto.GetDescription(),
		}
		resp.Stores = append(resp.Stores, detail)
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
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
	req := []*ptpb.ProductOrdered{}
	ctx.ShouldBindJSON(&req)
	if user_id == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_id))
		return
	}
	if len(req) == 0 {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_invalid_item_cart))
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
	log.Println("req:", req)
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
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_account_not_activated))
		return
	}
	req.IpAddress = ctx.ClientIP()
	if req.CodeId != "" {
		uv, err := r.voucherSer.GetUserVoucher(ctx, &vpb.UserVoucher{CodeId: req.CodeId, UserId: req.UserId})
		if err != nil {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_user_voucher))
			return
		}
		if uv.State != vpb.UserVoucher_got.String() {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_code_state_invalid))
			return
		}
		voucher, err := r.voucherSer.GetVoucher(c, &vpb.Voucher{Id: uv.VoucherId})
		if err != nil {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		if voucher.State != vpb.Voucher_active.String() {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_voucher_not_active))
			return
		}
		if voucher.EndAt < time.Now().Unix() {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_voucher_expired))
			return
		}
		req.Voucher = voucher
	}
	order, err := r.productSer.CreateOrder(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if order == nil {
		utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
		return
	}
	if req.MethodPayment == "cod" && req.CodeId != "" {
		if _, err := r.voucherSer.UpdateUserVoucher(c, &vpb.UserVoucher{
			UserId:    req.UserId,
			State:     vpb.UserVoucher_used.String(),
			CodeId:    req.CodeId,
			VoucherId: req.Voucher.Id,
		}); err != nil {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
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
	if req.CodeId != "" {
		if _, err := r.voucherSer.UpdateUserVoucher(c, &vpb.UserVoucher{
			UserId:    req.UserId,
			State:     vpb.UserVoucher_used.String(),
			CodeId:    req.CodeId,
			VoucherId: req.Voucher.Id,
		}); err != nil {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListOrder(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.OrderRequest{}
	utils.BindQuery(req, ctx)
	req.State = ctx.Query("state")
	req.UserId = claims.UserId
	log.Println("req, ", req)
	orders, err := r.productSer.ListOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: orders})
}

func (r *Router) handleListOrderAdmin(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.OrderRequest{}
	utils.BindQuery(req, ctx)
	req.State = ctx.Query("state")
	if err := r.isCanBeAccess(c, ctx, "order", "r"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	pid := claims.PartnerId
	req.PartnerId = pid
	orders, err := r.productSer.ListOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, order := range orders.Orders {
		if order.UserAddressId != "" {
			userAddress, err := r.userSer.GetUserAddress(c, &upb.UserAddress{Id: order.UserAddressId})
			if err != nil {
				log.Println("err ", err)
				utils.HandleError(LangMappingErr, ctx, err)
				return
			}
			order.UserAddress = userAddress
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: orders})
}

func (r *Router) handleGetOrder(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &ptpb.OrderRequest{}
	id := ctx.Param("id")
	if id == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_order_id))
		return
	}
	req.Id = id
	req.UserId = claims.UserId
	order, err := r.productSer.GetOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: order})
}

func (r *Router) handleGetOrderAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &ptpb.OrderRequest{}
	id := ctx.Param("id")
	if id == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_order_id))
		return
	}
	req.Id = id
	req.PartnerId = claims.PartnerId
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
	id := ctx.Param("id")
	req.State = ptpb.Order_cancelled.String()
	req.Id = id
	_, err := r.productSer.UpdateStateOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateStateOrderAdmin(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	acceptPartnerType := []string{"seller", "admin"}
	if !slices.Contains(acceptPartnerType, claims.PartnerType) {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_access_denied))
		return
	}
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Order{}
	ctx.ShouldBindJSON(&req)
	if err := r.isCanBeAccess(c, ctx, "order", "u"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	id := ctx.Param("id")
	req.Id = id
	_, err := r.productSer.UpdateStateOrder(c, req)
	if err != nil {
		log.Println("err ", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

// func (r *Router) handleUpdateStateOrderAdmin(ctx *gin.Context) {
// 	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
// 	defer cancel()
// 	req := &ptpb.Order{}
// 	ctx.ShouldBindJSON(&req)
// 	if err := r.isCanBeAccess(c, ctx, "order", "r"); err != nil {
// 		utils.HandleError(LangMappingErr, ctx, err)
// 		return
// 	}
// 	id := ctx.Param("id")
// 	if id == "" {
// 		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_order_id))
// 		return
// 	}
// 	req.Id = id
// 	order, err := r.productSer.UpdateStateOrder(c, req)
// 	if err != nil {
// 		log.Println("err ", err)
// 		utils.HandleError(LangMappingErr, ctx, err)
// 		return
// 	}
// 	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: order})
// }
