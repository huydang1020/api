package main

import (
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
	userpb "github.com/huyshop/header/user"
)

// Partner
func (r *Router) handleListPartner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.PartnerRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "partner", "r"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	partners, err := r.userSer.ListPartner(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: partners})
}

func (r *Router) handleGetPartner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "partner", "r"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	partner, err := r.userSer.GetPartner(c, &userpb.PartnerRequest{Id: id})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: partner})
}

func (r *Router) handleCreatePartner(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.Partner{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "partner", "c"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Type = claims.PartnerType
	_, err := r.userSer.CreatePartner(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdatePartner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.Partner{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "partner", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	_, err := r.userSer.UpdatePartner(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeletePartner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "partner", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.userSer.DeletePartner(c, &userpb.Partner{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

// Store
func (r *Router) handleListStore(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.StoreRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "store", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType != userpb.Partner_admin.String() {
		req.PartnerId = claims.PartnerId
	}
	stores, err := r.userSer.ListStore(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: stores})
}

func (r *Router) handleGetStore(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "store", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	store, err := r.userSer.GetStore(c, &userpb.StoreRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: store})
}

func (r *Router) handleCreateStore(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.Store{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "store", "c"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.PartnerId = claims.PartnerId
	log.Println("req", req)
	log.Println("claims", claims)
	var limitSto int32
	if claims.PartnerType == userpb.Partner_seller.String() {
		part, err := r.userSer.GetPartner(c, &userpb.PartnerRequest{Id: claims.PartnerId})
		if err != nil {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		log.Println("part", part)
		if part.PlanExpiredAt < time.Now().Unix() {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_plan_expired))
			return
		}
		stores, err := r.userSer.CountStore(c, &userpb.StoreRequest{PartnerId: claims.PartnerId})
		if err != nil {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		log.Println("stores", stores)
		if int32(stores.Count) >= part.MaxStoresAllowed {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_store_limit_reached))
			return
		}
		limitSto = part.CurrentStoresCount
	}
	_, err := r.userSer.CreateStore(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	// update partner
	if claims.PartnerType == userpb.Partner_seller.String() {
		_, err = r.userSer.UpdatePartner(c, &userpb.Partner{
			Id:                 claims.PartnerId,
			CurrentStoresCount: limitSto + 1,
			UpdatedAt:          time.Now().Unix(),
		})
		if err != nil {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateStore(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.Store{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "store", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	req.PartnerId = claims.PartnerId
	_, err := r.userSer.UpdateStore(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteStore(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "store", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.userSer.DeleteStore(c, &userpb.Store{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

type StoreCustomer struct {
	Store        *userpb.Store      `json:"store"`
	ProductTypes *ptpb.ProductTypes `json:"product_types"`
}

func (r *Router) handleGetStoreCustomer(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	slug := ctx.Param("slug")
	store, err := r.userSer.GetStore(c, &userpb.StoreRequest{Slug: slug})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	productTypes, err := r.productSer.ListProductType(c, &ptpb.ProductTypeRequest{StoreId: store.Id, State: ptpb.ProductType_active.String()})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: &StoreCustomer{
		Store:        store,
		ProductTypes: productTypes,
	}})
}
