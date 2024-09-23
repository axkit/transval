package transval

import (
	"testing"
)

func TestNew(t *testing.T) {
	trans := New()
	if trans == nil {
		t.Fatal("Expected New() to return a non-nil TransVal instance")
	}
}

func TestSetAndTransitions(t *testing.T) {
	trans := New()
	err := trans.Set("test", "1=>2,3;2=>3")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	transitions := trans.Transitions("test")
	if len(transitions) != 2 {
		t.Fatalf("Expected 2 transitions, but got %d", len(transitions))
	}

	if len(transitions[1]) != 2 || transitions[1][0] != 2 || transitions[1][1] != 3 {
		t.Errorf("Unexpected transitions for 1, got: %v", transitions[1])
	}

	if len(transitions[2]) != 1 || transitions[2][0] != 3 {
		t.Errorf("Unexpected transitions for 2, got: %v", transitions[2])
	}
}

func TestSetInvalidInput(t *testing.T) {
	trans := New()
	err := trans.Set("test", "invalid=>input")
	if err != ErrWrongInput {
		t.Fatalf("Expected ErrWrongInput, but got %v", err)
	}
}

func TestDel(t *testing.T) {
	trans := New()
	err := trans.Set("test", "1=>2,3;2=>3")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	trans.Del("test")
	transitions := trans.Transitions("test")
	if transitions != nil {
		t.Fatalf("Expected nil after deletion, but got %v", transitions)
	}
}

func TestIsValidTransition(t *testing.T) {
	trans := New()
	err := trans.Set("test", "1=>2,3;2=>3")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if !trans.IsTransitionValid("test", 1, 2) {
		t.Errorf("Expected transition from 1 to 2 to be valid")
	}

	if trans.IsTransitionValid("test", 1, 4) {
		t.Errorf("Expected transition from 1 to 4 to be invalid")
	}

	if trans.IsTransitionValid("unknown", 1, 2) {
		t.Errorf("Expected transition for unknown name to be invalid")
	}
}

func TestAllowedTo(t *testing.T) {
	trans := New()
	err := trans.Set("test", "1=>2,3;2=>3")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	allowed := trans.AllowedTo("test", 1)
	if len(allowed) != 2 || allowed[0] != 2 || allowed[1] != 3 {
		t.Errorf("Expected allowed transitions from 1 to be [2, 3], got %v", allowed)
	}

	allowed = trans.AllowedTo("test", 4)
	if allowed != nil {
		t.Errorf("Expected no allowed transitions for 4, got %v", allowed)
	}

	allowed = trans.AllowedTo("unknown", 1)
	if allowed != nil {
		t.Errorf("Expected no allowed transitions for unknown name, got %v", allowed)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		rules    string
		expected map[int][]int
		err      error
	}{
		{"1=>2,3;2=>3", map[int][]int{1: {2, 3}, 2: {3}}, nil},
		{"", nil, nil},
		{";", nil, nil},
		{"1=>,", nil, ErrTargetEmpty},
		{"1", nil, ErrWrongInput},
		{"1=>", nil, ErrTargetEmpty},
		{"invalid=>rule", nil, ErrWrongInput},
	}

	for _, tt := range tests {
		result, err := parse(tt.rules)
		if err != tt.err {
			t.Errorf("Expected error %v, but got %v", tt.err, err)
		}
		if len(result) != len(tt.expected) {
			t.Errorf("Expected result size %d, but got %d", len(tt.expected), len(result))
		}
	}
}

func TestStringToIntSlice(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
		err      error
	}{
		{"1,2,3", []int{1, 2, 3}, nil},
		{"", []int{}, nil},
		{"1, ,2", []int{1, 2}, nil},
		{"a,b,c", nil, ErrWrongInput},
	}

	for _, tt := range tests {
		result, err := stringToIntSlice(tt.input)
		if err != tt.err {
			t.Errorf("Expected error %v, but got %v", tt.err, err)
		}
		if len(result) != len(tt.expected) {
			t.Errorf("Expected result size %d, but got %d", len(tt.expected), len(result))
		}
	}
}
