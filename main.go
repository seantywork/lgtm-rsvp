package main

import (
	"log"
	"os"
	pkgdb "our-wedding-rsvp/pkg/db"
	pkgglob "our-wedding-rsvp/pkg/glob"
	pkgserver "our-wedding-rsvp/pkg/server"
)

func main() {

	err := pkgglob.LoadConfig()

	if err != nil {

		log.Printf("failed to load config: %v\n", err)

		os.Exit(-1)
	}

	srv, err := pkgserver.CreateServerFromConfig()

	if err != nil {

		log.Printf("failed to create server: %v\n", err)

		os.Exit(-1)

	}

	log.Printf("server running at: %s\n", pkgglob.G_CONF.ServeAddr)

	if pkgglob.G_CONF.Test > -1 {

		log.Printf("test mode: case: %d\n", pkgglob.G_CONF.Test)

		if err := test(pkgglob.G_CONF.Test); err != nil {

			log.Printf("test failed: %v\n", err)

			os.Exit(-1)

		} else {

			log.Printf("test success\n")
		}

		os.Exit(0)

	}

	err = pkgdb.OpenDB(pkgglob.G_CONF.Db.Addr)

	if err != nil {

		log.Printf("failed to open db: %s\n", err.Error())

		os.Exit(-1)

	}

	err = pkgdb.Init(pkgglob.G_CONF.Db.InitFile, pkgglob.G_CONF.Admin.Id, pkgglob.G_CONF.Admin.Pw)

	if err != nil {

		log.Printf("failed to open db: %s\n", err.Error())

		os.Exit(-1)
	}

	if err := srv.Run(pkgglob.G_CONF.ServeAddr); err != nil {

		log.Printf("failed to run: %v\n", err)

		os.Exit(-1)

	} else {

		log.Printf("success\n")
	}

	os.Exit(0)

}
