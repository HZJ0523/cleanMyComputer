# 电脑垃圾清理工程 - 设计文档

**版本**: 1.0  
**日期**: 2026-04-23  
**状态**: 设计阶段

---

## 目录

1. [项目概述](#项目概述)
2. [需求总结](#需求总结)
3. [技术方案](#技术方案)
4. [整体架构](#整体架构)
5. [清理对象分层和规则体系](#清理对象分层和规则体系)
6. [界面设计和用户交互流程](#界面设计和用户交互流程)
7. [技术实现细节](#技术实现细节)
8. [错误处理和测试策略](#错误处理和测试策略)
9. [部署和分发](#部署和分发)
10. [后续迭代计划](#后续迭代计划)

---

## 项目概述

### 项目目标

开发一套完备的 Windows 电脑垃圾清理工具，帮助用户安全、高效地清理系统中的各类垃圾文件，释放磁盘空间，提升系统性能。

### 核心价值

- **完整覆盖**: 支持三级清理（安全、深度、高级），覆盖系统、浏览器、应用、用户文件等全方位垃圾
- **安全可控**: 智能风险分级 + 预览确认机制，避免误删重要文件
- **透明可信**: 详细的清理报告和历史记录，用户清楚知道清理了什么
- **高性能**: 使用 Go 语言开发，资源占用低，扫描速度快

---

## 需求总结

### 功能需求

| 需求项 | 描述 |
|--------|------|
| **使用场景** | 个人日常使用 |
| **界面形式** | 图形界面应用（GUI） |
| **清理范围** | 三级完整覆盖（安全、深度、高级） |
| **安全策略** | 预览确认模式 + 智能分级 |
| **触发方式** | 手动触发 + 可选的自动清理计划 |
| **目标平台** | Windows 优先，架构上保留跨平台扩展空间 |
| **结果展示** | 清理前后对比、详细清理报告、历史记录 |
| **项目范围** | 功能完整版 |

### 非功能需求

| 需求项 | 描述 |
|--------|------|
| **性能** | 扫描速度快，资源占用低 |
| **安全** | 避免误删，提供恢复机制 |
| **可靠** | 错误处理完善，不会崩溃 |
| **易用** | 界面清晰，操作简单 |

---

## 技术方案

### 选定方案：Go + Fyne

**技术栈：**
- 后端：Go 1.21+（核心清理逻辑）
- UI框架：Fyne v2.4+（Go 原生 GUI 框架）
- 数据库：SQLite（历史记录、配置）

**选择理由：**
- ✅ **纯 Go 实现**：单一语言栈，代码统一
- ✅ **性能好**：Go 的并发模型适合文件扫描任务
- ✅ **学习曲线平缓**：Go 语法简洁，容易上手
- ✅ **原生渲染**：Fyne 使用 OpenGL 渲染，性能稳定
- ✅ **跨平台能力**：架构上支持未来扩展到 macOS/Linux

---

## 整体架构

采用 **Go + Fyne + 本地持久化** 的分层结构：

```text
┌─────────────────────────────────────┐
│            界面层 UI                 │
│  Fyne Desktop App                   │
│  - 首页仪表盘                        │
│  - 扫描结果预览页                    │
│  - 清理确认页                        │
│  - 历史记录页                        │
│  - 设置页                            │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│            应用层 App               │
│  - 扫描任务编排                      │
│  - 清理任务编排                      │
│  - 计划任务管理                      │
│  - 权限/风险确认流程                 │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│            核心层 Core              │
│  - 规则引擎                          │
│  - 垃圾识别器                        │
│  - 文件扫描器                        │
│  - 风险分级器                        │
│  - 清理执行器                        │
│  - 报告生成器                        │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│         平台适配层 Platform         │
│  - Windows 路径发现                  │
│  - Windows API / 命令调用            │
│  - 回收站/系统目录处理               │
│  - 计划任务接入                      │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│           数据层 Storage            │
│  - SQLite / BoltDB                  │
│  - 清理历史                          │
│  - 用户配置                          │
│  - 清理规则快照                      │
└─────────────────────────────────────┘
```

### 架构原则

1. **界面层只负责交互，不直接碰文件系统**
2. **核心层按“识别 → 评估 → 预览 → 清理”拆分**
3. **平台适配层单独隔离 Windows 细节**
4. **数据层单独保存历史和配置**

---

## 清理对象分层和规则体系

### 三级清理规则

#### 级别 1：安全清理（RiskScore: 0-30）
- 系统临时文件（%TEMP%, C:\Windows\Temp）
- 浏览器缓存（不含密码/Cookie）
- 回收站
- 缩略图缓存
- Windows 更新下载缓存
- 错误报告文件（WER）
- 小型内存转储文件

#### 级别 2：深度清理（RiskScore: 31-60）
- Windows.old 文件夹
- 旧的系统还原点（保留最近2个）
- Microsoft Office 缓存
- Adobe 产品缓存
- 游戏平台缓存（Steam/Epic/Origin）
- 开发工具缓存（npm/pip/Maven/Gradle）
- 下载文件夹中超过30天的文件
- 重复文件扫描
- 大文件扫描（>1GB）
- Windows Installer 缓存（部分）
- 字体缓存
- Microsoft Store 缓存

#### 级别 3：高级清理（RiskScore: 61-100）
- 休眠文件（hiberfil.sys）
- Windows Installer 完整缓存
- 旧驱动程序
- 系统日志文件
- 项目中的 node_modules
- build/dist 输出目录
- Docker 未使用的镜像/容器
- 页面文件优化（不删除，仅建议）

### 规则数据结构

```go
type CleanRule struct {
    ID            string
    Name          string
    Category      string
    Level         int
    Description   string
    Targets       []Target
    RiskScore     int
    RequiresAdmin bool
    Enabled       bool
}

type Target struct {
    Type        string
    Path        string
    Pattern     string
    Recursive   bool
    MaxDepth    int
    ExcludeList []string
}
```

### 智能分级展示

- **绿色**：可安全清理，直接勾选
- **黄色**：建议确认后清理，默认勾选并展示详情
- **红色**：高级清理，默认不勾选，需要二次确认

---

## 界面设计和用户交互流程

### 界面结构（5个页面）

**1. 首页仪表盘**
- 当前磁盘占用概览
- 上次清理时间
- 本次可扫描的垃圾类别数量
- 两个主按钮："快速扫描" 和 "完整扫描"

**2. 扫描结果页**
- 左侧：垃圾分类树（系统/浏览器/应用/用户/高级）
- 右侧：当前分类的详细列表（文件数量、占用空间、风险等级）

**3. 清理确认页**
- 显示本次总释放空间
- 显示高风险项数量
- 显示需要管理员权限的任务
- 提供："开始清理"、"导出清单"、"返回修改"

**4. 历史记录页**
- 每次清理的时间、释放空间、耗时
- 点开查看详细报告
- 支持删除历史记录和导出报告

**5. 设置页**
- 启用/禁用清理规则
- 设置自动清理计划
- 设置旧文件判定天数
- 设置高风险项目隔离时长
- 设置开机提醒扫描

### 核心交互流程

**快速扫描 vs 完整扫描**
- 快速扫描：只扫级别 1 + 部分级别 2，1-3 分钟内出结果
- 完整扫描：扫全部三级规则，包含重复文件、大文件等

**高风险项确认流程**
```
用户勾选高级清理项
  ↓
弹出警告说明
  ↓
展示影响范围
  ↓
要求再次确认
  ↓
加入最终清理队列
```

### 恢复中心

- 高风险清理项先移入隔离区
- 默认保留 24 小时（可配置为 3 天 / 7 天）
- 用户可以查看、恢复或立即永久删除

---

## 技术实现细节

### 项目结构

```text
cleanMyComputer/
├── cmd/
│   └── cleaner/
│       └── main.go                 # 应用入口
├── internal/
│   ├── ui/                         # 界面层
│   │   ├── app.go                  # Fyne 应用初始化
│   │   ├── dashboard.go            # 首页仪表盘
│   │   ├── scanner.go              # 扫描结果页
│   │   ├── confirm.go              # 清理确认页
│   │   ├── history.go              # 历史记录页
│   │   ├── settings.go             # 设置页
│   │   └── components/             # 可复用组件
│   │       ├── progress.go
│   │       ├── file_tree.go
│   │       └── risk_badge.go
│   ├── app/                        # 应用层
│   │   ├── orchestrator.go         # 任务编排器
│   │   ├── scan_task.go            # 扫描任务
│   │   ├── clean_task.go           # 清理任务
│   │   └── scheduler.go            # 计划任务管理
│   ├── core/                       # 核心层
│   │   ├── rule/
│   │   │   ├── engine.go           # 规则引擎
│   │   │   ├── loader.go           # 规则加载器
│   │   │   └── validator.go        # 规则验证器
│   │   ├── scanner/
│   │   │   ├── scanner.go          # 文件扫描器
│   │   │   ├── parallel.go         # 并行扫描
│   │   │   └── filter.go           # 文件过滤器
│   │   ├── analyzer/
│   │   │   ├── risk.go             # 风险评估器
│   │   │   ├── duplicate.go        # 重复文件检测
│   │   │   └── size.go             # 大文件分析
│   │   ├── cleaner/
│   │   │   ├── executor.go         # 清理执行器
│   │   │   ├── quarantine.go       # 隔离区管理
│   │   │   └── recovery.go         # 恢复功能
│   │   └── report/
│   │       ├── generator.go        # 报告生成器
│   │       └── exporter.go         # 报告导出器
│   ├── platform/                   # 平台适配层
│   │   ├── windows/
│   │   │   ├── paths.go            # Windows 路径发现
│   │   │   ├── registry.go         # 注册表操作
│   │   │   ├── recyclebin.go       # 回收站处理
│   │   │   └── admin.go            # 管理员权限检测
│   │   └── common/
│   │       └── interface.go        # 平台接口定义
│   ├── storage/                    # 数据层
│   │   ├── db.go                   # 数据库初始化
│   │   ├── history.go              # 历史记录操作
│   │   ├── config.go               # 配置管理
│   │   └── migrations/             # 数据库迁移
│   │       └── 001_init.sql
│   └── models/                     # 数据模型
│       ├── rule.go
│       ├── scan_result.go
│       ├── clean_record.go
│       └── config.go
├── pkg/                            # 可复用包
│   ├── fileutil/                   # 文件工具
│   │   ├── hash.go                 # 文件哈希计算
│   │   ├── size.go                 # 大小格式化
│   │   └── permission.go           # 权限检查
│   └── logger/                     # 日志工具
│       └── logger.go
├── configs/                        # 配置文件
│   ├── rules/                      # 清理规则定义
│   │   ├── level1_safe.json
│   │   ├── level2_deep.json
│   │   └── level3_advanced.json
│   └── app.yaml                    # 应用配置
├── assets/                         # 资源文件
│   ├── icons/
│   └── i18n/                       # 国际化
│       ├── zh-CN.json
│       └── en-US.json
├── scripts/                        # 构建脚本
│   ├── build.sh
│   └── package.sh
├── tests/                          # 测试
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── go.mod
├── go.sum
└── README.md
```

### 核心模块设计

#### 1. 规则引擎（Rule Engine）

规则引擎负责加载、验证和管理清理规则。

```go
// internal/core/rule/engine.go
package rule

import (
    "context"
    "sync"
)

type Engine struct {
    rules  map[string]*CleanRule
    mu     sync.RWMutex
    loader *Loader
}

func NewEngine(loader *Loader) *Engine {
    return &Engine{
        rules:  make(map[string]*CleanRule),
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
        if err := e.validateRule(rule); err != nil {
            continue
        }
        e.rules[rule.ID] = rule
    }
    return nil
}

func (e *Engine) GetEnabledRules(level int) []*CleanRule {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    var enabled []*CleanRule
    for _, rule := range e.rules {
        if rule.Enabled && rule.Level <= level {
            enabled = append(enabled, rule)
        }
    }
    return enabled
}
```

#### 2. 文件扫描器（Scanner）

扫描器使用并行策略快速扫描文件系统。

```go
// internal/core/scanner/scanner.go
package scanner

import (
    "context"
    "os"
    "path/filepath"
    "sync"
)

type Scanner struct {
    workers   int
    filter    *Filter
    results   chan *ScanResult
}

type ScanResult struct {
    Path      string
    Size      int64
    ModTime   time.Time
    RuleID    string
    RiskScore int
}

func NewScanner(workers int) *Scanner {
    return &Scanner{
        workers: workers,
        filter:  NewFilter(),
        results: make(chan *ScanResult, 1000),
    }
}

func (s *Scanner) ScanTargets(ctx context.Context, targets []Target) ([]*ScanResult, error) {
    var wg sync.WaitGroup
    taskChan := make(chan Target, len(targets))
    
    // 启动工作协程
    for i := 0; i < s.workers; i++ {
        wg.Add(1)
        go s.worker(ctx, taskChan, &wg)
    }
    
    // 分发任务
    for _, target := range targets {
        select {
        case taskChan <- target:
        case <-ctx.Done():
            close(taskChan)
            return nil, ctx.Err()
        }
    }
    close(taskChan)
    
    // 等待完成
    go func() {
        wg.Wait()
        close(s.results)
    }()
    
    // 收集结果
    var results []*ScanResult
    for result := range s.results {
        results = append(results, result)
    }
    return results, nil
}

func (s *Scanner) worker(ctx context.Context, tasks <-chan Target, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for target := range tasks {
        if err := s.scanTarget(ctx, target); err != nil {
            continue
        }
    }
}
```

#### 3. 风险评估器（Risk Analyzer）

风险评估器根据文件类型、位置、大小等因素计算风险分数。

```go
// internal/core/analyzer/risk.go
package analyzer

import (
    "path/filepath"
    "strings"
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

func (r *RiskAnalyzer) CalculateRisk(result *ScanResult) int {
    score := result.RiskScore
    
    // 系统路径增加风险
    if r.isSystemPath(result.Path) {
        score += 30
    }
    
    // 受保护的扩展名增加风险
    ext := strings.ToLower(filepath.Ext(result.Path))
    if r.protectedExts[ext] {
        score += 20
    }
    
    // 大文件增加风险（可能是重要数据）
    if result.Size > 100*1024*1024 { // >100MB
        score += 15
    }
    
    // 最近修改的文件增加风险
    if time.Since(result.ModTime) < 7*24*time.Hour {
        score += 10
    }
    
    // 限制在 0-100 范围内
    if score > 100 {
        score = 100
    }
    return score
}

func (r *RiskAnalyzer) GetRiskLevel(score int) string {
    switch {
    case score <= 30:
        return "safe"
    case score <= 60:
        return "moderate"
    default:
        return "high"
    }
}
```

#### 4. 清理执行器（Cleaner Executor）

清理执行器负责安全地删除文件，高风险项先移入隔离区。

```go
// internal/core/cleaner/executor.go
package cleaner

import (
    "context"
    "os"
    "path/filepath"
)

type Executor struct {
    quarantine *QuarantineManager
    dryRun     bool
}

type CleanTask struct {
    Files     []*FileItem
    TotalSize int64
    RiskLevel string
}

type FileItem struct {
    Path      string
    Size      int64
    RiskScore int
}

func NewExecutor(quarantine *QuarantineManager) *Executor {
    return &Executor{
        quarantine: quarantine,
        dryRun:     false,
    }
}

func (e *Executor) Execute(ctx context.Context, task *CleanTask) (*CleanResult, error) {
    result := &CleanResult{
        StartTime: time.Now(),
    }
    
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
    
    // 高风险文件先移入隔离区
    if file.RiskScore > 60 {
        return e.quarantine.Quarantine(file.Path)
    }
    
    // 低风险文件直接删除
    return os.Remove(file.Path)
}
```

#### 5. 隔离区管理（Quarantine Manager）

隔离区采用独立目录 + 元数据索引，支持恢复与过期清理。

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
    repo    QuarantineRepository
}

func (q *QuarantineManager) Quarantine(srcPath string) error {
    fileName := filepath.Base(srcPath)
    quarantineName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileName)
    dstPath := filepath.Join(q.baseDir, quarantineName)
    
    if err := os.Rename(srcPath, dstPath); err != nil {
        return err
    }
    
    record := &QuarantineRecord{
        OriginalPath:   srcPath,
        QuarantinePath: dstPath,
        QuarantinedAt:  time.Now(),
        ExpiresAt:      time.Now().Add(24 * time.Hour),
    }
    return q.repo.Save(record)
}

func (q *QuarantineManager) Restore(id int64) error {
    record, err := q.repo.GetByID(id)
    if err != nil {
        return err
    }
    return os.Rename(record.QuarantinePath, record.OriginalPath)
}
```

### 数据库设计

数据库采用 SQLite，本地存储于 `%LOCALAPPDATA%/CleanMyComputer/data/cleaner.db`。

#### 1. clean_history

记录每次清理任务的汇总信息。

```sql
CREATE TABLE clean_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    scan_level INTEGER NOT NULL,           -- 1=安全, 2=深度, 3=高级
    total_files INTEGER NOT NULL,
    total_size INTEGER NOT NULL,           -- 字节
    freed_size INTEGER NOT NULL,
    failed_count INTEGER DEFAULT 0,
    status TEXT NOT NULL,                  -- success, partial, failed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### 2. clean_details

记录每次清理的详细文件列表。

```sql
CREATE TABLE clean_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    history_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    rule_id TEXT NOT NULL,
    risk_score INTEGER NOT NULL,
    action TEXT NOT NULL,                  -- deleted, quarantined, failed
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (history_id) REFERENCES clean_history(id) ON DELETE CASCADE
);

CREATE INDEX idx_clean_details_history ON clean_details(history_id);
CREATE INDEX idx_clean_details_rule ON clean_details(rule_id);
```

#### 3. quarantine

隔离区文件记录。

```sql
CREATE TABLE quarantine (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_path TEXT NOT NULL,
    quarantine_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    risk_score INTEGER NOT NULL,
    quarantined_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    restored BOOLEAN DEFAULT 0,
    restored_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quarantine_expires ON quarantine(expires_at);
CREATE INDEX idx_quarantine_restored ON quarantine(restored);
```

#### 4. user_config

用户配置项。

```sql
CREATE TABLE user_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 默认配置
INSERT INTO user_config (key, value) VALUES
    ('quarantine_retention_hours', '24'),
    ('auto_clean_enabled', 'false'),
    ('auto_clean_schedule', ''),
    ('old_file_days', '30'),
    ('scan_workers', '4'),
    ('language', 'zh-CN');
```

#### 5. rule_status

规则启用状态（用户可自定义）。

```sql
CREATE TABLE rule_status (
    rule_id TEXT PRIMARY KEY,
    enabled BOOLEAN DEFAULT 1,
    last_used DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 性能优化策略

#### 1. 并行扫描

使用 Go 协程池并行扫描多个目录，默认 4 个工作协程。

```go
// 根据 CPU 核心数动态调整
workers := runtime.NumCPU()
if workers > 8 {
    workers = 8
}
scanner := NewScanner(workers)
```

#### 2. 增量扫描

缓存上次扫描结果，仅扫描变化的目录。

```go
type ScanCache struct {
    LastScan  time.Time
    DirHashes map[string]string  // 目录路径 -> 修改时间哈希
}
```

#### 3. 分块处理大目录

对包含大量文件的目录采用分页扫描，避免一次性加载过多元数据。

- 单批次最多处理 1000 个文件
- 扫描结果实时刷新到 UI，避免长时间无响应
- 重复文件计算采用文件大小预分组，再计算哈希值

#### 4. 数据库优化

- 为历史记录、隔离区过期时间建立索引
- 清理详情表使用批量插入
- 历史明细超过阈值时归档或压缩
- 读操作使用只读事务，写操作使用短事务

#### 5. IO 优化

- 先按目录分组，减少磁盘随机寻址
- 跳过已知系统保护目录和无权限目录
- 大文件采用流式读取，不整体载入内存
- 哈希计算仅用于重复文件候选集，不对全部文件执行

---

## 错误处理和测试策略

### 错误分类和处理原则

#### 1. 错误分类

| 错误类型 | 典型场景 | 处理方式 |
|----------|----------|----------|
| **权限错误** | 无法访问系统目录、注册表、回收站 | 记录并跳过，提示需要管理员权限 |
| **路径错误** | 路径不存在、路径被占用、路径格式非法 | 标记失败项，不中断整体任务 |
| **文件状态错误** | 文件被锁定、正在使用、已被删除 | 重试一次，失败则记录 |
| **数据库错误** | SQLite 锁冲突、迁移失败、写入失败 | 回滚事务，保留内存态结果 |
| **规则错误** | 规则配置缺失、字段非法、路径模板错误 | 禁用该规则并记录告警 |
| **系统错误** | 磁盘不可用、空间不足、系统 API 调用失败 | 终止当前阶段，提示用户处理 |

#### 2. 处理原则

1. **局部失败不影响整体**：单个文件失败不应导致整个清理任务终止
2. **可恢复优先**：高风险文件先隔离，避免永久删除
3. **错误可追踪**：所有失败项必须记录错误码、错误信息、发生时间
4. **用户可理解**：界面提示用业务语言，不直接暴露底层异常堆栈
5. **日志可诊断**：日志保留原始错误链，便于开发调试

#### 3. 错误处理示例

```go
// internal/core/cleaner/executor.go
func (e *Executor) cleanFile(file *FileItem) error {
    const maxRetries = 2
    
    for i := 0; i < maxRetries; i++ {
        err := e.doClean(file)
        if err == nil {
            return nil
        }
        
        // 权限错误不重试
        if errors.Is(err, os.ErrPermission) {
            return fmt.Errorf("权限不足，需要管理员权限: %w", err)
        }
        
        // 文件被占用，等待后重试
        if isFileLocked(err) && i < maxRetries-1 {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        
        return err
    }
    return fmt.Errorf("清理失败，已重试 %d 次", maxRetries)
}
```

### 测试策略

#### 1. 单元测试（Unit Tests）

覆盖核心逻辑模块，目标覆盖率 80%+。

```go
// tests/unit/scanner_test.go
func TestScanner_ScanTargets(t *testing.T) {
    tmpDir := t.TempDir()
    createTestFiles(tmpDir, 100)
    
    scanner := NewScanner(4)
    targets := []Target{{Path: tmpDir, Recursive: true}}
    
    results, err := scanner.ScanTargets(context.Background(), targets)
    assert.NoError(t, err)
    assert.Equal(t, 100, len(results))
}

func TestRiskAnalyzer_CalculateRisk(t *testing.T) {
    analyzer := NewRiskAnalyzer()
    
    tests := []struct {
        name     string
        result   *ScanResult
        expected int
    }{
        {
            name: "系统路径高风险",
            result: &ScanResult{
                Path: "C:\\Windows\\System32\\test.dll",
                RiskScore: 20,
            },
            expected: 70, // 20 + 30(系统路径) + 20(dll扩展名)
        },
        {
            name: "临时文件低风险",
            result: &ScanResult{
                Path: "C:\\Users\\Test\\AppData\\Local\\Temp\\cache.tmp",
                RiskScore: 10,
            },
            expected: 10,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            score := analyzer.CalculateRisk(tt.result)
            assert.Equal(t, tt.expected, score)
        })
    }
}
```

#### 2. 集成测试（Integration Tests）

测试模块间协作，使用真实文件系统和数据库。

```go
// tests/integration/clean_workflow_test.go
func TestCleanWorkflow_EndToEnd(t *testing.T) {
    // 准备测试环境
    tmpDir := t.TempDir()
    db := setupTestDB(t)
    defer db.Close()
    
    // 创建测试文件
    createTestFiles(tmpDir, map[string]int64{
        "temp/cache1.tmp": 1024,
        "temp/cache2.tmp": 2048,
        "downloads/old_file.zip": 5000,
    })
    
    // 执行完整流程
    engine := rule.NewEngine(rule.NewLoader("configs/rules"))
    scanner := scanner.NewScanner(4)
    analyzer := analyzer.NewRiskAnalyzer()
    executor := cleaner.NewExecutor(quarantine.NewManager(tmpDir + "/quarantine"))
    
    // 1. 加载规则
    err := engine.LoadRules(context.Background(), 1)
    assert.NoError(t, err)
    
    // 2. 扫描文件
    rules := engine.GetEnabledRules(1)
    results, err := scanner.ScanTargets(context.Background(), extractTargets(rules))
    assert.NoError(t, err)
    assert.Greater(t, len(results), 0)
    
    // 3. 风险评估
    for _, result := range results {
        result.RiskScore = analyzer.CalculateRisk(result)
    }
    
    // 4. 执行清理
    task := buildCleanTask(results)
    cleanResult, err := executor.Execute(context.Background(), task)
    assert.NoError(t, err)
    assert.Greater(t, cleanResult.FreedSize, int64(0))
}
```

#### 3. E2E 测试（End-to-End Tests）

测试完整用户场景，包括 UI 交互。

```go
// tests/e2e/app_test.go
func TestApp_QuickScanAndClean(t *testing.T) {
    app := test.NewApp()
    window := app.NewWindow("Test")
    
    // 模拟用户点击"快速扫描"
    scanButton := findButton(window, "快速扫描")
    test.Tap(scanButton)
    
    // 等待扫描完成
    waitForScanComplete(t, 30*time.Second)
    
    // 验证扫描结果显示
    resultList := findList(window, "scan-results")
    assert.Greater(t, resultList.Length(), 0)
    
    // 模拟用户点击"开始清理"
    cleanButton := findButton(window, "开始清理")
    test.Tap(cleanButton)
    
    // 等待清理完成
    waitForCleanComplete(t, 30*time.Second)
    
    // 验证清理报告
    report := findLabel(window, "clean-report")
    assert.Contains(t, report.Text, "清理完成")
}
```

#### 4. 性能测试（Performance Tests）

测试大规模文件扫描和清理的性能表现。

```go
// tests/performance/benchmark_test.go
func BenchmarkScanner_LargeDirectory(b *testing.B) {
    tmpDir := b.TempDir()
    createTestFiles(tmpDir, 10000) // 创建 10000 个文件
    
    scanner := NewScanner(4)
    targets := []Target{{Path: tmpDir, Recursive: true}}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := scanner.ScanTargets(context.Background(), targets)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkDuplicateFinder_HashCalculation(b *testing.B) {
    tmpDir := b.TempDir()
    files := createDuplicateFiles(tmpDir, 1000, 1024*1024) // 1000个1MB文件
    
    finder := NewDuplicateFinder()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := finder.FindDuplicates(files)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 安全考虑

#### 1. 权限管理

```go
// internal/platform/windows/admin.go
package windows

import (
    "golang.org/x/sys/windows"
)

func IsAdmin() bool {
    var sid *windows.SID
    err := windows.AllocateAndInitializeSid(
        &windows.SECURITY_NT_AUTHORITY,
        2,
        windows.SECURITY_BUILTIN_DOMAIN_RID,
        windows.DOMAIN_ALIAS_RID_ADMINS,
        0, 0, 0, 0, 0, 0,
        &sid)
    if err != nil {
        return false
    }
    defer windows.FreeSid(sid)
    
    token := windows.Token(0)
    member, err := token.IsMember(sid)
    return err == nil && member
}

func RequireAdmin(operation string) error {
    if !IsAdmin() {
        return fmt.Errorf("操作 '%s' 需要管理员权限", operation)
    }
    return nil
}
```

#### 2. 路径验证

防止路径遍历攻击和误删系统关键文件。

```go
// pkg/fileutil/permission.go
package fileutil

import (
    "path/filepath"
    "strings"
)

var protectedPaths = []string{
    "C:\\Windows\\System32",
    "C:\\Windows\\SysWOW64",
    "C:\\Program Files\\WindowsApps",
    "C:\\ProgramData\\Microsoft\\Windows\\Start Menu",
}

func IsPathSafe(path string) error {
    // 规范化路径
    absPath, err := filepath.Abs(path)
    if err != nil {
        return fmt.Errorf("无效路径: %w", err)
    }
    
    // 检查是否在保护路径内
    for _, protected := range protectedPaths {
        if strings.HasPrefix(strings.ToLower(absPath), strings.ToLower(protected)) {
            return fmt.Errorf("禁止操作受保护路径: %s", protected)
        }
    }
    
    // 检查路径遍历
    if strings.Contains(absPath, "..") {
        return fmt.Errorf("路径包含非法字符: ..")
    }
    
    return nil
}
```

#### 3. 数据隐私

- 不收集用户文件内容，仅记录路径和元数据
- 历史记录仅本地存储，不上传云端
- 导出报告时脱敏用户名和敏感路径
- 清理规则不扫描用户文档、图片、视频等个人文件（除非用户主动选择）

```go
// internal/core/report/exporter.go
func (e *Exporter) SanitizePath(path string) string {
    // 替换用户名
    username := os.Getenv("USERNAME")
    if username != "" {
        path = strings.ReplaceAll(path, username, "<USER>")
    }
    
    // 替换常见敏感路径
    path = strings.ReplaceAll(path, os.Getenv("USERPROFILE"), "%USERPROFILE%")
    path = strings.ReplaceAll(path, os.Getenv("APPDATA"), "%APPDATA%")
    
    return path
}
```

---

## 部署和分发

### 构建配置

使用 Go 交叉编译和 Fyne 打包工具生成 Windows 可执行文件。

```bash
# scripts/build.sh
#!/bin/bash

VERSION="1.0.0"
APP_NAME="CleanMyComputer"

# 构建 Windows 64位版本
GOOS=windows GOARCH=amd64 go build \
    -ldflags="-s -w -H windowsgui -X main.Version=${VERSION}" \
    -o dist/${APP_NAME}.exe \
    cmd/cleaner/main.go

# 使用 Fyne 打包工具生成安装包
fyne package \
    -os windows \
    -name ${APP_NAME} \
    -icon assets/icons/app.png \
    -appVersion ${VERSION} \
    -release
```

### 安装包内容

```text
CleanMyComputer_Setup_v1.0.0.exe
├── CleanMyComputer.exe           # 主程序
├── configs/
│   ├── rules/                    # 清理规则配置
│   │   ├── level1_safe.json
│   │   ├── level2_deep.json
│   │   └── level3_advanced.json
│   └── app.yaml                  # 应用配置
├── assets/
│   ├── icons/                    # 图标资源
│   └── i18n/                     # 语言包
│       ├── zh-CN.json
│       └── en-US.json
├── LICENSE.txt                   # 许可证
└── README.txt                    # 使用说明
```

**安装位置**：
- 程序文件：`C:\Program Files\CleanMyComputer\`
- 用户数据：`%LOCALAPPDATA%\CleanMyComputer\`
  - `data/cleaner.db`：数据库
  - `quarantine/`：隔离区
  - `logs/`：日志文件

**安装器功能**：
- 检测 Windows 版本（要求 Windows 10 1809+）
- 创建桌面快捷方式
- 添加到开始菜单
- 可选：开机自启动
- 可选：添加右键菜单"使用 CleanMyComputer 清理"

---

## 后续迭代计划

### 第一阶段：核心功能（v1.0）

- ✅ 三级清理规则体系
- ✅ 智能风险评估
- ✅ 隔离区和恢复机制
- ✅ 历史记录和报告导出
- ✅ 基础设置和配置管理

### 第二阶段：体验增强（v1.1 - v1.2）

#### 1. 跨平台支持
- 扩展平台适配层，增加 macOS 支持
- 增加 Linux 支持（Ubuntu / Debian 优先）
- 抽象平台相关路径、缓存目录和系统 API
- 提供统一规则描述，平台适配层做最终映射

#### 2. 智能推荐系统
- 基于用户清理历史推荐常用规则组合
- 根据磁盘使用趋势提示高收益清理项
- 自动识别长期未访问的大文件和重复文件
- 提供"推荐清理"模式，默认仅展示高收益低风险项

### 迭代原则

1. **安全优先**：任何新增能力都不能降低现有误删保护级别
2. **本地优先**：默认离线可用，联网能力均为可选增强
3. **渐进扩展**：先完善 Windows 核心体验，再考虑跨平台和 AI 能力
4. **可解释性**：所有清理建议都要能向用户说明依据
5. **可回退**：关键功能上线前必须具备回滚和恢复方案

