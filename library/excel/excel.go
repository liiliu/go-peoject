package excel

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// ========================================
// Excel 读取
// ========================================

// ReadExcel 读取 Excel 文件
func ReadExcel(filePath string) (*excelize.File, error) {
	return excelize.OpenFile(filePath)
}

// ReadSheet 读取指定 Sheet 的所有行
func ReadSheet(filePath, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// ReadFirstSheet 读取第一个 Sheet
func ReadFirstSheet(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	return f.GetRows(sheetName)
}

// ReadSheetToMap 读取 Sheet 并转为 Map 切片（第一行作为键）
func ReadSheetToMap(filePath, sheetName string) ([]map[string]string, error) {
	rows, err := ReadSheet(filePath, sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet is empty")
	}

	// 第一行作为表头
	headers := rows[0]
	result := make([]map[string]string, 0, len(rows)-1)

	// 从第二行开始读取数据
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

// GetCellValue 获取单元格值
func GetCellValue(filePath, sheetName, cell string) (string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return f.GetCellValue(sheetName, cell)
}

// ========================================
// Excel 写入
// ========================================

// NewExcel 创建新的 Excel 文件
func NewExcel() *excelize.File {
	return excelize.NewFile()
}

// WriteData 写入数据到 Excel
func WriteData(filePath string, sheetName string, data [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	// 创建或获取 Sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 写入数据
	for rowIndex, row := range data {
		for colIndex, value := range row {
			cell := getCellName(rowIndex+1, colIndex+1)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置活动 Sheet
	f.SetActiveSheet(index)

	// 保存文件
	return f.SaveAs(filePath)
}

// WriteMapData 写入 Map 数据到 Excel（自动生成表头）
func WriteMapData(filePath string, sheetName string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}

	f := excelize.NewFile()
	defer f.Close()

	// 创建 Sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 获取所有键作为表头
	headers := make([]string, 0)
	for key := range data[0] {
		headers = append(headers, key)
	}

	// 写入表头
	for i, header := range headers {
		cell := getCellName(1, i+1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for rowIndex, record := range data {
		for colIndex, header := range headers {
			cell := getCellName(rowIndex+2, colIndex+1)
			f.SetCellValue(sheetName, cell, record[header])
		}
	}

	// 设置活动 Sheet
	f.SetActiveSheet(index)

	return f.SaveAs(filePath)
}

// AppendRow 追加行到现有 Excel
func AppendRow(filePath, sheetName string, row []interface{}) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	nextRow := len(rows) + 1
	for i, value := range row {
		cell := getCellName(nextRow, i+1)
		f.SetCellValue(sheetName, cell, value)
	}

	return f.Save()
}

// WriteLargeData 高性能写入大量数据（推荐用于 >10000 行数据）
// 使用流式写入，性能比 WriteData 快 10-50 倍
func WriteLargeData(filePath string, sheetName string, data [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	// 创建 Sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 创建流式写入器
	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	// 逐行写入数据
	for rowIndex, row := range data {
		// 转换为 excelize 需要的格式
		cells := make([]interface{}, len(row))
		for i, value := range row {
			cells[i] = excelize.Cell{Value: value}
		}

		// 写入行
		cell := getCellName(rowIndex+1, 1)
		if err := streamWriter.SetRow(cell, cells); err != nil {
			return err
		}
	}

	// 刷新流式写入器
	if err := streamWriter.Flush(); err != nil {
		return err
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filePath)
}

// WriteLargeDataWithHeader 高性能写入大量数据（带表头样式）
func WriteLargeDataWithHeader(filePath, sheetName string, headers []string, data [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 创建流式写入器
	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	// 定义表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// 写入表头
	headerCells := make([]interface{}, len(headers))
	for i, header := range headers {
		headerCells[i] = excelize.Cell{
			StyleID: headerStyle,
			Value:   header,
		}
	}
	if err := streamWriter.SetRow("A1", headerCells); err != nil {
		return err
	}

	// 写入数据行
	for rowIndex, row := range data {
		cells := make([]interface{}, len(row))
		for i, value := range row {
			cells[i] = excelize.Cell{Value: value}
		}

		cell := getCellName(rowIndex+2, 1) // 从第二行开始
		if err := streamWriter.SetRow(cell, cells); err != nil {
			return err
		}
	}

	if err := streamWriter.Flush(); err != nil {
		return err
	}

	// 设置列宽
	for i := range headers {
		col := GetColumnName(i + 1)
		f.SetColWidth(sheetName, col, col, 15)
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filePath)
}

// WriteDataBatch 批量写入（适合中等数据量：1000-10000行）
// 使用 SetSheetRow 批量设置，比逐个 SetCellValue 快
func WriteDataBatch(filePath string, sheetName string, data [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 批量写入行
	for rowIndex, row := range data {
		cell := getCellName(rowIndex+1, 1)
		if err := f.SetSheetRow(sheetName, cell, &row); err != nil {
			return err
		}
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filePath)
}

// ========================================
// 流式写入（边查询边写入，真正节省内存）
// ========================================

// StreamWriter 流式写入器（用于边查询边写入）
type StreamWriter struct {
	file         *excelize.File
	streamWriter *excelize.StreamWriter
	sheetName    string
	rowIndex     int
}

// NewStreamWriter 创建流式写入器
func NewStreamWriter(filePath, sheetName string) (*StreamWriter, error) {
	f := excelize.NewFile()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)

	streamWriter, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return nil, err
	}

	return &StreamWriter{
		file:         f,
		streamWriter: streamWriter,
		sheetName:    sheetName,
		rowIndex:     1,
	}, nil
}

// WriteHeader 写入表头
func (sw *StreamWriter) WriteHeader(headers []string) error {
	// 创建表头样式
	headerStyle, err := sw.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// 写入表头
	cells := make([]interface{}, len(headers))
	for i, header := range headers {
		cells[i] = excelize.Cell{
			StyleID: headerStyle,
			Value:   header,
		}
	}

	if err := sw.streamWriter.SetRow("A1", cells); err != nil {
		return err
	}

	sw.rowIndex = 2 // 下一行从第2行开始
	return nil
}

// WriteRow 写入一行数据
func (sw *StreamWriter) WriteRow(row []interface{}) error {
	cells := make([]interface{}, len(row))
	for i, value := range row {
		cells[i] = excelize.Cell{Value: value}
	}

	cell := getCellName(sw.rowIndex, 1)
	if err := sw.streamWriter.SetRow(cell, cells); err != nil {
		return err
	}

	sw.rowIndex++
	return nil
}

// WriteRows 批量写入多行
func (sw *StreamWriter) WriteRows(rows [][]interface{}) error {
	for _, row := range rows {
		if err := sw.WriteRow(row); err != nil {
			return err
		}
	}
	return nil
}

// SetColumnWidth 设置列宽
func (sw *StreamWriter) SetColumnWidth(startCol, endCol string, width float64) error {
	return sw.file.SetColWidth(sw.sheetName, startCol, endCol, width)
}

// Close 关闭并保存文件
func (sw *StreamWriter) Close(filePath string) error {
	if err := sw.streamWriter.Flush(); err != nil {
		return err
	}

	if err := sw.file.SaveAs(filePath); err != nil {
		return err
	}

	return sw.file.Close()
}

// WriteFromCallback 使用回调函数流式写入（推荐用于数据库分页查询）
// callback 函数每次返回一批数据，返回空切片表示结束
func WriteFromCallback(filePath, sheetName string, headers []string, callback func(page int) [][]interface{}) error {
	writer, err := NewStreamWriter(filePath, sheetName)
	if err != nil {
		return err
	}
	defer writer.Close(filePath)

	// 写入表头
	if len(headers) > 0 {
		if err := writer.WriteHeader(headers); err != nil {
			return err
		}

		// 设置列宽
		for i := range headers {
			col := GetColumnName(i + 1)
			writer.SetColumnWidth(col, col, 15)
		}
	}

	// 分页查询并写入
	page := 1
	for {
		rows := callback(page)
		if len(rows) == 0 {
			break
		}

		if err := writer.WriteRows(rows); err != nil {
			return err
		}

		page++
	}

	return nil
}

// WriteFromChannel 从通道读取数据并流式写入（推荐用于并发场景）
func WriteFromChannel(filePath, sheetName string, headers []string, dataChan <-chan []interface{}) error {
	writer, err := NewStreamWriter(filePath, sheetName)
	if err != nil {
		return err
	}
	defer writer.Close(filePath)

	// 写入表头
	if len(headers) > 0 {
		if err := writer.WriteHeader(headers); err != nil {
			return err
		}

		// 设置列宽
		for i := range headers {
			col := GetColumnName(i + 1)
			writer.SetColumnWidth(col, col, 15)
		}
	}

	// 从通道读取数据并写入
	for row := range dataChan {
		if err := writer.WriteRow(row); err != nil {
			return err
		}
	}

	return nil
}

// ========================================
// Excel 样式设置
// ========================================

// SetCellStyle 设置单元格样式
func SetCellStyle(f *excelize.File, sheetName, startCell, endCell string, style *excelize.Style) error {
	styleID, err := f.NewStyle(style)
	if err != nil {
		return err
	}
	return f.SetCellStyle(sheetName, startCell, endCell, styleID)
}

// SetHeaderStyle 设置表头样式
func SetHeaderStyle(f *excelize.File, sheetName string, endCol int) error {
	style := &excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "FFFFFF",
			Family: "微软雅黑",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}

	startCell := "A1"
	endCell := getCellName(1, endCol)
	return SetCellStyle(f, sheetName, startCell, endCell, style)
}

// SetColumnWidth 设置列宽
func SetColumnWidth(f *excelize.File, sheetName, startCol, endCol string, width float64) error {
	return f.SetColWidth(sheetName, startCol, endCol, width)
}

// AutoFilterTable 设置自动筛选
func AutoFilterTable(f *excelize.File, sheetName, startCell, endCell string) error {
	return f.AutoFilter(sheetName, startCell, endCell, "")
}

// ========================================
// 辅助函数
// ========================================

// getCellName 根据行列号获取单元格名称 (1,1 -> A1)
func getCellName(row, col int) string {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	return cell
}

// GetColumnName 获取列名 (1->A, 2->B, 27->AA)
func GetColumnName(col int) string {
	name := ""
	for col > 0 {
		col--
		name = string(rune('A'+col%26)) + name
		col /= 26
	}
	return name
}

// ========================================
// 高级功能
// ========================================

// ExportData 导出数据到 Excel（带样式）
func ExportData(filePath, sheetName string, headers []string, data [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	// 创建 Sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 写入表头
	for i, header := range headers {
		cell := getCellName(1, i+1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置表头样式
	if err := SetHeaderStyle(f, sheetName, len(headers)); err != nil {
		return err
	}

	// 写入数据
	for rowIndex, row := range data {
		for colIndex, value := range row {
			cell := getCellName(rowIndex+2, colIndex+1)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	for i := range headers {
		col := GetColumnName(i + 1)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// 设置自动筛选
	lastCell := getCellName(len(data)+1, len(headers))
	f.AutoFilter(sheetName, "A1", lastCell, "")

	// 冻结首行
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	f.SetActiveSheet(index)
	return f.SaveAs(filePath)
}

// MergeExcelFiles 合并多个 Excel 文件
func MergeExcelFiles(outputPath string, inputPaths []string) error {
	f := excelize.NewFile()
	defer f.Close()

	for i, inputPath := range inputPaths {
		source, err := excelize.OpenFile(inputPath)
		if err != nil {
			return err
		}

		sheetName := fmt.Sprintf("Sheet%d", i+1)
		if i == 0 {
			sheetName = "Sheet1"
		} else {
			f.NewSheet(sheetName)
		}

		// 复制数据
		sourceSheet := source.GetSheetName(0)
		rows, _ := source.GetRows(sourceSheet)
		for rowIdx, row := range rows {
			for colIdx, value := range row {
				cell := getCellName(rowIdx+1, colIdx+1)
				f.SetCellValue(sheetName, cell, value)
			}
		}
		source.Close()
	}

	return f.SaveAs(outputPath)
}

// ConvertToInt 字符串转整数
func ConvertToInt(value string) int {
	num, _ := strconv.Atoi(value)
	return num
}

// ConvertToFloat 字符串转浮点数
func ConvertToFloat(value string) float64 {
	num, _ := strconv.ParseFloat(value, 64)
	return num
}

// ========================================
// 使用示例（注释）
// ========================================
/*
// 1. 读取 Excel
rows, err := excel.ReadSheet("users.xlsx", "Sheet1")
for _, row := range rows {
    fmt.Println(row)
}

// 2. 读取为 Map
data, err := excel.ReadSheetToMap("users.xlsx", "Sheet1")
for _, record := range data {
    fmt.Printf("Name: %s, Age: %s\n", record["姓名"], record["年龄"])
}

// 3. 写入数据
data := [][]interface{}{
    {"姓名", "年龄", "城市"},
    {"张三", 25, "北京"},
    {"李四", 30, "上海"},
}
err := excel.WriteData("output.xlsx", "用户数据", data)

// 4. 写入 Map 数据
users := []map[string]interface{}{
    {"name": "张三", "age": 25, "city": "北京"},
    {"name": "李四", "age": 30, "city": "上海"},
}
err := excel.WriteMapData("users.xlsx", "Sheet1", users)

// 5. 导出带样式的 Excel
headers := []string{"姓名", "年龄", "城市", "职位"}
data := [][]interface{}{
    {"张三", 25, "北京", "工程师"},
    {"李四", 30, "上海", "经理"},
}
err := excel.ExportData("export.xlsx", "员工表", headers, data)

// 6. 追加行
row := []interface{}{"王五", 28, "广州"}
err := excel.AppendRow("users.xlsx", "Sheet1", row)

// 7. 合并多个 Excel
files := []string{"file1.xlsx", "file2.xlsx", "file3.xlsx"}
err := excel.MergeExcelFiles("merged.xlsx", files)
*/
