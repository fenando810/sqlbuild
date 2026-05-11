package sqlbuilder

import (
	"fmt"
	"strings"
	"time"
)

type QueryBuilder struct {
	baseQuery  string
	conditions []string
	args       []any
	suffix     string
	suffixArgs []any
}

func NewQueryBuilder(base string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:  base,
		conditions: []string{},
		args:       []any{},
	}
}

func (qb *QueryBuilder) AddCondition(cond string, val any) *QueryBuilder {
	if val == nil {
		return qb
	}
	if s, ok := val.(string); ok && s == "" {
		return qb
	}
	if t, ok := val.(time.Time); ok && t.IsZero() {
		return qb
	}

	qb.conditions = append(qb.conditions, cond)
	qb.args = append(qb.args, val)
	return qb
}

func (qb *QueryBuilder) WhereEq(col string, val any) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s = ?", col), val)
}

func (qb *QueryBuilder) WhereLike(col string, pattern string) *QueryBuilder {
	if pattern == "" {
		return qb
	}
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s LIKE ?", col))
	qb.args = append(qb.args, "%"+pattern+"%")
	return qb
}

func (qb *QueryBuilder) WhereGt(col string, val any) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s > ?", col), val)
}

func (qb *QueryBuilder) WhereLt(col string, val any) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s < ?", col), val)
}

func (qb *QueryBuilder) WhereGte(col string, val any) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s >= ?", col), val)
}

func (qb *QueryBuilder) WhereLte(col string, val any) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s <= ?", col), val)
}

func WhereIn[T any](qb *QueryBuilder, col string, values []T) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1]
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IN (%s)", col, placeholders))

	for _, v := range values {
		qb.args = append(qb.args, v)
	}
	return qb
}

func (qb *QueryBuilder) WhereInAny(col string, values []any) *QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	placeholders := strings.Repeat("?,", len(values))
	placeholders = placeholders[:len(placeholders)-1]
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IN (%s)", col, placeholders))
	qb.args = append(qb.args, values...)
	return qb
}

var timeFormats = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006/01/02 15:04:05",
	"2006/01/02",
}

// 东八区时区（北京时间）
var cstZone = time.FixedZone("CST", 8*60*60)

// 自带时区的按时区解析，如果没有时区信息，默认按东八区（北京时间）解析
func parseTimeString(s string) (time.Time, error) {
	// RFC3339 格式包含时区信息，直接解析
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// 其他格式默认按东八区（北京时间）解析
	for _, layout := range timeFormats[1:] {
		if t, err := time.ParseInLocation(layout, s, cstZone); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析时间字符串: %s", s)
}

func (qb *QueryBuilder) WhereTimeRange(col string, start, end string) (*QueryBuilder, error) {
	if start != "" {
		t, err := parseTimeString(start)
		if err != nil {
			return qb, fmt.Errorf("解析开始时间失败: %w", err)
		}
		qb.conditions = append(qb.conditions, fmt.Sprintf("%s >= ?", col))
		qb.args = append(qb.args, t)
	}

	if end != "" {
		t, err := parseTimeString(end)
		if err != nil {
			return qb, fmt.Errorf("解析结束时间失败: %w", err)
		}
		qb.conditions = append(qb.conditions, fmt.Sprintf("%s <= ?", col))
		qb.args = append(qb.args, t)
	}
	return qb, nil
}

func (qb *QueryBuilder) WhereRaw(cond string, args ...any) *QueryBuilder {
	if cond != "" {
		qb.conditions = append(qb.conditions, cond)
		qb.args = append(qb.args, args...)
	}
	return qb
}

func (qb *QueryBuilder) WhereOr(conds []string, args ...any) *QueryBuilder {
	if len(conds) == 0 {
		return qb
	}
	groupCond := "(" + strings.Join(conds, " OR ") + ")"
	qb.conditions = append(qb.conditions, groupCond)
	qb.args = append(qb.args, args...)
	return qb
}

func (qb *QueryBuilder) WhereNotEq(col string, val any) *QueryBuilder {
	return qb.AddCondition(fmt.Sprintf("%s != ?", col), val)
}

func (qb *QueryBuilder) WhereBetween(col string, start, end any) *QueryBuilder {
	if start == nil || end == nil {
		return qb
	}
	if s, ok := start.(string); ok && s == "" {
		return qb
	}
	if e, ok := end.(string); ok && e == "" {
		return qb
	}
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s BETWEEN ? AND ?", col))
	qb.args = append(qb.args, start, end)
	return qb
}

func (qb *QueryBuilder) WhereIsNull(col string) *QueryBuilder {
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IS NULL", col))
	return qb
}

func (qb *QueryBuilder) WhereIsNotNull(col string) *QueryBuilder {
	qb.conditions = append(qb.conditions, fmt.Sprintf("%s IS NOT NULL", col))
	return qb
}

func (qb *QueryBuilder) OrderBy(col string, direction string) *QueryBuilder {
	dir := strings.ToUpper(direction)
	if dir != "ASC" && dir != "DESC" {
		dir = "ASC"
	}
	if qb.suffix != "" {
		qb.suffix += " "
	}
	qb.suffix += fmt.Sprintf("ORDER BY %s %s", col, dir)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	if qb.suffix != "" {
		qb.suffix += " "
	}
	qb.suffix += "LIMIT ?"
	qb.suffixArgs = append(qb.suffixArgs, limit)
	return qb
}

func (qb *QueryBuilder) LimitOffset(limit, offset int) *QueryBuilder {
	if qb.suffix != "" {
		qb.suffix += " "
	}
	qb.suffix += "LIMIT ? OFFSET ?"
	qb.suffixArgs = append(qb.suffixArgs, limit, offset)
	return qb
}

func (qb *QueryBuilder) GroupBy(cols ...string) *QueryBuilder {
	if len(cols) == 0 {
		return qb
	}
	if qb.suffix != "" {
		qb.suffix += " "
	}
	qb.suffix += "GROUP BY " + strings.Join(cols, ", ")
	return qb
}

func (qb *QueryBuilder) SetSuffix(suffix string, args ...any) *QueryBuilder {
	if qb.suffix != "" {
		qb.suffix += " "
	}
	qb.suffix += suffix
	qb.suffixArgs = append(qb.suffixArgs, args...)
	return qb
}

func (qb *QueryBuilder) Build() (string, []any) {
	query := qb.baseQuery
	if len(qb.conditions) > 0 {
		query += " WHERE " + strings.Join(qb.conditions, " AND ")
	}
	if qb.suffix != "" {
		query += " " + qb.suffix
	}

	args := make([]any, 0, len(qb.args)+len(qb.suffixArgs))
	args = append(args, qb.args...)
	args = append(args, qb.suffixArgs...)
	return query, args
}
