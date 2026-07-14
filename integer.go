package syntax

/*
nf.go contains methods and types for expressing X.680 number forms.
*/

import (
	"errors"
	"math"
	"math/big"
	"strconv"
)

/*
Integer implements the unbounded ASN.1 INTEGER type for use in the
context of X.680 number forms, which are present within [NameAndInteger],
[DotNotation] and [ASN1Notation] type instances.

Note that *[big.Int] is used internally ONLY if the number overflows uint64.

For safety reasons (with respect to ambiguity of default values), a zero
instance of this type is bogus.  Users MUST use the [NewInteger] or
[MustNewInteger] constructor to obtain valid instances of this type.
*/
type Integer struct {
	big, ok bool
	native  int64    // Stores native signed integer values when possible
	bigInt  *big.Int // Stores big.Int values only when necessary
}

/*
NewInteger returns an instance of [Integer] alongside an error
following an attempt to marshal x as an unbounded ASN.1 integer.

Input types may be int, int32, int64, uint64, string or *[big.Int].

Any signed magnitude is permitted. Values which underflow or overflow
int64 are stored as *[big.Int].
*/
func NewInteger[T any](x T) (i Integer, err error) {
	i, err = assertInteger(x)
	return
}

func assertInteger[T any](v T) (i Integer, err error) {
	switch value := any(v).(type) {
	case int:
		i = Integer{native: int64(value)}
	case int64:
		i = Integer{native: value}
	case uint64:
		i = uint64ToInteger(value)
	case *big.Int:
		i = bigToInteger(value)
	case int32:
		i = Integer{native: int64(value)}
	case string:
		i, err = strToInteger(value)
	case Integer:
		if !value.ok {
			err = errorIntNil
		}
		i = value
	default:
		err = errorIntBadType
	}

	if err == nil {
		i.ok = true
	}

	return
}

func integer(x any) (result bool, err error) {
	_, err = NewInteger(x)
	result = err == nil
	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r Integer) IsZero() bool { return &r == nil }

/*
String returns the string representation of the receiver instance.
*/
func (r Integer) String() string {
	var s string
	if r.big {
		s = r.bigInt.String()
	} else {
		s = strconv.FormatInt(r.native, 10)
	}

	return s
}

/*
IsBig returns a Boolean value indicative of the underlying value
overflowing uint64.
*/
func (r Integer) IsBig() bool { return r.big }

/*
Native returns the underlying int64 value found within the receiver
instance. Note that this method should not be used unless a call of
[Integer.IsBig] beforehand returns false.
*/
func (r Integer) Native() int64 { return r.native }

/*
Valid returns a Boolean value indicative of the receiver instance
being properly initialized via the [NewInteger] or [MustNewInteger]
constructor with an unambiguous (non-default) value.
*/
func (r Integer) Valid() bool { return r.ok }

/*
Big returns the *[big.Int] form of the receiver instance.

Note that use of this method constructs an entirely new instance of
*[big.Int] if the underlying value is an int64.  Thus, this method
should only usually be needed if a call to [Integer.IsBig] returns
true. In that case, the preexisting *[big.Int] value is returned, as
opposed to being generated on the fly.

When [Integer.IsBig] returns false, the return instance of *[big.Int]
is entirely independent of the receiver and does not replace the
underlying value. This can be useful, though potentially costly, in
cases where methods extended by *[big.Int] that are not wrapped in
this package directly need to be accessed for some reason.
*/
func (r Integer) Big() (i *big.Int) {
	if r.big {
		i = r.bigInt
	} else {
		i = newBigInt(0).SetInt64(r.native)
	}

	return
}

/*
Eq returns a bool indicative of an equality match between the
receiver instance and x.
*/
func (r Integer) Eq(x any) bool { return r.cmpAny(x) == 0 }

/*
Ne returns a bool indicative of a negative equality match between
the receiver instance and x.
*/
func (r Integer) Ne(x any) bool { return r.cmpAny(x) != 0 }

/*
Gt returns a bool indicative of r being greater than x.
*/
func (r Integer) Gt(x any) bool { return r.cmpAny(x) > 0 }

