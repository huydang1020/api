package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	permpb "github.com/huyshop/header/permission"
	userpb "github.com/huyshop/header/user"
	"go.elastic.co/apm/module/apmgin"
	"google.golang.org/grpc"
)

type Router struct {
	route   *gin.Engine
	permSer permpb.PermissionServiceClient
	userSer userpb.UserServiceClient
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (r *Router) dialPerm(target string) error {
	client, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		return err
	}
	r.permSer = permpb.NewPermissionServiceClient(client)
	return nil
}

func (r *Router) dialUser(target string) error {
	client, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		return err
	}
	r.userSer = userpb.NewUserServiceClient(client)
	return nil
}

func NewRouter(cf *Configs) error {
	r := &Router{}
	if err := r.dialPerm(cf.PermGrpcServer); err != nil {
		log.Print(err)
	}
	if err := r.dialUser(cf.UserGrpcServer); err != nil {
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
	r.router()
	r.route.Run(fmt.Sprintf(":%v", cf.Port))
	return nil
}
