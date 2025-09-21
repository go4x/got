package got_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go4x/got"
)

func TestNew(t *testing.T) {
	tr := got.New(t, "test new runner")
	tr.Case("new runner from testing.T")
	tr.Require(tr != nil, "new runner should be success")
}

func TestNewWithTitle(t *testing.T) {
	tr := got.New(t, "test new runner with title")
	tr.Case("new runner with title")
	tr.Require(tr != nil, "new runner with title should be success")
	tr.Run("test title", func(t *testing.T) {
		t.Log("test title")
	})
}

func TestCase(t *testing.T) {
	tr := got.New(t, "test case")
	ret := tr.Case("case description: %d", 1)
	if ret != tr {
		t.Error("Case should return the same runner instance")
	}
}

func TestRunMethod(t *testing.T) {
	tr := got.New(t, "test run")
	called := false
	tr.Run("test run", func(tt *testing.T) {
		called = true
	})
	if !called {
		t.Error("Run should execute the provided function")
	}
}

func TestCasesMethod(t *testing.T) {
	tr := got.New(t, "test cases")
	cases := []got.Case{
		got.NewCase("case1", nil, nil, false, nil),
		got.NewCase("case2", nil, nil, false, nil),
	}
	count := 0
	tr.Cases(cases, func(c got.Case, tt *testing.T) {
		count++
	})
	if count != 2 {
		t.Errorf("Cases should run function for each case, got %d", count)
	}
}

func TestPass(t *testing.T) {
	tr := got.New(t, "test pass")
	tr.Pass("pass message: %s", "ok")
}

func TestRequire(t *testing.T) {
	tr := got.New(t, "test require")
	tr.Require(true, "should pass")
}

func TestFailNow(t *testing.T) {
	tr := got.New(t, "test fail now")
	tr.FailNow(true, "should pass")
}

func TestNoErr(t *testing.T) {
	tr := got.New(t, "test no err")
	tr.AssertNoErr(nil)
}

func TestNoErrf(t *testing.T) {
	tr := got.New(t, "test no errf")
	tr.AssertNoErrf(nil, "no error expected")
}

func TestErr(t *testing.T) {
	tr := got.New(t, "test err")
	err := &struct{ error }{}
	defer func() {
		if r := recover(); r == nil {
			// NoErrf will call t.FailNow if error is not nil
		}
	}()
	tr.AssertErr(err)
}

func TestRun(t *testing.T) {
	cases := []got.Case{
		got.NewCase("t1", "a", "a", false, nil),
		got.NewCase("t2", 1, 1, false, nil),
	}
	equals := func(a, b any) bool {
		return reflect.DeepEqual(a, b)
	}
	tr := got.New(t, "test run")

	tr.Cases(cases, func(c got.Case, tt *testing.T) {
		if equals(c.Input(), c.Want()) {
			tr.Pass("the same")
		} else {
			tr.Fail("not the same")
		}
	})
}

func TestCases(t *testing.T) {
	var err = errors.New("div by zero")
	cases := []got.Case{
		got.NewCase("correct test", []int{1, 2}, 0, false, nil),
		got.NewCase("error test, div by zero", []int{1, 0}, 0, true, err),
	}
	tr := got.New(t, "test cases")
	div := func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("%w, b is zero", err)
		}
		return a / b, nil
	}
	tr.Cases(cases, func(c got.Case, tt *testing.T) {
		ret, err := div(c.Input().([]int)[0], c.Input().([]int)[1])
		if c.WantErr() {
			tr.Require(err != nil && errors.Is(err, c.Err()), "error expected")
		} else {
			tr.Require(err == nil && reflect.DeepEqual(ret, c.Want()), "the same")
		}
	})
}

