package user

import (
	"your_project/app/model"
	"your_project/app/view"
	"your_project/library/database"
	"your_project/library/logger"

	"github.com/gofiber/fiber/v2"
)

// InitialUserRoutes 注册用户路由
func InitialUserRoutes(app *fiber.App) {
	user := app.Group("/v1/user")
	user.Get("/info", GetUserInfo)       // 获取用户信息
	user.Put("/update", UpdateUserInfo)  // 更新用户信息
	user.Get("/list", GetUserList)       // 获取用户列表
}

// ========================================
// 接口处理函数
// ========================================

// GetUserInfo 获取当前用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户的详细信息（需要Token）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  view.Result{data=model.User}  "成功"
// @Failure      1006 {object}  view.Result  "身份鉴权失败"
// @Failure      1007 {object}  view.Result  "TOKEN失效"
// @Failure      1001 {object}  view.Result  "用户不存在"
// @Router       /v1/user/info [get]
func GetUserInfo(c *fiber.Ctx) error {
	logger.InfoWithTrace(c, "user", "获取用户信息请求")

	// 从上下文获取用户ID（由 CheckToken 中间件设置）
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		logger.WarnWithTrace(c, "user", "未找到用户ID")
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeNoAuth, "身份鉴权失败"))
	}

	// 查询用户信息
	db := database.NewEngine()
	var user model.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		logger.WarnWithTrace(c, "user", "用户不存在: %s", userID)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeFailed, "用户不存在"))
	}

	// 隐藏敏感信息
	user.Password = ""

	logger.InfoWithTrace(c, "user", "用户 %s 获取信息成功", user.Username)
	return c.JSON(view.SuccessWithCtx(c, user))
}

// UpdateUserInfo 更新用户信息
// @Summary      更新用户信息
// @Description  更新当前登录用户的信息（需要Token）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      object{email=string,phone=string}  true  "更新参数"
// @Success      200      {object}  view.Result  "更新成功"
// @Failure      1006     {object}  view.Result  "身份鉴权失败"
// @Failure      1007     {object}  view.Result  "TOKEN失效"
// @Failure      1002     {object}  view.Result  "非法参数"
// @Router       /v1/user/update [put]
func UpdateUserInfo(c *fiber.Ctx) error {
	logger.InfoWithTrace(c, "user", "更新用户信息请求")

	// 从上下文获取用户ID
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		logger.WarnWithTrace(c, "user", "未找到用户ID")
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeNoAuth, "身份鉴权失败"))
	}

	// 解析请求参数
	var req struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil {
		logger.WarnWithTrace(c, "user", "参数解析失败: %v", err)
		return c.JSON(view.ErrorWithCtx(c, view.CodeInvalidBody))
	}

	// 更新用户信息
	db := database.NewEngine()
	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if len(updates) == 0 {
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeInvalidParams, "没有可更新的字段"))
	}

	if err := db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		logger.ErrorWithTrace(c, "user", "更新用户信息失败: %v", err)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeSystemError, "更新失败"))
	}

	logger.InfoWithTrace(c, "user", "用户 %s 更新信息成功", userID)
	return c.JSON(view.SuccessWithoutDataCtx(c))
}

// GetUserList 获取用户列表（示例：管理员接口）
// @Summary      获取用户列表
// @Description  获取所有用户列表（需要Token和管理员权限）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page      query     int     false  "页码"  default(1)
// @Param        page_size query     int     false  "每页数量"  default(10)
// @Success      200       {object}  view.Result{data=object{list=[]model.User,total=int64}}  "成功"
// @Failure      1006      {object}  view.Result  "身份鉴权失败"
// @Failure      1007      {object}  view.Result  "TOKEN失效"
// @Router       /v1/user/list [get]
func GetUserList(c *fiber.Ctx) error {
	logger.InfoWithTrace(c, "user", "获取用户列表请求")

	// 从上下文获取用户信息
	username, ok := c.Locals("username").(string)
	if !ok {
		logger.WarnWithTrace(c, "user", "未找到用户信息")
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeNoAuth, "身份鉴权失败"))
	}

	// 获取分页参数
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	offset := (page - 1) * pageSize

	// 查询用户列表
	db := database.NewEngine()
	var users []model.User
	var total int64

	db.Model(&model.User{}).Count(&total)
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		logger.ErrorWithTrace(c, "user", "查询用户列表失败: %v", err)
		return c.JSON(view.ErrorWithMsgCtx(c, view.CodeSystemError, "查询失败"))
	}

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	logger.InfoWithTrace(c, "user", "用户 %s 获取用户列表成功，共 %d 条", username, total)
	return c.JSON(view.SuccessWithCtx(c, fiber.Map{
		"list":  users,
		"total": total,
		"page":  page,
	}))
}
