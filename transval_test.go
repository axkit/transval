package transval

import (
	"testing"
)

func TestTransVal_Set(t *testing.T) {

	var tests = []struct {
		name  string
		rules string
		ok    bool
	}{
		{"a", "", false},
		{"b", "1", false},
		{"c", "1;", false},
		{"d", "1;;", false},
		{"e", ";1;", false},
		{"f", ";1;;", false},
		{"g", ";;1;;", false},
		{"j", ";;;1;;;", false},
		{"e1", "1=>", false},
		{"e2", "1=>;", false},
		{"e3", "=>;", false},
		{"e4", "=>", false},
		{"e5", "=> ", false},
		{"e6", "=>;=>", false},
		{"e7", "=>=>", false},
		{"e8", "1=>2=>3", false},
		{"e8", "1=>2;=>3", false},
		{"ok1", "1=>2", true},
		{"ok2", "1=>2,3", true},
		{"ok3", "1=>2,3;", true},
		{"ok3", "1=>2, 3", true},
		{"ok3", "1=> 2,3;", true},
		{"ok3", "1=> 2, 3 ;", true},
		{"ok3", "1=> 2, 3, 4; 2=>1,3;6=>7,8;7=>6,8", true},
		{"ok4", "1=>2,3,4;2=>1,3;6=>7,8;7=>6,8", true},
	}

	tr := New()
	for i := range tests {
		err := tr.Set(tests[i].name, tests[i].rules)

		if (err != nil && tests[i].ok == false) || (err == nil && tests[i].ok == true) {
			continue
		}
		t.Errorf("case %d (%s) failed: %v", i, tests[i].rules, err)
	}

}

func TestTranVal_IsTransitionValid(t *testing.T) {

	rules := "1=>2,3,4;2=>1,3;6=>7,8;7=>6,8"
	var tests = []struct {
		from int
		to   int
		ok   bool
	}{
		{1, 2, true},
		{1, 1, false},
		{2, 1, true},
		{1, 10, false},
		{10, 1, false},
		{6, 7, true},
		{6, 8, true},
	}

	tr := New()
	err := tr.Set("name", rules)
	if err != nil {
		t.Fatal(err)
	}

	for i := range tests {
		if tr.IsTransitionValid("name", tests[i].from, tests[i].to) != tests[i].ok {
			t.Errorf("case %d failed", i)
		}
	}
}
