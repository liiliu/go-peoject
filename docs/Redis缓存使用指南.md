# Redis 缓存使用指南

## 📋 概述

项目使用 **go-redis** 作为 Redis 客户端，提供缓存、分布式锁等功能。

---

## 🚀 快速开始

### 全局 Redis 实例

```go
import "your_project/library/cache"

// 使用全局实例
cache.Set("key", "value", 3600)
value, err := cache.Get("key")
```

---

## 📝 基础操作

### 1. 字符串操作

```go
// 设置值（永不过期）
cache.Set("username", "张三", 0)

// 设置值（带过期时间，秒）
cache.Set("code", "123456", 300)  // 5分钟

// 获取值
value, err := cache.Get("username")
if err != nil {
    logger.Error("cache", "获取失败: %v", err)
}

// 判断 key 是否存在
exists, _ := cache.Exists("username")
if exists {
    logger.Info("cache", "key 存在")
}

// 删除 key
cache.Del("username")

// 批量删除
cache.Del("key1", "key2", "key3")

// 设置过期时间
cache.Expire("username", 3600)  // 1小时

// 获取剩余过期时间（秒）
ttl, _ := cache.TTL("username")
logger.Info("cache", "剩余时间: %d 秒", ttl)
```

### 2. 数字操作

```go
// 自增
newValue, _ := cache.Incr("counter")
logger.Info("cache", "计数器: %d", newValue)

// 自增指定值
cache.IncrBy("counter", 10)

// 自减
cache.Decr("counter")

// 自减指定值
cache.DecrBy("counter", 5)
```

### 3. Hash 操作

```go
// 设置单个字段
cache.HSet("user:1", "name", "张三")
cache.HSet("user:1", "age", 25)

// 批量设置
cache.HMSet("user:1", map[string]interface{}{
    "name":  "张三",
    "age":   25,
    "email": "zhangsan@example.com",
})

// 获取单个字段
name, _ := cache.HGet("user:1", "name")

// 获取所有字段
user, _ := cache.HGetAll("user:1")
logger.Info("cache", "用户信息: %v", user)

// 删除字段
cache.HDel("user:1", "email")

// 判断字段是否存在
exists, _ := cache.HExists("user:1", "name")

// 获取所有键
keys, _ := cache.HKeys("user:1")

// 获取所有值
values, _ := cache.HVals("user:1")

// 字段数量
count, _ := cache.HLen("user:1")
```

### 4. 列表操作

```go
// 左侧推入
cache.LPush("messages", "消息1", "消息2")

// 右侧推入
cache.RPush("messages", "消息3")

// 左侧弹出
msg, _ := cache.LPop("messages")

// 右侧弹出
msg, _ := cache.RPop("messages")

// 获取列表长度
length, _ := cache.LLen("messages")

// 获取列表范围
messages, _ := cache.LRange("messages", 0, -1)  // 获取全部

// 获取指定索引的元素
msg, _ := cache.LIndex("messages", 0)

// 修剪列表（只保留指定范围）
cache.LTrim("messages", 0, 99)  // 只保留前100条
```

### 5. 集合操作

```go
// 添加成员
cache.SAdd("tags", "Go", "Redis", "MySQL")

// 获取所有成员
members, _ := cache.SMembers("tags")

// 判断是否存在
exists, _ := cache.SIsMember("tags", "Go")

// 删除成员
cache.SRem("tags", "MySQL")

// 集合数量
count, _ := cache.SCard("tags")

// 随机获取成员
member, _ := cache.SRandMember("tags")

// 集合运算
// 交集
inter, _ := cache.SInter("tags1", "tags2")
// 并集
union, _ := cache.SUnion("tags1", "tags2")
// 差集
diff, _ := cache.SDiff("tags1", "tags2")
```

### 6. 有序集合操作

```go
// 添加成员（带分数）
cache.ZAdd("rank", map[string]float64{
    "user1": 100,
    "user2": 200,
    "user3": 150,
})

// 获取排名（从小到大）
rank, _ := cache.ZRank("rank", "user1")

// 获取排名（从大到小）
rank, _ := cache.ZRevRank("rank", "user1")

// 获取分数
score, _ := cache.ZScore("rank", "user1")

// 增加分数
newScore, _ := cache.ZIncrBy("rank", 10, "user1")

// 获取范围（按排名）
users, _ := cache.ZRange("rank", 0, 9)  // 前10名

// 获取范围（按分数）
users, _ := cache.ZRangeByScore("rank", "100", "200")

// 获取数量
count, _ := cache.ZCard("rank")

// 删除成员
cache.ZRem("rank", "user1")
```

