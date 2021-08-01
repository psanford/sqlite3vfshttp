package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/psanford/sqlite3vfs"
	"github.com/psanford/sqlite3vfshttp"
)

var url = flag.String("url", "", "URL of sqlite db")
var query = flag.String("query", "", "Query to run")

func main() {
	flag.Parse()

	if *url == "" || *query == "" {
		log.Printf("-url and -query are required")
		flag.Usage()
		os.Exit(1)
	}

	vfs := sqlite3vfshttp.HttpVFS{URL: *url}

	err := sqlite3vfs.RegisterVFS("httpvfs", &vfs)
	if err != nil {
		log.Fatalf("register vfs err: %s", err)
	}

	db, err := sql.Open("sqlite3", "not_a_read_name.db?vfs=httpvfs&mode=ro")
	if err != nil {
		log.Fatalf("open db err: %s", err)
	}

	rows, err := db.Query(*query)
	if err != nil {
		log.Fatalf("query err: %s", err)
	}

	cols, _ := rows.Columns()

	for rows.Next() {
		rows.Columns()

		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err = rows.Scan(columnPointers...)
		if err != nil {
			log.Fatalf("query rows err: %s", err)
		}

		fmt.Printf("row: %+v\n", columns)
	}
	err = rows.Close()
	if err != nil {
		log.Fatalf("query rows err: %s", err)
	}

}
