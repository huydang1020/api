package main

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	permpb "github.com/huyshop/header/permission"
	userpb "github.com/huyshop/header/user"
)

type Response struct {
	Code    int         `json:"code"` // 0: success, -1: error
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Router) handleSignInAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	resp, err := r.userSer.SignIn(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	role, err := r.permSer.GetRole(c, &permpb.RoleRequest{Id: resp.GetUser().GetRoleId()})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	pages, err := r.permSer.ListPages(c, &permpb.PageRequest{RoleId: role.GetId()})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	// Tạo map để tra cứu nhanh
	mapMenu := make(map[int32]*permpb.Page)
	for _, p := range pages.Pages {
		mapMenu[p.Id] = p
	}

	// Xây dựng cây menu
	var menu []*permpb.Page
	for _, p := range pages.Pages {
		if p.ParentId == 0 {
			menu = append(menu, p) // Trang chính
		} else if parent, ok := mapMenu[p.ParentId]; ok {
			parent.Children = append(parent.Children, p) // Thêm vào trang cha
		}
	}
	// Sắp xếp menu chính và menu con
	sortPages(menu)

	// Gán lại quyền cho user
	role.Permission = menu
	resp.User.Role = role
	resp.User.Permissions = menu
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: resp})
}

func sortPages(pages []*permpb.Page) {
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Id < pages[j].Id // Sắp xếp theo ID tăng dần (hoặc theo Name, Order, ...)
	})
	for _, p := range pages {
		if len(p.Children) > 0 {
			sortPages(p.Children) // Đệ quy để sắp xếp các menu con
		}
	}
}

func (r *Router) handleGetListUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.UserRequest{}
	utils.BindQuery(req, ctx)
	users, err := r.userSer.ListUsers(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: users})
}

func (r *Router) handleGetUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: id})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success", Data: user})
}

func (r *Router) handleCreateUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	_, err := r.userSer.CreateUser(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	req.Id = id
	_, err := r.userSer.UpdateUser(c, req)
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	_, err := r.userSer.DeleteUser(c, &userpb.User{Id: id})
	if err != nil {
		ctx.JSON(500, &Response{Code: -1, Message: err.Error()})
		return
	}
	ctx.JSON(200, &Response{Code: 0, Message: "success"})
}
