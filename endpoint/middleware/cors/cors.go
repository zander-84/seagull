package cors

import (
	"context"
	"errors"
	"github.com/zander-84/seagull/transport"
	"net/http"
	"strings"

	"github.com/zander-84/seagull/endpoint"
	zhttp "github.com/zander-84/seagull/transport/http"
)

type cors struct {
	allowAllOrigins  bool
	allowCredentials bool
	allowOriginFunc  func(string) bool
	allowOrigins     []string
	exposeHeaders    []string
	normalHeaders    http.Header
	preflightHeaders http.Header
	wildcardOrigins  [][]string
}

var (
	DefaultSchemas = []string{
		"http://",
		"https://",
	}
	ExtensionSchemas = []string{
		"chrome-extension://",
		"safari-extension://",
		"moz-extension://",
		"ms-browser-extension://",
	}
	FileSchemas = []string{
		"file://",
	}
	WebSocketSchemas = []string{
		"ws://",
		"wss://",
	}
)

func newCors(config Config) *cors {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}

	return &cors{
		allowOriginFunc:  config.AllowOriginFunc,
		allowAllOrigins:  config.AllowAllOrigins,
		allowCredentials: config.AllowCredentials,
		allowOrigins:     normalize(config.AllowOrigins),
		normalHeaders:    generateNormalHeaders(config),
		preflightHeaders: generatePreflightHeaders(config),
		wildcardOrigins:  config.parseWildcardRules(),
	}
}

func (cors *cors) applyCors(c zhttp.Context) error {
	origin := c.Header().Get("Origin")
	if len(origin) == 0 {
		// request is not a CORS request
		return nil
	}
	host := c.Request().Host

	if origin == "http://"+host || origin == "https://"+host {
		// request is not a CORS request but have origin header.
		// for example, use fetch api
		return nil
	}

	if !cors.validateOrigin(origin) {
		c.Response().WriteHeader(http.StatusForbidden)
		return errors.New("forbidden")
	}

	if c.Request().Method == "OPTIONS" {
		cors.handlePreflight(c)
		defer c.Response().WriteHeader(http.StatusNoContent) // Using 204 is better than 200 when the request status is OPTIONS
		return errors.New("OPTIONS")
	} else {
		cors.handleNormal(c)
	}

	if !cors.allowAllOrigins {
		c.Response().Header().Set("Access-Control-Allow-Origin", origin)
	}
	return nil
}

func (cors *cors) validateWildcardOrigin(origin string) bool {
	for _, w := range cors.wildcardOrigins {
		if w[0] == "*" && strings.HasSuffix(origin, w[1]) {
			return true
		}
		if w[1] == "*" && strings.HasPrefix(origin, w[0]) {
			return true
		}
		if strings.HasPrefix(origin, w[0]) && strings.HasSuffix(origin, w[1]) {
			return true
		}
	}

	return false
}

func (cors *cors) validateOrigin(origin string) bool {
	if cors.allowAllOrigins {
		return true
	}
	for _, value := range cors.allowOrigins {
		if value == origin {
			return true
		}
	}
	if len(cors.wildcardOrigins) > 0 && cors.validateWildcardOrigin(origin) {
		return true
	}
	if cors.allowOriginFunc != nil {
		return cors.allowOriginFunc(origin)
	}
	return false
}

func (cors *cors) handlePreflight(c zhttp.Context) {
	header := c.Response().Header()
	for key, value := range cors.preflightHeaders {
		header[key] = value
	}
}

func (cors *cors) handleNormal(c zhttp.Context) {
	header := c.Response().Header()
	for key, value := range cors.normalHeaders {
		header[key] = value
	}
}

// New returns the location middleware with user-defined custom configuration.
func New(config Config) endpoint.Middleware {
	cors1 := newCors(config)

	return func(next endpoint.HandlerFunc) endpoint.HandlerFunc {
		return func(ctx context.Context, request interface{}) (out interface{}, err error) {
			transporter := transport.GetTransporter(ctx)
			if transporter.Kind().IsHttp() {
				httpCtx := ctx.(zhttp.Context)
				err = cors1.applyCors(httpCtx)
				if err != nil {
					return nil, nil
				}
			}
			return next(ctx, request)
		}
	}
}
