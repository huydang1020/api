package main

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	permpb "github.com/huyshop/header/permission"
	userpb "github.com/huyshop/header/user"
)

func (r *Router) handleSignInAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	resp, err := r.userSer.SignIn(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	role, err := r.permSer.GetRole(c, &permpb.RoleRequest{Id: resp.GetUser().GetRoleId()})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	pages, err := r.permSer.ListPages(c, &permpb.PageRequest{RoleId: role.GetId()})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	menu := SortPage(pages)

	// Gán lại quyền cho user
	role.Page = menu
	resp.User.Role = role
	resp.User.Pages = menu
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleGetListUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.UserRequest{}
	utils.BindQuery(req, ctx)
	users, err := r.userSer.ListUsers(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	rids := make([]string, 0)
	for _, user := range users.Users {
		rids = append(rids, user.RoleId)
	}
	roles, err := r.permSer.ListRoles(c, &permpb.RoleRequest{RoleIds: rids})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	// Gán lại quyền cho user
	for _, user := range users.Users {
		for _, role := range roles.Roles {
			if user.RoleId == role.Id {
				user.Role = role
				break
			}
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: users})
}

func (r *Router) handleGetUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: user})
}

func (r *Router) handleCreateUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	fullName := ctx.PostForm("full_name")
	if fullName == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New("full_name is required"))
		return
	}
	req.FullName = fullName

	email := ctx.PostForm("email")
	if email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_email))
		return
	}
	req.Email = email

	username := ctx.PostForm("username")
	if username == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_username))
		return
	}
	req.Username = username

	password := ctx.PostForm("password")
	if password == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_password))
		return
	}
	req.Password = password

	roleId := ctx.PostForm("role_id")
	if roleId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_role_id))
		return
	}
	req.RoleId = roleId
	req.Province = ctx.PostForm("province")
	req.District = ctx.PostForm("district")
	req.Ward = ctx.PostForm("ward")
	req.Address = ctx.PostForm("address")
	if birthday := ctx.PostForm("birthday"); birthday != "" {
		birth, err := strconv.Atoi(birthday)
		if err != nil {
			log.Println("birthday err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		req.Birthday = int64(birth)
	}

	// Xử lý file avatar (nếu có)
	form, err := ctx.MultipartForm()
	if err == nil && form.File["avatar"] != nil {
		avatar := []string{}
		files := form.File["avatar"]
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
			avatar = append(avatar, imageUrl)
		}
		if len(avatar) > 0 {
			req.Avatar = avatar[0]
		}
	}
	log.Println("req:", req)
	_, err = r.userSer.CreateUser(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.User{}
	req.Id = id
	req.FullName = ctx.PostForm("full_name")
	req.Email = ctx.PostForm("email")
	req.PhoneNumber = ctx.PostForm("phone_number")
	req.Username = ctx.PostForm("username")
	req.Password = ctx.PostForm("password")
	req.RoleId = ctx.PostForm("role_id")
	req.Province = ctx.PostForm("province")
	req.District = ctx.PostForm("district")
	req.Ward = ctx.PostForm("ward")
	req.State = ctx.PostForm("state")
	req.Address = ctx.PostForm("address")
	if birthday := ctx.PostForm("birthday"); birthday != "" {
		birth, err := strconv.Atoi(birthday)
		if err != nil {
			log.Println("birthday err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		req.Birthday = int64(birth)
	}

	// Xử lý file avatar (nếu có)
	form, err := ctx.MultipartForm()
	if err == nil && form.File["avatar"] != nil {
		avatar := []string{}
		files := form.File["avatar"]
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
			avatar = append(avatar, imageUrl)
		}
		if len(avatar) > 0 {
			req.Avatar = avatar[0]
		}
	}
	log.Println("req:", req)
	_, err = r.userSer.UpdateUser(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteUser(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	_, err := r.userSer.DeleteUser(c, &userpb.User{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}
