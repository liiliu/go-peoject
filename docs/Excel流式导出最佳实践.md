# Excel 流式导出最佳实践

## 💡 问题分析

### ❌ 错误做法：一次性加载所有数据

```go
// 问题：100 万条数据会占用 GB 级内存！
var users []model.User
database.DB.Find(&users)  // 一次性查询 100 万条

data := make([][]interface{}, len(users))
for i, user := range users {
    data[i] = []interface{}{user.ID, user.Name}
}

excel.WriteLargeData("users.xlsx", "用户", data)
```

**问题**：
- ❌ 占用大量内存（100万条 ≈ 1-2GB）
- ❌ 数据库查询慢
- ❌ 容易内存溢出

---

## ✅ 正确做法：边查询边写入

### 方法 1：使用回调函数（推荐）

```go
func ExportUsers() error {
    headers := []string{"ID", "用户名", "邮箱", "注册时间"}
    
    // 使用回调函数，每次返回一批数据
    return excel.WriteFromCallback("users.xlsx", "用户数据", headers, func(page int) [][]interface{} {
        const pageSize = 1000  // 每次查 1000 条
        
        var users []model.User
        result := database.DB.
            Limit(pageSize).
            Offset((page - 1) * pageSize).
            Find(&users)
        
        // 没有数据了，返回空切片表示结束
        if result.RowsAffected == 0 {
            return nil
        }
        
        // 转换为 Excel 数据
        data := make([][]interface{}, len(users))
        for i, user := range users {
            data[i] = []interface{}{
                user.ID,
                user.Username,
                user.Email,
                user.CreatedAt.Format("2006-01-02"),
            }
        }
        
        logger.Info("excel", "导出第 %d 批，共 %d 条", page, len(users))
        return data
    })
}
```

**优势**：
- ✅ **内存占用极低** - 每次只加载 1000 条数据
- ✅ **自动分页** - 框架自动循环调用
- ✅ **代码简洁** - 业务逻辑清晰

---

### 方法 2：使用通道（并发场景）

```go
func ExportUsersWithChannel() error {
    headers := []string{"ID", "用户名", "邮箱", "注册时间"}
    
    // 创建数据通道
    dataChan := make(chan []interface{}, 100)
    
    // 启动查询协程
    go func() {
        defer close(dataChan)
        
        const pageSize = 1000
        page := 1
        
        for {
            var users []model.User
            result := database.DB.
                Limit(pageSize).
                Offset((page - 1) * pageSize).
                Find(&users)
            
            if result.RowsAffected == 0 {
                break
            }
            
            // 发送到通道
            for _, user := range users {
                dataChan <- []interface{}{
                    user.ID,
                    user.Username,
                    user.Email,
                    user.CreatedAt.Format("2006-01-02"),
                }
            }
            
            page++
        }
    }()
    
    // 从通道读取数据并写入 Excel
    return excel.WriteFromChannel("users.xlsx", "用户数据", headers, dataChan)
}
```

**优势**：
- ✅ 支持并发
- ✅ 生产者消费者模式
- ✅ 更灵活的控制

---

### 方法 3：手动控制（最灵活）

```go
func ExportUsersManual() error {
    // 创建流式写入器
    writer, err := excel.NewStreamWriter("users.xlsx", "用户数据")
    if err != nil {
        return err
    }
    defer writer.Close("users.xlsx")
    
    // 写入表头
    headers := []string{"ID", "用户名", "邮箱", "注册时间"}
    if err := writer.WriteHeader(headers); err != nil {
        return err
    }
    
    // 设置列宽
    writer.SetColumnWidth("A", "D", 15)
    
    // 分页查询并写入
    const pageSize = 1000
    page := 1
    
    for {
        var users []model.User
        result := database.DB.
            Limit(pageSize).
            Offset((page - 1) * pageSize).
            Find(&users)
        
        if result.RowsAffected == 0 {
            break
        }
        
        // 逐行写入
        for _, user := range users {
            row := []interface{}{
                user.ID,
                user.Username,
                user.Email,
                user.CreatedAt.Format("2006-01-02"),
            }
            
            if err := writer.WriteRow(row); err != nil {
                return err
            }
        }
        
        logger.Info("excel", "已导出 %d 页，共 %d 条", page, page*pageSize)
        page++
    }
    
    return nil
}
```

