package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	jt "github.com/huyshop/api/jwt"
	"github.com/huyshop/api/utils"
)

func authMiddleware(r *Router) gin.HandlerFunc {
	var jwtKey = []byte(config.JwtSecretKey)
	return func(c *gin.Context) {
		tokenString := c.GetHeader("access-token")
		if tokenString == "" {
			log.Println("123")
			utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_unauthorized))
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &jt.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			// log.Println("token.Header", token.Header)
			return jwtKey, nil
		})
		claims, ok := token.Claims.(*jt.JWTClaim)
		if !ok {
			utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_invalid_token))
			c.Abort()
			return
		}
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					keyRedis := fmt.Sprintf("refresh_token_user_id_%s", claims.UserId)
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					refreshToken, err := r.cache.Get(ctx, keyRedis).Result()
					if err == redis.Nil || err != nil {
						c.JSON(200, utils.ErrMsg{Code: 1, Message: utils.E_invalid_token})
						c.Abort()
						return
					}
					if refreshToken == "" {
						c.JSON(200, utils.ErrMsg{Code: 1, Message: utils.E_invalid_token})
						c.Abort()
						return
					}
					validRefreshToken, err := jt.ValidateRefreshToken(refreshToken, config.JwtSecretKey)
					if err != nil {
						utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_invalid_token))
						c.Abort()
						return
					}
					exprAct, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_ACCESS_TOKEN"))
					newToken, err := jt.GenerateAccessToken(validRefreshToken, time.Duration(exprAct), config.JwtSecretKey)
					if err != nil {
						utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_invalid_token))
						c.Abort()
						return
					}
					c.Header("access-token", newToken)
					c.Set("claims", claims)
					c.Next()
					return
				}
			}
			if err == jwt.ErrSignatureInvalid {
				utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_unauthorized))
				c.Abort()
				return
			}
			utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_bad_request))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jt.JWTClaim); ok && token.Valid {
			c.Set("claims", claims)
			c.Next()
		} else {
			utils.HandleError(LangMappingErr, c, fmt.Errorf(utils.E_unauthorized))
			c.Abort()
			return
		}
	}
}
