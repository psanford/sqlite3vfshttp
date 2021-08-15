# sqlite3vfshttp: a Go sqlite VFS for querying databases over http(s)

sqlite3vfshttp is a sqlite3 VFS for querying remote databases over http(s).
This allows you to perform queries without needing to download the complete database
first.

Your database must be hosted on a webserver that supports HTTP range requests (such as Amazon S3).

## Example

See [sqlitehttpcli/sqlitehttpcli.go](sqlitehttpcli/sqlitehttpcli.go) for a simple CLI
tool that is able to query a remotely hosted sqlite database.

## Usage


```
	vfs := sqlite3vfshttp.New(*url)

	err := sqlite3vfs.RegisterVFS("httpvfs", vfs)
	if err != nil {
		log.Fatalf("register vfs err: %s", err)
	}

	db, err := sql.Open("sqlite3", "not_a_read_name.db?vfs=httpvfs&mode=ro")
	if err != nil {
		log.Fatalf("open db err: %s", err)
	}
```

## Querying a database in S3

The original purpose of this library was to allow an AWS Lambda function to be able to query a sqlite database stored in S3 without downloading the entire database.

This is possible even for private files stored in S3 by generating a presigned URL and passing that to this library. That allows the client to make HTTP Get range requests without it needing to know how to sign S3 requests.

## Building a loadable extension for the sqlite3 cli

The `sqlite3` cli supports runtime loadable extensions. We can build `sqlite3vfshttp` as a shared library, and then load it from the `sqlite3` cli to interactively query sqlite databases over http connections. The shared library code is located in the `sqlite3http-ext` directory. See the sqlite3http-ext/README.md for more details.

## Demo

I've uploaded a 30MB sqlite database to a publicly accessible webserver for testing, based on "Balance of payments international investment position: March 2021 quarter – CSV" from https://www.stats.govt.nz/large-datasets/csv-files-for-download/. The schema is:

```
CREATE TABLE csv (series_reference,
period,
data_value,
suppressed,
status,
units,
magntude,
subject,
grp,
series_title_1);

```

You can query this dataset from the `sqlite3` cli tool using the shared library extension:
```
$ cd sqlite3http-ext

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

Alternatively, you can query this dataset using the `sqlitehttpcli` example tool:

```
# query the sqlite schema table
$ ./sqlitehttpcli -url 'https://www.sanford.io/demo.db' -query 'select * from main.sqlite_master'
row: [table csv csv 2 CREATE TABLE csv (series_reference,
period,
data_value,
suppressed,
status,
units,
magntude,
subject,
grp,
series_title_1)]

# get 10 rows from the dataset
./sqlitehttpcli -url 'https://www.sanford.io/demo.db' -query "select * from csv limit 10"
row: [BOPQ.S06AC000000000A 1971.06 426  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1971.09 435  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1971.12 360  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1972.03 417  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1972.06 528  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1972.09 471  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1972.12 437  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1973.03 607  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1973.06 666  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 1973.09 578  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]

# get 10 rows where the period is after 2010
$ ./sqlitehttpcli -url 'https://www.sanford.io/demo.db' -query "select * from csv where period > '2010' limit 10"
row: [BOPQ.S06AC000000000A 2010.03 17463  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2010.06 17260  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2010.09 15419  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2010.12 17088  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2011.03 18516  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2011.06 18835  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2011.09 16390  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2011.12 18748  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2012.03 18477  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]
row: [BOPQ.S06AC000000000A 2012.06 18270  F Dollars 6 Balance of Payments - BOP BPM6 Quarterly, Balance of payments major components Actual]

```
