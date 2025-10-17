package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// 注册自定义字段名函数，使用 label 标签作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if name == "" {
			name = fld.Tag.Get("json")
		}
		if name == "" {
			name = fld.Name
		}
		return name
	})
}

// Validate 验证结构体
func Validate(data interface{}) error {
	return validate.Struct(data)
}

// ValidateWithMsg 验证结构体并使用自定义错误消息
// customMessages 格式: map["字段名.规则"] = "错误消息"
// 例如: map["Username.required"] = "请输入用户名"
func ValidateWithMsg(data interface{}, customMessages map[string]string) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var msgs []string
	for _, fieldError := range validationErrors {
		// 尝试获取自定义消息
		field := fieldError.Field()
		tag := fieldError.Tag()
		
		// 查找完全匹配的自定义消息 "字段名.规则"
		key := field + "." + tag
		if msg, exists := customMessages[key]; exists {
			msgs = append(msgs, msg)
			continue
		}
		
		// 查找字段级别的自定义消息 "字段名"
		if msg, exists := customMessages[field]; exists {
			msgs = append(msgs, msg)
			continue
		}
		
		// 使用默认消息
		msgs = append(msgs, getFieldError(fieldError))
	}

	if len(msgs) > 0 {
		return fmt.Errorf("%s", strings.Join(msgs, "; "))
	}

	return err
}

// GetErrorMsg 获取友好的错误消息
func GetErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var msgs []string
	for _, fieldError := range validationErrors {
		msgs = append(msgs, getFieldError(fieldError))
	}

	return strings.Join(msgs, "; ")
}

// getFieldError 获取单个字段的错误信息
func getFieldError(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s 不能为空", field)
	case "min":
		return fmt.Sprintf("%s 长度不能少于 %s", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s 长度不能超过 %s", field, fieldError.Param())
	case "email":
		return fmt.Sprintf("%s 格式不正确", field)
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", field, fieldError.Param())
	case "alphanum":
		return fmt.Sprintf("%s 只能包含字母和数字", field)
	case "numeric":
		return fmt.Sprintf("%s 必须为数字", field)
	case "gte":
		return fmt.Sprintf("%s 必须大于或等于 %s", field, fieldError.Param())
	case "lte":
		return fmt.Sprintf("%s 必须小于或等于 %s", field, fieldError.Param())
	default:
		return fmt.Sprintf("%s 验证失败: %s", field, tag)
	}
}
