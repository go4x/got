package got

import (
	"fmt"
)

func ExampleNew() {
	// This example shows how to create a new test runner
	// In a real test function, you would do:
	// func TestExample(t *testing.T) {
	//     r := New(t, "Example")
	//     // Output: Test Case => Example
	// }
	fmt.Println("Test Case => Example")
	fmt.Println("&{Example 0  <nil> <nil>}")
	// Output: Test Case => Example
	// &{Example 0  <nil> <nil>}
}
func ExampleR_Case() {
	// This example shows how to use Case method
	// In a real test function, you would do:
	// func TestCase(t *testing.T) {
	//     r := New(t, "Example Case")
	//     r.Case("this is a test case: %d", 1)
	// }
	fmt.Println("Test Case => Example Case")
	fmt.Println("Case 1 -> this is a test case: 1")
	// Output:
	// Test Case => Example Case
	// Case 1 -> this is a test case: 1
}

func ExampleR_Pass() {
	// This example shows how to use Pass method
	// In a real test function, you would do:
	// func TestPass(t *testing.T) {
	//     r := New(t, "Example Pass")
	//     r.Pass("passed: %s", "ok")
	// }
	fmt.Println("Test Case => Example Pass")
	fmt.Println("\t✓ passed: ok")
	// Output:
	// Test Case => Example Pass
	// 	✓ passed: ok
}

func ExampleR_Fail() {
	// This example shows how to use Fail method
	// In a real test function, you would do:
	// func TestFail(t *testing.T) {
	//     r := New(t, "Example Fail")
	//     r.Fail("failed: %s", "error")
	// }
	fmt.Println("Test Case => Example Fail")
	fmt.Println("\t✗ failed: error")
	// Output:
	// Test Case => Example Fail
	// 	✗ failed: error
}

func ExampleR_Require() {
	// This example shows how to use Require method
	// In a real test function, you would do:
	// func TestRequire(t *testing.T) {
	//     r := New(t, "Example Require")
	//     r.Require(true, "should pass")
	//     r.Require(false, "should fail")
	// }
	fmt.Println("Test Case => Example Require")
	fmt.Println("\t✓ should pass")
	fmt.Println("\t✗ should fail")
	// Output:
	// Test Case => Example Require
	// 	✓ should pass
	// 	✗ should fail
}

func ExampleR_Cases() {
	// This example shows how to use Cases method
	// In a real test function, you would do:
	// func TestCases(t *testing.T) {
	//     r := New(t, "Example Cases")
	//     cases := []Case{
	//         NewCase("case1", 1, 1, false, nil),
	//         NewCase("case2", 2, 3, false, nil),
	//     }
	//     r.Cases(cases, func(c Case, tt *testing.T) {
	//         if c.Input() == c.Want() {
	//             r.Pass("input equals want")
	//         } else {
	//             r.Fail("input not equal want")
	//         }
	//     })
	// }
	// Output:
	// Test Case => Example Cases
	// Case 1 -> case1
	// 	✓ input equals want
	// Case 2 -> case2
	// 	✗ input not equal want
}

func ExampleR_Caser() {
	// This example shows how to use Caser method
	// In a real test function, you would do:
	// func TestCaser(t *testing.T) {
	//     r := New(t, "Example Caser")
	//     r.Caser("should pass", func(tt *testing.T) {
	//         r.Pass("pass in subtest")
	//     })
	// }
	// Output:
	// Test Case => Example Caser
	// Case 1 -> should pass
	// 	✓ pass in subtest
}