**优势**：
- ✅ 完全控制写入过程
- ✅ 可以添加进度显示
- ✅ 可以中途暂停/取消

---

## 📊 性能对比

### 测试：导出 100 万条数据

| 方法 | 内存占用 | 耗时 | 说明 |
|------|---------|------|------|
| 一次性加载 | ~2GB | 30-60秒 | ❌ 容易 OOM |
| 回调函数 | ~50MB | 15-25秒 | ✅ 推荐 |
| 通道方式 | ~50MB | 15-25秒 | ✅ 并发推荐 |
| 手动控制 | ~50MB | 15-25秒 | ✅ 灵活控制 |

---

## 💡 实战示例

### 示例 1：导出订单（带条件查询）

```go
func ExportOrders(startDate, endDate time.Time) error {
    headers := []string{"订单号", "用户", "金额", "状态", "时间"}
    
    return excel.WriteFromCallback("orders.xlsx", "订单数据", headers, func(page int) [][]interface{} {
        const pageSize = 1000
        
        var orders []model.Order
        database.DB.
            Where("created_at BETWEEN ? AND ?", startDate, endDate).
            Preload("User").  // 关联查询
            Limit(pageSize).
            Offset((page - 1) * pageSize).
            Find(&orders)
        
        if len(orders) == 0 {
            return nil
        }
        
        data := make([][]interface{}, len(orders))
        for i, order := range orders {
            data[i] = []interface{}{
                order.OrderNo,
                order.User.Username,
                order.Amount,
                order.GetStatusText(),
                order.CreatedAt.Format("2006-01-02 15:04:05"),
            }
        }
        
        return data
    })
}
```

---

### 示例 2：导出并添加进度显示

```go
func ExportWithProgress(c *fiber.Ctx) error {
    // 先统计总数
    var total int64
    database.DB.Model(&model.User{}).Count(&total)
    
    writer, _ := excel.NewStreamWriter("users.xlsx", "用户数据")
    defer writer.Close("users.xlsx")
    
    headers := []string{"ID", "用户名", "邮箱"}
    writer.WriteHeader(headers)
    
    const pageSize = 1000
    exported := 0
    
    for page := 1; ; page++ {
        var users []model.User
        result := database.DB.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users)
        
        if result.RowsAffected == 0 {
            break
        }
        
        for _, user := range users {
            writer.WriteRow([]interface{}{user.ID, user.Username, user.Email})
            exported++
        }
        
        // 发送进度（WebSocket 或 SSE）
        progress := float64(exported) / float64(total) * 100
        logger.Info("export", "进度: %.2f%% (%d/%d)", progress, exported, total)
        
        // 可以通过 WebSocket 实时推送进度给前端
        // websocket.SendProgress(c, progress)
    }
    
    return nil
}
```

---

### 示例 3：多表关联查询导出

```go
func ExportUserOrders() error {
    headers := []string{"用户ID", "用户名", "订单总数", "总金额"}
    
    return excel.WriteFromCallback("user_orders.xlsx", "用户订单统计", headers, func(page int) [][]interface{} {
        const pageSize = 500
        
        type UserStat struct {
            UserID      uint
            Username    string
            OrderCount  int
            TotalAmount float64
        }
        
        var stats []UserStat
        database.DB.Raw(`
            SELECT 
                u.id as user_id,
                u.username,
                COUNT(o.id) as order_count,
                COALESCE(SUM(o.amount), 0) as total_amount
            FROM users u
            LEFT JOIN orders o ON u.id = o.user_id
            GROUP BY u.id, u.username
            LIMIT ? OFFSET ?
        `, pageSize, (page-1)*pageSize).Scan(&stats)
        
        if len(stats) == 0 {
            return nil
        }
        
        data := make([][]interface{}, len(stats))
        for i, stat := range stats {
            data[i] = []interface{}{
                stat.UserID,
                stat.Username,
                stat.OrderCount,
                stat.TotalAmount,
            }
        }
        
        return data
    })
}
```