---

## 🎯 实战示例

### 示例 1：用户Token缓存

```go
// 生成并缓存 Token
func SaveUserToken(userID uint, token string) error {
    key := fmt.Sprintf("token:%d", userID)
    return cache.Set(key, token, 86400)  // 24小时
}

// 获取 Token
func GetUserToken(userID uint) (string, error) {
    key := fmt.Sprintf("token:%d", userID)
    return cache.Get(key)
}

// 删除 Token（登出）
func DeleteUserToken(userID uint) error {
    key := fmt.Sprintf("token:%d", userID)
    return cache.Del(key)
}
```

### 示例 2：验证码缓存

```go
// 发送验证码
func SendVerifyCode(mobile string) error {
    // 生成验证码
    code := util.GenerateCode(6)
    
    // 缓存验证码（5分钟）
    key := fmt.Sprintf("code:%s", mobile)
    if err := cache.Set(key, code, 300); err != nil {
        return err
    }
    
    // 发送短信...
    return nil
}

// 验证验证码
func VerifyCode(mobile, code string) bool {
    key := fmt.Sprintf("code:%s", mobile)
    
    savedCode, err := cache.Get(key)
    if err != nil {
        return false
    }
    
    if savedCode != code {
        return false
    }
    
    // 验证成功，删除验证码
    cache.Del(key)
    return true
}
```

### 示例 3：缓存查询结果

```go
// 获取用户信息（带缓存）
func GetUserWithCache(userID uint) (*model.User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    // 先从缓存读取
    cached, err := cache.Get(key)
    if err == nil && cached != "" {
        var user model.User
        if err := util.FromJSON(cached, &user); err == nil {
            logger.Info("cache", "从缓存读取用户: %d", userID)
            return &user, nil
        }
    }
    
    // 缓存未命中，从数据库查询
    var user model.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        return nil, err
    }
    
    // 写入缓存（1小时）
    userJSON, _ := util.ToJSON(user)
    cache.Set(key, userJSON, 3600)
    
    logger.Info("cache", "从数据库读取用户: %d", userID)
    return &user, nil
}

// 更新用户时删除缓存
func UpdateUser(userID uint, updates map[string]interface{}) error {
    if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
        return err
    }
    
    // 删除缓存
    key := fmt.Sprintf("user:%d", userID)
    cache.Del(key)
    
    return nil
}
```

### 示例 4：接口限流

```go
// 限流检查（每分钟最多 10 次）
func CheckRateLimit(userID uint) bool {
    key := fmt.Sprintf("ratelimit:%d", userID)
    
    // 获取当前计数
    count, _ := cache.Incr(key)
    
    // 首次访问，设置过期时间
    if count == 1 {
        cache.Expire(key, 60)  // 1分钟
    }
    
    // 超过限制
    if count > 10 {
        logger.Warn("ratelimit", "用户 %d 请求过于频繁", userID)
        return false
    }
    
    return true
}
```

### 示例 5：排行榜

```go
// 更新用户积分
func UpdateUserScore(userID uint, score float64) error {
    return cache.ZAdd("rank:score", map[string]float64{
        fmt.Sprintf("%d", userID): score,
    })
}

// 获取排行榜（前10名）
func GetTopUsers(limit int) ([]model.User, error) {
    // 从 Redis 获取用户 ID
    userIDs, err := cache.ZRevRange("rank:score", 0, int64(limit-1))
    if err != nil {
        return nil, err
    }
    
    // 查询用户信息
    var users []model.User
    for _, userIDStr := range userIDs {
        userID, _ := strconv.ParseUint(userIDStr, 10, 32)
        var user model.User
        if err := database.DB.First(&user, userID).Error; err == nil {
            // 获取积分
            score, _ := cache.ZScore("rank:score", userIDStr)
            user.Score = int(score)
            users = append(users, user)
        }
    }
    
    return users, nil
}

// 获取用户排名
func GetUserRank(userID uint) (int64, error) {
    rank, err := cache.ZRevRank("rank:score", fmt.Sprintf("%d", userID))
    return rank + 1, err  // 排名从1开始
}
```

---

