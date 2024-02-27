//nolint:gocritic,lll // Example code, ignore output comment line and log.Fatal.
package zaptool_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/na4ma4/go-zaptool"
	"go.uber.org/zap"
)

func ExampleLoggingHTTPHandler() {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()

	loggedRouter := zaptool.LoggingHTTPHandler(
		logger,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			// do nothing
		}),
		zaptool.LoggingOptionTimestamp(false),
		zaptool.LoggingOptionTiming(false),
	)

	ts := httptest.NewServer(loggedRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
	// Output:
	// {"level":"info","msg":"Request","http":{"host":"127.0.0.1","username":"-","method":"GET","uri":"/","proto":"HTTP/1.1","status":200,"size":0,"referer":"","user-agent":"Go-http-client/1.1"}}
}
