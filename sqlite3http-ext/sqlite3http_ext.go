//go:build SQLITE3VFS_LOADABLE_EXT
// +build SQLITE3VFS_LOADABLE_EXT

package main

// import C is necessary for us to export SqliteHTTPRegister in the c-archive .a file

import "C"
import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/psanford/sqlite3vfs"
	"github.com/psanford/sqlite3vfshttp"
)

func main() {
}

//export Sqlite3HTTPRegister
func Sqlite3HTTPRegister() {
	url := os.Getenv("SQLITE3VFSHTTP_URL")
	if url == "" {
		log.Fatal("SQLITE3VFSHTTP_URL environment variable not defined")
	}

	referer := os.Getenv("SQLITE3VFSHTTP_REFERER")
	userAgent := os.Getenv("SQLITE3VFSHTTP_USER_AGENT")

	vfs := sqlite3vfshttp.HttpVFS{
		URL: url,
		RoundTripper: &roundTripper{
			referer:   referer,
			userAgent: userAgent,
		},
	}

	err := sqlite3vfs.RegisterVFS("httpvfs", &vfs)
	if err != nil {
		log.Fatalf("register vfs err: %s", err)
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
