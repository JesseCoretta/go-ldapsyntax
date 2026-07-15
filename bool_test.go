package syntax

import (
	"testing"
)

func TestNewBoolean(t *testing.T) {
	for _, b := range []any{
		true, false,
		`true`, `TRUE`, `false`, `FALSE`, `True`, `False`,
		0, 1,
		byte(0x00), byte(0x01),
	} {
		if _, err := NewBoolean(b); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		}
	}

	// coverage
	_, _ = NewBoolean(struct{}{})
	_, _ = boolean(true)
	_, _ = boolean(`falsch`)
	_, _ = boolean(nil)
	_, _ = boolean(byte(0x02))
	_, _ = boolean(-2)

	_, _ = booleanMatch(struct{}{}, true)
	_, _ = booleanMatch(true, struct{}{})
	_, _ = booleanMatch(false, true)
	_, _ = booleanMatch(nil, true)

	if result, err := booleanMatch(`TRUE`, false); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	} else if result {
		t.Errorf("%s failed:\nwant: %t\ngot:  %t", t.Name(), false, result)
		return
	}
}
