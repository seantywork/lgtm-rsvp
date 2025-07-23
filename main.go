package main

import (
	"encoding/json"
	"fmt"
	pkgdb "lgtm-rsvp/pkg/db"
	pkgglob "lgtm-rsvp/pkg/glob"
	pkgserver "lgtm-rsvp/pkg/server"
	pkgserverapi "lgtm-rsvp/pkg/server/api"
	pkgutils "lgtm-rsvp/pkg/utils"
	"log"
	"os"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

func printHelp() {
	fmt.Println("help: print help")
	fmt.Println("coli y: allow all comment from comment-list.json, force overide")
	fmt.Println("coli n: block all comments approved from comment-list.json, key is title")
}

var _doColi bool = false
var _doColiAct bool = false

func doColi() error {

	err := pkgdb.OpenDB(pkgglob.G_CONF.Db.Addr)

	if err != nil {

		return fmt.Errorf("failed to open db: %s", err.Error())

	}

	err = pkgdb.Init(pkgglob.G_CONF.Db.InitFile, pkgglob.G_CONF.Admin.Id, pkgglob.G_CONF.Admin.Pw)

	if err != nil {

		return fmt.Errorf("failed to init db: %s", err.Error())

	}

	fb, err := os.ReadFile(pkgglob.G_COMMENT_LIST_PATH)

	if err != nil {
		return fmt.Errorf("failed to read comment list: %v", err)
	}

	var coli pkgserverapi.CommntDataList

	err = json.Unmarshal(fb, &coli)

	if err != nil {

		return fmt.Errorf("failed to unmarshal comment list: %v", err)
	}

	clen := len(coli)

	if _doColiAct {

		log.Printf("allowing all data in comment list...\n")

		for i := 0; i < clen; i++ {

			c, err := pkgdb.GetCommentById(coli[i].CommentId)

			id := ""
			title := ""

			if err != nil {
				log.Printf("  - comment by id doesn't exit: %v\n", err)
				comment_id, _ := pkgutils.GetRandomHex(32)
				now := time.Now().UTC()
				p := bluemonday.UGCPolicy()
				title_san := p.Sanitize(coli[i].Title)
				content_san := p.Sanitize(coli[i].Content)
				if title_san == "" {
					log.Printf("  - register comment: invalid title\n")
					continue
				}
				if content_san == "" {
					log.Printf("  - register comment: invalid content\n")
					continue
				}
				timeRegistered := now.Format("2006-01-02-15-04-05")
				err = pkgdb.RegisterComment(comment_id, title_san, content_san, timeRegistered)
				if err != nil {
					log.Printf("  - register comment failed: %v\n", err)
				} else {
					log.Printf("  - register comment success: %s\n", title_san)
				}
				id = comment_id
				title = title_san
			} else {
				id = c.Id
				title = c.Title
			}

			now := time.Now().UTC()
			timeApproved := now.Format("2006-01-02-15-04-05")

			err = pkgdb.ApproveComment(id, timeApproved)

			if err != nil {
				log.Printf("  - approve comment failed: %v\n", err)
			} else {
				log.Printf("  - approve comment success: %s\n", title)
			}
		}

	} else {

		log.Printf("blocking all data in comment list...\n")

		for i := 0; i < clen; i++ {
			err := pkgdb.DisapproveCommentByTitle(coli[i].Title)
			if err != nil {
				log.Printf("  - disapprove by title failed: %v\n", err)
			} else {
				log.Printf("  - disapprove by title success: %s\n", coli[i].Title)
			}

		}

	}

	return nil
}

func main() {

	err := pkgglob.LoadConfig()

	if err != nil {

		log.Printf("failed to load config: %v\n", err)

		os.Exit(-1)
	}

	if len(os.Args) > 1 {

		if os.Args[1] == "help" {
			printHelp()
			os.Exit(0)
		}
		if os.Args[1] == "coli" {
			if len(os.Args) != 3 {
				fmt.Printf("invalid coli command count: %v\n", os.Args[1:])
				printHelp()
				os.Exit(-1)
			} else if os.Args[2] == "y" {
				_doColi = true
				_doColiAct = true
			} else if os.Args[2] == "n" {
				_doColi = true
				_doColiAct = false

			} else {
				fmt.Printf("invalid command coli: %v\n", os.Args[1:])
				printHelp()
				os.Exit(-1)
			}

		} else {
			fmt.Printf("invalid command: %v\n", os.Args[1:])
			printHelp()
			os.Exit(-1)
		}

	}

	if _doColi {
		err := doColi()
		if err != nil {
			log.Printf("failed to do comment list action: %v\n", err)
			os.Exit(-1)
		} else {
			os.Exit(0)
		}
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
