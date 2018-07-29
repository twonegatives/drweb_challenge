[![CircleCI](https://circleci.com/gh/twonegatives/drweb_challenge.svg?style=svg)](https://circleci.com/gh/twonegatives/drweb_challenge)
[![codecov](https://codecov.io/gh/twonegatives/drweb_challenge/branch/master/graph/badge.svg)](https://codecov.io/gh/twonegatives/drweb_challenge)

# DrWeb challenge

Golang code challenge was to implement a filestore with http access:

1. User may upload a file to server: file should be saved locally at `/store/ab/abcdef12345`, where `abcdef12345` is a hash of file contents, `ab` is a couple of first hash letters;
2. It should be possible to call callbacks for pre-processing and post-processing of a file;
3. User may download a file by requesting it with its hash;
4. User may remove a file with its hash aswell.

## Key endpoints

<table>
  <thead>
    <tr>
      <th>Method</th>
      <th>Path</th>
      <th>Body params</th>
      <th>Code</th>
      <th>Response body</th>
      <th>Explanation</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th>POST</th>
      <th>/files</th>
      <th>file: form</th>
      <th>201</th>
      <th>{hashstring: string}</th>
      <th>Created succesfully</th>
    </tr>
    <tr>
      <th></th>
      <th></th>
      <th></th>
      <th>400</th>
      <th>{error: string}</th>
      <th>Could not parse form data</th>
    </tr>
    <tr>
      <th></th>
      <th></th>
      <th></th>
      <th>500</th>
      <th>{error: string}</th>
      <th>Server error</th>
    </tr>
    <tr>
      <th>GET</th>
      <th>/files/filename</th>
      <th></th>
      <th>200</th>
      <th>File contents</th>
      <th>Successfull download</th>
    </tr>
    <tr>
      <th></th>
      <th></th>
      <th></th>
      <th>404</th>
      <th></th>
      <th>Requested file was not found</th>
    </tr>
    <tr>
      <th></th>
      <th></th>
      <th></th>
      <th>500</th>
      <th>{error: string}</th>
      <th>Server error</th>
    </tr>
    <tr>
      <th>DELETE</th>
      <th>/files/filename</th>
      <th></th>
      <th>200</th>
      <th></th>
      <th>Deleted succesfully</th>
    </tr>
    <tr>
      <th></th>
      <th></th>
      <th></th>
      <th>404</th>
      <th></th>
      <th>Requested file was not found</th>
    </tr>
    <tr>
      <th></th>
      <th></th>
      <th></th>
      <th>500</th>
      <th>{error: string}</th>
      <th>Server error</th>
    </tr>
  </tbody>
</table>

## Configuration settings

Configuration settings might be passed to application via environment variables.

* `LISTEN` - `host:port` for server. Default: `:80`
* `WRITE_TIMEOUT` - Duration within which the whole request must be written back to the client (seconds). Default: `15`
* `READ_TIMEOUT` - Duration within which the whole request must be read from the client (seconds). Default: `15`
* `PATH_NESTED_LEVELS` - How many levels of nesting should be used when storing a file. Default: `2`
* `PATH_NESTED_FOLDERS_LENGTH` - How many characters should each folder's name consist of. Default: `2`
* `PATH_BASE` - Where to store files and corresponding folders. Default: `.`
* `STORAGE_FILE_MODE` - What filemode to use when creating files and folders. Default: `0755`

## Firing up

```
dep ensure
go install ./...
LISTEN=:3001 $GOPATH/bin/drweb
```

Local development would also require you to have mockgen for test mocks generation
```
go install github.com/golang/mock/mockgen
```

## Tests

```
go test -race ./...
```
