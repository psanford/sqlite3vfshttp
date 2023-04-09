package sqlite3vfshttp

import (
	"errors"
	"net/http"
	"strings"

	"github.com/psanford/httpreadat"
	"github.com/psanford/sqlite3vfs"
)

type HttpVFS struct {
	URL          string
	CacheHandler httpreadat.CacheHandler
	RoundTripper http.RoundTripper
}

func (vfs *HttpVFS) Open(name string, flags sqlite3vfs.OpenFlag) (sqlite3vfs.File, sqlite3vfs.OpenFlag, error) {
	var opts []httpreadat.Option
	if vfs.CacheHandler != nil {
		opts = append(opts, httpreadat.WithCacheHandler(vfs.CacheHandler))
	}
	if vfs.RoundTripper != nil {
		opts = append(opts, httpreadat.WithRoundTripper(vfs.RoundTripper))
	}
	tf := &httpFile{
		rr: httpreadat.New(vfs.URL, opts...),
	}

	return tf, flags, nil
}

func (vfs *HttpVFS) Delete(name string, dirSync bool) error {
	return sqlite3vfs.ReadOnlyError
}

func (vfs *HttpVFS) Access(name string, flag sqlite3vfs.AccessFlag) (bool, error) {
	if strings.HasSuffix(name, "-wal") || strings.HasSuffix(name, "-journal") {
		return false, nil
	}

	return true, nil
}

func (vfs *HttpVFS) FullPathname(name string) string {
	return name
}

type httpFile struct {
	rr *httpreadat.RangeReader
}

func (tf *httpFile) Close() error {
	return nil
}

func (tf *httpFile) ReadAt(p []byte, off int64) (int, error) {
	return tf.rr.ReadAt(p, off)
}

func (tf *httpFile) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, sqlite3vfs.ReadOnlyError
}

func (tf *httpFile) Truncate(size int64) error {
	return sqlite3vfs.ReadOnlyError
}

func (tf *httpFile) Sync(flag sqlite3vfs.SyncType) error {
	return nil
}

var invalidContentRangeErr = errors.New("invalid Content-Range response")

func (tf *httpFile) FileSize() (int64, error) {
	return tf.rr.Size()
}

func (tf *httpFile) Lock(elock sqlite3vfs.LockType) error {
	return nil
}

func (tf *httpFile) Unlock(elock sqlite3vfs.LockType) error {
	return nil
}

func (tf *httpFile) CheckReservedLock() (bool, error) {
	return false, nil
}

func (tf *httpFile) SectorSize() int64 {
	return 0
}

func (tf *httpFile) DeviceCharacteristics() sqlite3vfs.DeviceCharacteristic {
	return sqlite3vfs.IocapImmutable
}
