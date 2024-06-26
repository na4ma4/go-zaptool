# go-zaptool

[![CI](https://github.com/na4ma4/go-zaptool/workflows/CI/badge.svg)](https://github.com/na4ma4/go-zaptool/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/na4ma4/go-zaptool/?status.svg)](https://godoc.org/github.com/na4ma4/go-zaptool)
[![GitHub issues](https://img.shields.io/github/issues/na4ma4/go-zaptool)](https://github.com/na4ma4/go-zaptool/issues)
[![GitHub forks](https://img.shields.io/github/forks/na4ma4/go-zaptool)](https://github.com/na4ma4/go-zaptool/network)
[![GitHub stars](https://img.shields.io/github/stars/na4ma4/go-zaptool)](https://github.com/na4ma4/go-zaptool/stargazers)
[![GitHub license](https://img.shields.io/github/license/na4ma4/go-zaptool)](https://github.com/na4ma4/go-zaptool/blob/main/LICENSE)

[uber-go/zap](https://github.com/uber-go/zap) wrappers and tools.

## Install

```shell
go get -u github.com/na4ma4/go-zaptool
```

## Tools

### LogLevels

```golang
logger, _ := cfg.ZapConfig().Build()
ll := zaptool.NewLogLevels(logger)

processOne := server.NewProcess(ll.Named("Server.Process"))

// somewhere else.

ll.SetLevel("Server.Process", "debug")

// and triggered somewhere else again.

ll.SetLevel("Server.Process", "info")
```

### HTTP Logging Handler

```golang
logger, _ := zap.NewProduction()
defer logger.Sync() // flushes buffer, if any
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is a catch-all route"))
})

loggedRouter := zaptool.LoggingHTTPHandler(logger, r)
http.ListenAndServe(":1123", loggedRouter)
```
