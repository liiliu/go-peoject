-- 数据库初始化脚本
-- 使用方法: mysql -u root -p < scripts/init_database.sql

-- 创建数据库
CREATE DATABASE IF NOT EXISTS your_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE your_db;

-- ========================================
-- 用户表示例
-- ========================================
CREATE TABLE IF NOT EXISTS `app_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（MD5或bcrypt加密）',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1-正常 0-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ========================================
-- 示例：角色表
-- ========================================
CREATE TABLE IF NOT EXISTS `app_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1-启用 0-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ========================================
-- 示例：用户角色关联表
-- ========================================
CREATE TABLE IF NOT EXISTS `app_user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- ========================================
-- 插入初始数据
-- ========================================

-- 插入默认管理员用户（密码: admin123 的MD5: 0192023a7bbd73250516f069df18b500）
INSERT INTO `app_users` (`username`, `password`, `email`, `status`) 
VALUES ('admin', '0192023a7bbd73250516f069df18b500', 'admin@example.com', 1)
ON DUPLICATE KEY UPDATE `username` = `username`;

-- 插入默认角色
INSERT INTO `app_roles` (`name`, `code`, `description`, `status`) 
VALUES 
('管理员', 'admin', '系统管理员角色', 1),
('普通用户', 'user', '普通用户角色', 1)
ON DUPLICATE KEY UPDATE `name` = `name`;

-- ========================================
-- 完成提示
-- ========================================
SELECT '✅ 数据库初始化完成！' AS message;
SELECT CONCAT('数据库: ', DATABASE()) AS info;
SELECT CONCAT('表数量: ', COUNT(*)) AS table_count FROM information_schema.tables WHERE table_schema = DATABASE();
