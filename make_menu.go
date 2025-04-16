package main

import (
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/huyshop/api/utils"
	ppb "github.com/huyshop/header/permission"
)

type MenuData struct {
	Path     string      `json:"path"`
	Handle   *Handle     `json:"handle"`
	Children []*MenuData `json:"children,omitempty"`
}

type Handle struct {
	Icon         string   `json:"icon"`
	IconType     string   `json:"iconType"`
	Title        string   `json:"title"`
	Order        int      `json:"order"`
	IframeLink   string   `json:"iframeLink,omitempty"`
	ExternalLink string   `json:"externalLink,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
	KeepAlive    bool     `json:"keepAlive,omitempty"`
	HideInMenu   bool     `json:"hideInMenu"`
}

func BuildMenuTree(pages []*ppb.Page) []*MenuData {
	idToMenu := make(map[string]*MenuData)
	var roots []*MenuData

	// Bước 1: Tạo map để ánh xạ từ ID sang MenuData
	for _, page := range pages {
		// Lấy danh sách roles và permissions
		var roles []string
		var perms []string
		for _, ra := range page.RoleActions {
			roles = AppendIfMissing(roles, ra.Role.Name)
			for _, act := range ra.Actions {
				perm := fmt.Sprintf("permission:%s", ConvertAction(act))
				perms = AppendIfMissing(perms, perm)
			}
		}

		menu := &MenuData{
			Path: "/" + page.Route,
			Handle: &Handle{
				Title:       page.Name,
				Icon:        page.Icon.IconName,
				IconType:    page.Icon.IconType,
				Order:       int(page.Order),
				Roles:       roles,
				Permissions: perms,
				// KeepAlive:   page.KeepAlive,
				// HideInMenu:  page.HideMenu,
			},
			Children: []*MenuData{},
		}

		idToMenu[page.Id] = menu
	}

	// Bước 2: Gắn các node con vào cha
	for _, page := range pages {
		menu := idToMenu[page.Id]
		if page.ParentId != "" {
			parentMenu, ok := idToMenu[page.ParentId]
			if ok {
				parentMenu.Children = append(parentMenu.Children, menu)
			}
		} else {
			roots = append(roots, menu)
		}
	}

	// Bước 3: Sắp xếp các menu theo thứ tự
	sortMenuDataByOrder(roots)
	return roots
}

func AppendIfMissing(slice []string, value string) []string {
	for _, v := range slice {
		if v == value {
			return slice
		}
	}
	return append(slice, value)
}

func ConvertAction(code string) string {
	switch code {
	case "c":
		return "button:add"
	case "r":
		return "button:get"
	case "u":
		return "button:update"
	case "d":
		return "button:delete"
	default:
		return "unknown"
	}
}

func sortMenuDataByOrder(menu []*MenuData) {
	sort.Slice(menu, func(i, j int) bool {
		return menu[i].Handle.Order > menu[j].Handle.Order
	})
	for _, m := range menu {
		if len(m.Children) > 0 {
			sortMenuDataByOrder(m.Children)
		}
	}
}

func (r *Router) ListMenu(ctx *gin.Context) {
	c, cancel := utils.MakeContext(MAXTIMEREQ, nil)
	defer cancel()
	req := &ppb.PageRole{}
	utils.BindQuery(req, ctx)
	pages, err := r.permSer.ListPages(c, &ppb.PageRequest{})
	if err != nil {
		utils.HandleError(LangMappingErr, ctx, err)
		return
	}
	menuTree := BuildMenuTree(pages.Pages)
	utils.HandleSuccess(LangMappingSuccess, ctx, &utils.Response{Code: 0, Message: "success", Data: menuTree})
}
