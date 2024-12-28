package auth

import (
	"strings"
	"unicode"
)

func VerifyMediaKey(mediakey string) bool {

	mklist := strings.SplitN(mediakey, ".", 2)

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
