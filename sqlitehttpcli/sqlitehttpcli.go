package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/psanford/sqlite3vfs"
	"github.com/psanford/sqlite3vfshttp"
)

var url = flag.String("url", "", "URL of sqlite db")
var query = flag.String("query", "", "Query to run")
var referer = flag.String("referer", "", "HTTP Referer")
var userAgent = flag.String("user-agent", "", "HTTP User agent")

func main() {
	flag.Parse()

	if *url == "" || *query == "" {
		log.Printf("-url and -query are required")
		flag.Usage()
		os.Exit(1)
	}

	vfs := sqlite3vfshttp.HttpVFS{
		URL: *url,
		RoundTripper: &roundTripper{
			referer:   *referer,
			userAgent: *userAgent,
		},
	}

	err := sqlite3vfs.RegisterVFS("httpvfs", &vfs)
	if err != nil {
		log.Fatalf("register vfs err: %s", err)
	}

	db, err := sql.Open("sqlite3", "not_a_real_name.db?vfs=httpvfs&mode=ro")
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

		columns := make([]*string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err = rows.Scan(columnPointers...)
		if err != nil {
			log.Fatalf("query rows err: %s", err)
		}

		names := make([]string, 0, len(columns))
		for _, col := range columns {
			if col == nil {
				names = append(names, "NULL")
			} else {
				names = append(names, *col)
			}
		}
		fmt.Printf("row: %+v\n", names)
	}
	err = rows.Close()
	if err != nil {
		log.Fatalf("query rows err: %s", err)
	}

}

type roundTripper struct {
	referer   string
	userAgent string
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.referer != "" {
		req.Header.Set("Referer", rt.referer)
	}

	if rt.userAgent != "" {
		req.Header.Set("User-Agent", rt.userAgent)
	}

	tr := http.DefaultTransport

	if req.URL.Scheme == "file" {
		path := req.URL.Path
		root := filepath.Dir(path)
		base := filepath.Base(path)
		tr = http.NewFileTransport(http.Dir(root))
		req.URL.Path = base
	}

	return tr.RoundTrip(req)
}
