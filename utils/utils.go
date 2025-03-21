package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"google.golang.org/grpc/metadata"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.SetAliasTag("json")
}

type Map = map[string]interface{}

type MapString = map[string]string

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
