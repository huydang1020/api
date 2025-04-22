package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	userpb "github.com/huyshop/header/user"
)

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
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
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
	req.Type = claims.PartnerType
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
