package main

import (
	"errors"
	"log"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
	permpb "github.com/huyshop/header/permission"
	"github.com/huyshop/header/product"
	userpb "github.com/huyshop/header/user"
	"google.golang.org/grpc/status"
)

const (
	ROLE_CUSTOMER = "roled0di17m9ipf12jq5ndlg"
)

func (r *Router) handleSignInAdmin(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	resp, err := r.userSer.SignIn(c, req)
	if err != nil || resp == nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	role, err := r.permSer.GetRole(c, &permpb.RoleRequest{Id: resp.User.RoleId})
	if err != nil {
		log.Println("GetRole err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	resp.User.Role = role
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleSignOutAdmin(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	_, err := r.userSer.SignOut(c, &userpb.User{Id: claims.UserId})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleGetMe(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: claims.UserId})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	user.Password = ""
	// lấy số lượng sản phẩm trong giỏ hàng
	cart, err := r.productSer.ListCart(c, &product.Cart{UserId: claims.UserId})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	var total int
	for _, item := range cart.GetItem() {
		total += int(item.Quantity)
	}
	user.CartQuantity = int32(total)
	// lấy đối tác
	if user.PartnerId != "" {
		part, err := r.userSer.GetPartner(c, &userpb.PartnerRequest{Id: user.PartnerId})
		if err != nil {
			log.Println("GetOart err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		user.Partner = part
	}
	// lấy sản phẩm yêu thích
	key := "favorites:" + user.Id
	productTypeIds, err := r.cache.SMembers(c, key).Result()
	if err != nil {
		log.Println("Redis SMembers error:", err)
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_invalid_user_id))
		return
	}
	user.FavoriteQuantity = int32(len(productTypeIds))
	// lấy số lượng đơn hàng đã đặt
	order, err := r.productSer.ListOrder(c, &product.OrderRequest{UserId: user.Id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	user.TotalOrders = int32(len(order.Orders))
	// lấy số lượng đơn hàng đã đặt
	for _, ord := range order.Orders {
		if ord.State == product.Order_completed.String() {
			user.TotalAmountSpent += int64(ord.TotalMoney)
		}
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: user})
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
	ctx.ShouldBindJSON(req)
	if req.FullName == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_full_name_cannot_empty))
		return
	}
	if req.Email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	if req.Username == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_username_cannot_empty))
		return
	}
	if req.Password == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_password_cannot_empty))
		return
	}
	if req.RoleId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_role_cannot_empty))
		return
	}
	if req.PhoneNumber == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_phone_number_cannot_empty))
		return
	}
	check, _ := r.userSer.IsExistUser(c, &userpb.User{Email: req.Email, Username: req.Username, PhoneNumber: req.PhoneNumber})
	if check.Exist {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_user_existed))
		return
	}
	log.Println("req:", req)
	_, err := r.userSer.CreateUser(c, req)
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
	ctx.ShouldBindJSON(req)
	req.Id = id
	log.Println("req:", req)
	user, err := r.userSer.GetUser(c, &userpb.UserRequest{Id: id})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	newUsername := req.Username
	if newUsername == "" {
		newUsername = user.Username
	}
	newEmail := req.Email
	if newEmail == "" {
		newEmail = user.Email
	}
	newPhone := req.PhoneNumber
	if newPhone == "" {
		newPhone = user.PhoneNumber
	}
	if !reflect.DeepEqual(
		&userpb.User{Username: user.Username, Email: user.Email, PhoneNumber: user.PhoneNumber},
		&userpb.User{Username: newUsername, Email: newEmail, PhoneNumber: newPhone},
	) {
		check, _ := r.userSer.IsExistUser(c, &userpb.User{
			Email: newEmail, Username: newUsername, PhoneNumber: newPhone,
		})
		if check.Exist {
			utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_user_existed))
			return
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

func (r *Router) handleRegisterSellerAccount(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	fullName := ctx.PostForm("full_name")
	if fullName == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_full_name_cannot_empty))
		return
	}
	req.FullName = fullName

	email := ctx.PostForm("email")
	if email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	req.Email = email

	username := ctx.PostForm("username")
	if username == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_username_cannot_empty))
		return
	}
	req.Username = username

	password := ctx.PostForm("password")
	if password == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_password_cannot_empty))
		return
	}
	req.Password = password

	roleId := ctx.PostForm("role_id")
	if roleId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_role_cannot_empty))
		return
	}
	req.PhoneNumber = ctx.PostForm("phone_number")
	if req.PhoneNumber == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_phone_number_cannot_empty))
		return
	}
	req.RoleId = roleId
	// req.Province = ctx.PostForm("province")
	// req.District = ctx.PostForm("district")
	// req.Ward = ctx.PostForm("ward")
	// req.Address = ctx.PostForm("address")
	if birthday := ctx.PostForm("birthday"); birthday != "" {
		birth, err := strconv.Atoi(birthday)
		if err != nil {
			log.Println("birthday err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		req.Birthday = int64(birth)
	}
	check, _ := r.userSer.IsExistUser(c, &userpb.User{Email: req.Email, Username: req.Username, PhoneNumber: req.PhoneNumber})
	if check.Exist {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_user_existed))
		return
	}
	log.Println("req:", req)
	// resp, err := r.userSer.RegisterSellerAccount(c, req)
	// if err != nil {
	// 	utils.HandleError(LangMappingErr, ctx, err)
	// 	return
	// }

	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

