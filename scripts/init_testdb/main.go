package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/common/l"
)

var (
	flDrop   = flag.Bool("drop", false, "Drop database before executing")
	flDBName = flag.String("dbname", "test", "Database name")
	ll       = l.New()
)

type Content struct {
	Path string
	Body []byte
}

func main() {
	flag.Parse()

	projectPath := os.Getenv("ETOPDIR") + "/backend"
	sqlPath := filepath.Join(projectPath, "/db/main")

	var contents []Content
	err := filepath.Walk(sqlPath, func(path string, info os.FileInfo, err error) error {
		baseName := filepath.Base(path)
		if strings.HasPrefix(baseName, "_") {
			log.Println("Skipped file", baseName)
			return nil
		}
		if !strings.HasSuffix(baseName, ".sql") {
			log.Println("Skipped file", baseName)
			return nil
		}
		body, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		contents = append(contents, Content{path, body})
		return nil
	})

	cfg := config.DefaultTest().Postgres
	cfg.Database = *flDBName
	db, err := cmsql.Connect(cfg)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)

	if *flDrop {
		ll.Warn("Drop database: " + cfg.Database)
		db.MustExec(`
			DROP SCHEMA public CASCADE;
			DROP SCHEMA IF EXISTS history CASCADE;
			CREATE SCHEMA public;
			GRANT ALL ON SCHEMA public TO postgres;
			GRANT ALL ON SCHEMA public TO public;
`)
	}

	err = db.InTransaction(bus.Ctx(), func(tx cmsql.QueryInterface) error {
		for _, content := range contents {
			log.Println("--- Executing", content.Path)
			_, err := tx.SQL(string(content.Body)).Exec()
			if err != nil {
				ll.Error("Error while executing", l.String("script", content.Path), l.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		os.Exit(1)
	}
	log.Println("Initialized database for testing")
}
