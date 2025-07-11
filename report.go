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
	"github.com/huyshop/header/voucher"
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
		resp.TotalUsers = int32(len(listUser.Users))
		listPartner, err := r.userSer.ListPartner(c, &userpb.PartnerRequest{})
		if err != nil {
			log.Println("err", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		resp.TotalPartners = int32(len(listPartner.Partners))
		listVoucher, err := r.voucherSer.ListVouchers(c, &voucher.VoucherRequest{})
		if err != nil {
			log.Println("err", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		resp.TotalVouchers = int32(len(listVoucher.Vouchers))
		listCodeUsed, err := r.voucherSer.ListUserVouchers(c, &voucher.UserVoucherRequest{
			State: voucher.UserVoucher_used.String(),
		})
		if err != nil {
			log.Println("err", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		resp.TotalCodeUsed = listCodeUsed.Total
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

func (r *Router) handleGetReportProduct(ctx *gin.Context) {
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
	if req.OrderBy == "" {
		req.OrderBy = "sold"
	}
	resp, err := r.productSer.ListProductType(c, &ptpb.ProductTypeRequest{Limit: 10, State: ptpb.ProductType_active.String(), PartnerId: req.PartnerId, OrderBy: req.OrderBy})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	dataProduct := make(map[string]int32, 0)
	for _, product := range resp.ProductTypes {
		if req.OrderBy == "views" {
			dataProduct[product.Name] = product.GetViews()
		} else {
			dataProduct[product.Name] = product.GetQuantitySold()
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: dataProduct})
}

func (r *Router) handleGetReportTopProducts(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &ptpb.ReportRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "home", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}

	// Set default time range to current month
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

	// Get all completed orders in the time range
	orderResp, err := r.productSer.ListOrder(c, &ptpb.OrderRequest{
		State: ptpb.Order_completed.String(),
	})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}

	// Count products sold in the time range
	productCount := make(map[string]int32)
	productInfo := make(map[string]*ptpb.Product)

	for _, order := range orderResp.Orders {
		// Check if order is within time range
		if order.TimeOrder >= req.StartDate && order.TimeOrder <= req.EndDate {
			// Count products in this order
			for _, item := range order.ProductOrdered {
				if item.Product != nil && item.Product.Id != "" {
					productCount[item.Product.Id] += item.Quantity

					// Get product info if not already cached
					if _, exists := productInfo[item.Product.Id]; !exists {
						product, err := r.productSer.GetProduct(c, &ptpb.ProductRequest{Id: item.Product.Id})
						if err == nil && product != nil {
							productInfo[item.Product.Id] = product
						}
					}
				}
			}
		}
	}

	// Create labels and values for ReportRevenue
	var labels []string
	var values []int32

	// Convert map to slice for sorting
	type ProductSale struct {
		ProductId string
		Quantity  int32
		Product   *ptpb.Product
	}

	var productSales []ProductSale

	for productId, quantity := range productCount {
		if productInfo[productId] != nil {
			product := productInfo[productId]
			productSale := ProductSale{
				ProductId: productId,
				Quantity:  quantity,
				Product:   product,
			}
			productSales = append(productSales, productSale)
		}
	}

	// Sort by quantity sold (descending) and take top 10
	sort.Slice(productSales, func(i, j int) bool {
		return productSales[i].Quantity > productSales[j].Quantity
	})

	// Limit to top 10 and create labels/values
	for i, productSale := range productSales {
		if i >= 10 {
			break
		}

		// Create label: product name
		label := productSale.Product.Name

		labels = append(labels, label)
		values = append(values, productSale.Quantity)
	}

	// Create ReportRevenue response
	reportRevenue := &ptpb.ReportRevenue{
		Labels: labels,
		Values: make([]int64, len(values)),
	}

	// Convert int32 to int64
	for i, v := range values {
		reportRevenue.Values[i] = int64(v)
	}

	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: reportRevenue})
}
