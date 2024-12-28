package main

import (
	"fmt"
)

func test(tc int) error {

	var reterr error = nil

	switch tc {

	case 0:

		break

	default:

		reterr = fmt.Errorf("invalid test case: %d", tc)

	}

	return reterr

}
