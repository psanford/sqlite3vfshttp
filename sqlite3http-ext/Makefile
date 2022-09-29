httpvfs.so: sqlite3http_ext.c sqlite3http_ext.a
	gcc -g -fPIC -shared -o $@ $^

httpvfs.dylib: sqlite3http_ext.c sqlite3http_ext.a
	gcc -framework CoreServices -framework Security -dynamiclib -o $@ $^

sqlite3http_ext.a: sqlite3http_ext.go
	go build -tags SQLITE3VFS_LOADABLE_EXT -o $@ -buildmode=c-archive $^
	rm -f sqlite3http_ext.h

.PHONY: clean
clean:
	rm -f sqlite3http_ext.a httpvfs.so httpvfs.dylib
