package utils

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	ppb "github.com/huyshop/header/permission"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.SetAliasTag("json")
}

type Map = map[string]interface{}

type MapString = map[string]string

type LangCode struct {
	Vi string `json:"vi"`
	En string `json:"en"`
}

type ErrMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Code    int         `json:"code"` // 0: success, -1: error
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func MakeContext(sec int, claims interface{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sec)*time.Second)
	if claims != nil {
		bin, err := json.Marshal(claims)
		if err != nil {
			log.Print(err)
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "ctx", string(bin))
		return ctx, cancel
	}
	return ctx, cancel
}

func BindQuery(in interface{}, ctx *gin.Context) error {
	err := decoder.Decode(in, ctx.Request.URL.Query())
	return err
}

func Include(slice []string, in string) bool {
	for _, item := range slice {
		if item == in {
			return true
		}
	}
	return false
}

func ConvertUnixToDateTime(format string, t int64) (string, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Println("load location err:", err)
		return "", err
	}
	formattedDate := time.Unix(t, 0).In(location).Format(format)
	return formattedDate, nil
}

func HandleError(mLangs map[string]LangCode, ctx *gin.Context, err error) {
	s := status.Convert(err)
	statusCode := 200
	lang := ctx.GetHeader("Accept-Language")
	if strings.Contains(lang, "vi-VN") {
		if data, ok := mLangs[s.Message()]; ok {
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: data.Vi})
			return
		} else {
			log.Println("error:", err)
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: "Có lỗi xảy ra"})
			return
		}
	} else {
		if data, ok := mLangs[s.Message()]; ok {
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: data.En})
			return
		} else {
			log.Println("error:", err)
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: "An error occurred"})
			return
		}
	}
}

func HandleSuccess(mLangs map[string]LangCode, ctx *gin.Context, resp *Response) {
	statusCode := 200
	lang := ctx.GetHeader("Accept-Language")
	if strings.Contains(lang, "vi-VN") {
		if data, ok := mLangs[resp.Message]; ok {
			resp.Message = data.Vi
			ctx.JSON(statusCode, resp)
			return
		}
	} else {
		if data, ok := mLangs[resp.Message]; ok {
			resp.Message = data.En
			ctx.JSON(statusCode, resp)
			return
		}
	}
}

func BuildMenuTree(pages []*ppb.Page) []*ppb.MenuData {
	idToMenu := make(map[string]*ppb.MenuData)
	var roots []*ppb.MenuData
	// Bước 1: Tạo map idToMenu và danh sách root
	for _, page := range pages {
		menu := &ppb.MenuData{
			Path: page.Path,
			Handle: &ppb.Handle{
				Icon:       page.Handle.Icon,
				IconType:   page.Handle.IconType,
				Title:      page.Handle.Title,
				Order:      page.Handle.Order,
				KeepAlive:  page.Handle.KeepAlive,
				HideInMenu: page.Handle.HideInMenu,
			},
			Children: []*ppb.MenuData{},
		}

		idToMenu[page.Id] = menu
	}
	// Bước 2: Gắn các node con vào node cha
	for _, page := range pages {
		menu := idToMenu[page.Id]
		if page.ParentId != "" {
			if parentMenu, ok := idToMenu[page.ParentId]; ok {
				parentMenu.Children = append(parentMenu.Children, menu)
			}
		} else {
			roots = append(roots, menu)
		}
	}
	return roots
}

func AppendIfMissing(slice []string, value string) []string {
	if slices.Contains(slice, value) {
		return slice
	}
	return append(slice, value)
}
