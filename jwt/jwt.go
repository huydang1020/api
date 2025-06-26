package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/huyshop/api/utils"
)

type JWTClaim struct {
	UserId      string `json:"user_id"`
	RoleId      string `json:"role_id"`
	PartnerId   string `json:"partner_id"`
	PartnerType string `json:"partner_type"`
	jwt.StandardClaims
}

func GenerateAccessToken(newClaims *JWTClaim, expireTime time.Duration, secretKey string) (string, error) {
	expirationTime := time.Now().Add(expireTime * time.Minute)
	log.Println("newClaims:", newClaims)
	claims := &JWTClaim{
		UserId:      newClaims.UserId,
		RoleId:      newClaims.RoleId,
		PartnerId:   newClaims.PartnerId,
		PartnerType: newClaims.PartnerType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("generate access token err:", err)
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(newClaims *JWTClaim, expireTime time.Duration, secretKey string) (string, error) {
	expirationTime := time.Now().Add(expireTime * time.Minute)
	log.Println("newClaims:", newClaims)
	claims := &JWTClaim{
		UserId:      newClaims.UserId,
		RoleId:      newClaims.RoleId,
		PartnerId:   newClaims.PartnerId,
		PartnerType: newClaims.PartnerType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("generate refresh token err:", err)
		return "", err
	}
	return tokenString, nil
}

func ValidateRefreshToken(refreshToken, secretKey string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New(utils.E_invalid_sigature + err.Error())
		}
		return nil, errors.New(utils.E_invalid_token)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New(utils.E_invalid_token)
	}

	return claims, nil
}
