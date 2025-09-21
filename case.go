// Package got provides a comprehensive testing framework for Go applications.
package got

// Namer is an interface for objects that have a name.
// This is used by the Case interface to provide naming functionality.
type Namer interface {
	Name() string
}

// Case defines the structure of a test case for table-driven tests.
// It provides a standardized way to define test inputs, expected outputs,
// and error conditions for multiple test scenarios.
//
// The interface includes:
//   - Name(): A descriptive name for the test case
//   - Input(): The input data for the test
//   - Want(): The expected output/result
//   - WantErr(): Whether the test case should produce an error
//   - Err(): The specific error expected (if WantErr() is true)
//
// Example:
//
//	case := got.NewCase("Valid Login", "user@example.com", true, false, nil)
//	// case.Name() returns "Valid Login"
//	// case.Input() returns "user@example.com"
//	// case.Want() returns true
//	// case.WantErr() returns false
//	// case.Err() returns nil
type Case interface {
	Namer

	Input() any    // the input of the test case
	Want() any     // the expected result of the test case
	WantErr() bool // whether the test case should return an error
	Err() error    // the error of the test case
}

// caseImpl is the default implementation of the Case interface.
// It stores all the test case data in a simple struct format.
type caseImpl struct {
	name    string
	input   any
	want    any
	wantErr bool
	err     error
}

// Name returns the name of the test case.
func (c *caseImpl) Name() string {
	return c.name
}

// Input returns the input data for the test case.
func (c *caseImpl) Input() any {
	return c.input
}

// Want returns the expected output for the test case.
func (c *caseImpl) Want() any {
	return c.want
}

// WantErr returns whether the test case should produce an error.
func (c *caseImpl) WantErr() bool {
	return c.wantErr
}

// Err returns the specific error expected for this test case.
func (c *caseImpl) Err() error {
	return c.err
}

// NewCase creates a new test case with the provided parameters.
// This is the simplest way to create a test case for table-driven tests.
//
// Parameters:
//   - name: A descriptive name for the test case
//   - input: The input data for the test
//   - want: The expected output/result
//   - wantErr: Whether the test case should produce an error
//   - err: The specific error expected (if wantErr is true)
//
// Returns:
//   - Case: A new test case instance
//
// Example:
//
//	case := got.NewCase("Valid Input", "hello", 5, false, nil)
//	case := got.NewCase("Invalid Input", "", 0, true, errors.New("empty input"))
func NewCase(name string, input any, want any, wantErr bool, err error) Case {
	return &caseImpl{name: name, input: input, want: want, wantErr: wantErr, err: err}
}

// CaseBuilder creates a new case builder for fluent test case construction.
// This provides a builder pattern for creating test cases with method chaining.
//
// Parameters:
//   - name: The initial name for the test case
//
// Returns:
//   - *caseBuilder: A new case builder instance
//
// Example:
//
//	case := got.CaseBuilder("Test Case").
//		Input("hello").
//		Want(5).
//		WantErr(false).
//		Err(nil).
//		Build()
func CaseBuilder(name string) *caseBuilder {
	return &caseBuilder{caseImpl: caseImpl{name: name}}
}

// caseBuilder provides a fluent interface for building test cases.
// It implements the builder pattern to allow method chaining.
type caseBuilder struct {
	caseImpl
}

// Name sets the name of the test case and returns the builder for chaining.
func (b *caseBuilder) Name(name string) *caseBuilder {
	b.name = name
	return b
}

// Input sets the input data for the test case and returns the builder for chaining.
func (b *caseBuilder) Input(input any) *caseBuilder {
	b.input = input
	return b
}

// Want sets the expected output for the test case and returns the builder for chaining.
func (b *caseBuilder) Want(want any) *caseBuilder {
	b.want = want
	return b
}

// WantErr sets whether the test case should produce an error and returns the builder for chaining.
func (b *caseBuilder) WantErr(wantErr bool) *caseBuilder {
	b.wantErr = wantErr
	return b
}

// Err sets the specific error expected for the test case and returns the builder for chaining.
func (b *caseBuilder) Err(err error) *caseBuilder {
	b.err = err
	return b
}

// Build creates the final Case instance from the builder.
// This method should be called at the end of the builder chain.
//
// Returns:
//   - Case: The completed test case instance
func (b *caseBuilder) Build() Case {
	return &b.caseImpl
}
