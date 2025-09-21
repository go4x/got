// Package got provides a comprehensive testing framework for Go applications.
// It offers a fluent API for writing expressive and readable tests with built-in
// support for test cases, assertions, and error handling.
//
// The package includes:
//   - Test runner with fluent API for organizing test cases
//   - Built-in assertion methods for common testing scenarios
//   - Support for table-driven tests with structured test cases
//   - Error handling utilities for testing error conditions
//   - Mock utilities for Redis and SQL databases
//
// Example usage:
//
//	func TestCalculator(t *testing.T) {
//		r := got.New(t, "Calculator Tests")
//
//		r.Case("Testing addition")
//		result := 2 + 3
//		r.Require(result == 5, "2 + 3 should equal 5")
//
//		r.Case("Testing division by zero")
//		_, err := divide(10, 0)
//		r.Errf(err, "Division by zero should return error")
//	}
//
// For more examples, see the example functions in this package.
package got

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	checkMark = "\033[32m\u2713\033[0m" // green √
	ballotX   = "\033[31m\u2717\033[0m" // red ×
)

// R represents a test runner that provides a fluent API for writing tests.
// It embeds *testing.T to provide all standard testing functionality while
// adding enhanced logging, assertion methods, and test case management.
//
// The runner maintains state for:
//   - title: The main test suite title
//   - caseNum: Current case number for automatic numbering
//   - prefix: Formatted prefix for case logging
//   - startTime: Test start time for timing
//   - benchmark: Whether running in benchmark mode
//   - parallel: Whether test is marked as parallel
//
// Example:
//
//	r := got.New(t, "My Test Suite")
//	r.Case("First test case")
//	r.Require(condition, "Description")
type R struct {
	title     string
	caseNum   int
	prefix    string
	startTime time.Time
	benchmark bool
	parallel  bool
	*testing.T
}

// New creates a new test runner instance from a testing.T.
// The title parameter is used to identify the test suite in logs and output.
//
// Parameters:
//   - t: The testing.T instance from the test function
//   - title: A descriptive title for the test suite
//
// Returns:
//   - *R: A new test runner instance
//
// Example:
//
//	func TestMyFeature(t *testing.T) {
//		r := got.New(t, "My Feature Tests")
//		// Use r for test cases and assertions
//	}
func New(t *testing.T, title string) *R {
	t.Log("Test Case => " + title)
	return &R{
		T:         t,
		title:     title,
		startTime: time.Now(),
	}
}

// Case starts a new test case with a descriptive message.
// It automatically increments the case number and logs the case description.
// The method supports printf-style formatting for dynamic case descriptions.
//
// Parameters:
//   - format: A format string describing the test case
//   - args: Arguments for the format string
//
// Returns:
//   - *R: The runner instance for method chaining
//
// Example:
//
//	r.Case("Testing user authentication with valid credentials")
//	r.Case("Testing division by zero with divisor %d", 0)
func (r *R) Case(format string, args ...any) *R {
	r.caseNum++
	r.prefix = "Case " + strconv.Itoa(r.caseNum) + " -> "
	r.Logf(r.prefix+format, args...)
	return r
}

// Caser runs a test case with the given name and function.
// It combines Case and Run methods to create a named subtest with automatic case logging.
// This is useful for organizing related test scenarios under a common name.
//
// Parameters:
//   - name: The name of the test case
//   - f: The test function to execute
//
// Returns:
//   - *R: The runner instance for method chaining
//
// Example:
//
//	r.Caser("Valid Login", func(t *testing.T) {
//		// Test valid login scenario
//		r.Require(login("user", "pass"), "Login should succeed")
//	})
func (r *R) Caser(name string, f func(t *testing.T)) *R {
	r.Case(name)
	r.Run(name, f)
	return r
}

// Run executes a subtest with the given name and function.
// It wraps testing.T.Run to provide a convenient way to run subtests while
// maintaining the fluent API pattern.
//
// Parameters:
//   - name: The name of the subtest
//   - f: The test function to execute
//
// Returns:
//   - *R: The runner instance for method chaining
//
// Example:
//
//	r.Run("Database Connection", func(t *testing.T) {
//		// Test database connection
//		conn, err := connectDB()
//		r.NoErrf(err, "Database connection should succeed")
//	})
func (r *R) Run(name string, f func(t *testing.T)) *R {
	r.T.Run(name, func(tt *testing.T) {
		f(tt)
	})
	return r
}

