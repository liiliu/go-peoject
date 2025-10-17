package auth

import (
	"your_project/app/model"
	"your_project/app/request"
	"your_project/app/view"
	"your_project/library/database"
	"your_project/library/jwt"
	"your_project/library/util"
	"your_project/library/validator"

	"github.com/gofiber/fiber/v2"
)

// Login 用户登录（示例）
func Login(c *fiber.Ctx) error {
	// 解析请求参数
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(view.ErrorResult(view.CodeInvalidBody))
	}

	// 验证参数
	if err := validator.Validate(&req); err != nil {
		return c.JSON(view.ErrorResultWithMsg(view.CodeInvalidParams, validator.GetErrorMsg(err)))
	}

	// 查询用户
	db := database.NewEngine()
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.JSON(view.ErrorResultWithMsg(view.CodeFailed, "用户名或密码错误"))
	}

	// 验证密码（这里简化处理，实际应使用bcrypt等加密）
	if util.Md5(req.Password) != user.Password {
		return c.JSON(view.ErrorResultWithMsg(view.CodeFailed, "用户名或密码错误"))
	}

	// 检查用户状态
	if user.Status != 1 {
		return c.JSON(view.ErrorResultWithMsg(view.CodeFailed, "用户已被禁用"))
	}

	// 生成Token
	j := jwt.NewJWT()
	token := j.IssueToken(string(rune(user.ID)), user.Username)

	return c.JSON(view.SuccessResult(fiber.Map{
		"token":    token,
		"username": user.Username,
		"user_id":  user.ID,
	}))
}

// Register 用户注册（示例）
func Register(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(view.ErrorResult(view.CodeInvalidBody))
	}

	// 验证参数
	if err := validator.Validate(&req); err != nil {
		return c.JSON(view.ErrorResultWithMsg(view.CodeInvalidParams, validator.GetErrorMsg(err)))
	}

	// 检查用户名是否已存在
	db := database.NewEngine()
	var count int64
	db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return c.JSON(view.ErrorResultWithMsg(view.CodeFailed, "用户名已存在"))
	}

	// 创建用户（密码应使用bcrypt等加密，这里简化处理）
	user := model.User{
		Username: req.Username,
		Password: util.Md5(req.Password),
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(view.ErrorResultWithMsg(view.CodeSystemError, "注册失败"))
	}

	return c.JSON(view.SuccessResult(fiber.Map{
		"user_id":  user.ID,
		"username": user.Username,
	}))
}

// InitialAuthRoutes 注册认证路由
func InitialAuthRoutes(app *fiber.App) {
	auth := app.Group("/v1/auth")
	auth.Post("/login", Login)
	auth.Post("/register", Register)
}
