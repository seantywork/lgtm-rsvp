package auth

import (
	"strings"
	"unicode"
)

func VerifyDefaultValue(raw string) bool {

	for _, c := range raw {

		if unicode.IsLower(c) {

			continue

		} else if unicode.IsDigit(c) {

			continue

		} else if c == '-' {

			continue
		} else if c == '@' {
			continue
		} else if c == '.' {
			continue
		} else {

			return false
		}

	}

	return true
}

func VerifyMediaKey(mediakey string) bool {

	mklist := strings.SplitN(mediakey, ".", 2)

	if len(mklist) != 2 {
		return false
	}

	for _, c := range mklist[0] {

		if unicode.IsLetter(c) {

			continue

		} else if unicode.IsDigit(c) {

			continue

		} else {

			return false
		}

	}

	for _, c := range mklist[1] {

		if unicode.IsLetter(c) {

			continue

		} else if unicode.IsDigit(c) {

			continue

		} else {

			return false
		}

	}

	return true
}
