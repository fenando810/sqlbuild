# SQL Builder 使用说明

动态 SQL 查询构建器，支持链式调用、参数化查询、自动空值过滤。

## 导入

```go
import "github.com/fenando810/sqlbuilder"
```

## 快速开始

```go
qb := sqlbuilder.NewQueryBuilder("SELECT * FROM users").
    WhereEq("status", "active").
    WhereGt("age", 18).
    OrderBy("name", "ASC").
    Limit(10)

sql, args := qb.Build()
// sql:  SELECT * FROM users WHERE status = ? AND age > ? ORDER BY name ASC LIMIT ?
// args: ["active", 18, 10]
```

调用 `Build()` 返回两个值：拼接好的 SQL 语句和对应的参数切片，可直接传给 `db.Query(sql, args...)` 等方法。

---

## 条件方法

所有条件方法自动过滤空值（`nil`、空字符串 `""`、零值 `time.Time`），不会生成无效条件。

### 等于 / 不等于

```go
qb.WhereEq("name", "alice")      // name = ?
qb.WhereNotEq("status", "deleted") // status != ?
```

空值自动忽略：

```go
qb.WhereEq("name", "")  // 不生成条件
qb.WhereEq("name", nil) // 不生成条件
```

### 比较运算

```go
qb.WhereGt("age", 18)   // age > ?
qb.WhereLt("age", 60)   // age < ?
qb.WhereGte("age", 18)  // age >= ?
qb.WhereLte("age", 60)  // age <= ?
```

### 模糊匹配

自动包裹 `%` 通配符，生成 `%pattern%` 形式：

```go
qb.WhereLike("name", "ali") // name LIKE ?  参数: "%ali%"
```

空字符串不生成条件。

### IN 查询

**泛型函数**（类型安全，需断开链式调用）：

```go
qb := sqlbuilder.NewQueryBuilder("SELECT * FROM users")
sqlbuilder.WhereIn(qb, "id", []int{1, 2, 3})
// id IN (?,?,?)  参数: [1, 2, 3]

sqlbuilder.WhereIn(qb, "status", []string{"active", "pending"})
// status IN (?,?)  参数: ["active", "pending"]
```

**方法版本**（保持链式调用）：

```go
qb := sqlbuilder.NewQueryBuilder("SELECT * FROM users").
    WhereInAny("id", []any{1, 2, 3}).
    WhereEq("status", "active")
// WHERE id IN (?,?,?) AND status = ?
```

空切片不生成条件。

### BETWEEN 范围查询

```go
qb.WhereBetween("age", 18, 60) // age BETWEEN ? AND ?
```

`nil` 或空字符串自动忽略：

```go
qb.WhereBetween("age", nil, 60)  // 不生成条件
qb.WhereBetween("age", "", 60)   // 不生成条件
```

### NULL 判断

```go
qb.WhereIsNull("deleted_at")     // deleted_at IS NULL
qb.WhereIsNotNull("email")       // email IS NOT NULL
```

### 时间范围查询

支持多种时间格式，返回 `(*QueryBuilder, error)`：

```go
qb, err := qb.WhereTimeRange("created_at", "2024-01-01", "2024-12-31")
// created_at >= ? AND created_at <= ?

qb, err := qb.WhereTimeRange("created_at", "2024-01-01", "")  // 仅开始时间
// created_at >= ?

qb, err := qb.WhereTimeRange("created_at", "", "2024-12-31")  // 仅结束时间
// created_at <= ?
```

支持的时间格式：

| 格式 | 示例 |
|------|------|
| RFC3339 | `2024-01-15T10:30:00Z` |
| 日期时间 | `2024-01-15 10:30:00` |
| 仅日期 | `2024-01-15` |
| 斜杠日期时间 | `2024/01/15 10:30:00` |
| 斜杠日期 | `2024/01/15` |

解析失败时返回 error：

```go
qb, err := qb.WhereTimeRange("created_at", "invalid-date", "2024-12-31")
// err: 解析开始时间失败: 无法解析时间字符串: invalid-date
```

### OR 条件组

```go
qb.WhereOr(
    []string{"status = ?", "role = ?"},
    "active", "admin",
)
// (status = ? OR role = ?)  参数: ["active", "admin"]
```

参数按顺序对应条件中的 `?` 占位符。

### 原始条件

支持带参数的原始 SQL 条件：

```go
qb.WhereRaw("1 = 1")                    // 无参数
qb.WhereRaw("status != ?", "banned")    // 带参数
qb.WhereRaw("id > ? AND id < ?", 10, 20) // 多参数
```

> ⚠️ 谨慎使用，不要拼接用户输入，存在 SQL 注入风险。

---

## 后缀方法

后缀方法支持追加调用，不会互相覆盖。

### 排序

自动校验排序方向，非法值默认为 `ASC`：

