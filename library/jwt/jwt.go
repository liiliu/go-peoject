package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"your_project/library/cache"
	"your_project/library/common"
	"your_project/library/config"
	"your_project/library/util"

	"github.com/gofiber/fiber/v2"
	jwtPkg "github.com/golang-jwt/jwt/v4"
)

// 错误定义
var (
	ErrTokenExpired           = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         = errors.New("请求令牌格式有误")
	ErrTokenInvalid           = errors.New("请求令牌无效")
	ErrHeaderEmpty            = errors.New("需要认证才能访问")
	ErrHeaderMalformed        = errors.New("请求头中 Authorization 格式有误")
	ErrTokenBlocked           = errors.New("令牌已加入黑名单")
)

// JWT JWT对象
type JWT struct {
	SignKey    []byte        // 签名密钥
	MaxRefresh time.Duration // 最大刷新时间
}

// CustomClaims 自定义载荷
type CustomClaims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	ExpireAtTime int64  `json:"expire_time"`
	jwtPkg.RegisteredClaims
}

// NewJWT 创建JWT实例
func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetConfig().JWT.Secret),
		MaxRefresh: config.GetConfig().JWT.MaxRefreshTime * time.Minute,
	}
}

// IssueToken 生成Token
func (j *JWT) IssueToken(userID, username string) string {
	expireAtTime := j.expireAtTime()
	claims := CustomClaims{
		UserID:       userID,
		Username:     username,
		ExpireAtTime: expireAtTime.Unix(),
		RegisteredClaims: jwtPkg.RegisteredClaims{
			ExpiresAt: jwtPkg.NewNumericDate(expireAtTime),
			IssuedAt:  jwtPkg.NewNumericDate(time.Now()),
			NotBefore: jwtPkg.NewNumericDate(time.Now()),
			Issuer:    config.GetConfig().Server.Name,
			ID:        userID,
		},
	}

	token, err := j.createToken(claims)
	if err != nil {
		return ""
	}

	return token
}

// createToken 创建Token
func (j *JWT) createToken(claims CustomClaims) (string, error) {
	token := jwtPkg.NewWithClaims(jwtPkg.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

// expireAtTime 计算过期时间
func (j *JWT) expireAtTime() time.Time {
	expireTime := int64(config.GetConfig().JWT.ExpireTime)
	expire := time.Duration(expireTime) * time.Minute
	return time.Now().Add(expire)
}

// ParserToken 解析Token
func (j *JWT) ParserToken(c *fiber.Ctx) (*CustomClaims, error) {
	tokenString, parseErr := j.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	token, err := j.parseTokenString(tokenString)
	if err != nil {
		validationErr, ok := err.(*jwtPkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtPkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtPkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	// 检查黑名单
	if j.tokenInBlacklist(tokenString) {
		return nil, ErrTokenBlocked
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新Token
func (j *JWT) RefreshToken(c *fiber.Ctx) (string, error) {
	tokenString, parseErr := j.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	token, err := j.parseTokenString(tokenString)
	if err != nil {
		validationErr, ok := err.(*jwtPkg.ValidationError)
		if !ok || validationErr.Errors != jwtPkg.ValidationErrorExpired {
			return "", err
		}
	}

	if j.tokenInBlacklist(tokenString) {
		return "", ErrTokenBlocked
	}

	claims := token.Claims.(*CustomClaims)

	// 检查是否超过最大刷新时间
	x := time.Now().Add(-j.MaxRefresh).Unix()
	if claims.IssuedAt.Unix() > x {
		// 更新过期时间
		claims.RegisteredClaims.ExpiresAt = jwtPkg.NewNumericDate(j.expireAtTime())

		// 旧Token加入黑名单
		expiresIn := claims.IssuedAt.Add(j.MaxRefresh).Unix() - time.Now().Unix()
		j.blacklistAdd(tokenString, expiresIn)

		return j.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// getTokenFromHeader 从请求头获取Token
func (j *JWT) getTokenFromHeader(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == common.TokenType) {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}

// parseTokenString 解析Token字符串
func (j *JWT) parseTokenString(tokenString string) (*jwtPkg.Token, error) {
	return jwtPkg.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtPkg.Token) (interface{}, error) {
		return j.SignKey, nil
	})
}

// blacklistAdd 添加Token到黑名单
func (j *JWT) blacklistAdd(token string, expiresIn int64) error {
	key := fmt.Sprintf("%s_%s", common.CacheTokenPrefix, util.Md5(token))
	cache.SaveString(key, token, time.Duration(expiresIn)*time.Second)
	return nil
}

// tokenInBlacklist 检查Token是否在黑名单
func (j *JWT) tokenInBlacklist(token string) bool {
	key := fmt.Sprintf("%s_%s", common.CacheTokenPrefix, util.Md5(token))
	return cache.GetString(key) != ""
}
