# Got - A Fluent Go Testing Framework

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)](https://github.com/go4x/got)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[English](README.md) | [中文](README_zh.md)

---

## Overview

Got is a comprehensive testing framework for Go applications that provides a fluent API for writing expressive and readable tests. It offers built-in support for test cases, assertions, error handling, and mock utilities for Redis and SQL databases.

## Features

- 🚀 **Fluent API** - Write expressive and readable tests with method chaining
- 🎯 **Built-in Assertions** - Comprehensive assertion methods for common testing scenarios
- 📊 **Table-driven Tests** - Support for structured test cases with the Case interface
- 🔧 **Error Handling** - Utilities for testing error conditions
- 🗄️ **Mock Support** - Built-in mock utilities for Redis and SQL databases
- 🎨 **Smart Output** - Color-coded terminal output with graceful fallback
- ⚡ **Performance** - Lightweight and fast execution

## Installation

```bash
go get github.com/go4x/got
```

## Quick Start

### Basic Usage

```go
package main

import (
    "testing"
    "github.com/go4x/got"
)

func TestCalculator(t *testing.T) {
    r := got.New(t, "Calculator Tests")

    r.Case("Testing addition")
    result := 2 + 3
    r.Require(result == 5, "2 + 3 should equal 5")

    r.Case("Testing division by zero")
    _, err := divide(10, 0)
    r.AssertErrf(err, "Division by zero should return error")
}
```

### Table-driven Tests

```go
func TestStringLength(t *testing.T) {
    r := got.New(t, "String Length Tests")
    
    cases := []got.Case{
        got.NewCase("Valid Input", "hello", 5, false, nil),
        got.NewCase("Empty Input", "", 0, false, nil),
    }
    
    r.Cases(cases, func(c got.Case, tt *testing.T) {
        result := len(c.Input().(string))
        r.Require(result == c.Want().(int), "Length should match expected")
    })
}
```

### Enhanced Assertions

```go
func TestEnhancedAssertions(t *testing.T) {
    r := got.New(t, "Enhanced Assertions")
    
    // Equality assertions
    r.AssertEqual(5, 5, "Numbers should be equal")
    r.AssertNotEqual(5, 6, "Numbers should not be equal")
    
    // Nil assertions
    r.AssertNil(nil, "Value should be nil")
    r.AssertNotNil("hello", "Value should not be nil")
    
    // Boolean assertions
    r.AssertTrue(true, "Condition should be true")
    r.AssertFalse(false, "Condition should be false")
    
    // Contains assertions
    r.AssertContains("hello world", "world", "String should contain substring")
    r.AssertNotContains("hello world", "foo", "String should not contain substring")
}
```

## API Reference

### Core Methods

#### Test Runner
- `New(t *testing.T, title string) *R` - Create a new test runner
- `Case(format string, args ...any) *R` - Start a new test case
- `Run(name string, f func(t *testing.T)) *R` - Execute a subtest
- `Cases(cases []Case, f func(c Case, tt *testing.T))` - Run table-driven tests

#### Assertions
- `Require(cond bool, desc string, args ...any)` - Basic boolean assertion
- `FailNow(cond bool, desc string, args ...any)` - Critical assertion that stops on failure
- `AssertEqual(expected, actual any, msg ...string) *R` - Equality assertion
- `AssertNotEqual(expected, actual any, msg ...string) *R` - Inequality assertion
- `AssertNil(value any, msg ...string) *R` - Nil assertion
- `AssertNotNil(value any, msg ...string) *R` - Non-nil assertion
- `AssertTrue(condition bool, msg ...string) *R` - True assertion
- `AssertFalse(condition bool, msg ...string) *R` - False assertion
- `AssertContains(container, item any, msg ...string) *R` - Contains assertion
- `AssertNotContains(container, item any, msg ...string) *R` - Not contains assertion

#### Error Handling
- `AssertNoErr(err error)` - Assert no error
- `AssertNoErrf(err error, desc string, args ...any)` - Assert no error with description
- `AssertErr(err error)` - Assert error exists
- `AssertErrf(err error, desc string, args ...any)` - Assert error with description

#### Utility Methods
- `Pass(format string, args ...any)` - Log success message
- `Fail(format string, args ...any)` - Log failure message
- `Fatal(format string, args ...any)` - Log fatal error and stop
- `StartTimer() *R` - Start timing
- `StopTimer() *R` - Stop timing and log duration
- `Parallel() *R` - Mark test as parallel
- `Skip(reason string, args ...any) *R` - Skip test
- `Cleanup(fn func()) *R` - Register cleanup function

### Mock Utilities

#### Redis Mock
```go
import "github.com/go4x/got/redist"

// Mock Redis client
client, mock := redist.MockRedis()

// Mini Redis for testing
client, err := redist.NewMiniRedis()
```

#### SQL Mock
```go
import "github.com/go4x/got/sqlt"

// Create SQL mock
mockDB, err := sqlt.NewSqlmock()

// Create GORM mock
gormMock, err := mockDB.Gorm()
```

## Advanced Features

### Environment Variables
- `NO_COLOR=1` - Disable color output
- `TERM=dumb` - Use text-only output

### Color Support
The framework automatically detects terminal color support and provides:
- Color-coded output in supported terminals
- Graceful fallback to text labels in unsupported environments

### Performance Monitoring
```go
func TestPerformance(t *testing.T) {
    r := got.New(t, "Performance Test")
    
    r.StartTimer()
    // Your test code here
    r.StopTimer()
    
    r.MemoryUsage()
    r.GoroutineCount()
    r.TestInfo()
}
```

## Examples

Check out the [examples](example_test.go) for more comprehensive usage patterns.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Fluent API for Go testing
- Built-in assertion methods
- Redis and SQL mock utilities
- Smart color output with fallback
- Comprehensive test coverage