// Handle customer
func (r *Router) handleSignInCustomer(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	resp, err := r.userSer.SignInCustomer(c, req)
	if err != nil {
		s := status.Convert(err)
		if s.Message() != utils.E_account_not_activated {
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 1, Message: utils.E_account_not_activated})
		return
	}
	cart, err := r.productSer.ListCart(c, &product.Cart{UserId: resp.User.Id})
	if err != nil {
		log.Println("ListCart err:", err)
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	var total int
	for _, item := range cart.GetItem() {
		total += int(item.Quantity)
	}
	resp.User.CartQuantity = int32(total)
	if resp.User.PartnerId != "" {
		part, err := r.userSer.GetPartner(c, &userpb.PartnerRequest{Id: resp.User.PartnerId})
		if err != nil {
			log.Println("GetOart err:", err)
			utils.HandleError(LangMappingErr, ctx, err)
			return
		}
		resp.User.Partner = part
	}
	key := "favorites:" + resp.User.Id
	productTypeIds, err := r.cache.SMembers(c, key).Result()
	if err != nil {
		log.Println("Redis SMembers error:", err)
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_invalid_user_id))
		return
	}
	resp.User.FavoriteQuantity = int32(len(productTypeIds))
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: resp})
}

func (r *Router) handleSignOutCustomer(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	_, err := r.userSer.SignOut(c, &userpb.User{Id: claims.UserId})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleSignUpCustomer(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	if req.PhoneNumber == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_phone_number_cannot_empty))
		return
	}
	if req.Email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	if req.Password == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_password_cannot_empty))
		return
	}
	req.RoleId = ROLE_CUSTOMER
	_, err := r.userSer.CreateCustomer(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleVerifyOtp(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	if req.Email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	if req.VerifyOtp == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_not_found_code))
		return
	}
	_, err := r.userSer.VerifyEmail(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleSendOtp(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	if req.Email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	ttl, err := r.userSer.SendVerifyOtp(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: ttl})
}

func (r *Router) handleSendResetPasswordOtp(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	if req.Email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	ttl, err := r.userSer.SendResetPasswordOtp(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: ttl})
}

func (r *Router) handleUpdateProfile(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	req.Id = claims.UserId
	_, err := r.userSer.UpdateUser(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleResetPassword(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.User{}
	ctx.ShouldBindJSON(req)
	if req.Email == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_email_cannot_empty))
		return
	}
	_, err := r.userSer.ResetPassword(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleCreatePointTransaction(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	defer cancel()
	req := &userpb.PointExchange{}
	ctx.ShouldBindJSON(req)
	if claims.RoleId != config.AdminRole {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_invalid_role))
		return
	}
	if req.ReceiverId == "" {
		utils.HandleError(LangMappingErr, ctx, errors.New(utils.E_receiver_id_cannot_empty))
		return
	}
	_, err := r.userSer.CreatePointExchange(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}

	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleListUserAddress(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.UserAddressRequest{}
	utils.BindQuery(req, ctx)
	req.UserId = claims.UserId
	userAddress, err := r.userSer.ListUserAddress(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: userAddress})
}

func (r *Router) handleCreateUserAddress(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.UserAddress{}
	ctx.ShouldBindJSON(req)
	req.Id = ctx.Param("id")
	req.UserId = claims.UserId
	_, err := r.userSer.CreateUserAddress(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleUpdateUserAddress(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &userpb.UserAddress{}
	ctx.ShouldBindJSON(req)
	req.Id = ctx.Param("id")
	req.UserId = claims.UserId
	log.Println("req:", req)
	_, err := r.userSer.UpdateUserAddress(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}

func (r *Router) handleDeleteUserAddress(ctx *gin.Context) {
	claims, _ := ctx.MustGet("claims").(*jwt.JWTClaim)
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	id := ctx.Param("id")
	req := &userpb.UserAddress{Id: id, UserId: claims.UserId}
	_, err := r.userSer.DeleteUserAddress(c, req)
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success"})
}
