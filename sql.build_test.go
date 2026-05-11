package sqlbuild

import (
	"strings"
	"testing"
	"time"
)

func TestNewQueryBuilder(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users")
	sql, args := qb.Build()
	if sql != "SELECT * FROM users" {
		t.Errorf("expected 'SELECT * FROM users', got '%s'", sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereEq(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereEq("name", "alice")
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE name = ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != "alice" {
		t.Errorf("expected args [alice], got %v", args)
	}
}

func TestWhereEq_EmptyValue(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereEq("name", "")
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereEq_NilValue(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereEq("name", nil)
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereLike(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereLike("name", "ali")
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE name LIKE ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != "%ali%" {
		t.Errorf("expected args [%%ali%%], got %v", args)
	}
}

func TestWhereLike_EmptyPattern(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereLike("name", "")
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereGt(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereGt("age", 18)
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE age > ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != 18 {
		t.Errorf("expected args [18], got %v", args)
	}
}

func TestWhereLt(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereLt("age", 60)
	sql, _ := qb.Build()
	expected := "SELECT * FROM users WHERE age < ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereGte(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereGte("age", 18)
	sql, _ := qb.Build()
	expected := "SELECT * FROM users WHERE age >= ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereLte(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereLte("age", 60)
	sql, _ := qb.Build()
	expected := "SELECT * FROM users WHERE age <= ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereNotEq(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereNotEq("status", "deleted")
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE status != ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != "deleted" {
		t.Errorf("expected args [deleted], got %v", args)
	}
}

func TestWhereIn_Generic(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users")
	WhereIn(qb, "id", []int{1, 2, 3})
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE id IN (?,?,?)"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d", len(args))
	}
}

func TestWhereIn_Empty(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users")
	WhereIn(qb, "id", []int{})
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereInAny(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereInAny("id", []any{1, 2, 3})
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE id IN (?,?,?)"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d", len(args))
	}
}

func TestWhereInAny_Empty(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereInAny("id", []any{})
	sql, _ := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereBetween(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereBetween("age", 18, 60)
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE age BETWEEN ? AND ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 2 || args[0] != 18 || args[1] != 60 {
		t.Errorf("expected args [18, 60], got %v", args)
	}
}

func TestWhereBetween_NilStart(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereBetween("age", nil, 60)
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereBetween_NilEnd(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereBetween("age", 18, nil)
	sql, _ := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereBetween_EmptyStringStart(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereBetween("age", "", 60)
	sql, _ := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereIsNull(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereIsNull("deleted_at")
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE deleted_at IS NULL"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereIsNotNull(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereIsNotNull("email")
	sql, _ := qb.Build()
	expected := "SELECT * FROM users WHERE email IS NOT NULL"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereOr(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereOr(
		[]string{"status = ?", "role = ?"},
		"active", "admin",
	)
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE (status = ? OR role = ?)"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 2 || args[0] != "active" || args[1] != "admin" {
		t.Errorf("expected args [active, admin], got %v", args)
	}
}

func TestWhereOr_Empty(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereOr([]string{})
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereRaw(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereRaw("1 = 1")
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE 1 = 1"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereRaw_WithArgs(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereRaw("status != ?", "banned")
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE status != ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != "banned" {
		t.Errorf("expected args [banned], got %v", args)
	}
}

func TestWhereRaw_Empty(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereRaw("")
	sql, _ := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestWhereTimeRange(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	qb, err := qb.WhereTimeRange("created_at", "2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, args := qb.Build()
	expected := "SELECT * FROM orders WHERE created_at >= ? AND created_at <= ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

func TestWhereTimeRange_OnlyStart(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	qb, err := qb.WhereTimeRange("created_at", "2024-01-01", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, args := qb.Build()
	expected := "SELECT * FROM orders WHERE created_at >= ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(args))
	}
}

func TestWhereTimeRange_OnlyEnd(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	qb, err := qb.WhereTimeRange("created_at", "", "2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, args := qb.Build()
	expected := "SELECT * FROM orders WHERE created_at <= ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(args))
	}
}

func TestWhereTimeRange_BothEmpty(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	qb, err := qb.WhereTimeRange("created_at", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, args := qb.Build()
	expected := "SELECT * FROM orders"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestWhereTimeRange_InvalidStart(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	_, err := qb.WhereTimeRange("created_at", "not-a-date", "2024-12-31")
	if err == nil {
		t.Error("expected error for invalid start time, got nil")
	}
	if !strings.Contains(err.Error(), "解析开始时间失败") {
		t.Errorf("error should contain '解析开始时间失败', got: %v", err)
	}
}

func TestWhereTimeRange_InvalidEnd(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	_, err := qb.WhereTimeRange("created_at", "2024-01-01", "not-a-date")
	if err == nil {
		t.Error("expected error for invalid end time, got nil")
	}
	if !strings.Contains(err.Error(), "解析结束时间失败") {
		t.Errorf("error should contain '解析结束时间失败', got: %v", err)
	}
}

func TestWhereTimeRange_RFC3339(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders")
	qb, err := qb.WhereTimeRange("created_at", "2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, args := qb.Build()
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
	for _, arg := range args {
		if _, ok := arg.(time.Time); !ok {
			t.Errorf("expected time.Time arg, got %T", arg)
		}
	}
	_ = sql
}

func TestOrderBy(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").OrderBy("name", "ASC")
	sql, args := qb.Build()
	expected := "SELECT * FROM users ORDER BY name ASC"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestOrderBy_Desc(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").OrderBy("age", "desc")
	sql, _ := qb.Build()
	expected := "SELECT * FROM users ORDER BY age DESC"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestOrderBy_InvalidDirection(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").OrderBy("name", "invalid")
	sql, _ := qb.Build()
	expected := "SELECT * FROM users ORDER BY name ASC"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestLimit(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").Limit(10)
	sql, args := qb.Build()
	expected := "SELECT * FROM users LIMIT ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != 10 {
		t.Errorf("expected args [10], got %v", args)
	}
}

func TestLimitOffset(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").LimitOffset(10, 20)
	sql, args := qb.Build()
	expected := "SELECT * FROM users LIMIT ? OFFSET ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 2 || args[0] != 10 || args[1] != 20 {
		t.Errorf("expected args [10, 20], got %v", args)
	}
}

func TestGroupBy(t *testing.T) {
	qb := NewQueryBuilder("SELECT department, COUNT(*) FROM users").GroupBy("department")
	sql, args := qb.Build()
	expected := "SELECT department, COUNT(*) FROM users GROUP BY department"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestGroupBy_MultipleCols(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").GroupBy("department", "role")
	sql, _ := qb.Build()
	expected := "SELECT * FROM users GROUP BY department, role"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestGroupBy_Empty(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").GroupBy()
	sql, _ := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
}

func TestSetSuffix_Append(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").
		SetSuffix("ORDER BY name ASC").
		SetSuffix("LIMIT ?", 10)
	sql, args := qb.Build()
	expected := "SELECT * FROM users ORDER BY name ASC LIMIT ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 || args[0] != 10 {
		t.Errorf("expected args [10], got %v", args)
	}
}

func TestChainedConditions(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").
		WhereEq("status", "active").
		WhereGt("age", 18).
		WhereLike("name", "ali").
		OrderBy("name", "ASC").
		Limit(10)

	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE status = ? AND age > ? AND name LIKE ? ORDER BY name ASC LIMIT ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 4 {
		t.Errorf("expected 4 args, got %d: %v", len(args), args)
	}
}

func TestComplexQuery(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM orders").
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
	expected := "SELECT * FROM orders WHERE status = ? AND shipped_at IS NOT NULL AND (priority = ? OR vip = ?) AND amount BETWEEN ? AND ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 7 {
		t.Errorf("expected 7 args, got %d: %v", len(args), args)
	}
}

func TestZeroTimeFiltered(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").WhereEq("created_at", time.Time{})
	sql, args := qb.Build()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestValidTimeNotFiltered(t *testing.T) {
	now := time.Now()
	qb := NewQueryBuilder("SELECT * FROM users").WhereEq("created_at", now)
	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE created_at = ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(args))
	}
}

func TestParseTimeString_Formats(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"2024-01-15T10:30:00Z", true},
		{"2024-01-15 10:30:00", true},
		{"2024-01-15", true},
		{"2024/01/15 10:30:00", true},
		{"2024/01/15", true},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range cases {
		_, err := parseTimeString(tc.input)
		if tc.valid && err != nil {
			t.Errorf("expected '%s' to parse, got error: %v", tc.input, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("expected '%s' to fail parsing, but it succeeded", tc.input)
		}
	}
}

func TestOrderByWithLimit(t *testing.T) {
	qb := NewQueryBuilder("SELECT * FROM users").
		WhereEq("status", "active").
		OrderBy("id", "DESC").
		LimitOffset(10, 5)

	sql, args := qb.Build()
	expected := "SELECT * FROM users WHERE status = ? ORDER BY id DESC LIMIT ? OFFSET ?"
	if sql != expected {
		t.Errorf("expected '%s', got '%s'", expected, sql)
	}
	if len(args) != 3 {
		t.Errorf("expected 3 args, got %d: %v", len(args), args)
	}
	if args[0] != "active" || args[1] != 10 || args[2] != 5 {
		t.Errorf("expected args [active, 10, 5], got %v", args)
	}
}
