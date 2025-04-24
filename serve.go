package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/huyshop/api/utils"
	permpb "github.com/huyshop/header/permission"
	propb "github.com/huyshop/header/product"
	userpb "github.com/huyshop/header/user"
	voupb "github.com/huyshop/header/voucher"
	"go.elastic.co/apm/module/apmgin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Router struct {
	route      *gin.Engine
	cache      *redis.Client
	permSer    permpb.PermissionServiceClient
	userSer    userpb.UserServiceClient
	voucherSer voupb.VoucherServiceClient
	productSer propb.ProductServiceClient
	
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	loadErrorCodes()
	loadSuccessCodes()
}

var LangMappingErr map[string]utils.LangCode
var LangMappingSuccess map[string]utils.LangCode

func loadErrorCodes() {
	bin, err := os.ReadFile("assets/errors.json")
	if err != nil {
		panic(err)
	}
	errorCodes := make(map[string]utils.LangCode)
	if err := json.Unmarshal(bin, &errorCodes); err != nil {
		panic(err)
	}
	LangMappingErr = errorCodes
}

func loadSuccessCodes() {
	bin, err := os.ReadFile("assets/success.json")
	if err != nil {
		panic(err)
	}
	successCodes := make(map[string]utils.LangCode)
	if err := json.Unmarshal(bin, &successCodes); err != nil {
		panic(err)
	}
	LangMappingSuccess = successCodes
}

func (r *Router) dialPerm(target string) error {
	client, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return err
	}
	r.permSer = permpb.NewPermissionServiceClient(client)
	return nil
}

func (r *Router) dialUser(target string) error {
	client, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return err
	}
	r.userSer = userpb.NewUserServiceClient(client)
	return nil
}

// func (r *Router) dialPartner(target string) error {
// 	client, err := grpc.Dial(target,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	r.partnerSer = partpb.NewPartnerServiceClient(client)
// 	return nil
// }

func (r *Router) dialVoucher(target string) error {
	client, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return err
	}
	r.voucherSer = voupb.NewVoucherServiceClient(client)
	return nil
}

func (r *Router) dialProduct(target string) error {
	client, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return err
	}
	r.productSer = propb.NewProductServiceClient(client)
	return nil
}

func NewRedisCache(addr, pw string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
	tick := time.NewTicker(10 * time.Minute)
	ctx := context.Background()
	go func(client *redis.Client) {
		for {
			select {
			case <-tick.C:
				if err := client.Ping(ctx).Err(); err != nil {
					panic(err)
				}
			}
		}
	}(client)
	return client
}

func NewRouter(cf *Configs) error {
	r := &Router{}
	if err := r.dialPerm(cf.PermGrpcServer); err != nil {
		log.Print(err)
	}
	if err := r.dialUser(cf.UserGrpcServer); err != nil {
		log.Print(err)
	}
	// if err := r.dialPartner(cf.PartnerGrpcServer); err != nil {
	// 	log.Print(err)
	// }
	if err := r.dialVoucher(cf.VoucherGrpcServer); err != nil {
		log.Print(err)
	}
	if err := r.dialProduct(cf.ProductGrpcServer); err != nil {
		log.Print(err)
	}
	r.route = gin.New()
	r.route.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/"},
		Formatter: func(param gin.LogFormatterParams) string {
			if param.Method == "OPTIONS" {
				return ""
			}
			var statusColor, methodColor, resetColor string
			if param.IsOutputColor() {
				statusColor = param.StatusCodeColor()
				methodColor = param.MethodColor()
				resetColor = param.ResetColor()
			}
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		},
	}))
	r.route.Use(gin.Recovery())
	r.route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "access-token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}), func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	})
	r.route.Use(apmgin.Middleware(r.route))
	redisDb, _ := strconv.Atoi(config.RedisDb)
	rd := NewRedisCache(config.RedisAddr, config.RedisPassword, redisDb)
	r.cache = rd
	r.router()
	r.route.Run(fmt.Sprintf(":%v", cf.Port))
	return nil
}