/*
Ge returns a bool indicative of r being greater than or equal to x.
*/
func (r Integer) Ge(x any) bool { return r.cmpAny(x) >= 0 }

/*
Lt returns a bool indicative of r being less than x.
*/
func (r Integer) Lt(x any) bool { return r.cmpAny(x) < 0 }

/*
Le returns a bool indicative of r being less than or equal to x.
*/
func (r Integer) Le(x any) bool { return r.cmpAny(x) <= 0 }

func (r Integer) cmpAny(x any) (result int) {
	switch t := x.(type) {
	case Integer:
		result = cmpInteger(r, t)

	case int:
		result = r.cmpInt64(int64(t))

	case int32:
		result = r.cmpInt64(int64(t))

	case int64:
		result = r.cmpInt64(t)

	case uint64:
		result = r.cmpUint64(t)

	case string:
		result = r.cmpIntegerStr(t)

	case *big.Int:
		result = r.cmpBig(t)

	default:
		panic("Integer: unsupported type for comparison")
	}

	return
}

func (r Integer) cmpIntegerStr(v string) int {
	nf, err := NewInteger(v)
	if err != nil {
		panic(err)
	}
	return cmpInteger(r, nf)
}

func cmpInteger(a, b Integer) int {
	if !a.big && !b.big {
		switch {
		case a.native < b.native:
			return -1
		case a.native > b.native:
			return +1
		default:
			return 0
		}
	}
	return a.Big().Cmp(b.Big())
}

func (r Integer) cmpInt64(v int64) int {
	if !r.big {
		switch {
		case r.native < v:
			return -1
		case r.native > v:
			return +1
		default:
			return 0
		}
	}
	return r.Big().Cmp(big.NewInt(v))
}

func (r Integer) cmpUint64(u uint64) int {
	if !r.big && u <= math.MaxInt64 {
		return r.cmpInt64(int64(u))
	}
	b := newBigInt(0).SetUint64(u)
	return r.Big().Cmp(b)
}

func (r Integer) cmpBig(b *big.Int) int {
	if !r.big {
		return newBigInt(0).SetInt64(r.native).Cmp(b)
	}
	return r.bigInt.Cmp(b)
}

/*
IsInteger returns a Boolean value indicative of the nf string
input value representing a valid [Integer], in that:

  - The number is one (1) or more valid digits, and ...
  - The number is base10 (e.g.: not octal)

Assuming the above requirements are satisfied, any unsigned magnitude
is considered valid.
*/
func IsInteger(i string) bool {
	return integerCheck(i) == nil
}

func integerCheck(num string) (err error) {
	if len(num) == 0 {
		err = errorIntNoInput
		return
	}

	if num[0] == '-' {
		num = num[1:]
	}

	if len(num) > 1 && num[0] == '0' {
		err = errorIntOctal
		return
	}

	for i := 0; i < len(num); i++ {
		if ch := num[i]; !('0' <= ch && ch <= '9') {
			err = errorIntNaN
			break
		}
	}

	return
}

func strToInteger(num string) (i Integer, err error) {
	if err = integerCheck(num); err != nil {
		return
	}

	_i, _ := newBigInt(0).SetString(num, 10)
	if _i.IsInt64() {
		i = Integer{native: _i.Int64()}
	} else {
		i = Integer{big: true, bigInt: _i}
	}

	return
}

func bigToInteger(num *big.Int) (i Integer) {
	if i.big = !num.IsInt64(); i.big {
		i.bigInt = num
	} else {
		i.native = num.Int64()
	}

	return
}

func uint64ToInteger(num uint64) (i Integer) {
	if i.big = num > uint64(math.MaxInt64); i.big {
		i.bigInt = newBigInt(0).SetUint64(num)
	} else {
		i.native = int64(num)
	}

	return
}

var (
	errorIntNil     = errors.New("INTEGER: nil or bogus instance")
	errorIntBadType = errors.New("INTEGER: unsupported input type")
	errorIntNoInput = errors.New("INTEGER: nil or zero input")
	errorIntOctal   = errors.New("INTEGER: leading zeroes (octal numbers) prohibited")
	errorIntNaN     = errors.New("INTEGER: non numeric character found")
)