func TestCaser(t *testing.T) {
	tr := got.New(t, "test Caser")
	tr.Caser("should pass when values are equal", func(tt *testing.T) {
		a, b := 2, 2
		if a == b {
			tr.Pass("values are equal: %d == %d", a, b)
		} else {
			tr.Fail("values are not equal: %d != %d", a, b)
		}
	})

	tr.Caser("should fail when values are not equal", func(tt *testing.T) {
		a, b := 2, 3
		if a == b {
			tr.Fail("values are equal: %d == %d", a, b)
		} else {
			tr.Pass("values are not equal: %d != %d", a, b)
		}
	})
}

func TestCase1(t *testing.T) {
	r := got.New(t, "test case1")
	r.Case("case1").Run("test case1", func(tt *testing.T) {
		r.Pass("test case1")
	}).Run("test case2", func(tt *testing.T) {
		r.Pass("test case2")
	})
	r.Case("case2").Run("test case3", func(tt *testing.T) {
		r.Pass("test case3")
	}).Run("test case4", func(tt *testing.T) {
		r.Pass("test case4")
	})
}

// TestEnhancedRunner demonstrates the enhanced runner capabilities
func TestEnhancedRunner(t *testing.T) {
	er := got.New(t, "Enhanced Runner Test")

	er.StartTimer()
	er.Case("Testing enhanced runner capabilities")

	// Test basic assertions
	er.AssertEqual(5, 5, "Numbers should be equal")
	er.AssertNotEqual(5, 6, "Numbers should not be equal")
	er.AssertTrue(true, "Condition should be true")
	er.AssertFalse(false, "Condition should be false")
	er.AssertNil(nil, "Value should be nil")
	er.AssertNotNil("hello", "Value should not be nil")

	// Test contains assertions
	er.AssertContains("hello world", "world", "String should contain substring")
	er.AssertNotContains("hello world", "foo", "String should not contain substring")

	// Test panic assertions
	er.AssertPanics(func() {
		panic("test panic")
	}, "Function should panic")

	er.AssertNotPanics(func() {
		// Do nothing
	}, "Function should not panic")

	// Test environment variables
	er.Setenv("TEST_VAR", "test_value")
	value := er.Getenv("TEST_VAR")
	er.AssertEqual("test_value", value, "Environment variable should be set")

	// Test cleanup
	er.Cleanup(func() {
		er.Log("Cleanup function called")
	})

	// Test helper function
	er.Helper()

	// Test temp directory
	tempDir := er.TempDir()
	er.AssertNotEqual("", tempDir, "Temp directory should be created")

	// Test deadline
	deadline, ok := er.Deadline()
	if ok {
		er.Logf("Test deadline: %v", deadline)
	}

	// Test info
	er.TestInfo()

	er.StopTimer()
}

// TestEnhancedRunnerParallel demonstrates parallel execution
func TestEnhancedRunnerParallel(t *testing.T) {
	er := got.New(t, "Enhanced Runner Parallel Test")

	er.Parallel()
	er.Case("Testing parallel execution")

	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	er.AssertTrue(true, "Parallel test should pass")
}

// TestEnhancedRunnerSkip demonstrates test skipping
func TestEnhancedRunnerSkip(t *testing.T) {
	er := got.New(t, "Enhanced Runner Skip Test")

	er.Case("Testing skip functionality")
	er.SkipIf(true, "Skipping test because condition is true")
	er.SkipUnless(false, "Skipping test because condition is false")

	// This should not be reached
	er.AssertTrue(false, "This should not be executed")
}

// TestEnhancedRunnerBenchmark demonstrates benchmarking
func TestEnhancedRunnerBenchmark(t *testing.T) {
	er := got.New(t, "Enhanced Runner Benchmark Test")

	er.Case("Testing benchmark functionality")
	er.Benchmark("Test Benchmark", func(b *testing.B) {
		// Benchmark implementation would go here
		for i := 0; i < b.N; i++ {
			// Do some work
			_ = i * i
		}
	})
}

