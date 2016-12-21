### go-test

Example project to show off unit testing with go.

Tests with coverage:
```
$ go test -cover -tags test
```

Run with:
```
$ go build && ./go-test
```

Curl it!
```
$ curl http://localhost:8080/?location=Jacksonville,+FL
{"sunset":"5:30 pm","timestamp":"2016-12-21T03:34:40Z"}
```