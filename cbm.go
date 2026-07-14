package syntax

import (
	"strings"
)

func caseIgnoreMatch(a, b any) (result bool, err error) {
	result, err = caseBasedMatch(a, b, false)
	return
}

func caseExactMatch(a, b any) (result bool, err error) {
	result, err = caseBasedMatch(a, b, true)
	return
}

func caseBasedMatch(a, b any, caseExact bool) (result bool, err error) {
	var str1, str2 string
	str1, err = assertString(a, 1, "string")
	if err != nil {
		return
	}

	str2, err = assertString(b, 1, "string")
	if err != nil {
		return
	}

	if caseExact {
		result = str1 == str2
	} else {
		result = strings.EqualFold(str1, str2)
	}

	return
}

func caseIgnoreOrderingMatch(a any, operator byte, b any) (bool, error) {
	return caseBasedOrderingMatch(a, b, false, operator)
}

func caseExactOrderingMatch(a any, operator byte, b any) (bool, error) {
	return caseBasedOrderingMatch(a, b, true, operator)
}

func caseBasedOrderingMatch(a, b any, caseExact bool, operator byte) (result bool, err error) {
	var str1, str2 string
	if str1, str2, err = prepareNumericStringAssertion(a, b); err == nil {
		if caseExact {
			if operator == GreaterOrEqual {
				result = str1 >= str2
			} else {
				result = str1 <= str2
			}
		} else {
			lc := strings.ToLower
			if operator == GreaterOrEqual {
				result = lc(str1) >= lc(str2)
			} else {
				result = lc(str1) <= lc(str2)
			}
		}
	}

	return
}

/*
caseIgnoreSubstringsMatch implements [§ 4.2.13 of RFC 4517].

OID: 2.5.13.4.

[§ 4.2.13 of RFC 4517]: https://datatracker.ietf.org/doc/html/rfc4517#section-4.2.13
*/
func caseIgnoreSubstringsMatch(a, b any) (result bool, err error) {
	result, err = substringsMatch(a, b, true)
	return
}

/*
caseIgnoreSubstringsMatch implements [§ 4.2.6 of RFC 4517].

OID: 2.5.13.7.

[§ 4.2.6 of RFC 4517]: https://datatracker.ietf.org/doc/html/rfc4517#section-4.2.6
*/
func caseExactSubstringsMatch(a, b any) (result bool, err error) {
	result, err = substringsMatch(a, b, false)
	return
}

func substringsMatch(a, b any, caseIgnore ...bool) (result bool, err error) {
	var value string
	if value, err = assertString(a, 1, "actual value"); err != nil {
		return
	}

	var B SubstringAssertion
	if B, err = marshalSubstringAssertion(b); err != nil {
		return
	}

	caseHandler := func(val string) string { return val }

	if len(caseIgnore) > 0 {
		if caseIgnore[0] {
			caseHandler = strings.ToLower
		}
	}

	value = caseHandler(value)

	if B.Any == nil {
		err = errorBadType("Missing SubstringAssertion.Any")
		return
	}

	if B.Initial != nil {
		initialStr := caseHandler(string(B.Initial))

		if !strings.HasPrefix(value, initialStr) {
			return
		}
		value = strings.TrimPrefix(value, initialStr)
	}

	anyStr := `*` + strings.Trim(caseHandler(string(B.Any)), `*`) + `*`
	substrings := strings.Split(anyStr, "*")
	for _, substr := range substrings {
		index := strings.Index(value, substr)
		if index == -1 {
			return
		}
		value = value[index+len(substr):]
	}

	if B.Final != nil {
		finalStr := caseHandler(string(B.Final))
		result = strings.HasSuffix(value, finalStr)
		return
	}

	result = true

	return
}

func prepareStringListAssertion(a, b any) (str1, str2 string, err error) {
	assertSubstringsList := func(x any) (list string, err error) {
		var ok bool
		var slices []string
		if slices, ok = x.([]string); ok {
			list = strings.Join(slices, ``)
			list = strings.ReplaceAll(list, `\\`, ``)
			list = strings.ReplaceAll(list, `$`, ``)
		} else {
			errorBadType("substringslist")
		}
		return
	}

	if str1, err = assertSubstringsList(a); err == nil {
		str2, err = assertSubstringsList(b)
	}

	return
}

func caseIgnoreListSubstringsMatch(a, b any) (result bool, err error) {
	var str1, str2 string
	if str1, str2, err = prepareStringListAssertion(a, b); err == nil {
		result, err = caseIgnoreSubstringsMatch(str1, str2)
	}

	return
}

func caseIgnoreListMatch(a, b any) (result bool, err error) {
	var strs1, strs2 []string
	if strs1, strs2, err = assertLists(a, b); err != nil {
		return
	}

	if len(strs1) != len(strs2) {
		return
	}

	for idx, slice := range strs1 {
		if !strings.EqualFold(slice, strs2[idx]) || slice == "" {
			return
		}
	}

	result = true
	return
}

func assertLists(a, b any) (strs1, strs2 []string, err error) {
	var ok bool

	if strs1, ok = a.([]string); !ok {
		err = errorBadType("list")
		return
	}

	if strs2, ok = b.([]string); !ok {
		err = errorBadType("list")
	}

	return
}

func caseExactIA5Match(a, b any) (bool, error) {
	return caseBasedIA5Match(a, b, true)
}

func caseIgnoreIA5Match(a, b any) (bool, error) {
	return caseBasedIA5Match(a, b, false)
}

func caseBasedIA5Match(a, b any, caseExact bool) (result bool, err error) {
	var str1, str2 string
	if str1, err = assertString(a, 1, "ia5String"); err != nil {
		return
	}

	if str2, err = assertString(b, 1, "ia5String"); err != nil {
		return
	}

	if err = checkIA5String(str1); err == nil {
		if err = checkIA5String(str2); err == nil {
			if caseExact {
				result = str1 == str2
			} else {
				result = strings.EqualFold(str1, str2)
			}
		}
	}

	return
}

func prepareIA5StringAssertion(a, b any) (str1, str2 string, err error) {
	assertIA5 := func(x any) (i string, err error) {
		var raw string
		if raw, err = assertString(x, 1, "IA5String"); err == nil {
			if err = checkIA5String(raw); err == nil {
				i = raw
			}
		}
		return
	}

	if str1, err = assertIA5(a); err == nil {
		str2, err = assertIA5(b)
	}

	return
}

func caseIgnoreIA5SubstringsMatch(a, b any) (result bool, err error) {
	var str1, str2 string
	if str1, str2, err = prepareIA5StringAssertion(a, b); err == nil {
		result, err = caseIgnoreSubstringsMatch(str1, str2)
	}

	return
}
