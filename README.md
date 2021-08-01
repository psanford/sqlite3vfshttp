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
