package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
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
	if strings.Contains(lang, "vi_VN") {
		if data, ok := mLangs[s.Message()]; ok {
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: data.En})
			return
		} else {
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: "An error occurred"})
			return
		}
	} else {
		if data, ok := mLangs[s.Message()]; ok {
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: data.Vi})
			return
		} else {
			ctx.JSON(statusCode, ErrMsg{Code: -1, Message: "Có lỗi xảy ra"})
			return
		}
	}
}

func HandleSuccess(mLangs map[string]LangCode, ctx *gin.Context, resp *Response) {
	statusCode := 200
	lang := ctx.GetHeader("Accept-Language")
	if strings.Contains(lang, "vi_VN") {
		if data, ok := mLangs[resp.Message]; ok {
			resp.Message = data.En
			ctx.JSON(statusCode, resp)
			return
		}
	} else {
		if data, ok := mLangs[resp.Message]; ok {
			resp.Message = data.Vi
			ctx.JSON(statusCode, resp)
			return
		}
	}
}

func GetUserIdByToken(ctx *gin.Context) (string, error) {
	uid, exist := ctx.Get("user_id")
	if !exist {
		return "", errors.New(E_invalid_user_id)
	}
	userid, ok := uid.(string)
	if !ok {
		return "", errors.New(E_internal_errors)
	}

	return userid, nil
}

func GetRoleByToken(ctx *gin.Context) (string, error) {
	rid, exist := ctx.Get("user_id")
	if !exist {
		return "", errors.New(E_invalid_user_id)
	}
	userid, ok := rid.(string)
	if !ok {
		return "", errors.New(E_internal_errors)
	}

	return userid, nil
}
