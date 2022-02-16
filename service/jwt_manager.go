package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//
// JWTManager
//  @Description:
//
type JWTManager struct {
	//用于签署和验证访问令牌的密钥
	secretKey string
	//令牌有效期
	tokenDuration time.Duration
}

//
// UserClaims
//  @Description:
//
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

//
// NewJWTManager
//  @Description:
//  @param secretKey
//  @param tokenDuration
//  @return *JWTManager
//
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

//
// Generate
//  @Description:
//  @receiver manager
//  @param user
//  @return string
//  @return error
//
func (manager *JWTManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role,
	}

	//生成令牌对象,这里使用基于HMAC的签名方法HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//用密钥签名生成令牌
	return token.SignedString([]byte(manager.secretKey))
}

//
// Verify
//  @Description:
//  @receiver manager
//  @param accessToken
//  @return *UserClaims
//  @return error
//
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{}, //空用户声称
		func(token *jwt.Token) (interface{}, error) { //自定义key功能
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil //匹配则返回用于对令牌进行签名的密钥
		},
	)

	//错误不为nil, 返回无效令牌
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	//否则从token中获取claims,并转换为user声明对象
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
