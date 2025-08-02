package got

import (
	"testing"
)

func TestCaseBuilder(t *testing.T) {
	// Test building a case with all fields set
	c := CaseBuilder("test case").
		Input(123).
		Want(456).
		WantErr(true).
		Err(nil).
		Build()

	if c.Name() != "test case" {
		t.Errorf("expected Name to be 'test case', got '%s'", c.Name())
	}
	if c.Input() != 123 {
		t.Errorf("expected Input to be 123, got %v", c.Input())
	}
	if c.Want() != 456 {
		t.Errorf("expected Want to be 456, got %v", c.Want())
	}
	if !c.WantErr() {
		t.Errorf("expected WantErr to be true, got false")
	}
	if c.Err() != nil {
		t.Errorf("expected Err to be nil, got %v", c.Err())
	}

	// Test default values when not set
	c2 := CaseBuilder("default case").Build()
	if c2.Name() != "default case" {
		t.Errorf("expected Name to be 'default case', got '%s'", c2.Name())
	}
	if c2.Input() != nil {
		t.Errorf("expected Input to be nil, got %v", c2.Input())
	}
	if c2.Want() != nil {
		t.Errorf("expected Want to be nil, got %v", c2.Want())
	}
	if c2.WantErr() {
		t.Errorf("expected WantErr to be false, got true")
	}
	if c2.Err() != nil {
		t.Errorf("expected Err to be nil, got %v", c2.Err())
	}

	// Test setting error
	testErr := &struct{ error }{}
	c3 := CaseBuilder("error case").Err(testErr).Build()
	if c3.Err() != testErr {
		t.Errorf("expected Err to be %v, got %v", testErr, c3.Err())
	}

	// Test changing fields after initial set
	cb := CaseBuilder("change case").Input(1).Want(2)
	cb.Input(10).Want(20).WantErr(true)
	c4 := cb.Build()
	if c4.Input() != 10 {
		t.Errorf("expected Input to be 10, got %v", c4.Input())
	}
	if c4.Want() != 20 {
		t.Errorf("expected Want to be 20, got %v", c4.Want())
	}
	if !c4.WantErr() {
		t.Errorf("expected WantErr to be true, got false")
	}

	// Test Name method on builder
	cb2 := CaseBuilder("init").Name("renamed")
	c5 := cb2.Build()
	if c5.Name() != "renamed" {
		t.Errorf("expected Name to be 'renamed', got '%s'", c5.Name())
	}
}