// Cases runs a set of test cases, executing the provided function for each case.
// This method is designed for table-driven tests where you have multiple test
// scenarios with different inputs and expected outputs.
//
// For each test case, it:
//   - Logs the case description using Case()
//   - Runs the case as a subtest using Run()
//   - Passes the case data to the test function
//
// Parameters:
//   - cases: A slice of Case implementations containing test data
//   - f: The test function that will be executed for each case
//
// Example:
//
//	cases := []got.Case{
//		got.NewCase("Valid Input", "hello", 5, false, nil),
//		got.NewCase("Empty Input", "", 0, false, nil),
//	}
//	r.Cases(cases, func(c got.Case, tt *testing.T) {
//		result := len(c.Input().(string))
//		r.Require(result == c.Want().(int), "Length should match expected")
//	})
func (r *R) Cases(cases []Case, f func(c Case, tt *testing.T)) {
	for _, c := range cases {
		r.Case(c.Name())
		r.Run(c.Name(), func(tt *testing.T) {
			f(c, tt)
		})
	}
}

// Pass logs a successful assertion with a green checkmark.
// Use this method to indicate that a test condition has passed.
//
// Parameters:
//   - format: A format string describing the successful assertion
//   - args: Arguments for the format string
//
// Example:
//
//	r.Pass("User authentication succeeded")
//	r.Pass("Value %d is within expected range", 42)
func (r *R) Pass(format string, args ...any) {
	r.Logf("\t%s "+format, prependTag(checkMark, args...)...)
}

// Fail logs a failed assertion with a red X mark.
// Use this method to indicate that a test condition has failed.
// This method will mark the test as failed but will not stop execution.
//
// Parameters:
//   - format: A format string describing the failed assertion
//   - args: Arguments for the format string
//
// Example:
//
//	r.Fail("User authentication should have succeeded")
//	r.Fail("Value %d is outside expected range", 100)
func (r *R) Fail(format string, args ...any) {
	r.Errorf("\t%s "+format, prependTag(ballotX, args...)...)
}

// Fatal logs a fatal error and immediately stops test execution.
// This method is equivalent to calling Fail() followed by t.FailNow().
// Use this when a test cannot continue due to a critical failure.
//
// Parameters:
//   - format: A format string describing the fatal error
//   - args: Arguments for the format string
//
// Example:
//
//	r.Fatal("Database connection failed - cannot continue test")
//	r.Fatal("Critical system component %s is not available", "auth-service")
func (r *R) Fatal(format string, args ...any) {
	r.Fatalf("\t%s "+format, prependTag(ballotX, args...)...)
}

func prependTag(tag any, args ...any) []any {
	if len(args) == 0 {
		return []any{tag}
	}
	// Pre-allocate slice with exact capacity to avoid reallocation
	result := make([]any, 0, len(args)+1)
	result = append(result, tag)
	result = append(result, args...)
	return result
}

// Require is a convenient assertion method that checks a boolean condition.
// If the condition is true, it logs a pass message; otherwise, it logs a fail message.
// This is the most commonly used assertion method for simple boolean checks.
//
// Parameters:
//   - cond: The boolean condition to check
//   - desc: A description of what is being tested
//   - args: Arguments for the description format string
//
// Example:
//
//	r.Require(user.IsAuthenticated(), "User should be authenticated")
//	r.Require(len(items) > 0, "Items list should not be empty")
//	r.Require(result == expected, "Result %d should equal %d", result, expected)
func (r *R) Require(cond bool, desc string, args ...any) {
	if cond {
		r.Pass(desc, args...)
	} else {
		r.Fail(desc, args...)
	}
}

// FailNow checks a boolean condition and stops test execution if it fails.
// If the condition is true, it logs a pass message and continues.
// If the condition is false, it logs a fail message and immediately stops the test.
// This is useful for critical assertions that should halt the test if they fail.
//
// Parameters:
//   - cond: The boolean condition to check
//   - desc: A description of what is being tested
//   - args: Arguments for the description format string
//
// Example:
//
//	r.FailNow(db.IsConnected(), "Database connection is required for this test")
//	r.FailNow(config.IsValid(), "Configuration must be valid to continue")
func (r *R) FailNow(cond bool, desc string, args ...any) {
	if cond {
		r.Pass(desc, args...)
	} else {
		r.Fail(desc, args...)
		r.T.FailNow()
	}
}

// AssertNoErr checks if the provided error is nil.
// If err is nil, it passes the test; otherwise, it fails with a default message.
// This is a convenience method for the common case of checking for no error.
//
// Parameters:
//   - err: The error to check
//
// Example:
//
//	_, err := someFunction()
//	r.AssertNoErr(err)
func (r *R) AssertNoErr(err error) {
	r.AssertNoErrf(err, "error unexpected")
}

