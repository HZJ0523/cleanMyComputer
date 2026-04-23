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
