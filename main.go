package main

import (
	pkgdb "lgtm-rsvp/pkg/db"
	pkgglob "lgtm-rsvp/pkg/glob"
	pkgserver "lgtm-rsvp/pkg/server"
	"log"
	"os"
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

	err = pkgdb.OpenDB(pkgglob.G_CONF.Db.Addr)

	if err != nil {

		log.Printf("failed to open db: %s\n", err.Error())

		os.Exit(-1)

	}

	err = pkgdb.Init(pkgglob.G_CONF.Db.InitFile, pkgglob.G_CONF.Admin.Id, pkgglob.G_CONF.Admin.Pw)

	if err != nil {

		log.Printf("failed to init db: %s\n", err.Error())

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
