[![GoDoc](https://godoc.org/github.com/gernotegerlog15-gelf?status.svg)](https://godoc.org/github.com/gernoteger/log15-gelf)
[![Go Report Card](https://goreportcard.com/badge/gernoteger/log15-gelf)](https://goreportcard.com/report/gernoteger/log15-gelf)
[![Build Status](https://travis-ci.org/gernoteger/log15-gelf.svg?branch=master)](https://travis-ci.org/gernoteger/log15-gelf)

# Gelf Handler for log15
Adds the [GELF](http://docs.graylog.org/en/2.1/pages/gelf.html) format for Graylog-based logging to the [log15](https://github.com/inconshreveable/log15) logging library.
GELF can be udp+tcp based, and supports chunking with udp, thus avoiding reconnection- and performance issues.

# Usage

create a Handler with:

     h,err:=gelf.GelfHandler("myhost:12201")
Currently only udp is transport implemented. You can also use [log15-config)[https://github.com/gernoteger/log15-config] 
with the [config](config) package from this repo.


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