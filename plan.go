package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	userpb "github.com/huyshop/header/user"
)

func (r *Router) handleListPlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.PlansRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "plan", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp, err := r.userSer.ListPlans(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetPlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "plan", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp, err := r.userSer.GetPlan(c, &userpb.Plan{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleListPlan(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.PlansRequest{}
	utils.BindQuery(req, ctx)
	req.State = userpb.Plan_active.String()
	resp, err := r.userSer.ListPlans(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetPlan(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	resp, err := r.userSer.GetPlan(c, &userpb.Plan{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleCreatePlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.Plan{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "plan", "c"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	if req.Name == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_name_cannot_empty))
		return
	}
	log.Println("req:", req)
	_, err := r.userSer.CreatePlan(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdatePlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.Plan{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "plan", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	log.Println("req:", req)
	_, err := r.userSer.UpdatePlan(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeletePlanAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "plan", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.userSer.DeletePlan(c, &userpb.Plan{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}
