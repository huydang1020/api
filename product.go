package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	ptpb "github.com/huyshop/header/product"
	userpb "github.com/huyshop/header/user"
)

func (r *Router) handleListProduct(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ProductRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "product_type", "r"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	products, err := r.productSer.ListProduct(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: products})
}

func (r *Router) handleGetProduct(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "product_type", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	product, err := r.productSer.GetProduct(c, &ptpb.ProductRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: product})
}

func (r *Router) handleCreateProductType(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ProductType{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "product_type", "c"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType == userpb.Partner_admin.String() {
		req.State = ptpb.ProductType_active.String()
	} else {
		req.State = ptpb.ProductType_pending.String()
	}
	req.PartnerId = claims.PartnerId
	_, err := r.productSer.CreateProductType(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateProductType(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.ProductType{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	if err := r.isCanBeAccess(c, ctx, "product_type", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType != userpb.Partner_admin.String() {
		req.State = ptpb.ProductType_active.String()
	} else {
		req.State = ptpb.ProductType_pending.String()
	}
	_, err := r.productSer.UpdateProductType(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteProductType(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "product_type", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.productSer.DeleteProductType(c, &ptpb.ProductType{Id: id})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateStateProductType(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.ProductType{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	if err := r.isCanBeAccess(c, ctx, "product_type", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType != userpb.Partner_admin.String() && req.State == ptpb.ProductType_active.String() {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_access_denied))
		return
	}
	log.Println("req ", req)
	_, err := r.productSer.UpdateStateProductType(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListProductType(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ProductTypeRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "product_type", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if claims.PartnerType != userpb.Partner_admin.String() {
		req.PartnerId = claims.PartnerId
	}
	productTypes, err := r.productSer.ListProductType(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, pty := range productTypes.GetProductTypes() {
		if pty.PartnerId != "" {
			partner, err := r.userSer.GetPartner(c, &userpb.PartnerRequest{Id: pty.PartnerId})
			if err != nil {
				log.Println("err", err)
				continue
			}
			pty.Partner = partner
		}
		if pty.StoreId != "" {
			store, err := r.userSer.GetStore(c, &userpb.StoreRequest{Id: pty.StoreId})
			if err != nil {
				log.Println("err", err)
				continue
			}
			pty.Store = store
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: productTypes})
}

func (r *Router) handleListProductTypeCustomer(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.ProductTypeRequest{}
	utils.BindQuery(req, ctx)
	productTypes, err := r.productSer.ListProductType(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	for _, pty := range productTypes.GetProductTypes() {
		if pty.PartnerId != "" {
			partner, err := r.userSer.GetPartner(c, &userpb.PartnerRequest{Id: pty.PartnerId})
			if err != nil {
				log.Println("err", err)
				continue
			}
			pty.Partner = partner
		}
		if pty.StoreId != "" {
			store, err := r.userSer.GetStore(c, &userpb.StoreRequest{Id: pty.StoreId})
			if err != nil {
				log.Println("err", err)
				continue
			}
			pty.Store = store
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: productTypes})
}

func (r *Router) handleGetProductType(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "product_type", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	productType, err := r.productSer.GetProductType(c, &ptpb.ProductTypeRequest{Id: id})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: productType})
}

func (r *Router) handleGetProductTypeCustomer(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	pty, err := r.productSer.GetProductType(c, &ptpb.ProductTypeRequest{Id: id})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if pty.StoreId != "" {
		store, err := r.userSer.GetStore(c, &userpb.StoreRequest{Id: pty.StoreId})
		if err != nil {
			log.Println("err", err)
		}
		pty.Store = store
	}
	countpty, err := r.productSer.CountProductType(ctx, &ptpb.ProductTypeRequest{StoreId: pty.StoreId})
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if countpty.GetCount() > 0 {
		pty.Store.QuantityProduct = int32(countpty.GetCount())
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: pty})
}

func (r *Router) handleListCategory(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.CategoryRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "category", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	productCategories, err := r.productSer.ListCategory(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: productCategories})
}

func (r *Router) handleGetCategory(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "category", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	Category, err := r.productSer.GetCategory(c, &ptpb.CategoryRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: Category})
}

func (r *Router) handleCreateCategory(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Category{}
	if err := r.isCanBeAccess(c, ctx, "category", "c"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Name = ctx.PostForm("name")
	// Xử lý file logo (nếu có)
	form, err := ctx.MultipartForm()
	if err == nil && form.File["logo"] != nil {
		logo := []string{}
		files := form.File["logo"]
		for _, file := range files {
			imageName := file.Filename
			image, err := file.Open()
			if err != nil {
				log.Println("file open err:", err)
				continue
			}
			defer image.Close()

			imageUrl, err := UploadImageToCloudinary(c, image, imageName)
			if err != nil {
				log.Println("upload img err:", err)
				continue
			}
			logo = append(logo, imageUrl)
		}
		if len(logo) > 0 {
			req.Logo = logo[0]
		}
	}
	_, err = r.productSer.CreateCategory(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateCategory(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.Category{Id: id}
	if err := r.isCanBeAccess(c, ctx, "category", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Name = ctx.PostForm("name")
	// Xử lý file logo (nếu có)
	form, err := ctx.MultipartForm()
	if err == nil && form.File["logo"] != nil {
		logo := []string{}
		files := form.File["logo"]
		for _, file := range files {
			imageName := file.Filename
			image, err := file.Open()
			if err != nil {
				log.Println("file open err:", err)
				continue
			}
			defer image.Close()

			imageUrl, err := UploadImageToCloudinary(c, image, imageName)
			if err != nil {
				log.Println("upload img err:", err)
				continue
			}
			logo = append(logo, imageUrl)
		}
		if len(logo) > 0 {
			req.Logo = logo[0]
		}
	}
	req.Id = id
	log.Println("req:", req)
	_, err = r.productSer.UpdateCategory(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteCategory(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "category", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.productSer.DeleteCategory(c, &ptpb.Category{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleCreateBanner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.Banner{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "banner", "c"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.productSer.CreateBanner(c, req)
	if err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateBanner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ptpb.Banner{Id: id}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "banner", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.productSer.UpdateBanner(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteBanner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "banner", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.productSer.DeleteBanner(c, &ptpb.Banner{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListBanner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ptpb.BannerRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "banner", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	banners, err := r.productSer.ListBanner(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: banners})
}

func (r *Router) handleGetBanner(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "banner", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	banner, err := r.productSer.GetBanner(c, &ptpb.BannerRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: banner})
}

type Home struct {
	Categories []*ptpb.Category `json:"categories"`
	Banners    []*ptpb.Banner   `json:"banners"`
}

type HomeRequest struct {
	CategoryId string `json:"category_id"`
}

func (r *Router) handleHome(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &HomeRequest{}
	utils.BindQuery(req, ctx)
	log.Println("req", req)
	cates, err := r.productSer.ListCategory(c, &ptpb.CategoryRequest{State: ptpb.Category_active.String()})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	banners, err := r.productSer.ListBanner(c, &ptpb.BannerRequest{State: ptpb.Banner_active.String()})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	home := &Home{
		Categories: cates.GetCategories(),
		Banners:    banners.GetBanners(),
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: home})
}