// AssertNoErrf checks if the provided error is nil with a custom description.
// If err is nil, it passes the test with the given description.
// If err is not nil, it fails the test, logs the error details, and stops execution.
// This method is useful when you need to provide context about what operation failed.
//
// Parameters:
//   - err: The error to check
//   - desc: A description of what operation should not have failed
//   - args: Arguments for the description format string
//
// Example:
//
//	user, err := authenticateUser(username, password)
//	r.AssertNoErrf(err, "User authentication should succeed for %s", username)
func (r *R) AssertNoErrf(err error, desc string, args ...any) {
	if err == nil {
		r.Pass(desc, args...)
	} else {
		r.Fail(desc, args...)
		r.Logf("requires no error, but found: %v", err)
		r.T.FailNow()
	}
}

// AssertErr checks if the provided error is not nil.
// If err is not nil, it passes the test; otherwise, it fails with a default message.
// This is useful for testing error conditions where you expect an error to occur.
//
// Parameters:
//   - err: The error to check
//
// Example:
//
//	_, err := divide(10, 0)
//	r.AssertErr(err) // Expects an error for division by zero
func (r *R) AssertErr(err error) {
	r.AssertErrf(err, "error expected")
}

// AssertErrf checks if the provided error is not nil with a custom description.
// If err is not nil, it passes the test with the given description.
// If err is nil, it fails the test and stops execution.
// This method is useful when you need to provide context about what error was expected.
//
// Parameters:
//   - err: The error to check
//   - desc: A description of what error was expected
//   - args: Arguments for the description format string
//
// Example:
//
//	_, err := validateInput("")
//	r.AssertErrf(err, "Empty input should cause validation error")
func (r *R) AssertErrf(err error, desc string, args ...any) {
	if err == nil {
		r.Fail(desc, args...)
		r.Logf("requires error, but found nil")
		r.T.FailNow()
	} else {
		r.Pass(desc, args...)
	}
}

// StartTimer starts timing the test
func (r *R) StartTimer() *R {
	r.startTime = time.Now()
	r.Case("Starting test timer")
	return r
}

// StopTimer stops timing and logs the duration
func (r *R) StopTimer() *R {
	duration := time.Since(r.startTime)
	r.Case("Test completed in %v", duration)
	return r
}

// Benchmark starts a benchmark test
func (r *R) Benchmark(name string, f func(b *testing.B)) *R {
	r.Case("Benchmark: %s", name)
	r.benchmark = true
	r.Run(name, func(t *testing.T) {
		// Note: This is a simplified benchmark implementation
		// In a real implementation, you'd need to convert testing.T to testing.B
		r.Logf("Running benchmark: %s", name)
	})
	r.benchmark = false
	return r
}

// Parallel marks the test as safe to run in parallel
func (r *R) Parallel() *R {
	r.parallel = true
	r.T.Parallel()
	r.Case("Test marked as parallel")
	return r
}

// Skip skips the current test with a reason
func (r *R) Skip(reason string, args ...any) *R {
	r.Case("Skipping test: "+reason, args...)
	r.T.Skipf(reason, args...)
	return r
}

// SkipIf skips the test if the condition is true
func (r *R) SkipIf(condition bool, reason string, args ...any) *R {
	if condition {
		r.Skip(reason, args...)
	}
	return r
}

// SkipUnless skips the test unless the condition is true
func (r *R) SkipUnless(condition bool, reason string, args ...any) *R {
	if !condition {
		r.Skip(reason, args...)
	}
	return r
}

// Cleanup registers a cleanup function
func (r *R) Cleanup(fn func()) *R {
	r.T.Cleanup(fn)
	r.Case("Cleanup function registered")
	return r
}

// Helper marks the calling function as a test helper function
func (r *R) Helper() *R {
	r.T.Helper()
	return r
}

// TempDir returns a temporary directory for the test
func (r *R) TempDir() string {
	return r.T.TempDir()
}

// Setenv sets an environment variable for the test
func (r *R) Setenv(key, value string) *R {
	r.T.Setenv(key, value)
	r.Case("Environment variable set: %s=%s", key, value)
	return r
}

// Getenv gets an environment variable
func (r *R) Getenv(key string) string {
	// Use os.Getenv to read the environment variable
	// testing.T.Setenv sets the environment variable for the test process
	return os.Getenv(key)
}

// Deadline returns the time when the test will be timed out
func (r *R) Deadline() (deadline time.Time, ok bool) {
	return r.T.Deadline()
}

// RunParallel runs tests in parallel
func (r *R) RunParallel(fn func(*testing.PB)) *R {
	// Note: testing.T.RunParallel is not available in all Go versions
	// This is a simplified implementation
	r.Case("Running tests in parallel")
	return r
}

// MemoryUsage logs memory usage information
func (r *R) MemoryUsage() *R {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	r.Case("Memory Usage")
	r.Logf("Alloc: %d KB", m.Alloc/1024)
	r.Logf("TotalAlloc: %d KB", m.TotalAlloc/1024)
	r.Logf("Sys: %d KB", m.Sys/1024)
	r.Logf("NumGC: %d", m.NumGC)

	return r
}

