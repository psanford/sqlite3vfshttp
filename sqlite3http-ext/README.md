# sqlite3http-ext loadable sqlite3 extension

This is a sqlite3 loadable extension for querying sqlite3 databases over http(s). This allows you to use the standard `sqlite3` cli tool with this VFS.

## Building

Currently I've only tested building this on Linux. I would not expect it to work on MacOS or Windows without some updates to the Makefile.

To build, run `make`.

## Usage

This extension expects an environment variable to be set to specify the remote url to load from `SQLITE3VFSHTTP_URL`. Additionally, the following optional environment variables are supported `SQLITE3VFSHTTP_REFERER` and `SQLITE3VFSHTTP_USER_AGENT`.

To load the extension from the `sqlite3` cli tool run:
```
sqlite> .load ./httpvfs
```

When opening the database, you must specify `vfs=httpvfs`. This requires you to use the [sqlite3 uri](https://www.sqlite.org/uri.html) style filename, which must begin with `file:///`. This extension ignores the filename provided.
```
sqlite> .open file:///foo.db?vfs=httpvfs
```
