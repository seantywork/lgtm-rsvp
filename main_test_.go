package main

import (
	"fmt"
	pkgdb "our-wedding-rsvp/pkg/db"
	pkgglob "our-wedding-rsvp/pkg/glob"
)

func test(tc int) error {

	var reterr error = nil

	switch tc {

	case 0:

		reterr = test_db()

		break

	default:

		reterr = fmt.Errorf("invalid test case: %d", tc)

	}

	return reterr

}

func test_db() error {

	err := pkgdb.OpenDB(pkgglob.G_CONF.Db.Addr)

	if err != nil {
		return err
	}

	err = pkgdb.Init(pkgglob.G_CONF.Db.InitFile, pkgglob.G_CONF.Admin.Id, pkgglob.G_CONF.Admin.Pw)

	if err != nil {
		return err
	}

	return nil
}