// TestEnhancedRunnerComplex demonstrates complex scenarios
func TestEnhancedRunnerComplex(t *testing.T) {
	er := got.New(t, "Complex Enhanced Test")

	er.StartTimer()
	er.Case("Testing complex scenarios")

	// Test with environment setup
	er.Setenv("COMPLEX_TEST", "true")
	er.Setenv("TEST_NUM", "42")

	// Test assertions with complex data
	testData := []string{"hello", "world", "test"}
	er.AssertContains(testData, "hello", "Slice should contain 'hello'")
	er.AssertNotContains(testData, "foo", "Slice should not contain 'foo'")

	// Test error handling
	err := errors.New("test error")
	er.AssertErrf(err, "Expected error should be present")

	// Test memory and goroutine info
	er.MemoryUsage()
	er.GoroutineCount()

	// Test info
	er.TestInfo()

	er.StopTimer()
}

// TestEnhancedRunnerRunParallel demonstrates parallel test execution
func TestEnhancedRunnerRunParallel(t *testing.T) {
	er := got.New(t, "Enhanced Runner RunParallel Test")

	er.Case("Testing RunParallel functionality")

	// This would run in parallel
	er.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Do some work
			time.Sleep(1 * time.Millisecond)
		}
	})
}

// TestEnhancedRunnerCleanup demonstrates cleanup functionality
func TestEnhancedRunnerCleanup(t *testing.T) {
	er := got.New(t, "Enhanced Runner Cleanup Test")

	cleanupCalled := false
	er.Cleanup(func() {
		cleanupCalled = true
		er.Log("Cleanup function called")
	})

	er.Case("Testing cleanup functionality")
	er.AssertTrue(true, "Test should pass")

	// Cleanup will be called automatically when test ends
	er.AssertFalse(cleanupCalled, "Cleanup should not be called yet")
}

// TestEnhancedRunnerHelper demonstrates helper functionality
func TestEnhancedRunnerHelper(t *testing.T) {
	er := got.New(t, "Enhanced Runner Helper Test")

	er.Case("Testing helper functionality")
	er.Helper() // Mark as helper function

	helperFunction := func() {
		er.Helper() // Mark helper function
		er.AssertTrue(true, "Helper function should work")
	}

	helperFunction()
}

// TestEnhancedRunnerTempDir demonstrates temporary directory usage
func TestEnhancedRunnerTempDir(t *testing.T) {
	er := got.New(t, "Enhanced Runner TempDir Test")

	er.Case("Testing temporary directory")
	tempDir := er.TempDir()
	er.AssertNotEqual("", tempDir, "Temp directory should be created")
	er.Logf("Temp directory: %s", tempDir)
}

// TestEnhancedRunnerSetenv demonstrates environment variable setting
func TestEnhancedRunnerSetenv(t *testing.T) {
	er := got.New(t, "Enhanced Runner Setenv Test")

	er.Case("Testing environment variable setting")
	er.Setenv("TEST_VAR", "test_value")

	value := er.Getenv("TEST_VAR")
	er.AssertEqual("test_value", value, "Environment variable should be set")
}

// TestEnhancedRunnerDeadline demonstrates test deadline
func TestEnhancedRunnerDeadline(t *testing.T) {
	er := got.New(t, "Enhanced Runner Deadline Test")

	er.Case("Testing deadline functionality")
	deadline, ok := er.Deadline()

	if ok {
		er.Logf("Test deadline: %v", deadline)
		er.AssertTrue(deadline.After(time.Now()), "Deadline should be in the future")
	} else {
		er.Log("No deadline set")
	}
}

// TestFailNow_Failure tests the FailNow method when condition is false
func TestFailNow_Failure(t *testing.T) {
	r := got.New(t, "Test FailNow Failure")
	r.Case("Testing FailNow with false condition")

	// This should fail and stop the test
	r.FailNow(false, "This should fail and stop")
	// If we reach here, the test didn't stop as expected
	t.Error("FailNow should have stopped test execution")
}

