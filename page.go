package main

import (
	"errors"
	"log"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	ppb "github.com/huyshop/header/permission"
	upb "github.com/huyshop/header/user"
)

func (r *Router) handleListPageByPermission(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.PageRequest{}
	utils.BindQuery(req, ctx)
	uid, exist := ctx.Get("user_id")
	if exist {
		uidStr, ok := uid.(string)
		if !ok {
			utils.HandleError(LangMapping, ctx, errors.New(utils.E_invalid_user))
			return
		}
		user, err := r.userSer.GetUser(c, &upb.UserRequest{Id: uidStr})
		if err != nil {
			ctx.JSON(400, &Response{Code: -1, Message: utils.E_access_denied})
			return
		}
		req.RoleId = user.GetRoleId()
	}
	log.Println("req", req)
	pages, err := r.permSer.ListPages(c, req)
	if err != nil {
		utils.HandleError(LangMapping, ctx, err)
		return
	}
	// log.Println("pages", pages)
	menu := SortPage(pages)
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: &ppb.Pages{Pages: menu, Total: pages.Total}})
}

func (r *Router) handleListPage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.PageRequest{}
	utils.BindQuery(req, ctx)
	uid, exist := ctx.Get("user_id")
	if exist {
		uidStr, ok := uid.(string)
		if !ok {
			utils.HandleError(LangMapping, ctx, errors.New(utils.E_invalid_user))
			return
		}
		user, err := r.userSer.GetUser(c, &upb.UserRequest{Id: uidStr})
		if err != nil {
			utils.HandleError(LangMapping, ctx, errors.New(utils.E_access_denied))
			return
		}
		req.RoleId = user.GetRoleId()
	}
	log.Println("req", req)
	pages, err := r.permSer.ListPages(c, req)
	if err != nil {
		utils.HandleError(LangMapping, ctx, err)
		return
	}
	for _, p := range pages.Pages {
		p.Actions = nil
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: pages})
}

func SortPage(pages *ppb.Pages) []*ppb.Page {
	// Tạo map để tra cứu nhanh
	mapMenu := make(map[string]*ppb.Page)
	for _, p := range pages.Pages {
		mapMenu[p.Id] = p
	}

	// Xây dựng cây menu
	var menu []*ppb.Page
	for _, p := range pages.Pages {
		if p.ParentId == "" {
			menu = append(menu, p) // Trang chính
		} else if parent, ok := mapMenu[p.ParentId]; ok {
			parent.Children = append(parent.Children, p) // Thêm vào trang cha
		} else {
			// Trường hợp không tìm thấy cha, vẫn thêm vào menu để không mất dữ liệu
			menu = append(menu, p)
		}
	}
	// Sắp xếp menu chính và menu con
	sortPagesByOrder(menu)
	return menu
}

func sortPagesByOrder(pages []*ppb.Page) {
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Id < pages[j].Id // Sắp xếp theo ID tăng dần (hoặc theo Name, Order, ...)
	})
	for _, p := range pages {
		if len(p.Children) > 0 {
			sortPagesByOrder(p.Children) // Đệ quy để sắp xếp các menu con
		}
	}
}

func (r *Router) handleGetPage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	page, err := r.permSer.GetPage(c, &ppb.PageRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMapping, ctx, err)
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: page})
}

func (r *Router) handleCreatePage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.Page{}
	ctx.ShouldBindJSON(req)
	_, err := r.permSer.CreatePage(c, req)
	if err != nil {
		utils.HandleError(LangMapping, ctx, err)
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdatePage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &ppb.Page{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	_, err := r.permSer.UpdatePage(c, req)
	if err != nil {
		utils.HandleError(LangMapping, ctx, err)
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeletePage(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	_, err := r.permSer.DeletePage(c, &ppb.Page{Id: id})
	if err != nil {
		utils.HandleError(LangMapping, ctx, err)
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}
