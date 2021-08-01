package sqlite3vfshttp

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/psanford/sqlite3vfs"
)

func TestSqlite3vfshttp(t *testing.T) {
	dir, err := ioutil.TempDir("", "sqlite3vfshttptest")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite3", filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS foo (
id text NOT NULL PRIMARY KEY,
title text
)`)
	if err != nil {
		t.Fatal(err)
	}

	rows := []FooRow{
		{
			ID:    "415",
			Title: "romantic-swell",
		},
		{
			ID:    "610",
			Title: "ironically-gnarl",
		},
		{
			ID:    "768",
			Title: "biophysicist-straddled",
		},
	}

	for _, row := range rows {
		_, err = db.Exec(`INSERT INTO foo (id, title) values (?, ?)`, row.ID, row.Title)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	s := httptest.NewServer(http.FileServer(http.Dir(dir)))

	vfs := HttpVFS{
		URL: s.URL + "/test.db",
	}

	err = sqlite3vfs.RegisterVFS("httpvfs", &vfs)
	if err != nil {
		t.Fatal(err)
	}

	db, err = sql.Open("sqlite3", "meaningless_name.db?vfs=httpvfs&mode=ro")
	if err != nil {
		t.Fatal(err)
	}

	rowIter, err := db.Query(`SELECT id, title from foo order by id`)
	if err != nil {
		t.Fatal(err)
	}

	var gotRows []FooRow

	for rowIter.Next() {
		var row FooRow
		err = rowIter.Scan(&row.ID, &row.Title)
		if err != nil {
			t.Fatal(err)
		}
		gotRows = append(gotRows, row)
	}
	err = rowIter.Close()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(rows, gotRows) {
		t.Fatal(cmp.Diff(rows, gotRows))
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

}

type FooRow struct {
	ID    string
	Title string
}