// TestNoErrf_WithError tests NoErrf when an error is present
func TestNoErrf_WithError(t *testing.T) {
	r := got.New(t, "Test NoErrf With Error")
	r.Case("Testing NoErrf with error")

	// This should fail and stop the test
	err := errors.New("test error")
	r.AssertNoErrf(err, "This should fail because error is present")
	// If we reach here, the test didn't stop as expected
	t.Error("NoErrf should have stopped test execution when error is present")
}

// TestErrf_WithoutError tests Errf when no error is present
func TestErrf_WithoutError(t *testing.T) {
	r := got.New(t, "Test Errf Without Error")
	r.Case("Testing Errf without error")

	// This should fail and stop the test
	r.AssertErrf(nil, "This should fail because no error is present")
	// If we reach here, the test didn't stop as expected
	t.Error("Errf should have stopped test execution when no error is present")
}

// TestRequire_False tests Require when condition is false
func TestRequire_False(t *testing.T) {
	r := got.New(t, "Test Require False")
	r.Case("Testing Require with false condition")

	// This should fail but not stop the test
	r.Require(false, "This condition should be false")
}

// TestFail_Logging tests the Fail method logging
func TestFail_Logging(t *testing.T) {
	r := got.New(t, "Test Fail Logging")
	r.Case("Testing Fail method")

	// This should log a failure message
	r.Fail("This is a test failure message")
}

// TestFatal_Logging tests the Fatal method logging
func TestFatal_Logging(t *testing.T) {
	r := got.New(t, "Test Fatal Logging")
	r.Case("Testing Fatal method")

	// This should log a fatal message and stop the test
	r.Fatal("This is a fatal error message")
	// If we reach here, the test didn't stop as expected
	t.Error("Fatal should have stopped test execution")
}

// TestCase_Chaining tests method chaining with Case
func TestCase_Chaining(t *testing.T) {
	r := got.New(t, "Test Case Chaining")

	// Test that Case returns the same instance for chaining
	result := r.Case("First case").Case("Second case")
	if result != r {
		t.Error("Case should return the same runner instance for chaining")
	}
}

// TestCaser_Chaining tests method chaining with Caser
func TestCaser_Chaining(t *testing.T) {
	r := got.New(t, "Test Caser Chaining")

	// Test that Caser returns the same instance for chaining
	result := r.Caser("Test case", func(t *testing.T) {
		// Empty test function
	})
	if result != r {
		t.Error("Caser should return the same runner instance for chaining")
	}
}

// TestRun_Chaining tests method chaining with Run
func TestRun_Chaining(t *testing.T) {
	r := got.New(t, "Test Run Chaining")

	// Test that Run returns the same instance for chaining
	result := r.Run("Test subtest", func(t *testing.T) {
		// Empty test function
	})
	if result != r {
		t.Error("Run should return the same runner instance for chaining")
	}
}

// TestComplexTestCase demonstrates a complex test scenario
func TestComplexTestCase(t *testing.T) {
	r := got.New(t, "Complex Test Scenario")

	// Test multiple cases with different scenarios
	r.Case("Testing string operations")

	// Test string length
	str := "hello world"
	r.Require(len(str) == 11, "String length should be 11")

	// Test string concatenation
	str1 := "hello"
	str2 := " world"
	result := str1 + str2
	r.Require(result == "hello world", "Concatenated string should be 'hello world'")

	// Test error handling
	r.Case("Testing error handling")
	err := errors.New("test error")
	r.AssertErrf(err, "Error should be present")

	// Test successful operation
	r.Case("Testing successful operation")
	r.AssertNoErrf(nil, "No error should be present")
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	r := got.New(t, "Edge Cases Test")

	// Test with empty string
	r.Case("Testing empty string")
	r.Require(len("") == 0, "Empty string should have length 0")

	// Test with nil values
	r.Case("Testing nil handling")
	var nilSlice []string
	r.Require(nilSlice == nil, "Nil slice should be nil")

	// Test with zero values
	r.Case("Testing zero values")
	var zero int
	r.Require(zero == 0, "Zero value should be 0")
}
