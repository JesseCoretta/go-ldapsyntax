package syntax

import (
	"errors"
	"strconv"
)

/*
NewBoolean returns a Boolean value alongside an error following
an analysis of input argument x as a boolean.

If x is a bool, it is guaranteed to be valid and is returned as-is.

If x is a string, an underlying call to [strconv.ParseBool] is made.

If x is a byte, only values of zero (0x00) for false, or one (0x01)
for true, are considered valid. Any other byte value is an error.

If x is an int, only values of zero (0) for false, or one (1)
for true, are considered valid. Any other int value is an error.

Any other input type is an error.
*/
func NewBoolean(x any) (b bool, err error) {
	switch tv := x.(type) {
	case bool:
		b = tv
	case string:
		b, err = strconv.ParseBool(tv)
	case byte:
		if b = tv == 0x01; !b && tv != 0x00 {
			err = errors.New("Invalid bool byte; want 0x00 or 0x01")
		}
	case int:
		if b = tv == 1; !b && tv != 0 {
			err = errors.New("Invalid bool integer; want 0 or 1")
		}
	default:
		err = errorBadType("boolean")
	}
	return
}

func boolean(x any) (result bool, err error) {
	_, err = NewBoolean(x)
	result = err == nil
	return
}

func booleanMatch(realValue, assertionValue any) (result bool, err error) {
	var a, b bool
	if a, err = NewBoolean(realValue); err == nil {
		if b, err = NewBoolean(assertionValue); err == nil {
			result = a == b
		}
	}

	if err != nil {
		err = errors.New("UNDEFINED: " + err.Error())
	}

	return
}
