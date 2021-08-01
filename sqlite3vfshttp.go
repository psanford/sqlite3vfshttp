package sqlite3vfshttp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/psanford/sqlite3vfs"
)

func New(url string) sqlite3vfs.VFS {
	return &httpVFS{
		url: url,
	}
}

type httpVFS struct {
	url string
}

func (vfs *httpVFS) Open(name string, flags sqlite3vfs.OpenFlag) (sqlite3vfs.File, sqlite3vfs.OpenFlag, error) {
	tf := &httpFile{
		url: vfs.url,
	}

	return tf, flags, nil
}

func (vfs *httpVFS) Delete(name string, dirSync bool) error {
	return sqlite3vfs.ReadOnlyError
}

func (vfs *httpVFS) Access(name string, flag sqlite3vfs.AccessFlag) (bool, error) {
	if strings.HasSuffix(name, "-wal") || strings.HasSuffix(name, "-journal") {
		return false, nil
	}

	return true, nil
}

func (vfs *httpVFS) FullPathname(name string) string {
	return name
}

type httpFile struct {
	url string
}

func (tf *httpFile) Close() error {
	return nil
}

func (tf *httpFile) ReadAt(p []byte, off int64) (int, error) {
	req, err := http.NewRequest("GET", tf.url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", off, off+int64(len(p)-1)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return io.ReadFull(resp.Body, p)
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
	req, err := http.NewRequest("GET", tf.url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Range", "bytes=0-0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	rangeHeader := resp.Header.Get("Content-Range")
	rangeFields := strings.Fields(rangeHeader)
	if len(rangeFields) != 2 {
		return 0, invalidContentRangeErr
	}

	if strings.ToLower(rangeFields[0]) != "bytes" {
		return 0, invalidContentRangeErr
	}

	amts := strings.Split(rangeFields[1], "/")

	if len(amts) != 2 {
		return 0, invalidContentRangeErr
	}

	if amts[1] == "*" {
		return 0, invalidContentRangeErr
	}

	n, err := strconv.Atoi(amts[1])
	if err != nil {
		return 0, invalidContentRangeErr
	}

	return int64(n), nil
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
