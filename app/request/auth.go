package request

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" label:"用户名" validate:"required,min=3,max=50"`
	Password string `json:"password" label:"密码" validate:"required,min=6"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" label:"用户名" validate:"required,min=3,max=50"`
	Password string `json:"password" label:"密码" validate:"required,min=6"`
	Email    string `json:"email" label:"邮箱" validate:"omitempty,email"`
	Phone    string `json:"phone" label:"手机号" validate:"omitempty,len=11"`
}
