# DrWeb challange

Golang code challange was to implement a filestore with http access:

1. User may upload a file to server: file should be saved locally at `/store/ab/abcdef12345`, where `abcdef12345` is a hash of file contents, `ab` is a couple of first hash letters;
2. It should be possible to call callbacks for pre-processing and post-processing of a file;
3. User may download a file by requesting it with its hash;
4. User may remove a file with its hash aswell.

## Configuration settings

Configuration settings might be passed to application via environment variables.

* `LISTEN` - `host:port` for server. Default: `:80`
* `WRITE_TIMEOUT` - Duration within which the whole request must be written back to the client (seconds). Default: `15`
* `READ_TIMEOUT` - Duration within which the whole request must be read from the client (seconds). Default: `15`
* `PATH_NESTED_LEVELS` - How many levels of nesting should be used when storing a file. Default: `2`
* `PATH_NESTED_FOLDERS_LENGTH` - How many characters should each folder's name consist of. Default: `2`
* `PATH_BASE` - Where to store files and corresponding folders. Default: `.`
* `STORAGE_FILE_MODE` - What filemode to use when creating files and folders. Default: `0700`

## Firing up

```
go install github.com/golang/mock/mockgen
dep ensure
go install ./...
$GOPATH/bin/drweb
```

## Tests

```
go test -race ./...
```
