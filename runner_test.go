package got_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/gophero/got"
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
	tr.NoErr(nil)
}

func TestNoErrf(t *testing.T) {
	tr := got.New(t, "test no errf")
	tr.NoErrf(nil, "no error expected")
}

func TestErr(t *testing.T) {
	tr := got.New(t, "test err")
	err := &struct{ error }{}
	defer func() {
		if r := recover(); r == nil {
			// NoErrf will call t.FailNow if error is not nil
		}
	}()
	tr.Err(err)
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