---

### 示例 4：异步导出（后台任务）

```go
// 创建导出任务
func CreateExportTask(c *fiber.Ctx) error {
    // 创建任务记录
    task := &model.ExportTask{
        Type:   "users",
        Status: "pending",
        UserID: c.Locals("user_id").(uint),
    }
    database.DB.Create(task)
    
    // 异步执行导出
    go performExport(task)
    
    return c.JSON(fiber.Map{
        "code": 200,
        "data": task,
        "msg":  "导出任务已创建，请稍后下载",
    })
}

// 执行导出
func performExport(task *model.ExportTask) {
    defer func() {
        if err := recover(); err != nil {
            task.Status = "failed"
            task.Error = fmt.Sprintf("%v", err)
            database.DB.Save(task)
        }
    }()
    
    task.Status = "processing"
    database.DB.Save(task)
    
    filename := fmt.Sprintf("export_%d.xlsx", task.ID)
    filePath := filepath.Join("exports", filename)
    
    // 流式导出
    err := excel.WriteFromCallback(filePath, "数据", []string{"ID", "名称"}, func(page int) [][]interface{} {
        // 查询数据...
        var users []model.User
        database.DB.Limit(1000).Offset((page - 1) * 1000).Find(&users)
        
        if len(users) == 0 {
            return nil
        }
        
        // 更新进度
        task.Progress = page * 1000
        database.DB.Save(task)
        
        // 转换数据
        data := make([][]interface{}, len(users))
        for i, user := range users {
            data[i] = []interface{}{user.ID, user.Username}
        }
        return data
    })
    
    if err != nil {
        task.Status = "failed"
        task.Error = err.Error()
    } else {
        task.Status = "completed"
        task.FilePath = filePath
    }
    
    database.DB.Save(task)
}
```

---

## 🎯 最佳实践建议

### 1. 选择合适的方法

| 场景 | 推荐方法 |
|------|---------|
| 简单导出 | `WriteFromCallback` |
| 需要并发 | `WriteFromChannel` |
| 需要进度 | `NewStreamWriter` 手动控制 |
| 异步导出 | 后台任务 + 任意方法 |

### 2. 性能优化

```go
// ✅ 好的做法
const pageSize = 1000  // 每批 1000 条，平衡性能和内存

// ❌ 避免
const pageSize = 100   // 太小，查询次数太多
const pageSize = 10000 // 太大，内存占用高
```

### 3. 错误处理

```go
err := excel.WriteFromCallback(filePath, sheetName, headers, func(page int) [][]interface{} {
    var users []model.User
    if err := database.DB.Find(&users).Error; err != nil {
        logger.Error("export", "查询失败: %v", err)
        return nil  // 返回 nil 停止导出
    }
    // ...
})

if err != nil {
    logger.Error("export", "导出失败: %v", err)
    return err
}
```

### 4. 资源清理

```go
writer, err := excel.NewStreamWriter(filePath, sheetName)
if err != nil {
    return err
}
defer writer.Close(filePath)  // 确保文件正确关闭
```

---

## 📝 总结

| 对比项 | 一次性加载 | 流式写入 |
|--------|-----------|---------|
| **内存** | 2GB+ | 50MB |
| **速度** | 中等 | 快 |
| **稳定性** | 易崩溃 | 稳定 |
| **代码** | 简单 | 略复杂 |
| **推荐度** | ❌ | ✅ |

**推荐使用**：`WriteFromCallback` 或 `NewStreamWriter`

---

**更新日期**：2025-10-16  
**版本**：v2.0
