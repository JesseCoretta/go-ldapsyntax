package syntax

import (
	"encoding/hex"
	"strconv"
	"strings"
)

/*
AssertionValue implements an OCTET STRING value.
*/
type AssertionValue []uint8

/*
Set assigns x to the receiver instance.
*/
func (r *AssertionValue) Set(x any) {
        var s string
        switch tv := x.(type) {
        case string:
                s = tv
        case []byte:
                s = string(tv)
        default:
                return
        }

        *r = AssertionValue(escapeString(s))
}

/*
String returns the string representation of the receiver instance.
Note that this method is an alias of [AssertionValue.Escaped].
*/
func (r AssertionValue) String() string {
        return r.Escaped()
}

/*
Unescaped returns the unescaped receiver value. For example, "ジェシー"
is returned instead of "\e3\82\b8\e3\82\a7\e3\82\b7\e3\83\bc".
*/
func (r AssertionValue) Unescaped() string {
        var u string
        if len(r) > 0 {
                u = hexDecode(string(r))
        }

        return u
}

func (r AssertionValue) Escaped() (esc string) {
        if len(r) > 0 {
                esc = escapeString(string(r))
        }

        return
}

func escapeString(x string) (esc string) {
        if len(x) > 0 {
                bld := &strings.Builder{}
                for _, z := range x {
                        if z > maxASCII {
                                for _, c := range []byte(string(z)) {
                                        bld.WriteString(`\`)
                                        bld.WriteString(strconv.FormatUint(uint64(c), 16))
                                }
                        } else {
                                bld.WriteRune(z)
                        }
                }  

                esc = bld.String()
        }

        return     
}

func hexDecode(x any) string {
        var r string
        switch tv := x.(type) {
        case string:
                r = tv
        case []byte:
                r = string(tv)
        default:   
                return ``
        }

        d := &strings.Builder{}
        length := len(r)

        for i := 0; i < length; i++ {
                if r[i] == '\\' && i+3 <= length {
                        b, err := hex.DecodeString(r[i+1 : i+3])
                        if err != nil || !(isHex(rune(r[i+1])) || isHex(rune(r[i+2]))) {
                                return ``
                        }
                        d.Write(b)
                        i += 2
                } else {
                        d.WriteString(string(r[i]))
                }
        }

        return d.String()
}
