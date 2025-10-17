package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ========================================
// 字符串工具
// ========================================

// IsEmpty 判断字符串是否为空（去除空格后）
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNotEmpty 判断字符串是否不为空
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// TrimSpace 去除字符串两端空格
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// Contains 检查字符串是否包含子串
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// HasPrefix 检查字符串是否以指定前缀开头
func HasPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// HasSuffix 检查字符串是否以指定后缀结尾
func HasSuffix(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// Substring 截取字符串（支持中文）
func Substring(str string, start, length int) string {
	runes := []rune(str)
	if start < 0 || start >= len(runes) {
		return ""
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

// ReverseString 反转字符串
func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CamelToSnake 驼峰转蛇形 (UserName -> user_name)
func CamelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SnakeToCamel 蛇形转驼峰 (user_name -> UserName)
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// MaskString 遮罩字符串中间部分
// 例如: MaskString("13812345678", 3, 4, "*") -> "138****5678"
func MaskString(str string, start, end int, mask string) string {
	runes := []rune(str)
	if start < 0 || end < 0 || start+end >= len(runes) {
		return str
	}

	var result strings.Builder
	result.WriteString(string(runes[:start]))
	for i := start; i < len(runes)-end; i++ {
		result.WriteString(mask)
	}
	result.WriteString(string(runes[len(runes)-end:]))
	return result.String()
}

// ========================================
// 切片工具
// ========================================

// InStringSlice 判断字符串是否在切片中
func InStringSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// InIntSlice 判断整数是否在切片中
func InIntSlice(num int, list []int) bool {
	for _, v := range list {
		if v == num {
			return true
		}
	}
	return false
}

// UniqueStrings 字符串切片去重
func UniqueStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	keys := make(map[string]bool)
	list := make([]string, 0, len(input))
	for _, entry := range input {
		if !keys[entry] {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// UniqueInts 整数切片去重
func UniqueInts(input []int) []int {
	if len(input) == 0 {
		return input
	}
	keys := make(map[int]bool)
	list := make([]int, 0, len(input))
	for _, entry := range input {
		if !keys[entry] {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// SliceChunk 将切片分块
func SliceChunk(slice []interface{}, chunkSize int) [][]interface{} {
	if chunkSize <= 0 {
		return nil
	}
	var chunks [][]interface{}
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// StringSliceChunk 字符串切片分块
func StringSliceChunk(slice []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		return nil
	}
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// ========================================
// 加密工具
// ========================================

// Md5 MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Sha256 SHA256加密
func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// ========================================
// 随机工具
// ========================================

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randGen    = rand.New(randSource)
)

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randGen.Intn(len(charset))]
	}
	return string(b)
}

// RandomNumber 生成随机数字字符串
func RandomNumber(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randGen.Intn(len(charset))]
	}
	return string(b)
}

// RandomInt 生成指定范围的随机整数 [min, max]
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + randGen.Intn(max-min+1)
}

// GenerateCode 生成验证码（6位数字）
func GenerateCode(length int) string {
	if length <= 0 {
		length = 6
	}
	return RandomNumber(length)
}

// GenerateOrderNo 生成订单号（时间戳+随机数）
func GenerateOrderNo() string {
	return fmt.Sprintf("%d%s", time.Now().Unix(), RandomNumber(6))
}

// ========================================
// 类型转换工具
// ========================================

// ToString 转换为字符串
func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt 字符串转整数
func ToInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

// ToInt64 字符串转 int64
func ToInt64(str string) int64 {
	num, _ := strconv.ParseInt(str, 10, 64)
	return num
}

// ToFloat64 字符串转 float64
func ToFloat64(str string) float64 {
	num, _ := strconv.ParseFloat(str, 64)
	return num
}

// ToBool 字符串转布尔值
func ToBool(str string) bool {
	val, _ := strconv.ParseBool(str)
	return val
}

// ========================================
// JSON 工具
// ========================================

// ToJSON 对象转 JSON 字符串
func ToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// FromJSON JSON 字符串转对象
func FromJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

// ToJSONIndent 对象转格式化的 JSON 字符串
func ToJSONIndent(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}

// ========================================
// 验证工具
// ========================================

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsMobile 验证手机号（中国大陆）
func IsMobile(mobile string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, mobile)
	return matched
}

// IsURL 验证 URL 格式
func IsURL(url string) bool {
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(:\d+)?(/.*)?$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// IsIDCard 验证身份证号（中国大陆，简单验证）
func IsIDCard(idCard string) bool {
	pattern := `^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}

// RegexpMatchSlice 正则匹配切片中的任一项
func RegexpMatchSlice(str string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, str)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// ========================================
// 通用工具
// ========================================

// Ternary 三元运算符（泛型版本）
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// DefaultString 返回默认字符串（如果为空）
func DefaultString(str, defaultVal string) string {
	if IsEmpty(str) {
		return defaultVal
	}
	return str
}

// DefaultInt 返回默认整数（如果为0）
func DefaultInt(num, defaultVal int) int {
	if num == 0 {
		return defaultVal
	}
	return num
}

// Min 返回两个整数中的最小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max 返回两个整数中的最大值
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs 返回整数的绝对值
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
