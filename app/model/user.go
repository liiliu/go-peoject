package model

// User 用户模型（示例）
type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Email     string `gorm:"size:100" json:"email"`
	Phone     string `gorm:"size:20" json:"phone"`
	Status    int    `gorm:"default:1;comment:状态 1-正常 0-禁用" json:"status"`
	CreatedAt int64  `gorm:"autoCreateTime:milli;comment:创建时间(毫秒)" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli;comment:更新时间(毫秒)" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
