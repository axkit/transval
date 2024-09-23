// Package transval provides features to validate transitions between states.
package transval

import (
	"errors"
	"strconv"
	"strings"
)

// Splitter constants used to define the delimiters for parsing transition rules.
var (
	SingleRuleSplitter = ";"  // Delimiter for splitting multiple rules.
	TransitionSplitter = "=>" // Delimiter for separating 'from' and 'to' states.
	EndSplitter        = ","  // Delimiter for separating target states.
)

// transitionRules stores the original rule string and a map of transitions.
type transitionRules struct {
	original    string        // The original string of transition rules.
	transitions map[int][]int // A map where keys are 'from' states and values are slices of 'to' states.
}

// TransVal is the main structure that holds all the transition rules.
type TransVal struct {
	m map[string]transitionRules // Map of named transition rules, keyed by the name.
}

// Error variables for common error conditions.
var (
	ErrWrongInput  = errors.New("input is invalid")      // Error when the input is not formatted correctly.
	ErrEmptyInput  = errors.New("input is empty")        // Error when the input string is empty.
	ErrTargetEmpty = errors.New("target value is empty") // Error when the target state value is missing.
)

// New creates and returns a new TransVal instance.
// It initializes an empty map to store transition rules.
func New() *TransVal {
	return &TransVal{m: make(map[string]transitionRules)}
}

// Del deletes the transition rules associated with the given name.
// Parameters:
//   - name: The name of the transition rules to delete.
func (t *TransVal) Del(name string) {
	delete(t.m, name)
}

// Set adds or updates transition rules under the specified name.
// Parameters:
//   - name: The name under which the transition rules are stored.
//   - rules: The string representation of the transition rules to be parsed.
//
// Returns:
//   - error: Returns an error if the rule parsing fails, otherwise nil.
func (t *TransVal) Set(name string, rules string) error {

	pv, err := parse(rules)
	if err != nil {
		return err
	}
	r, ok := t.m[name]
	if !ok {
		r.original = rules
		r.transitions = make(map[int][]int)
	}

	for k, v := range pv {
		r.transitions[k] = append(r.transitions[k], v...)
	}
	t.m[name] = r
	return nil
}

// parse parses the input string of transition rules into a map.
// Parameters:
//   - rules: The string representing transition rules, formatted as "from=>to1,to2;from2=>to3"
//
// Returns:
//   - map[int][]int: A map where the key is the 'from' state, and the value is a slice of 'to' states.
//   - error: Returns an error if the input string is invalid.
func parse(rules string) (map[int][]int, error) {

	if len(rules) == 0 {
		return map[int][]int{}, nil
	}

	res := make(map[int][]int)

	parts := strings.Split(rules, SingleRuleSplitter)

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}

		kv := strings.Split(p, TransitionSplitter)
		if len(kv) != 2 {
			return nil, ErrWrongInput
		}

		kv[0] = strings.TrimSpace(kv[0])
		k, err := strconv.Atoi(kv[0])
		if err != nil {
			return nil, ErrWrongInput
		}

		kv[1] = strings.TrimSpace(kv[1])
		if len(kv[1]) == 0 {
			return nil, ErrTargetEmpty
		}

		v, err := stringToIntSlice(kv[1])
		if err != nil {
			return nil, err
		}
		if len(v) == 0 {
			return nil, ErrTargetEmpty
		}
		res[k] = append(res[k], v...)
	}
	return res, nil
}

// IsTransitionValid checks whether a transition from one state to another is allowed.
// Parameters:
//   - name: The name of the transition rule set.
//   - from: The starting state.
//   - to: The target state.
//
// Returns:
//   - bool: Returns true if the transition is allowed, false otherwise.
func (e *TransVal) IsTransitionValid(name string, from, to int) bool {

	r, ok := e.m[name]
	if !ok {
		return false
	}

	v, ok := r.transitions[from]
	if !ok {
		return false
	}

	for _, vto := range v {
		if vto == to {
			return true
		}
	}

	return false
}

// stringToIntSlice converts a comma-separated string into a slice of integers.
// Parameters:
//   - s: The string to convert, e.g., "1,2,3".
//
// Returns:
//   - []int: A slice of integers.
//   - error: Returns an error if any part of the string cannot be converted to an integer.
func stringToIntSlice(s string) ([]int, error) {
	res := []int{}

	nums := strings.Split(s, EndSplitter)

	for _, n := range nums {
		n = strings.TrimSpace(n)
		if len(n) == 0 {
			continue
		}

		v, err := strconv.Atoi(n)
		if err != nil {
			return nil, ErrWrongInput
		}
		res = append(res, v)
	}
	return res, nil
}

// AllowedTo returns the list of states that can be transitioned to from a given state.
// Parameters:
//   - name: The name of the transition rule set.
//   - from: The starting state.
//
// Returns:
//   - []int: A slice of valid target states, or nil if no transitions are found.
func (t *TransVal) AllowedTo(name string, from int) []int {
	r, ok := t.m[name]
	if !ok {
		return nil
	}

	v, ok := r.transitions[from]
	if !ok {
		return nil
	}

	return v
}

// Transitions returns all transition rules for the given name.
// Parameters:
//   - name: The name of the transition rule set.
//
// Returns:
//   - map[int][]int: A map of transitions, or nil if the name is not found.
func (t *TransVal) Transitions(name string) map[int][]int {
	r, ok := t.m[name]
	if !ok {
		return nil
	}

	return r.transitions
}
