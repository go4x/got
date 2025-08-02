package got

import (
	"fmt"
	"testing"
)

func ExampleNew() {
	t := testing.T{}
	r := New(&t, "Example")
	fmt.Println(r)
	// Output: Test Case => Example
	// &{Example 0  <nil> <nil>}
}
func ExampleR_Case() {
	t := testing.T{}
	r := New(&t, "Example Case")
	r.Case("this is a test case: %d", 1)
	// Output:
	// Test Case => Example Case
	// Case 1 -> this is a test case: 1
}

func ExampleR_Pass() {
	t := testing.T{}
	r := New(&t, "Example Pass")
	r.Pass("passed: %s", "ok")
	// Output:
	// Test Case => Example Pass
	// 	✓ passed: ok
}

func ExampleR_Fail() {
	t := testing.T{}
	r := New(&t, "Example Fail")
	r.Fail("failed: %s", "error")
	// Output:
	// Test Case => Example Fail
	// 	✗ failed: error
}

func ExampleR_Require() {
	t := testing.T{}
	r := New(&t, "Example Require")
	r.Require(true, "should pass")
	r.Require(false, "should fail")
	// Output:
	// Test Case => Example Require
	// 	✓ should pass
	// 	✗ should fail
}

func ExampleR_Cases() {
	t := testing.T{}
	r := New(&t, "Example Cases")
	cases := []Case{
		NewCase("case1", 1, 1, false, nil),
		NewCase("case2", 2, 3, false, nil),
	}
	r.Cases(cases, func(c Case, tt *testing.T) {
		if c.Input() == c.Want() {
			r.Pass("input equals want")
		} else {
			r.Fail("input not equal want")
		}
	})
	// Output:
	// Test Case => Example Cases
	// Case 1 -> case1
	// 	✓ input equals want
	// Case 2 -> case2
	// 	✗ input not equal want
}

func ExampleR_Caser() {
	t := testing.T{}
	r := New(&t, "Example Caser")
	r.Caser("should pass", func(tt *testing.T) {
		r.Pass("pass in subtest")
	})
	// Output:
	// Test Case => Example Caser
	// Case 1 -> should pass
	// 	✓ pass in subtest
}
