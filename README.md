# DrWeb challange

Golang code challange was to implement a filestore with http access:

1. User may upload a file to server: file should be saved locally at `/store/ab/abcdef12345`, where `abcdef12345` is a hash of file contents, `ab` is a couple of first hash letters;
2. It should be possible to call callbacks for pre-processing and post-processing of a file;
3. User may download a file by requesting it with its hash;
4. User may remove a file with its hash aswell.

## Firing up

```
dep ensure
go install ./...
$GOPATH/bin/drweb
```

## Tests

```
go test -race ./...
```
