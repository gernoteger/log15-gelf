[![GoDoc](https://godoc.org/github.com/gernoteger/mapstructure-hooks?status.svg)](https://godoc.org/github.com/gernoteger/log15-gelf)
[![Go Report Card](https://goreportcard.com/badge/gernoteger/mapstructure-hooks)](https://goreportcard.com/report/gernoteger/log15-gelf)
[![Build Status](https://travis-ci.org/inconshreveable/log15.svg?branch=master)](https://travis-ci.org/gernoteger/log15-gelf)

# Gelf Handler for log15
Adds the [GELF](http://docs.graylog.org/en/2.1/pages/gelf.html) format for Graylog-based logging to the [log15](https://github.com/inconshreveable/log15) logging library.
GELF can be udp+tcp based, and supports chunking with udp, thus avoiding reconnection- and performance issues.

# Duplicate keys
Currently log15 will append duplicate keys to the contect list. Gelf expects a map, therefore keys have to be unique.
This implementation assures that only the last value is used for this key.

```go
    l1:=log.New("foo","bar")
    l1.Info("a message","foo","baz")
    // Output: in GELF: msg: "a message", foo: "baz"
```

# Limitations
- only supports udp with gzip compression.
- buffer size not adjustable

# License

Released under the [MIT License](LICENSE).