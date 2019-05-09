package main

import (
	"flag"
	"log"
	"os/exec"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/up/gogen/pkg/gen"
)

var (
	flDrop   = flag.Bool("drop", false, "Drop database before executing")
	flDBName = flag.String("dbname", "test", "Database name")
	ll       = l.New()
)

func main() {
	flag.Parse()

	projectPath := gen.ProjectPath()
	sqlPath := projectPath + "/db/migrate/*.sql"
	cmd := exec.Command("bash", "-c", "cat "+sqlPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err, string(output))
	}

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
	db.MustExec(string(output))

	log.Println("Initialized database for testing")
}
