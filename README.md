# idgen

`idgen.ID` is an unsiged 64-bit pseudo unique ID consisting of a 42-bit Unix timestamp (in ms), a node ID (8 bits) and a sequence number (14 bits). Timestamp is placed first so that IDs generated are (roughly) sortable on creation. It is designed to support up to 256 nodes, and 16384 IDs unique within a millisecond.

The program is single-threaded, but is fast enough to generate 1 million IDs under a second when tested locally (2014 MBP).

## Setup

Download and place this project under your `$GOPATH`.
```
$GOPATH/src/github.com/microamp/idgen
```

## Build

```
go build ./...
```

## Tests

`stretchr/testify` is required to run the tests.
```
go get -u github.com/stretchr/testify
```

Then,
```
go test ./idgen/...
```

## Sample usage

`cmd/test_simple/test_simple.go` demonstrates how `idgen.ID` can be generated. e.g.
```
go run cmd/test_simple/test_simple.go
```
