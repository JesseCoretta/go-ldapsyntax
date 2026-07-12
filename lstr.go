package syntax

/*
lstr.go contains LDAPString types and methods.
*/

/*
LDAPString aliases [OctetString] to implement [§ 4.1.2 of RFC 4511].

[§ 4.1.2 of RFC 4511]: https://datatracker.ietf.org/doc/html/rfc4511#section-4.1.2
*/
type LDAPString OctetString

/*
LDAPString returns an instance of [LDAPString] alongside an error.
*/
func NewLDAPString(x ...any) (LDAPString, error) {
	return marshalLDAPString(x...)
}

func marshalLDAPString(x ...any) (ls LDAPString, err error) {
	if len(x) > 0 {
		switch tv := x[0].(type) {
		case OctetString, string, []byte:
			var o OctetString
			o, err = marshalOctetString(tv)
			ls = LDAPString(o)
		default:
			err = errorBadType("LDAPString")
		}
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r LDAPString) String() string { return string(r) }
