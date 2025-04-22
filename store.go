package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	userpb "github.com/huyshop/header/user"
)

func (r *Router) handleListStore(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.StoreRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "store", "r"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	stores, err := r.userSer.ListStore(c, req)
	if err != nil {
		log.Println("err", err)
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
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	store, err := r.userSer.GetStore(c, &userpb.StoreRequest{Id: id})
	if err != nil {
		log.Println("err", err)
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
	if err := r.isCanBeAccess(c, ctx, "store", "c"); err != nil {
		log.Println("err", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Name = ctx.PostForm("name")
	req.Province = ctx.PostForm("province")
	req.District = ctx.PostForm("district")
	req.Ward = ctx.PostForm("ward")
	req.Address = ctx.PostForm("address")
	req.Lat = ctx.PostForm("lat")
	req.Lng = ctx.PostForm("lng")
	req.PhoneNumber = ctx.PostForm("phone_number")
	req.Description = ctx.PostForm("description")
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
	req.PartnerId = claims.PartnerId
	_, err = r.userSer.CreateStore(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateStore(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.Store{}
	if err := r.isCanBeAccess(c, ctx, "store", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Name = ctx.PostForm("name")
	req.Province = ctx.PostForm("province")
	req.District = ctx.PostForm("district")
	req.Ward = ctx.PostForm("ward")
	req.Address = ctx.PostForm("address")
	req.Lat = ctx.PostForm("lat")
	req.Lng = ctx.PostForm("lng")
	req.PhoneNumber = ctx.PostForm("phone_number")
	req.Description = ctx.PostForm("description")
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
	req.PartnerId = claims.PartnerId
	_, err = r.userSer.UpdateStore(c, req)
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
