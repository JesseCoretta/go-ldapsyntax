package syntax

import (
	"testing"
)

func TestUniversalString(t *testing.T) {
	for _, raw := range []string{
		`平仮名`,
		`This is a UniversalString.`,
		`This is@~@@~~~ not UniversalString ﺝﺦﺕﺣﺛ^\^\rOH WAIT yes it is`,
	} {
		if _, err := NewUniversalString(raw); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		}
	}
}

func TestUniversalString_codecov(t *testing.T) {
	_ = universalString(`This is@~@@~~~ not UniversalString ﺝﺦﺕﺣﺛ^\^\rOH WAIT yes it is`)
}
