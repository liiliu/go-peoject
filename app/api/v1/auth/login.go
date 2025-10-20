package auth

import (
	"your_project/app/model"
	"your_project/app/request"
	"your_project/app/view"
	"your_project/library/database"
	"your_project/library/jwt"
	"your_project/library/logger"
	"your_project/library/util"
	"your_project/library/validator"

	"github.com/gofiber/fiber/v2"
)

// InitialAuthRoutes 注册认证路由
func InitialAuthRoutes(app *fiber.App) {
	auth := app.Group("/v1/auth")
	auth.Post("/login", Login)
	auth.Post("/register", Register)
}

// ========================================
// 接口处理函数
// ========================================

// Login 用户登录（示例）
// @Summary      用户登录
// @Description  用户通过用户名和密码登录，返回 JWT Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      request.LoginRequest  true  "登录参数"
// @Success      200      {object}  view.Result{data=object{token=string,username=string,user_id=int}}  "登录成功"
// @Failure      1002     {object}  view.Result  "非法参数"
// @Failure      1004     {object}  view.Result  "参数验证失败"
// @Failure      1001     {object}  view.Result  "用户名或密码错误"
// @Router       /v1/auth/login [post]
func Login(c *fiber.Ctx) error {
	logger.InfoWithTrace(c, "auth", "用户登录请求")
	
	// 解析请求参数
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		logger.WarnWithTrace(c, "auth", "登录参数解析失败: %v", err)
		return c.JSON(view.ErrorWithCtx(c, view.CodeInvalidBody))
	}

	// 验证参数
	if err := validator.Validate(&req); err != nil {
		logger.WarnWithTrace(c, "auth", "登录参数验证失败: %v", err)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeInvalidParams, validator.GetErrorMsg(err)))
	}

	// 查询用户
	db := database.NewEngine()
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		logger.WarnWithTrace(c, "auth", "用户不存在: %s", req.Username)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeFailed, "用户名或密码错误"))
	}

	// 验证密码（这里简化处理，实际应使用bcrypt等加密）
	if util.Md5(req.Password) != user.Password {
		logger.WarnWithTrace(c, "auth", "用户 %s 密码错误", req.Username)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeFailed, "用户名或密码错误"))
	}

	// 检查用户状态
	if user.Status != 1 {
		logger.WarnWithTrace(c, "auth", "用户 %s 已被禁用", req.Username)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeFailed, "用户已被禁用"))
	}

	// 生成Token
	j := jwt.NewJWT()
	token := j.IssueToken(string(rune(user.ID)), user.Username)

	logger.InfoWithTrace(c, "auth", "用户 %s 登录成功", req.Username)
	return c.JSON(view.SuccessWithCtx(c, fiber.Map{
		"token":    token,
		"username": user.Username,
		"user_id":  user.ID,
	}))
}

// Register 用户注册（示例）
// @Summary      用户注册
// @Description  新用户注册账号
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      request.RegisterRequest  true  "注册参数"
// @Success      200      {object}  view.Result{data=object{user_id=int,username=string}}  "注册成功"
// @Failure      1002     {object}  view.Result  "非法参数"
// @Failure      1004     {object}  view.Result  "参数验证失败"
// @Failure      1001     {object}  view.Result  "用户名已存在"
// @Failure      9999     {object}  view.Result  "注册失败"
// @Router       /v1/auth/register [post]
func Register(c *fiber.Ctx) error {
	logger.InfoWithTrace(c, "auth", "用户注册请求")
	
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		logger.WarnWithTrace(c, "auth", "注册参数解析失败: %v", err)
		return c.JSON(view.ErrorWithCtx(c, view.CodeInvalidBody))
	}

	// 验证参数
	if err := validator.Validate(&req); err != nil {
		logger.WarnWithTrace(c, "auth", "注册参数验证失败: %v", err)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeInvalidParams, validator.GetErrorMsg(err)))
	}

	// 检查用户名是否已存在
	db := database.NewEngine()
	var count int64
	db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		logger.WarnWithTrace(c, "auth", "用户名已存在: %s", req.Username)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeFailed, "用户名已存在"))
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
		logger.ErrorWithTrace(c, "auth", "创建用户失败: %v", err)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeSystemError, "注册失败"))
	}

	logger.InfoWithTrace(c, "auth", "用户 %s 注册成功", req.Username)
	return c.JSON(view.SuccessWithCtx(c, fiber.Map{
		"user_id":  user.ID,
		"username": user.Username,
	}))
}
