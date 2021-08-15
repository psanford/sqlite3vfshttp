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

## Example

```
# build httpvfs.so shared library
$ make
go build -tags SQLITE3VFS_LOADABLE_EXT -o sqlite3http_ext.a -buildmode=c-archive sqlite3http_ext.go
gcc -g -fPIC -shared -o httpvfs.so sqlite3http_ext.c sqlite3http_ext.a

# set url of sqlite3 db as environment variable SQLITE3VFSHTTP_URL:
$ export SQLITE3VFSHTTP_URL='https://www.sanford.io/demo.db'

$ sqlite3
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite> -- load extention
sqlite> .load ./httpvfs
sqlite> -- open db using vfs=httpvfs, note you must use the sqlite uri syntax which starts with file://
sqlite> .open file:///foo.db?vfs=httpvfs
sqlite> -- query from remote db
sqlite> SELECT * from csv where period > '2010' limit 10;
series_reference      period      data_value  suppressed  status      units       magntude    subject                    grp                                                   series_title_1
--------------------  ----------  ----------  ----------  ----------  ----------  ----------  -------------------------  ----------------------------------------------------  --------------
BOPQ.S06AC000000000A  2010.03     17463                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2010.06     17260                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2010.09     15419                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2010.12     17088                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2011.03     18516                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2011.06     18835                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2011.09     16390                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2011.12     18748                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2012.03     18477                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
BOPQ.S06AC000000000A  2012.06     18270                   F           Dollars     6           Balance of Payments - BOP  BPM6 Quarterly, Balance of payments major components  Actual
```
