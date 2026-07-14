package syntax

/*
ostr.go contains ASN.1 OCTET STRING types and methods.
*/

import (
	"errors"
	"strconv"
	"unicode"
)

/*
OctetString implements [§ 3.3.25 of RFC 4517]:

	OctetString = *OCTET

[§ 3.3.25 of RFC 4517]: https://datatracker.ietf.org/doc/html/rfc4517#section-3.3.25
*/
type OctetString []byte

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r OctetString) IsZero() bool { return r == nil }

/*
String returns the string representation of the receiver instance.
*/
func (r OctetString) String() string {
	return string(r)
}

/*
Len returns the integer length of the receiver instance.
*/
func (r OctetString) Len() int { return len(r) }

/*
OctetString returns an instance of [OctetString] alongside an error
following an analysis of x in the context of an Octet String.
*/
func NewOctetString(x any) (OctetString, error) {
	return marshalOctetString(x)
}

func octetString(x any) (result bool, err error) {
	_, err = marshalOctetString(x)
	result = err == nil
	return
}

func marshalOctetString(x any) (oct OctetString, err error) {
	var raw []byte
	if raw, err = assertOctetString(x); err != nil {
		return
	}

	runes := []rune(string(raw))
	for i := 0; i < len(runes) && err == nil; i++ {
		var char rune = runes[i]
		// octet range is simply IA5 chars
		if !unicode.In(char, iA5Range) {
			err = errors.New("Incompatible Octet String character: " + strconv.Itoa(int(char)))
		}
	}

	if err == nil {
		oct = OctetString(raw)
	}

	return
}

func assertOctetString(in any) (raw []byte, err error) {
	switch tv := in.(type) {
	case []byte:
		raw = tv
	case OctetString:
		raw = []byte(tv)
	case string:
		raw = []byte(tv)
	default:
		err = errorBadType("OctetStringMatch")
	}

	return
}

/*
octetStringMatch implements [§ 4.2.27 of RFC 4517].

OID: 2.5.13.17.

[§ 4.2.27 of RFC 4517]: https://datatracker.ietf.org/doc/html/rfc4517#section-4.2.27
*/
func octetStringMatch(a, b any) (result bool, err error) {

	var A, B []byte
	if A, err = assertOctetString(a); err != nil {
		return
	}

	if B, err = assertOctetString(b); err != nil {
		return
	}

	var res bool
	if res = len(A) == len(B); res {
		for i, ch := range B {
			if res = A[i] == ch; !res {
				break
			}
		}
	}

	result = res

	return
}

func octetStringOrderingMatch(a any, operator byte, b any) (result bool, err error) {
	var str1, str2 []byte

	if str1, err = assertOctetString(a); err != nil {
		return
	}

	if str2, err = assertOctetString(b); err != nil {
		return
	}

	mLen := len(str2)
	if len(str1) < mLen {
		mLen = len(str1)
	}

	// Compare octet strings from the first octet to the last
	for i := 0; i < mLen; i++ {
		if operator == GreaterOrEqual {
			if str2[i] < str1[i] {
				result = true
				return
			} else if str2[i] > str1[i] {
				result = false
				return
			}
		} else {
			if str1[i] < str2[i] {
				result = true
				return
			} else if str1[i] > str2[i] {
				result = false
				return
			}
		}
	}

	// If the strings are identical up to the length of the
	// shorter string, the shorter string precedes the longer
	// string
	if operator == GreaterOrEqual {
		result = len(str2) < len(str1)
	} else {
		result = len(str2) > len(str1)
	}

	return
}
