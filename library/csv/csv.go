package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ========================================
// CSV 读取
// ========================================

// Read 读取 CSV 文件
func Read(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

// ReadToMap 读取 CSV 并转为 Map 切片（第一行作为键）
func ReadToMap(filePath string) ([]map[string]string, error) {
	rows, err := Read(filePath)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	headers := rows[0]
	result := make([]map[string]string, 0, len(rows)-1)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		record := make(map[string]string)
		for j, header := range headers {
			if j < len(row) {
				record[header] = row[j]
			} else {
				record[header] = ""
			}
		}
		result = append(result, record)
	}

	return result, nil
}

// ReadWithCustomDelimiter 使用自定义分隔符读取
func ReadWithCustomDelimiter(filePath string, delimiter rune) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	return reader.ReadAll()
}

// ========================================
// CSV 写入
// ========================================

// Write 写入 CSV 文件
func Write(filePath string, data [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(data)
}

// WriteFromMap 从 Map 数据写入 CSV
func WriteFromMap(filePath string, headers []string, data []map[string]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 写入数据
	for _, record := range data {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = record[header]
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Append 追加数据到 CSV
func Append(filePath string, data [][]string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(data)
}

// WriteWithCustomDelimiter 使用自定义分隔符写入
func WriteWithCustomDelimiter(filePath string, data [][]string, delimiter rune) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = delimiter
	defer writer.Flush()

	return writer.WriteAll(data)
}

// ========================================
// 辅助函数
// ========================================

// ConvertToCSV 转换二维数组为 CSV 字符串
func ConvertToCSV(data [][]string) (string, error) {
	var result string
	for _, row := range data {
		for i, cell := range row {
			if i > 0 {
				result += ","
			}
			// 如果包含逗号或引号，需要转义
			if containsSpecialChar(cell) {
				result += fmt.Sprintf(`"%s"`, escapeQuotes(cell))
			} else {
				result += cell
			}
		}
		result += "\n"
	}
	return result, nil
}

func containsSpecialChar(s string) bool {
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' {
			return true
		}
	}
	return false
}

func escapeQuotes(s string) string {
	result := ""
	for _, c := range s {
		if c == '"' {
			result += `""`
		} else {
			result += string(c)
		}
	}
	return result
}

// ========================================
// 使用示例（注释）
// ========================================
/*
// 1. 读取 CSV
data, err := csv.Read("data.csv")
for _, row := range data {
    fmt.Println(row)
}

// 2. 读取为 Map
records, err := csv.ReadToMap("users.csv")
for _, record := range records {
    fmt.Printf("Name: %s, Age: %s\n", record["name"], record["age"])
}

// 3. 写入 CSV
data := [][]string{
    {"姓名", "年龄", "城市"},
    {"张三", "25", "北京"},
    {"李四", "30", "上海"},
}
err := csv.Write("output.csv", data)

// 4. 从 Map 写入
headers := []string{"name", "age", "city"}
users := []map[string]string{
    {"name": "张三", "age": "25", "city": "北京"},
    {"name": "李四", "age": "30", "city": "上海"},
}
err := csv.WriteFromMap("users.csv", headers, users)

// 5. 追加数据
newData := [][]string{
    {"王五", "28", "广州"},
}
err := csv.Append("users.csv", newData)

// 6. 使用 Tab 分隔符
data, err := csv.ReadWithCustomDelimiter("data.tsv", '\t')
err := csv.WriteWithCustomDelimiter("output.tsv", data, '\t')
*/
