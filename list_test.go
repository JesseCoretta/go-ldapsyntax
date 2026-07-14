package syntax

import (
	"testing"
)

func TestList(t *testing.T) {
	result, err := caseIgnoreListMatch(
		[]string{
			`this`, `is`, `a`, `list`,
		},
		[]string{
			`this`, `is`, `a`, `list`,
		})

	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	} else if !result {
		t.Errorf("%s failed:\nwant: %t\ngot:  %t",
			t.Name(), true, result)
	}

	result, err = caseIgnoreListMatch(
		[]string{
			`this`, `iz`, `a`, `list`,
		},
		[]string{
			`this`, `is`, `a`, `list`,
		})

	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	} else if result {
		t.Errorf("%s failed:\nwant: %t\ngot:  %t",
			t.Name(), false, result)
	}

	_, _ = caseIgnoreListMatch(nil, nil)
	_, _ = caseIgnoreListMatch([]string{}, nil)
	_, _ = caseIgnoreListMatch([]string{}, []string{`a`})
	_, _ = caseIgnoreListMatch(nil, struct{}{})
	_, _ = caseIgnoreListMatch(struct{}{}, nil)
}
