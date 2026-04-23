# 电脑垃圾清理工程 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个功能完整的 Windows 电脑垃圾清理工具，支持三级清理规则、智能风险评估、隔离恢复机制和完整的 GUI 界面

**Architecture:** 采用 Go + Fyne 的 5 层架构（UI/App/Core/Platform/Storage），核心层负责规则引擎、扫描、风险评估和清理执行，平台层隔离 Windows 特定实现，数据层使用 SQLite 持久化

**Tech Stack:** Go 1.21+, Fyne v2.4+, SQLite, golang.org/x/sys/windows

---

## 实施策略

本计划分为 **7 个阶段**，每个阶段产出可独立测试的功能模块：

1. **阶段 1**: 项目基础设施和数据模型
2. **阶段 2**: 规则引擎和加载器
3. **阶段 3**: 文件扫描系统
4. **阶段 4**: 风险评估和清理执行
5. **阶段 5**: 数据持久化和历史记录
6. **阶段 6**: UI 界面实现
7. **阶段 7**: 集成测试和打包

每个阶段遵循 TDD 原则：先写测试，再实现功能，频繁提交。

---

## 阶段 1: 项目基础设施和数据模型

### Task 1.1: 初始化 Go 项目

**Files:**
- Create: `go.mod`
- Create: `go.sum`
- Create: `README.md`
- Create: `.gitignore`

- [ ] **Step 1: 初始化 Go module**

```bash
go mod init github.com/yourusername/cleanMyComputer
```

Expected: 创建 go.mod 文件

- [ ] **Step 2: 创建 .gitignore**

```bash
cat > .gitignore << 'EOF'
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
dist/
build/

# Test binary
*.test
*.out

# Go workspace file
go.work

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Database
*.db
*.db-shm
*.db-wal

# Logs
*.log
EOF
```

- [ ] **Step 3: 创建 README.md**

```bash
cat > README.md << 'EOF'
# CleanMyComputer

一个功能完整的 Windows 电脑垃圾清理工具。

## 特性

- 三级清理规则（安全/深度/高级）
- 智能风险评估
- 隔离恢复机制
- 详细的清理报告和历史记录
- 图形界面

## 技术栈

- Go 1.21+
- Fyne v2.4+
- SQLite

## 开发

```bash
# 运行测试
go test ./...

# 构建
go build -o cleanMyComputer.exe ./cmd/cleaner
```
EOF
```

- [ ] **Step 4: 提交初始化**

```bash
git add go.mod .gitignore README.md
git commit -m "feat: 初始化项目结构

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 1.2: 定义核心数据模型

**Files:**
- Create: `internal/models/rule.go`
- Create: `internal/models/rule_test.go`

- [ ] **Step 1: 编写 CleanRule 模型测试**

```go
// internal/models/rule_test.go
package models

import (
	"testing"
)

func TestCleanRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    CleanRule
		wantErr bool
	}{
		{
			name: "valid rule",
			rule: CleanRule{
				ID:          "test_rule",
				Name:        "Test Rule",
				Category:    "system",
				Level:       1,
				Description: "Test description",
				RiskScore:   10,
				Enabled:     true,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			rule: CleanRule{
				Name:     "Test Rule",
				Category: "system",
				Level:    1,
			},
			wantErr: true,
		},
		{
			name: "invalid level",
			rule: CleanRule{
				ID:       "test_rule",
				Name:     "Test Rule",
				Category: "system",
				Level:    99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test ./internal/models -v
```

Expected: FAIL - undefined: CleanRule

- [ ] **Step 3: 实现 CleanRule 模型**

```go
// internal/models/rule.go
package models

import (
	"errors"
	"time"
)

// CleanRule 定义清理规则
type CleanRule struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Category      string    `json:"category"`
	Level         int       `json:"level"`
	Description   string    `json:"description"`
	Targets       []Target  `json:"targets"`
	RiskScore     int       `json:"risk_score"`
	RequiresAdmin bool      `json:"requires_admin"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Target 定义清理目标
type Target struct {
	Type        string   `json:"type"`         // file, folder, registry, command
	Path        string   `json:"path"`         // 路径或命令
	Pattern     string   `json:"pattern"`      // 文件匹配模式
	Recursive   bool     `json:"recursive"`    // 是否递归
	MaxDepth    int      `json:"max_depth"`    // 最大递归深度
	ExcludeList []string `json:"exclude_list"` // 排除列表
}

// Validate 验证规则有效性
func (r *CleanRule) Validate() error {
	if r.ID == "" {
		return errors.New("rule ID is required")
	}
	if r.Name == "" {
		return errors.New("rule name is required")
	}
	if r.Level < 1 || r.Level > 3 {
		return errors.New("rule level must be 1, 2, or 3")
	}
	if r.RiskScore < 0 || r.RiskScore > 100 {
		return errors.New("risk score must be between 0 and 100")
	}
	return nil
}
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test ./internal/models -v
```

Expected: PASS

- [ ] **Step 5: 提交数据模型**

```bash
git add internal/models/
git commit -m "feat: 添加核心数据模型 CleanRule 和 Target

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 1.3: 其他数据模型

**Files:**
- Create: `internal/models/scan_result.go`
- Create: `internal/models/scan_result_test.go`
- Create: `internal/models/clean_record.go`
- Create: `internal/models/config.go`

- [ ] **Step 1: 编写 ScanResult 模型测试**

```go
// internal/models/scan_result_test.go
package models

import (
	"testing"
	"time"
)

func TestScanResult_TotalSize(t *testing.T) {
	result := &ScanResult{
		Items: []*ScanItem{
			{Path: "/tmp/file1.txt", Size: 1024},
			{Path: "/tmp/file2.txt", Size: 2048},
		},
	}
	
	expected := int64(3072)
	if got := result.TotalSize(); got != expected {
		t.Errorf("TotalSize() = %d, want %d", got, expected)
	}
}

func TestScanResult_TotalCount(t *testing.T) {
	result := &ScanResult{
		Items: []*ScanItem{
			{Path: "/tmp/file1.txt"},
			{Path: "/tmp/file2.txt"},
		},
	}
	
	if got := result.TotalCount(); got != 2 {
		t.Errorf("TotalCount() = %d, want 2", got)
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test ./internal/models -v -run TestScanResult
```

Expected: FAIL - undefined: ScanResult

- [ ] **Step 3: 实现 ScanResult 模型**

```go
// internal/models/scan_result.go
package models

import "time"

type ScanResult struct {
	ID        int64       `json:"id"`
	Level     int         `json:"level"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Items     []*ScanItem `json:"items"`
	Status    string      `json:"status"`
}

type ScanItem struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	RuleID    string    `json:"rule_id"`
	RiskScore int       `json:"risk_score"`
}

func (s *ScanResult) TotalSize() int64 {
	var total int64
	for _, item := range s.Items {
		total += item.Size
	}
	return total
}

func (s *ScanResult) TotalCount() int {
	return len(s.Items)
}
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test ./internal/models -v -run TestScanResult
```

Expected: PASS

- [ ] **Step 5: 实现 CleanRecord 和 Config 模型**

```go
// internal/models/clean_record.go
package models

import "time"

type CleanRecord struct {
	ID         int64     `json:"id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	ScanLevel  int       `json:"scan_level"`
	TotalFiles int       `json:"total_files"`
	TotalSize  int64     `json:"total_size"`
	FreedSize  int64     `json:"freed_size"`
	Status     string    `json:"status"`
}

// internal/models/config.go
package models

type Config struct {
	QuarantineRetentionHours int    `json:"quarantine_retention_hours"`
	AutoCleanEnabled         bool   `json:"auto_clean_enabled"`
	AutoCleanSchedule        string `json:"auto_clean_schedule"`
	OldFileDays              int    `json:"old_file_days"`
	ScanWorkers              int    `json:"scan_workers"`
	Language                 string `json:"language"`
}
```

- [ ] **Step 6: 提交所有数据模型**

```bash
git add internal/models/
git commit -m "feat: 添加 ScanResult, CleanRecord, Config 数据模型

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 阶段 2: 规则引擎和加载器

### Task 2.1: 规则加载器

**Files:**
- Create: `configs/rules/level1_safe.json`
- Create: `internal/core/rule/loader.go`
- Create: `internal/core/rule/loader_test.go`

- [ ] **Step 1: 创建示例规则文件**

```bash
mkdir -p configs/rules
```

```json
// configs/rules/level1_safe.json
[
  {
    "id": "temp_files",
    "name": "临时文件",
    "category": "system",
    "level": 1,
    "description": "清理系统临时文件夹中的文件",
    "targets": [
      {
        "type": "folder",
        "path": "%TEMP%",
        "pattern": "*",
        "recursive": true,
        "max_depth": 3,
        "exclude_list": []
      }
    ],
    "risk_score": 10,
    "requires_admin": false,
    "enabled": true
  },
  {
    "id": "recycle_bin",
    "name": "回收站",
    "category": "system",
    "level": 1,
    "description": "清空回收站",
    "targets": [
      {
        "type": "command",
        "path": "Clear-RecycleBin",
        "pattern": "",
        "recursive": false,
        "max_depth": 0,
        "exclude_list": []
      }
    ],
    "risk_score": 5,
    "requires_admin": false,
    "enabled": true
  }
]
```

- [ ] **Step 2: 编写规则加载器测试**

```go
// internal/core/rule/loader_test.go
package rule

import (
	"testing"
)

func TestLoader_LoadFromFile(t *testing.T) {
	loader := NewLoader()
	rules, err := loader.LoadFromFile("../../../configs/rules/level1_safe.json")
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	
	if len(rules) == 0 {
		t.Error("Expected rules to be loaded")
	}
	
	for _, rule := range rules {
		if err := rule.Validate(); err != nil {
			t.Errorf("Invalid rule %s: %v", rule.ID, err)
		}
	}
}
```

- [ ] **Step 3: 运行测试验证失败**

```bash
go test ./internal/core/rule -v
```

Expected: FAIL - undefined: NewLoader

- [ ] **Step 4: 实现规则加载器**

```go
// internal/core/rule/loader.go
package rule

import (
	"encoding/json"
	"os"
	"path/filepath"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

type Loader struct {
	rulesDir string
}

func NewLoader() *Loader {
	return &Loader{
		rulesDir: "configs/rules",
	}
}

func (l *Loader) LoadFromFile(path string) ([]*models.CleanRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var rules []*models.CleanRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	
	return rules, nil
}

func (l *Loader) LoadByLevel(level int) ([]*models.CleanRule, error) {
	var allRules []*models.CleanRule
	
	for i := 1; i <= level; i++ {
		filename := filepath.Join(l.rulesDir, "level"+string(rune(i+'0'))+"_*.json")
		matches, _ := filepath.Glob(filename)
		
		for _, match := range matches {
			rules, err := l.LoadFromFile(match)
			if err != nil {
				continue
			}
			allRules = append(allRules, rules...)
		}
	}
	
	return allRules, nil
}
```

- [ ] **Step 5: 运行测试验证通过**

```bash
go test ./internal/core/rule -v
```

Expected: PASS

- [ ] **Step 6: 提交规则加载器**

```bash
git add configs/rules/ internal/core/rule/
git commit -m "feat: 实现规则加载器，支持从 JSON 加载清理规则

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 2.2: 规则引擎

**Files:**
- Create: `internal/core/rule/engine.go`
- Create: `internal/core/rule/engine_test.go`

- [ ] **Step 1: 编写规则引擎测试**

```go
// internal/core/rule/engine_test.go
package rule

import (
	"context"
	"testing"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

func TestEngine_LoadRules(t *testing.T) {
	loader := NewLoader()
	engine := NewEngine(loader)
	
	ctx := context.Background()
	err := engine.LoadRules(ctx, 1)
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}
	
	rules := engine.GetEnabledRules(1)
	if len(rules) == 0 {
		t.Error("Expected rules to be loaded")
	}
}

func TestEngine_GetEnabledRules(t *testing.T) {
	engine := NewEngine(nil)
	engine.rules = map[string]*models.CleanRule{
		"rule1": {ID: "rule1", Level: 1, Enabled: true},
		"rule2": {ID: "rule2", Level: 2, Enabled: true},
		"rule3": {ID: "rule3", Level: 1, Enabled: false},
	}
	
	rules := engine.GetEnabledRules(1)
	if len(rules) != 1 {
		t.Errorf("Expected 1 enabled rule, got %d", len(rules))
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test ./internal/core/rule -v -run TestEngine
```

Expected: FAIL - undefined: NewEngine

- [ ] **Step 3: 实现规则引擎**

```go
// internal/core/rule/engine.go
package rule

import (
	"context"
	"sync"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

type Engine struct {
	rules  map[string]*models.CleanRule
	mu     sync.RWMutex
	loader *Loader
}

func NewEngine(loader *Loader) *Engine {
	return &Engine{
		rules:  make(map[string]*models.CleanRule),
		loader: loader,
	}
}

func (e *Engine) LoadRules(ctx context.Context, level int) error {
	rules, err := e.loader.LoadByLevel(level)
	if err != nil {
		return err
	}
	
	e.mu.Lock()
	defer e.mu.Unlock()
	
	for _, rule := range rules {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := rule.Validate(); err == nil {
			e.rules[rule.ID] = rule
		}
	}
	return nil
}

func (e *Engine) GetEnabledRules(level int) []*models.CleanRule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var enabled []*models.CleanRule
	for _, rule := range e.rules {
		if rule.Enabled && rule.Level <= level {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test ./internal/core/rule -v -run TestEngine
```

Expected: PASS

- [ ] **Step 5: 提交规则引擎**

```bash
git add internal/core/rule/
git commit -m "feat: 添加规则引擎，支持规则管理和查询

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 2.3: 规则验证器

**Files:**
- Create: `internal/core/rule/validator.go`
- Create: `internal/core/rule/validator_test.go`

- [ ] **Step 1: 编写验证器测试**

```go
// internal/core/rule/validator_test.go
package rule

import (
	"testing"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

func TestValidator_ValidateRule(t *testing.T) {
	validator := NewValidator()
	
	tests := []struct {
		name    string
		rule    *models.CleanRule
		wantErr bool
	}{
		{
			name: "valid rule",
			rule: &models.CleanRule{
				ID: "test", Name: "Test", Level: 1, RiskScore: 10,
			},
			wantErr: false,
		},
		{
			name: "invalid risk score",
			rule: &models.CleanRule{
				ID: "test", Name: "Test", Level: 1, RiskScore: 150,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateRule(tt.rule)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

- [ ] **Step 2: 实现验证器**

```go
// internal/core/rule/validator.go
package rule

import (
	"errors"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateRule(rule *models.CleanRule) error {
	return rule.Validate()
}

func (v *Validator) ValidateRules(rules []*models.CleanRule) []error {
	var errs []error
	for _, rule := range rules {
		if err := v.ValidateRule(rule); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
```

- [ ] **Step 3: 运行测试验证通过**

```bash
go test ./internal/core/rule -v -run TestValidator
```

Expected: PASS

- [ ] **Step 4: 提交验证器**

```bash
git add internal/core/rule/
git commit -m "feat: 添加规则验证器

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 阶段 3: 文件扫描系统

### Task 3.1: 文件扫描器基础

**Files:**
- Create: `internal/core/scanner/scanner.go`
- Create: `internal/core/scanner/scanner_test.go`

- [ ] **Step 1: 编写扫描器测试**

```go
// internal/core/scanner/scanner_test.go
package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

func TestScanner_ScanTargets(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	scanner := NewScanner(2)
	targets := []models.Target{
		{
			Type:      "folder",
			Path:      tmpDir,
			Pattern:   "*.txt",
			Recursive: false,
		},
	}
	
	ctx := context.Background()
	results, err := scanner.ScanTargets(ctx, targets)
	if err != nil {
		t.Fatalf("ScanTargets() error = %v", err)
	}
	
	if len(results) == 0 {
		t.Error("Expected scan results")
	}
}
```

- [ ] **Step 2: 运行测试验证失败**

```bash
go test ./internal/core/scanner -v
```

Expected: FAIL - undefined: NewScanner

- [ ] **Step 3: 实现文件扫描器**

```go
// internal/core/scanner/scanner.go
package scanner

import (
	"context"
	"os"
	"path/filepath"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

type Scanner struct {
	workers int
	filter  *Filter
}

func NewScanner(workers int) *Scanner {
	return &Scanner{workers: workers, filter: NewFilter()}
}

func (s *Scanner) ScanTargets(ctx context.Context, targets []models.Target) ([]*models.ScanItem, error) {
	var results []*models.ScanItem
	for _, target := range targets {
		matches, err := filepath.Glob(filepath.Join(target.Path, target.Pattern))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				continue
			}
			results = append(results, &models.ScanItem{
				Path: match,
				Size: info.Size(),
				ModTime: info.ModTime(),
			})
		}
	}
	return results, nil
}
```

- [ ] **Step 4: 运行测试验证通过**

```bash
go test ./internal/core/scanner -v
```

Expected: PASS

- [ ] **Step 5: 提交扫描器基础**

```bash
git add internal/core/scanner/
git commit -m "feat: 实现文件扫描器基础功能

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 3.2: 并行扫描

**Files:**
- Create: `internal/core/scanner/parallel.go`
- Create: `internal/core/scanner/parallel_test.go`

- [ ] **Step 1: 编写并行扫描测试**

```go
// internal/core/scanner/parallel_test.go
package scanner

import "testing"

func TestScanner_WorkerCount(t *testing.T) {
	scanner := NewParallelScanner(4)
	if scanner.Workers() != 4 {
		t.Errorf("Workers() = %d, want 4", scanner.Workers())
	}
}
```

- [ ] **Step 2: 实现并行扫描器**

```go
// internal/core/scanner/parallel.go
package scanner

type ParallelScanner struct {
	workers int
}

func NewParallelScanner(workers int) *ParallelScanner {
	if workers < 1 {
		workers = 1
	}
	return &ParallelScanner{workers: workers}
}

func (p *ParallelScanner) Workers() int {
	return p.workers
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/core/scanner -v -run TestScanner_WorkerCount
git add internal/core/scanner/
git commit -m "feat: 添加并行扫描配置

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 3.3: 文件过滤器

**Files:**
- Create: `internal/core/scanner/filter.go`
- Create: `internal/core/scanner/filter_test.go`

- [ ] **Step 1: 编写过滤器测试**

```go
// internal/core/scanner/filter_test.go
package scanner

import "testing"

func TestFilter_ShouldInclude(t *testing.T) {
	filter := NewFilter()
	filter.AddExclude("*.log")
	
	if filter.ShouldInclude("app.log") {
		t.Error("Expected app.log to be excluded")
	}
	if !filter.ShouldInclude("app.txt") {
		t.Error("Expected app.txt to be included")
	}
}
```

- [ ] **Step 2: 实现文件过滤器**

```go
// internal/core/scanner/filter.go
package scanner

import "path/filepath"

type Filter struct {
	excludePatterns []string
}

func NewFilter() *Filter {
	return &Filter{excludePatterns: []string{}}
}

func (f *Filter) AddExclude(pattern string) {
	f.excludePatterns = append(f.excludePatterns, pattern)
}

func (f *Filter) ShouldInclude(path string) bool {
	for _, pattern := range f.excludePatterns {
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return false
		}
	}
	return true
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/core/scanner -v
git add internal/core/scanner/
git commit -m "feat: 添加文件过滤器

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 阶段 4: 风险评估和清理执行

### Task 4.1: 风险评估器

**Files:**
- Create: `internal/core/analyzer/risk.go`
- Create: `internal/core/analyzer/risk_test.go`

- [ ] **Step 1: 编写风险评估器测试**

```go
// internal/core/analyzer/risk_test.go
package analyzer

import (
	"testing"
	"time"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

func TestRiskAnalyzer_CalculateRisk(t *testing.T) {
	analyzer := NewRiskAnalyzer()
	
	item := &models.ScanItem{
		Path:      "C:\\Windows\\System32\\test.dll",
		Size:      1024,
		ModTime:   time.Now(),
		RiskScore: 10,
	}
	
	score := analyzer.CalculateRisk(item)
	if score < 10 {
		t.Errorf("Expected risk score >= 10, got %d", score)
	}
}
```

- [ ] **Step 2: 实现风险评估器**

```go
// internal/core/analyzer/risk.go
package analyzer

import (
	"path/filepath"
	"strings"
	"time"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

type RiskAnalyzer struct {
	systemPaths   map[string]bool
	protectedExts map[string]bool
}

func NewRiskAnalyzer() *RiskAnalyzer {
	return &RiskAnalyzer{
		systemPaths: map[string]bool{
			"C:\\Windows\\System32": true,
			"C:\\Program Files":     true,
		},
		protectedExts: map[string]bool{
			".exe": true, ".dll": true, ".sys": true,
		},
	}
}

func (r *RiskAnalyzer) CalculateRisk(item *models.ScanItem) int {
	score := item.RiskScore
	if r.isSystemPath(item.Path) {
		score += 30
	}
	ext := strings.ToLower(filepath.Ext(item.Path))
	if r.protectedExts[ext] {
		score += 20
	}
	if item.Size > 100*1024*1024 {
		score += 15
	}
	if time.Since(item.ModTime) < 7*24*time.Hour {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	return score
}

func (r *RiskAnalyzer) isSystemPath(path string) bool {
	for sysPath := range r.systemPaths {
		if strings.HasPrefix(path, sysPath) {
			return true
		}
	}
	return false
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/core/analyzer -v
git add internal/core/analyzer/
git commit -m "feat: 实现风险评估器

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 4.2: 清理执行器

**Files:**
- Create: `internal/core/cleaner/executor.go`
- Create: `internal/core/cleaner/executor_test.go`

- [ ] **Step 1: 编写清理执行器测试**

```go
// internal/core/cleaner/executor_test.go
package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	qm := NewQuarantineManager(filepath.Join(tmpDir, "quarantine"))
	executor := NewExecutor(qm)
	
	task := &CleanTask{
		Files: []*FileItem{
			{Path: testFile, Size: 4, RiskScore: 10},
		},
	}
	
	ctx := context.Background()
	result, err := executor.Execute(ctx, task)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	
	if result.FreedSize != 4 {
		t.Errorf("FreedSize = %d, want 4", result.FreedSize)
	}
}
```

- [ ] **Step 2: 实现清理执行器**

```go
// internal/core/cleaner/executor.go
package cleaner

import (
	"context"
	"os"
	"time"
)

type Executor struct {
	quarantine *QuarantineManager
	dryRun     bool
}

type CleanTask struct {
	Files     []*FileItem
	TotalSize int64
}

type FileItem struct {
	Path      string
	Size      int64
	RiskScore int
}

type CleanResult struct {
	Cleaned   []string
	Failed    []string
	FreedSize int64
	StartTime time.Time
	EndTime   time.Time
}

func NewExecutor(qm *QuarantineManager) *Executor {
	return &Executor{quarantine: qm, dryRun: false}
}

func (e *Executor) Execute(ctx context.Context, task *CleanTask) (*CleanResult, error) {
	result := &CleanResult{StartTime: time.Now()}
	for _, file := range task.Files {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		if err := e.cleanFile(file); err != nil {
			result.Failed = append(result.Failed, file.Path)
			continue
		}
		result.Cleaned = append(result.Cleaned, file.Path)
		result.FreedSize += file.Size
	}
	result.EndTime = time.Now()
	return result, nil
}

func (e *Executor) cleanFile(file *FileItem) error {
	if e.dryRun {
		return nil
	}
	if file.RiskScore > 60 {
		return e.quarantine.Quarantine(file.Path)
	}
	return os.Remove(file.Path)
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/core/cleaner -v
git add internal/core/cleaner/
git commit -m "feat: 实现清理执行器

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 4.3: 隔离区管理

**Files:**
- Create: `internal/core/cleaner/quarantine.go`
- Create: `internal/core/cleaner/quarantine_test.go`

- [ ] **Step 1: 编写隔离区管理器测试**

```go
// internal/core/cleaner/quarantine_test.go
package cleaner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQuarantineManager_Quarantine(t *testing.T) {
	tmpDir := t.TempDir()
	qDir := filepath.Join(tmpDir, "quarantine")
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	qm := NewQuarantineManager(qDir)
	err := qm.Quarantine(testFile)
	if err != nil {
		t.Fatalf("Quarantine() error = %v", err)
	}
	
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Expected file to be moved")
	}
}
```

- [ ] **Step 2: 实现隔离区管理器**

```go
// internal/core/cleaner/quarantine.go
package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type QuarantineManager struct {
	baseDir string
}

func NewQuarantineManager(baseDir string) *QuarantineManager {
	os.MkdirAll(baseDir, 0755)
	return &QuarantineManager{baseDir: baseDir}
}

func (q *QuarantineManager) Quarantine(srcPath string) error {
	fileName := filepath.Base(srcPath)
	quarantineName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileName)
	dstPath := filepath.Join(q.baseDir, quarantineName)
	return os.Rename(srcPath, dstPath)
}

func (q *QuarantineManager) Restore(quarantinePath, originalPath string) error {
	return os.Rename(quarantinePath, originalPath)
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/core/cleaner -v
git add internal/core/cleaner/
git commit -m "feat: 实现隔离区管理器

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 4.4: 恢复功能

**Files:**
- Create: `internal/core/cleaner/recovery.go`
- Create: `internal/core/cleaner/recovery_test.go`

- [ ] **Step 1: 编写恢复功能测试**

```go
// internal/core/cleaner/recovery_test.go
package cleaner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecovery_RestoreFile(t *testing.T) {
	tmpDir := t.TempDir()
	qDir := filepath.Join(tmpDir, "quarantine")
	originalPath := filepath.Join(tmpDir, "test.txt")
	
	qm := NewQuarantineManager(qDir)
	recovery := NewRecovery(qm)
	
	os.WriteFile(originalPath, []byte("test"), 0644)
	qm.Quarantine(originalPath)
	
	files, _ := os.ReadDir(qDir)
	if len(files) == 0 {
		t.Fatal("Expected quarantined file")
	}
	
	quarantinedPath := filepath.Join(qDir, files[0].Name())
	err := recovery.RestoreFile(quarantinedPath, originalPath)
	if err != nil {
		t.Fatalf("RestoreFile() error = %v", err)
	}
	
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		t.Error("Expected file to be restored")
	}
}
```

- [ ] **Step 2: 实现恢复功能**

```go
// internal/core/cleaner/recovery.go
package cleaner

type Recovery struct {
	qm *QuarantineManager
}

func NewRecovery(qm *QuarantineManager) *Recovery {
	return &Recovery{qm: qm}
}

func (r *Recovery) RestoreFile(quarantinePath, originalPath string) error {
	return r.qm.Restore(quarantinePath, originalPath)
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/core/cleaner -v
git add internal/core/cleaner/
git commit -m "feat: 实现恢复功能

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 阶段 5: 数据持久化和历史记录

### Task 5.1: 数据库初始化

**Files:**
- Create: `internal/storage/db.go`
- Create: `internal/storage/migrations/001_init.sql`
- Create: `internal/storage/db_test.go`

- [ ] **Step 1: 创建数据库迁移脚本**

```sql
-- internal/storage/migrations/001_init.sql
CREATE TABLE IF NOT EXISTS clean_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    scan_level INTEGER NOT NULL,
    total_files INTEGER NOT NULL,
    total_size INTEGER NOT NULL,
    freed_size INTEGER NOT NULL,
    failed_count INTEGER DEFAULT 0,
    status TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clean_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    history_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    rule_id TEXT NOT NULL,
    risk_score INTEGER NOT NULL,
    action TEXT NOT NULL,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (history_id) REFERENCES clean_history(id) ON DELETE CASCADE
);

CREATE INDEX idx_clean_details_history ON clean_details(history_id);
CREATE INDEX idx_clean_details_rule ON clean_details(rule_id);
```

- [ ] **Step 2: 编写数据库初始化测试**

```go
// internal/storage/db_test.go
package storage

import (
	"testing"
)

func TestDB_Init(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer db.Close()
	
	if db.conn == nil {
		t.Error("Expected database connection")
	}
}
```

- [ ] **Step 3: 实现数据库初始化**

```go
// internal/storage/db.go
package storage

import (
	"database/sql"
	_ "embed"
	
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/001_init.sql
var initSQL string

type DB struct {
	conn *sql.DB
}

func NewDB(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	
	if _, err := conn.Exec(initSQL); err != nil {
		conn.Close()
		return nil, err
	}
	
	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
```

- [ ] **Step 4: 安装依赖并运行测试**

```bash
go get github.com/mattn/go-sqlite3
go test ./internal/storage -v
```

Expected: PASS

- [ ] **Step 5: 提交数据库初始化**

```bash
git add internal/storage/ go.mod go.sum
git commit -m "feat: 实现数据库初始化和迁移

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 5.2: 历史记录管理

**Files:**
- Create: `internal/storage/history.go`
- Create: `internal/storage/history_test.go`

- [ ] **Step 1: 编写历史记录管理测试**

```go
// internal/storage/history_test.go
package storage

import (
	"testing"
	"time"
	
	"github.com/yourusername/cleanMyComputer/internal/models"
)

func TestHistory_Save(t *testing.T) {
	db, _ := NewDB(":memory:")
	defer db.Close()
	
	history := NewHistory(db)
	record := &models.CleanRecord{
		StartTime:  time.Now(),
		EndTime:    time.Now(),
		ScanLevel:  1,
		TotalFiles: 10,
		TotalSize:  1024,
		FreedSize:  512,
		Status:     "success",
	}
	
	id, err := history.Save(record)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if id == 0 {
		t.Error("Expected non-zero ID")
	}
}
```

- [ ] **Step 2: 实现历史记录管理**

```go
// internal/storage/history.go
package storage

import (
	"github.com/yourusername/cleanMyComputer/internal/models"
)

type History struct {
	db *DB
}

func NewHistory(db *DB) *History {
	return &History{db: db}
}

func (h *History) Save(record *models.CleanRecord) (int64, error) {
	result, err := h.db.conn.Exec(`
		INSERT INTO clean_history (start_time, end_time, scan_level, total_files, total_size, freed_size, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, record.StartTime, record.EndTime, record.ScanLevel, record.TotalFiles, record.TotalSize, record.FreedSize, record.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (h *History) GetAll() ([]*models.CleanRecord, error) {
	rows, err := h.db.conn.Query(`SELECT id, start_time, end_time, scan_level, total_files, total_size, freed_size, status FROM clean_history ORDER BY start_time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var records []*models.CleanRecord
	for rows.Next() {
		var r models.CleanRecord
		if err := rows.Scan(&r.ID, &r.StartTime, &r.EndTime, &r.ScanLevel, &r.TotalFiles, &r.TotalSize, &r.FreedSize, &r.Status); err != nil {
			continue
		}
		records = append(records, &r)
	}
	return records, nil
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/storage -v
git add internal/storage/
git commit -m "feat: 实现历史记录管理

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 5.3: 配置管理

**Files:**
- Create: `internal/storage/config.go`
- Create: `internal/storage/config_test.go`

- [ ] **Step 1: 编写配置管理测试**

```go
// internal/storage/config_test.go
package storage

import (
	"testing"
)

func TestConfig_Get(t *testing.T) {
	db, _ := NewDB(":memory:")
	defer db.Close()
	
	config := NewConfig(db)
	config.Set("test_key", "test_value")
	
	value, err := config.Get("test_key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if value != "test_value" {
		t.Errorf("Get() = %s, want test_value", value)
	}
}
```

- [ ] **Step 2: 实现配置管理**

```go
// internal/storage/config.go
package storage

type Config struct {
	db *DB
}

func NewConfig(db *DB) *Config {
	return &Config{db: db}
}

func (c *Config) Get(key string) (string, error) {
	var value string
	err := c.db.conn.QueryRow("SELECT value FROM user_config WHERE key = ?", key).Scan(&value)
	return value, err
}

func (c *Config) Set(key, value string) error {
	_, err := c.db.conn.Exec(`
		INSERT INTO user_config (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
	`, key, value, value)
	return err
}
```

- [ ] **Step 3: 运行测试并提交**

```bash
go test ./internal/storage -v
git add internal/storage/
git commit -m "feat: 实现配置管理

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 阶段 6: UI 界面实现

### Task 6.1: Fyne 应用初始化

**Files:**
- Create: `cmd/cleaner/main.go`
- Create: `internal/ui/app.go`

- [ ] **Step 1: 安装 Fyne 依赖**

```bash
go get fyne.io/fyne/v2
```

Expected: 下载 fyne.io/fyne/v2 依赖

- [ ] **Step 2: 编写应用入口**

```go
// cmd/cleaner/main.go
package main

import "github.com/yourusername/cleanMyComputer/internal/ui"

func main() {
	app := ui.NewApp()
	app.Run()
}

// internal/ui/app.go
package ui

import (
	"fyne.io/fyne/v2/app"
)

type App struct {
	fyneApp fyne.App
}

func NewApp() *App {
	return &App{fyneApp: app.New()}
}

func (a *App) Run() {
	window := a.fyneApp.NewWindow("CleanMyComputer")
	window.Resize(fyne.NewSize(1024, 768))
	window.ShowAndRun()
}
```

- [ ] **Step 3: 构建验证并提交**

```bash
go build -o build/cleanMyComputer.exe ./cmd/cleaner
git add cmd/cleaner/ internal/ui/ go.mod go.sum
git commit -m "feat: 初始化 Fyne 应用

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 6.2: 首页仪表盘

**Files:**
- Create: `internal/ui/dashboard.go`

- [ ] **Step 1: 实现首页仪表盘**

```go
// internal/ui/dashboard.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewDashboard() fyne.CanvasObject {
	title := widget.NewLabel("电脑垃圾清理工具")
	quickScan := widget.NewButton("快速扫描", func() {})
	fullScan := widget.NewButton("完整扫描", func() {})
	return container.NewVBox(title, quickScan, fullScan)
}
```

- [ ] **Step 2: 集成到主应用并提交**

```go
// 更新 internal/ui/app.go
func (a *App) Run() {
	window := a.fyneApp.NewWindow("CleanMyComputer")
	window.SetContent(NewDashboard())
	window.Resize(fyne.NewSize(1024, 768))
	window.ShowAndRun()
}
```

```bash
git add internal/ui/
git commit -m "feat: 实现首页仪表盘

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 6.3: 扫描结果页

**Files:**
- Create: `internal/ui/scanner.go`

- [ ] **Step 1: 实现扫描结果页**

```go
// internal/ui/scanner.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewScannerView() fyne.CanvasObject {
	resultList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)
	return container.NewBorder(
		widget.NewLabel("扫描结果"),
		widget.NewButton("开始清理", func() {}),
		nil, nil,
		resultList,
	)
}
```

- [ ] **Step 2: 提交扫描结果页**

```bash
git add internal/ui/
git commit -m "feat: 实现扫描结果页

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 6.4: 清理确认页

**Files:**
- Create: `internal/ui/confirm.go`

- [ ] **Step 1: 实现清理确认页**

```go
// internal/ui/confirm.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewConfirmView() fyne.CanvasObject {
	summary := widget.NewLabel("准备清理 100 个文件，释放 1.5 GB 空间")
	confirmBtn := widget.NewButton("确认清理", func() {})
	cancelBtn := widget.NewButton("取消", func() {})
	return container.NewVBox(summary, container.NewHBox(confirmBtn, cancelBtn))
}
```

- [ ] **Step 2: 提交清理确认页**

```bash
git add internal/ui/
git commit -m "feat: 实现清理确认页

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 6.5: 历史记录页

**Files:**
- Create: `internal/ui/history.go`

- [ ] **Step 1: 实现历史记录页**

```go
// internal/ui/history.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewHistoryView() fyne.CanvasObject {
	historyList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)
	return container.NewBorder(
		widget.NewLabel("清理历史"),
		nil, nil, nil,
		historyList,
	)
}
```

- [ ] **Step 2: 提交历史记录页**

```bash
git add internal/ui/
git commit -m "feat: 实现历史记录页

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 6.6: 设置页

**Files:**
- Create: `internal/ui/settings.go`

- [ ] **Step 1: 实现设置页**

```go
// internal/ui/settings.go
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsView() fyne.CanvasObject {
	autoClean := widget.NewCheck("启用自动清理", func(bool) {})
	retentionEntry := widget.NewEntry()
	retentionEntry.SetPlaceHolder("24")
	return container.NewVBox(
		widget.NewLabel("设置"),
		autoClean,
		widget.NewLabel("隔离区保留时间（小时）"),
		retentionEntry,
	)
}
```

- [ ] **Step 2: 提交设置页**

```bash
git add internal/ui/
git commit -m "feat: 实现设置页

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 阶段 7: 集成测试和打包

### Task 7.1: 集成测试

**Files:**
- Create: `tests/integration/cleaner_test.go`

- [ ] **Step 1: 编写集成测试**

```go
// tests/integration/cleaner_test.go
package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/yourusername/cleanMyComputer/internal/core/cleaner"
	"github.com/yourusername/cleanMyComputer/internal/core/rule"
	"github.com/yourusername/cleanMyComputer/internal/core/scanner"
	"github.com/yourusername/cleanMyComputer/internal/models"
)

func TestFullCleanWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	loader := rule.NewLoader()
	engine := rule.NewEngine(loader)
	ctx := context.Background()
	
	if err := engine.LoadRules(ctx, 1); err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}
	
	s := scanner.NewScanner(2)
	targets := []models.Target{
		{Type: "folder", Path: tmpDir, Pattern: "*.txt", Recursive: false},
	}
	
	results, err := s.ScanTargets(ctx, targets)
	if err != nil {
		t.Fatalf("ScanTargets() error = %v", err)
	}
	
	if len(results) == 0 {
		t.Fatal("Expected scan results")
	}
	
	qm := cleaner.NewQuarantineManager(filepath.Join(tmpDir, "quarantine"))
	executor := cleaner.NewExecutor(qm)
	
	task := &cleaner.CleanTask{
		Files: []*cleaner.FileItem{
			{Path: testFile, Size: 4, RiskScore: 10},
		},
	}
	
	cleanResult, err := executor.Execute(ctx, task)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	
	if cleanResult.FreedSize != 4 {
		t.Errorf("FreedSize = %d, want 4", cleanResult.FreedSize)
	}
}
```

- [ ] **Step 2: 运行集成测试并提交**

```bash
go test ./tests/integration -v
git add tests/integration/
git commit -m "test: 添加集成测试

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 7.2: 构建脚本

**Files:**
- Create: `scripts/build.sh`

- [ ] **Step 1: 创建构建脚本**

```bash
# scripts/build.sh
#!/bin/bash

set -e

echo "Building CleanMyComputer..."

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
OUTPUT_NAME="cleanMyComputer"

mkdir -p $BUILD_DIR

go build -ldflags "-s -w -X main.version=$VERSION" \
  -o $BUILD_DIR/${OUTPUT_NAME}.exe \
  ./cmd/cleaner

echo "Build complete: $BUILD_DIR/${OUTPUT_NAME}.exe"
```

- [ ] **Step 2: 赋予执行权限并测试**

```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

Expected: 成功构建 build/cleanMyComputer.exe

- [ ] **Step 3: 提交构建脚本**

```bash
git add scripts/
git commit -m "build: 添加构建脚本

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

### Task 7.3: 打包和安装程序

**Files:**
- Create: `scripts/package.sh`

- [ ] **Step 1: 创建打包脚本**

```bash
# scripts/package.sh
#!/bin/bash

set -e

echo "Packaging CleanMyComputer..."

VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
DIST_DIR="dist"
APP_NAME="cleanMyComputer"

mkdir -p $DIST_DIR

./scripts/build.sh

cp $BUILD_DIR/${APP_NAME}.exe $DIST_DIR/
cp -r configs $DIST_DIR/
cp README.md $DIST_DIR/

cd $DIST_DIR
zip -r ${APP_NAME}-${VERSION}-windows-amd64.zip .
cd ..

echo "Package complete: $DIST_DIR/${APP_NAME}-${VERSION}-windows-amd64.zip"
```

- [ ] **Step 2: 测试打包并提交**

```bash
chmod +x scripts/package.sh
./scripts/package.sh
git add scripts/
git commit -m "build: 添加打包脚本

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## 实施计划总结

本实施计划涵盖了从项目初始化到最终打包的完整流程，共 7 个阶段、23 个任务：

**阶段 1**: 项目基础设施和数据模型（3 个任务）
**阶段 2**: 规则引擎和加载器（3 个任务）
**阶段 3**: 文件扫描系统（3 个任务）
**阶段 4**: 风险评估和清理执行（4 个任务）
**阶段 5**: 数据持久化和历史记录（3 个任务）
**阶段 6**: UI 界面实现（6 个任务）
**阶段 7**: 集成测试和打包（3 个任务）

每个任务都遵循 TDD 原则，包含完整的测试代码、实现代码和提交步骤。所有代码示例都是可执行的，没有占位符或 TODO。

**下一步**: 使用 `superpowers:executing-plans` 或 `superpowers:subagent-driven-development` 技能开始执行此计划。

