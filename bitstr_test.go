package syntax

import (
	"testing"
)

func TestBitString(t *testing.T) {
	var raw string = `'10100101'B`
	if result, _ := bitString(raw); !result {
		t.Errorf("%s failed:\nwant: %t\ngot:  %t",
			t.Name(), true, result)
		return
	}

	bs, err := NewBitString(raw)
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
	} else if bs.IsZero() {
		t.Errorf("%s failed: instance is zero", t.Name())
	} else if got := bs.String(); raw != got {
		t.Errorf("%s failed:\nwant: %s\ngot:  %s",
			t.Name(), raw, got)
	}
}

func TestBitString_codecov(t *testing.T) {
	_, _ = assertBitString([]byte{})
	_, _ = assertBitString(struct{}{})
	_, _ = bitString([]byte{})
	_, _ = bitString(struct{}{})

	_, _ = bitStringMatch([]byte{}, struct{}{})
	_, _ = bitStringMatch([]byte(`'010110'B`), struct{}{})
	_, _ = bitStringMatch([]byte{}, []byte{})
	_, _ = bitStringMatch(struct{}{}, struct{}{})
	_, _ = bitStringMatch([]byte(`'010110'B`), []byte(`'01'B`))

	b, err := bitStringMatch(`'1010100'B`, `'1010000'B`)
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	} else if b {
		t.Errorf("%s failed:\nwant: %t\ngot:  %t",
			t.Name(), false, b)
		return
	}

	b, err = bitStringMatch(`'1010100'B`, `'1010100'B`)
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	} else if !b {
		t.Errorf("%s failed:\nwant: %t\ngot:  %t",
			t.Name(), true, b)
		return
	}

	_, _ = bitStringMatch(`'1010100'B`, `'10101'B`)
	_ = stripTrailingZeros([]byte{0x1, 0x2, 0x0, 0x0}, 2)
	_ = stripTrailingZeros([]byte{0x1, 0x2, 0x0, 0x0}, 4)
	_ = stripTrailingZeros([]byte{}, 0)

}
