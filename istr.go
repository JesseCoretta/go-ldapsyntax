package syntax

import (
	"errors"
	"unicode"
)

/*
IA5String implements [§ 3.2 of RFC 4517]:

	IA5 = 0x0000 through 0x00FF

[§ 3.2 of RFC 4517]: https://datatracker.ietf.org/doc/html/rfc4517#section-3.2
*/
type IA5String string

/*
String returns the string representation of the receiver instance.
*/
func (r IA5String) String() string {
	return string(r)
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r IA5String) IsZero() bool { return len(r) == 0 }

/*
IA5String returns an instance of [IA5String] alongside an error following
an analysis of x in the context of an IA5 String.
*/
func NewIA5String(x any) (ia5 IA5String, err error) {
	return marshalIA5String(x)
}

func iA5String(x any) (result bool, err error) {
	_, err = marshalIA5String(x)
	result = err == nil
	return
}

func marshalIA5String(x any) (ia5 IA5String, err error) {
	var raw string
	if raw, err = assertString(x, 1, "IA5String"); err == nil {
		if err = checkIA5String(raw); err == nil {
			ia5 = IA5String(raw)
		}
	}

	return
}

func checkIA5String(raw string) (err error) {
	if len(raw) == 0 {
		err = errors.New("Invalid IA5 String (zero)")
		return
	}

	runes := []rune(raw)
	for i := 0; i < len(runes) && err == nil; i++ {
		var char rune = runes[i]
		if !unicode.Is(iA5Range, char) {
			err = errors.New("Invalid IA5 String character: " + string(char))
		}
	}

	return
}

var iA5Range *unicode.RangeTable

func init() {
	iA5Range = &unicode.RangeTable{R16: []unicode.Range16{
		{0x0000, 0x00FF, 1},
	}}
}
