package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zxc7563598/oneadmin/internal/config"
)

var (
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
)

type Claims struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"` // access / refresh
	Role string `json:"role"` // admin / user
	jwt.RegisteredClaims
}

// Init 初始化 JWT
func Init(c config.JWTConfig) {
	secret = []byte(c.Secret)
	accessTokenTTL = time.Duration(c.AccessTTL) * time.Second
	refreshTokenTTL = time.Duration(c.RefreshTTL) * time.Second
}

// GenerateAccessToken 生成 AccessToken
func GenerateAccessToken(id uint64, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		ID:   id,
		Type: "access",
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// GenerateRefreshToken 生成 RefreshToken
func GenerateRefreshToken(id uint64, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		ID:   id,
		Type: "refresh",
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("无法验证的签名方式")
			}
			return secret, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 token")
	}
	return claims, nil
}

func AccessTTL() time.Duration {
	return accessTokenTTL
}

func RefreshTTL() time.Duration {
	return refreshTokenTTL
}
