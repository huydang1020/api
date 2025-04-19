package main

import (
	"context"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	ppb "github.com/huyshop/header/permission"
)

func (r *Router) handleListRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.RoleRequest{}
	utils.BindQuery(req, ctx)
	if err := r.isCanBeAccess(c, ctx, "role", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	roles, err := r.permSer.ListRoles(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: roles})
}

func (r *Router) handleGetRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "role", "r"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	role, err := r.permSer.GetRole(c, &ppb.RoleRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: role})
}

func (r *Router) handleCreateRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.Role{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "role", "c"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.permSer.CreateRole(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ppb.Role{}
	ctx.ShouldBindJSON(req)
	if err := r.isCanBeAccess(c, ctx, "role", "u"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	req.Id = id
	_, err := r.permSer.UpdateRole(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteRole(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	if err := r.isCanBeAccess(c, ctx, "role", "d"); err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	_, err := r.permSer.DeleteRole(c, &ppb.Role{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r Router) isCanBeAccess(c context.Context, ctx *gin.Context, group, action string) error {
	rid, exist := ctx.Get("role_id")
	if !exist {
		return errors.New(utils.E_invalid_role)
	}
	roleid, ok := rid.(string)
	if !ok {
		return errors.New(utils.E_invalid_role)
	}
	enforcer := &ppb.PolicyRequest{
		RoleId: roleid, Group: group, Action: action,
	}
	log.Println("enforcer", enforcer)
	_, err := r.permSer.CheckAccess(c, enforcer)
	return err
}
