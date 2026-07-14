package syntax

import (
	"strings"
)

func wordMatch(a, b any) (result bool, err error) {
	var str1, str2 string
	str1, err = assertString(a, 1, "word")
	if err != nil {
		return
	}

	str2, err = assertString(b, 1, "word")
	if err != nil {
		return
	}

	// Split the attribute value into words
	words := strings.Fields(str2)

	// Check if any word matches the assertion value
	for i := 0; i < len(words) && !result; i++ {
		result = strings.EqualFold(words[i], str1)
	}

	return
}

/*
TODO: dig deeper into other impls. to determine best (or most common)
practice to adopt.
*/
func keywordSplit(input string) (out []string) {
	bld := &strings.Builder{}

	for _, char := range input {
		if isWHSP(char) || isPunct(char) {
			if bld.Len() > 0 {
				out = append(out, bld.String())
				bld.Reset()
			}
		} else {
			bld.WriteRune(char)
		}
	}

	if bld.Len() > 0 {
		out = append(out, bld.String())
	}

	return
}

func keywordMatch(a, b any) (result bool, err error) {
	var str1, str2 string
	if str1, err = assertString(a, 1, "keyword"); err != nil {
		return
	}

	if str2, err = assertString(b, 1, "keyword"); err != nil {
		return
	}

	keys := keywordSplit(str2)
	for i := 0; i < len(keys) && !result; i++ {
		result = strings.EqualFold(keys[i], str1)
	}

	return
}