## 🔒 分布式锁

### 基础分布式锁

```go
// 获取锁
func AcquireLock(key string, expiration int) bool {
    // 使用 SETNX 实现
    success, err := cache.SetNX(key, "locked", expiration)
    if err != nil || !success {
        return false
    }
    return true
}

// 释放锁
func ReleaseLock(key string) {
    cache.Del(key)
}

// 使用示例
lockKey := "lock:order:123"
if !AcquireLock(lockKey, 10) {
    return errors.New("操作进行中，请稍后重试")
}
defer ReleaseLock(lockKey)

// 执行业务逻辑...
```

### 高级分布式锁（带重试）

```go
func ProcessOrderWithLock(orderID uint) error {
    lockKey := fmt.Sprintf("lock:order:%d", orderID)
    
    // 尝试获取锁（重试3次）
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        if AcquireLock(lockKey, 30) {
            defer ReleaseLock(lockKey)
            
            // 处理订单
            return processOrder(orderID)
        }
        
        // 等待一段时间后重试
        time.Sleep(time.Second)
    }
    
    return errors.New("系统繁忙，请稍后重试")
}
```

---

## 🔄 Pipeline 批量操作

```go
// 使用 Pipeline 批量操作
func BatchSetCache(data map[string]string) error {
    pipe := cache.Client.Pipeline()
    
    for key, value := range data {
        pipe.Set(context.Background(), key, value, 3600*time.Second)
    }
    
    _, err := pipe.Exec(context.Background())
    return err
}

// 批量获取
func BatchGetCache(keys []string) (map[string]string, error) {
    pipe := cache.Client.Pipeline()
    
    cmds := make([]*redis.StringCmd, len(keys))
    for i, key := range keys {
        cmds[i] = pipe.Get(context.Background(), key)
    }
    
    _, err := pipe.Exec(context.Background())
    if err != nil && err != redis.Nil {
        return nil, err
    }
    
    result := make(map[string]string)
    for i, cmd := range cmds {
        val, err := cmd.Result()
        if err == nil {
            result[keys[i]] = val
        }
    }
    
    return result, nil
}
```

---

## ⚠️ 注意事项

### 1. 缓存穿透

```go
// ❌ 问题：查询不存在的数据，缓存和数据库都没有
func GetUser(userID uint) (*model.User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    // 缓存未命中
    cached, _ := cache.Get(key)
    if cached == "" {
        // 查询数据库（不存在的数据每次都查数据库）
        var user model.User
        database.DB.First(&user, userID)
        return &user, nil
    }
    return nil, nil
}

// ✅ 解决：缓存空值
func GetUserSafe(userID uint) (*model.User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    cached, _ := cache.Get(key)
    if cached == "NULL" {
        // 缓存的空值
        return nil, errors.New("用户不存在")
    }
    
    if cached != "" {
        var user model.User
        util.FromJSON(cached, &user)
        return &user, nil
    }
    
    var user model.User
    err := database.DB.First(&user, userID).Error
    if err != nil {
        // 缓存空值（短时间）
        cache.Set(key, "NULL", 60)
        return nil, err
    }
    
    // 缓存正常值
    userJSON, _ := util.ToJSON(user)
    cache.Set(key, userJSON, 3600)
    return &user, nil
}
```

### 2. 缓存雪崩

```go
// ❌ 问题：大量缓存同时过期
for _, user := range users {
    key := fmt.Sprintf("user:%d", user.ID)
    cache.Set(key, userJSON, 3600)  // 都是1小时过期
}

// ✅ 解决：添加随机过期时间
for _, user := range users {
    key := fmt.Sprintf("user:%d", user.ID)
    // 1小时 + 随机0-300秒
    expiration := 3600 + rand.Intn(300)
    cache.Set(key, userJSON, expiration)
}
```

### 3. 缓存更新策略

```go
// 先更新数据库，再删除缓存
func UpdateUser(userID uint, updates map[string]interface{}) error {
    // 更新数据库
    if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
        return err
    }
    
    // 删除缓存（让下次查询重新加载）
    key := fmt.Sprintf("user:%d", userID)
    cache.Del(key)
    
    return nil
}
```

---

## 📚 相关资源

- **go-redis 文档**：https://redis.uptrace.dev/
- **Redis 命令参考**：https://redis.io/commands

---

**更新日期**：2025-10-16  
**版本**：v1.0
