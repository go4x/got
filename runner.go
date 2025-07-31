package got

import (
	"strconv"
	"testing"
)

const (
	checkMark = "\033[32m\u2713\033[0m" // green √
	ballotX   = "\033[31m\u2717\033[0m" // red ×
)

type R struct {
	title   string
	caseNum int
	prefix  string
	*testing.T
}

// New creates a new T instance from a testing.T.
func New(t *testing.T, title string) *R {
	t.Log("Test Case => " + title)
	return &R{T: t, title: title}
}

// Case is a method to start a new test case, the format parameter will describe the case's purpose,
// given condition, and expecting results.
func (r *R) Case(format string, args ...any) *R {
	r.caseNum++
	r.prefix = "Case " + strconv.Itoa(r.caseNum) + " -> "
	r.Logf(r.prefix+format, args...)
	return r
}

// Caser runs a test case with the given name and function.
// It first logs the case using Case, then executes the subtest using Run.
func (r *R) Caser(name string, f func(t *testing.T)) *R {
	r.Case(name)
	r.Run(name, f)
	return r
}

// Run executes a subtest with the given name and function.
// It wraps testing.T.Run to provide a convenient way to run subtests.
func (r *R) Run(name string, f func(t *testing.T)) *R {
	r.T.Run(name, func(tt *testing.T) {
		f(tt)
	})
	return r
}

// Cases runs a set of test cases, executing the provided function for each case.
// For each test case, it sets up the case description and runs the function as a subtest.
//
// The cases parameter must implement the Namer interface with Name() method.
func (r *R) Cases(cases []Case, f func(c Case, tt *testing.T)) {
	for _, c := range cases {
		r.Case(c.Name())
		r.Run(c.Name(), func(tt *testing.T) {
			f(c, tt)
		})
	}
}

// Pass is a method who invoking indicate that current condition matches the expecting result.
func (r *R) Pass(format string, args ...any) {
	r.Logf("\t%s "+format, prependTag(checkMark, args...)...)
}

// Fail is an opposite method to Pass, it indicate that current condition does not match the expecting result.
func (r *R) Fail(format string, args ...any) {
	r.Errorf("\t%s "+format, prependTag(ballotX, args...)...)
}

// Fatal is a method that immediately fails the test and stops its execution,
// logging the provided message with a failure mark.
func (r *R) Fatal(format string, args ...any) {
	prependTag(ballotX, args...)
	r.Fatalf("\t%s "+format, prependTag(ballotX, args...)...)
}

func prependTag(tag any, args ...any) []any {
	if args == nil {
		args = make([]any, 1)
		args[0] = tag
	} else {
		args = append([]any{tag}, args...)
	}
	return args
}

// Require is a convenient method for Pass and Fail, it requires the given bool parameter named "cond" to be true,
// if so it will invoke Pass, otherwise invoke Fail.
func (r *R) Require(cond bool, desc string, args ...any) {
	if cond {
		r.Pass(desc, args...)
	} else {
		r.Fail(desc, args...)
	}
}

// FailNow checks the given condition "cond".
// If true, it logs a pass message with the provided description and arguments.
// If false, it logs a fail message and immediately stops the test execution.
func (r *R) FailNow(cond bool, desc string, args ...any) {
	if cond {
		r.Pass(desc, args...)
	} else {
		r.Fail(desc, args...)
		r.T.FailNow()
	}
}

// NoErr checks if the provided error is nil.
// If err is nil, it passes the test; otherwise, it fails with a default message.
func (r *R) NoErr(err error) {
	r.NoErrf(err, "error unexpected")
}

// NoErrf checks if the provided error is nil.
// If err is nil, it passes the test with the given description and arguments.
// If err is not nil, it fails the test, logs the error, and immediately stops the test execution.
func (r *R) NoErrf(err error, desc string, args ...any) {
	if err == nil {
		r.Pass(desc, args...)
	} else {
		r.Fail(desc, args...)
		r.Logf("requires no error, but found: %v", err)
		r.T.FailNow()
	}
}

func (r *R) Err(err error) {
	r.Errf(err, "error expected")
}

func (r *R) Errf(err error, desc string, args ...any) {
	if err == nil {
		r.Fail(desc, args...)
		r.Logf("requires error, but found nil")
		r.T.FailNow()
	} else {
		r.Pass(desc, args...)
	}
}
