package http_gin

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestRouter(t *testing.T) {
	g := gin.New()
	gin.SetMode("release")
	g.GET("/a", func(c *gin.Context) {
		c.String(200, "hello")
	})
	g.Run("127.0.0.1:9009") //
}

/*
zander@macos ~ % wrk -t4 -c1000 -d10s --latency "http://127.0.0.1:9009/a"
Running 10s test @ http://127.0.0.1:9009/a
  4 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.19ms  754.22us   8.56ms   71.66%
    Req/Sec    27.28k     6.72k   42.07k    65.00%
  Latency Distribution
     50%    2.22ms
     75%    2.69ms
     90%    3.02ms
     99%    4.35ms
  1085539 requests in 10.01s, 125.27MB read
  Socket errors: connect 751, read 81, write 0, timeout 0
Requests/sec: 108499.17
Transfer/sec:     12.52MB
*/
