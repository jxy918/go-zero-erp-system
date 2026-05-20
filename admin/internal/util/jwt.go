package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构体
// 存储用户认证信息，用于token生成和解析
type Claims struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	RoleID   uint   `json:"role_id"`  // 角色ID
	IsAdmin  bool   `json:"is_admin"` // 是否管理员
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
// 参数:
//
//	userID - 用户ID
//	username - 用户名
//	roleID - 角色ID
//	isAdmin - 是否管理员
//	secret - 签名密钥（生产环境应从配置读取）
//	expire - token过期时间（秒）
//
// 返回:
//
//	token字符串和错误信息
func GenerateToken(userID uint, username string, roleID uint, isAdmin bool, secret string, expire int64) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                          // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                          // 生效时间
		},
	}

	// 使用HS256算法签名token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析并验证JWT token
// 参数:
//
//	tokenString - JWT token字符串
//	secret - 可选的签名密钥，不传则使用默认密钥
//
// 返回:
//
//	解析后的Claims结构体和错误信息
func ParseToken(tokenString string, secret ...string) (*Claims, error) {
	secretKey := ""
	if len(secret) > 0 {
		secretKey = secret[0]
	} else {
		secretKey = "your-secret-key" // 默认密钥（生产环境应替换）
	}

	// 解析token并验证签名
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token有效性并提取claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