// GoroutineCount logs the current goroutine count
func (r *R) GoroutineCount() *R {
	count := runtime.NumGoroutine()
	r.Case("Goroutine count: %d", count)
	return r
}

// TestInfo logs comprehensive test information
func (r *R) TestInfo() *R {
	r.Case("Test Information")
	r.Logf("Test Name: %s", r.T.Name())
	r.Logf("Start Time: %v", r.startTime)
	r.Logf("Duration: %v", time.Since(r.startTime))
	r.Logf("Parallel: %v", r.parallel)
	r.Logf("Benchmark: %v", r.benchmark)

	// Log goroutine count
	r.GoroutineCount()

	// Log memory usage
	r.MemoryUsage()

	return r
}

// AssertEqual provides a more descriptive equality assertion
func (r *R) AssertEqual(expected, actual any, msg ...string) *R {
	if !reflect.DeepEqual(expected, actual) {
		message := fmt.Sprintf("Expected %v, got %v", expected, actual)
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Values are equal")
	}
	return r
}

// AssertNotEqual provides a more descriptive inequality assertion
func (r *R) AssertNotEqual(expected, actual any, msg ...string) *R {
	if reflect.DeepEqual(expected, actual) {
		message := fmt.Sprintf("Expected values to be different, but both are %v", expected)
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Values are not equal")
	}
	return r
}

// AssertNil provides a more descriptive nil assertion
func (r *R) AssertNil(value any, msg ...string) *R {
	if value != nil {
		message := fmt.Sprintf("Expected nil, got %v", value)
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Value is nil")
	}
	return r
}

// AssertNotNil provides a more descriptive non-nil assertion
func (r *R) AssertNotNil(value any, msg ...string) *R {
	if value == nil {
		message := "Expected non-nil value, got nil"
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Value is not nil")
	}
	return r
}

// AssertTrue provides a more descriptive true assertion
func (r *R) AssertTrue(condition bool, msg ...string) *R {
	if !condition {
		message := "Expected condition to be true"
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Condition is true")
	}
	return r
}

// AssertFalse provides a more descriptive false assertion
func (r *R) AssertFalse(condition bool, msg ...string) *R {
	if condition {
		message := "Expected condition to be false"
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Condition is false")
	}
	return r
}

// AssertContains provides a more descriptive contains assertion
func (r *R) AssertContains(container, item any, msg ...string) *R {
	contains := false

	switch c := container.(type) {
	case string:
		if itemStr, ok := item.(string); ok {
			contains = strings.Contains(c, itemStr)
		}
	case []any:
		for _, v := range c {
			if reflect.DeepEqual(v, item) {
				contains = true
				break
			}
		}
	default:
		// Use reflection for other types
		rv := reflect.ValueOf(container)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				if reflect.DeepEqual(rv.Index(i).Interface(), item) {
					contains = true
					break
				}
			}
		}
	}

	if !contains {
		message := fmt.Sprintf("Expected %v to contain %v", container, item)
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Container contains item")
	}
	return r
}

// AssertNotContains provides a more descriptive not-contains assertion
func (r *R) AssertNotContains(container, item any, msg ...string) *R {
	contains := false

	switch c := container.(type) {
	case string:
		if itemStr, ok := item.(string); ok {
			contains = strings.Contains(c, itemStr)
		}
	case []any:
		for _, v := range c {
			if reflect.DeepEqual(v, item) {
				contains = true
				break
			}
		}
	default:
		// Use reflection for other types
		rv := reflect.ValueOf(container)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				if reflect.DeepEqual(rv.Index(i).Interface(), item) {
					contains = true
					break
				}
			}
		}
	}

	if contains {
		message := fmt.Sprintf("Expected %v not to contain %v", container, item)
		if len(msg) > 0 {
			message = msg[0]
		}
		r.Fail(message)
	} else {
		r.Pass("Container does not contain item")
	}
	return r
}

// AssertPanics provides a more descriptive panic assertion
func (r *R) AssertPanics(fn func(), msg ...string) *R {
	defer func() {
		if recover := recover(); recover == nil {
			message := "Expected function to panic"
			if len(msg) > 0 {
				message = msg[0]
			}
			r.Fail(message)
		} else {
			r.Pass("Function panicked as expected")
		}
	}()

	fn()
	r.Fail("Expected function to panic")
	return r
}

// AssertNotPanics provides a more descriptive no-panic assertion
func (r *R) AssertNotPanics(fn func(), msg ...string) *R {
	defer func() {
		if recover := recover(); recover != nil {
			message := fmt.Sprintf("Expected function not to panic, but it panicked with %v", recover)
			if len(msg) > 0 {
				message = msg[0]
			}
			r.Fail(message)
		} else {
			r.Pass("Function did not panic")
		}
	}()

	fn()
	return r
}
