package http_router

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	http2 "net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	g := httprouter.New()
	g.GET("/a", func(writer http2.ResponseWriter, request *http2.Request, params httprouter.Params) {
		_, _ = fmt.Fprint(writer, "hello")
	})

	log.Fatal(http2.ListenAndServe(":9009", g))
}

/*
zander@macos ~ % wrk -t4 -c1000 -d10s --latency "http://127.0.0.1:9009/a"
Running 10s test @ http://127.0.0.1:9009/a
  4 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.16ms  749.92us   8.84ms   72.38%
    Req/Sec    27.63k     4.62k   39.55k    63.61%
  Latency Distribution
     50%    2.17ms
     75%    2.66ms
     90%    2.97ms
     99%    4.39ms
  1110453 requests in 10.10s, 128.14MB read
  Socket errors: connect 751, read 93, write 0, timeout 0
Requests/sec: 109907.57
Transfer/sec:     12.68MB
*/
