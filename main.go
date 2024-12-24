package main

import (
	"log"
	"os"
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

	if err := srv.Run(pkgglob.G_CONF.ServeAddr); err != nil {

		log.Printf("failed to run: %v\n", err)

		os.Exit(-1)

	} else {

		log.Printf("success\n")
	}

	os.Exit(0)

}
