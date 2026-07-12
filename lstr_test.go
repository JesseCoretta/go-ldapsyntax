package syntax

import (
	"testing"
)

func TestLDAPString(t *testing.T) {
	for _, lstring := range []any{
		"cn",
		[]byte("cn"),
		OctetString("cn"),
	} {
		if _, err := NewLDAPString(lstring); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
			return
		}
	}
}
