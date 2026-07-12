package syntax

import (
	"testing"
)

func TestCountryString(t *testing.T) {

	for _, raw := range []string{
		`US`,
		`CA`,
		`UK`,
		`JP`,
	} {
		if cs, err := NewCountryString(raw); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
		} else if got := cs.String(); raw != got {
			t.Errorf("%s failed:\nwant: %s\ngot:  %s",
				t.Name(), raw, got)
		}
	}

	NewCountryString(nil)
	NewCountryString([]byte{})
	NewCountryString(``)
}

func TestCountryString_codecov(t *testing.T) {
	countryString(`US`)
	countryString(`FR`)
	countryString(`FRANCE`)
}
