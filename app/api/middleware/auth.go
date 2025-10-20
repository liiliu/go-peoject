package middleware

import (
	"strings"
	"your_project/app/view"
	"your_project/library/common"
	"your_project/library/jwt"
	"your_project/library/logger"
	"your_project/library/util"

	"github.com/gofiber/fiber/v2"
)

// NoAuthUrls 不需要权限校验的路由
var NoAuthUrls = []string{
	"/v1/auth/login",
	"/v1/auth/register",
	"/v1/health",
	"/swagger",  // Swagger 文档路径前缀
}

// CheckToken 验证Token中间件
func CheckToken(c *fiber.Ctx) error {
	urlPath := string(c.Request().URI().Path())

	// 不需要验证的接口（精确匹配）
	if util.InStringSlice(urlPath, NoAuthUrls) {
		return c.Next()
	}
	
	// Swagger 路径前缀匹配
	if strings.HasPrefix(urlPath, "/swagger") {
		return c.Next()
	}

	// 解析Token
	newJWT := jwt.NewJWT()
	claims, err := newJWT.ParserToken(c)
	if err != nil {
		logger.Error(common.LogTagAuthError, "Token verification failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(view.ErrorResult(view.CodeInvalidToken))
	}

	// 将用户信息存储到上下文
	c.Locals("user_id", claims.UserID)
	c.Locals("username", claims.Username)

	return c.Next()
}
