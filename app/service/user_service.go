package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"your_project/app/model"
	"your_project/library/database"
	"your_project/library/logger"
	"your_project/library/util"
)

// ========================================
// 用户相关业务逻辑
// 只放复杂业务逻辑和多处使用的公共逻辑
// ========================================

// CheckUsernameExists 检查用户名是否存在（公共逻辑，多个接口会用到）
func CheckUsernameExists(c *fiber.Ctx, username string) (bool, error) {
	logger.DebugWithTrace(c, "service", "检查用户名是否存在: %s", username)

	var count int64
	db := database.NewEngine()
	if err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		logger.ErrorWithTrace(c, "service", "查询用户名失败: %v", err)
		return false, err
	}

	exists := count > 0
	logger.DebugWithTrace(c, "service", "用户名 %s 是否存在: %v", username, exists)
	return exists, nil
}

// RegisterUser 用户注册（复杂业务：多个步骤、事务处理）
func RegisterUser(c *fiber.Ctx, username, password, email string) error {
	logger.InfoWithTrace(c, "service", "开始注册用户: %s", username)

	// 1. 检查用户名是否已存在
	exists, err := CheckUsernameExists(c, username)
	if err != nil {
		return err
	}
	if exists {
		logger.WarnWithTrace(c, "service", "用户名已存在: %s", username)
		return errors.New("用户名已存在")
	}

	// 2. 密码加密
	hashedPassword := util.Md5(password)
	logger.DebugWithTrace(c, "service", "密码已加密")

	// 3. 创建用户
	db := database.NewEngine()
	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Status:   1,
	}

	if err := db.Create(user).Error; err != nil {
		logger.ErrorWithTrace(c, "service", "创建用户失败: %v", err)
		return errors.New("注册失败")
	}

	logger.InfoWithTrace(c, "service", "用户注册成功: %s, ID: %d", username, user.ID)
	return nil
}

// BatchUpdateUserStatus 批量更新用户状态（复杂操作：事务、批量处理）
func BatchUpdateUserStatus(c *fiber.Ctx, userIDs []uint, status int) error {
	logger.InfoWithTrace(c, "service", "批量更新用户状态，用户数: %d, 状态: %d", len(userIDs), status)

	if len(userIDs) == 0 {
		return errors.New("用户ID列表不能为空")
	}

	db := database.NewEngine()

	// 开启事务
	tx := db.Begin()

	// 批量更新
	if err := tx.Model(&model.User{}).Where("id IN ?", userIDs).Update("status", status).Error; err != nil {
		tx.Rollback()
		logger.ErrorWithTrace(c, "service", "批量更新失败: %v", err)
		return errors.New("批量更新失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.ErrorWithTrace(c, "service", "事务提交失败: %v", err)
		return errors.New("操作失败")
	}

	logger.InfoWithTrace(c, "service", "批量更新成功，影响用户数: %d", len(userIDs))
	return nil
}

// GetUserStatistics 获取用户统计信息（复杂查询：多表关联、数据聚合）
func GetUserStatistics(c *fiber.Ctx, userID uint) (map[string]interface{}, error) {
	logger.InfoWithTrace(c, "service", "获取用户统计信息: %d", userID)

	db := database.NewEngine()

	// 查询用户基本信息
	var user model.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		logger.WarnWithTrace(c, "service", "用户不存在: %d", userID)
		return nil, errors.New("用户不存在")
	}

	// 这里可以添加更多复杂的统计逻辑
	// 例如：统计订单数、积分、等级等

	result := map[string]interface{}{
		"user_id":     user.ID,
		"username":    user.Username,
		"register_at": user.CreatedAt,
		"status":      user.Status,
		// 可以添加更多统计数据
	}

	logger.InfoWithTrace(c, "service", "统计信息获取成功: %s", user.Username)
	return result, nil
}
