package main

import (
	"errors"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
	userpb "github.com/huyshop/header/user"
)

func (r *Router) handleGetReportOverview(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &ptpb.ReportRequest{}
	utils.BindQuery(req, ctx)

	if claims.PartnerType != "admin" {
		req.PartnerId = claims.PartnerId
	}
	log.Println("req", req)
	resp, err := r.productSer.GetReportOverview(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	log.Println("resp", resp.OrderStatus)
	if claims.PartnerType == "admin" {
		listUser, err := r.userSer.ListUsers(c, &userpb.UserRequest{})
		if err != nil {
			log.Println("err", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		resp.NewUsers = int32(len(listUser.Users))
		listPartner, err := r.userSer.ListPartner(c, &userpb.PartnerRequest{})
		if err != nil {
			log.Println("err", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		resp.NewPartners = int32(len(listPartner.Partners))
	}
	listStore, err := r.userSer.ListStore(c, &userpb.StoreRequest{PartnerId: req.PartnerId})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp.TotalStores = listStore.Total
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetReportRevenue(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &ptpb.ReportRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "home", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType != "admin" {
		req.PartnerId = claims.PartnerId
	}
	log.Println("req", req)
	resp, err := r.productSer.GetReportRevenue(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetReportRevenueByStore(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &ptpb.ReportRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "home", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	now := time.Now()
	if req.StartDate == 0 {
		req.StartDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix()
	}
	if req.EndDate == 0 {
		req.EndDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1).Unix() + 86399
	}
	if claims.PartnerType != "admin" {
		req.PartnerId = claims.PartnerId
	}
	log.Println("req", req)
	resp, err := r.productSer.GetReportStoreRevenue(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	var storeIds []string
	for _, s := range resp.StoreRevenues {
		if s.Store == nil && s.StoreId != "" {
			storeIds = append(storeIds, s.StoreId)
		}
	}
	if len(storeIds) > 0 {
		storesResp, err := r.userSer.ListStore(c, &userpb.StoreRequest{Ids: storeIds})
		if err == nil && storesResp != nil {
			storeMap := map[string]*userpb.Store{}
			for _, st := range storesResp.Stores {
				storeMap[st.Id] = st
			}
			for _, s := range resp.StoreRevenues {
				if s.Store == nil {
					s.Store = storeMap[s.StoreId]
				}
			}
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetReportUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &userpb.ReportRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "home", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.RoleId != config.AdminRole {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_access_denied))
		return
	}
	log.Println("req", req)
	resp, err := r.userSer.GetReportUser(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, user := range resp.Users {
		orderResp, err := r.productSer.ListOrder(c, &ptpb.OrderRequest{
			UserId: user.UserId,
			State:  ptpb.Order_completed.String(),
		})
		if err == nil && orderResp != nil {
			totalOrders := int32(len(orderResp.Orders))
			totalAmount := int64(0)
			for _, ord := range orderResp.Orders {
				totalAmount += int64(ord.TotalMoney)
			}
			user.TotalOrders = totalOrders
			user.TotalSpent = totalAmount
		}
	}
	orderBy := req.OrderBy
	switch orderBy {
	case "total_spent":
		sort.Slice(resp.Users, func(i, j int) bool {
			return resp.Users[i].TotalSpent > resp.Users[j].TotalSpent
		})
	default:
		sort.Slice(resp.Users, func(i, j int) bool {
			if resp.Users[i].TotalOrders == resp.Users[j].TotalOrders {
				return resp.Users[i].TotalSpent > resp.Users[j].TotalSpent
			}
			return resp.Users[i].TotalOrders > resp.Users[j].TotalOrders
		})
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}
