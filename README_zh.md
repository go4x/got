# Got - 流畅的 Go 测试框架

[![Go 版本](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/)
[![版本](https://img.shields.io/badge/版本-1.0.0-green.svg)](https://github.com/go4x/got)
[![许可证](https://img.shields.io/badge/许可证-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[English](README.md) | [中文](README_zh.md)

---

## 概述

Got 是一个为 Go 应用程序设计的综合测试框架，提供流畅的 API 来编写表达性强、可读性好的测试。它内置支持测试用例、断言、错误处理，以及 Redis 和 SQL 数据库的模拟工具。

## 特性

- 🚀 **流畅 API** - 使用方法链编写表达性强、可读性好的测试
- 🎯 **内置断言** - 为常见测试场景提供全面的断言方法
- 📊 **表驱动测试** - 支持使用 Case 接口的结构化测试用例
- 🔧 **错误处理** - 测试错误条件的实用工具
- 🗄️ **模拟支持** - 内置 Redis 和 SQL 数据库模拟工具
- 🎨 **智能输出** - 带优雅降级的彩色终端输出
- ⚡ **性能** - 轻量级且快速执行

## 安装

```bash
go get github.com/go4x/got
```

## 快速开始

### 基本用法

```go
package main

import (
    "testing"
    "github.com/go4x/got"
)

func TestCalculator(t *testing.T) {
    r := got.New(t, "计算器测试")

    r.Case("测试加法")
    result := 2 + 3
    r.Require(result == 5, "2 + 3 应该等于 5")

    r.Case("测试除零错误")
    _, err := divide(10, 0)
    r.AssertErrf(err, "除零应该返回错误")
}
```

### 表驱动测试

```go
func TestStringLength(t *testing.T) {
    r := got.New(t, "字符串长度测试")
    
    cases := []got.Case{
        got.NewCase("有效输入", "hello", 5, false, nil),
        got.NewCase("空输入", "", 0, false, nil),
    }
    
    r.Cases(cases, func(c got.Case, tt *testing.T) {
        result := len(c.Input().(string))
        r.Require(result == c.Want().(int), "长度应该匹配期望值")
    })
}
```

### 增强断言

```go
func TestEnhancedAssertions(t *testing.T) {
    r := got.New(t, "增强断言")
    
    // 相等断言
    r.AssertEqual(5, 5, "数字应该相等")
    r.AssertNotEqual(5, 6, "数字应该不相等")
    
    // 空值断言
    r.AssertNil(nil, "值应该为空")
    r.AssertNotNil("hello", "值应该不为空")
    
    // 布尔断言
    r.AssertTrue(true, "条件应该为真")
    r.AssertFalse(false, "条件应该为假")
    
    // 包含断言
    r.AssertContains("hello world", "world", "字符串应该包含子字符串")
    r.AssertNotContains("hello world", "foo", "字符串不应该包含子字符串")
}
```

## API 参考

### 核心方法

#### 测试运行器
- `New(t *testing.T, title string) *R` - 创建新的测试运行器
- `Case(format string, args ...any) *R` - 开始新的测试用例
- `Run(name string, f func(t *testing.T)) *R` - 执行子测试
- `Cases(cases []Case, f func(c Case, tt *testing.T))` - 运行表驱动测试

#### 断言
- `Require(cond bool, desc string, args ...any)` - 基本布尔断言
- `FailNow(cond bool, desc string, args ...any)` - 失败时停止的关键断言
- `AssertEqual(expected, actual any, msg ...string) *R` - 相等断言
- `AssertNotEqual(expected, actual any, msg ...string) *R` - 不等断言
- `AssertNil(value any, msg ...string) *R` - 空值断言
- `AssertNotNil(value any, msg ...string) *R` - 非空断言
- `AssertTrue(condition bool, msg ...string) *R` - 真值断言
- `AssertFalse(condition bool, msg ...string) *R` - 假值断言
- `AssertContains(container, item any, msg ...string) *R` - 包含断言
- `AssertNotContains(container, item any, msg ...string) *R` - 不包含断言

#### 错误处理
- `AssertNoErr(err error)` - 断言无错误
- `AssertNoErrf(err error, desc string, args ...any)` - 带描述的无错误断言
- `AssertErr(err error)` - 断言有错误
- `AssertErrf(err error, desc string, args ...any)` - 带描述的错误断言

#### 实用方法
- `Pass(format string, args ...any)` - 记录成功消息
- `Fail(format string, args ...any)` - 记录失败消息
- `Fatal(format string, args ...any)` - 记录致命错误并停止
- `StartTimer() *R` - 开始计时
- `StopTimer() *R` - 停止计时并记录持续时间
- `Parallel() *R` - 标记测试为并行
- `Skip(reason string, args ...any) *R` - 跳过测试
- `Cleanup(fn func()) *R` - 注册清理函数

### 模拟工具

#### Redis 模拟
```go
import "github.com/go4x/got/redist"

// 模拟 Redis 客户端
client, mock := redist.MockRedis()

// 用于测试的 Mini Redis
client, err := redist.NewMiniRedis()
```

#### SQL 模拟
```go
import "github.com/go4x/got/sqlt"

// 创建 SQL 模拟
mockDB, err := sqlt.NewSqlmock()

// 创建 GORM 模拟
gormMock, err := mockDB.Gorm()
```

## 高级特性

### 环境变量
- `NO_COLOR=1` - 禁用彩色输出
- `TERM=dumb` - 使用纯文本输出

### 颜色支持
框架自动检测终端颜色支持并提供：
- 在支持的终端中显示彩色输出
- 在不支持的环境中优雅降级到文本标签

### 性能监控
```go
func TestPerformance(t *testing.T) {
    r := got.New(t, "性能测试")
    
    r.StartTimer()
    // 您的测试代码
    r.StopTimer()
    
    r.MemoryUsage()
    r.GoroutineCount()
    r.TestInfo()
}
```

## 示例

查看 [示例](example_test.go) 了解更多综合使用模式。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

本项目采用 Apache License 2.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 更新日志

### v1.0.0
- 初始版本
- Go 测试的流畅 API
- 内置断言方法
- Redis 和 SQL 模拟工具
- 带降级的智能彩色输出
- 全面的测试覆盖