```go
qb.OrderBy("name", "ASC")    // ORDER BY name ASC
qb.OrderBy("age", "desc")    // ORDER BY age DESC
qb.OrderBy("id", "invalid")  // ORDER BY id ASC (非法方向自动修正)
```

### 分页

```go
qb.Limit(10)              // LIMIT ?
qb.LimitOffset(10, 20)    // LIMIT ? OFFSET ?
```

### 分组

```go
qb.GroupBy("department")           // GROUP BY department
qb.GroupBy("department", "role")   // GROUP BY department, role
```

### 自定义后缀

```go
qb.SetSuffix("ORDER BY name ASC")
qb.SetSuffix("LIMIT ?", 10)
// 追加结果: ORDER BY name ASC LIMIT ?
```

---

## 完整示例

### 多条件组合查询

```go
qb := sqlbuilder.NewQueryBuilder("SELECT * FROM orders").
    WhereEq("status", "paid").
    WhereIsNotNull("shipped_at").
    WhereOr(
        []string{"priority = ?", "vip = ?"},
        "high", true,
    ).
    WhereBetween("amount", 100, 5000).
    OrderBy("created_at", "DESC").
    LimitOffset(20, 0)

sql, args := qb.Build()
// sql:  SELECT * FROM orders WHERE status = ? AND shipped_at IS NOT NULL
//       AND (priority = ? OR vip = ?) AND amount BETWEEN ? AND ?
//       ORDER BY created_at DESC LIMIT ? OFFSET ?
// args: ["paid", "high", true, 100, 5000, 20, 0]
```

### 动态条件查询

```go
func BuildUserQuery(name string, minAge int, status string) (string, []any) {
    qb := sqlbuilder.NewQueryBuilder("SELECT * FROM users")

    // 空值自动忽略，无需手动判断
    qb.WhereEq("name", name)
    qb.WhereGte("age", minAge)
    qb.WhereEq("status", status)

    return qb.Build()
}

// name="", minAge=0, status="active" 时:
// sql:  SELECT * FROM users WHERE status = ?
// args: ["active"]
```

### 时间范围查询

```go
qb := sqlbuilder.NewQueryBuilder("SELECT * FROM orders")

qb, err := qb.WhereTimeRange("created_at", "2024-01-01", "2024-12-31")
if err != nil {
    log.Fatal(err)
}

qb.OrderBy("id", "DESC").Limit(50)

sql, args := qb.Build()
// sql:  SELECT * FROM orders WHERE created_at >= ? AND created_at <= ? ORDER BY id DESC LIMIT ?
// args: [time.Time{2024-01-01}, time.Time{2024-12-31}, 50]
```

### 配合数据库使用

```go
qb := sqlbuilder.NewQueryBuilder("SELECT id, name FROM users").
    WhereEq("status", "active").
    OrderBy("id", "ASC").
    Limit(100)

sql, args := qb.Build()

rows, err := db.Query(sql, args...)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()
```

---

## API 速览

| 方法 | 签名 | 说明 |
|------|------|------|
| `NewQueryBuilder` | `(base string) *QueryBuilder` | 创建构建器 |
| `WhereEq` | `(col string, val any) *QueryBuilder` | 等于 |
| `WhereNotEq` | `(col string, val any) *QueryBuilder` | 不等于 |
| `WhereGt` | `(col string, val any) *QueryBuilder` | 大于 |
| `WhereLt` | `(col string, val any) *QueryBuilder` | 小于 |
| `WhereGte` | `(col string, val any) *QueryBuilder` | 大于等于 |
| `WhereLte` | `(col string, val any) *QueryBuilder` | 小于等于 |
| `WhereLike` | `(col, pattern string) *QueryBuilder` | 模糊匹配（自动加 `%`） |
| `WhereIn` | `[T any](qb, col string, values []T) *QueryBuilder` | IN 查询（泛型函数） |
| `WhereInAny` | `(col string, values []any) *QueryBuilder` | IN 查询（方法） |
| `WhereBetween` | `(col string, start, end any) *QueryBuilder` | 范围查询 |
| `WhereIsNull` | `(col string) *QueryBuilder` | IS NULL |
| `WhereIsNotNull` | `(col string) *QueryBuilder` | IS NOT NULL |
| `WhereTimeRange` | `(col, start, end string) (*QueryBuilder, error)` | 时间范围 |
| `WhereOr` | `(conds []string, args ...any) *QueryBuilder` | OR 条件组 |
| `WhereRaw` | `(cond string, args ...any) *QueryBuilder` | 原始条件 |
| `OrderBy` | `(col, direction string) *QueryBuilder` | 排序 |
| `Limit` | `(limit int) *QueryBuilder` | 分页 |
| `LimitOffset` | `(limit, offset int) *QueryBuilder` | 分页（带偏移） |
| `GroupBy` | `(cols ...string) *QueryBuilder` | 分组 |
| `SetSuffix` | `(suffix string, args ...any) *QueryBuilder` | 自定义后缀 |
| `Build` | `() (string, []any)` | 生成 SQL 和参数 |
