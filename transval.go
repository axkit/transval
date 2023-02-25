// The package provides features
package transval

import (
	"errors"
	"strconv"
	"strings"
)

var (
	SingleRuleSplitter = ";"
	TransitionSplitter = "=>"
	EndSplitter        = ","
)

type transitionRules struct {
	original    string
	transitions map[int][]int
}

type TransVal struct {
	m map[string]transitionRules
}

var (
	ErrWrongInput  = errors.New("input is invalid")
	ErrEmptyInput  = errors.New("input is empty")
	ErrTargetEmpty = errors.New("target value is empty")
)

// New creates Transiter instance.
func New() *TransVal {
	return &TransVal{m: make(map[string]transitionRules)}
}

func (t *TransVal) Del(name string) {
	delete(t.m, name)
}

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

func parse(rules string) (map[int][]int, error) {

	if len(rules) == 0 {
		return nil, nil
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
			return nil, err
		}

		kv[1] = strings.TrimSpace(kv[1])
		if len(kv[1]) == 0 {
			return nil, ErrTargetEmpty
		}

		v, err := stringToIntSlice(kv[1])
		if err != nil {
			return nil, err
		}
		res[k] = append(res[k], v...)
	}
	return res, nil
}

// IsValidTransition returns true if transition valid.
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

// stringToIntSlice converts string containing comma separated numbers
// to []int.
func stringToIntSlice(s string) ([]int, error) {
	res := []int{}

	nums := strings.Split(s, ",")

	for _, n := range nums {
		n = strings.TrimSpace(n)
		if len(n) == 0 {
			continue
		}

		v, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

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

func (t *TransVal) Transitions(name string) map[int][]int {
	r, ok := t.m[name]
	if !ok {
		return nil
	}

	return r.transitions
}
