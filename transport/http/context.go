package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin/binding"
	"github.com/zander-84/seagull/internal/host"
	"github.com/zander-84/seagull/think"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultMultipartMemory = 32 << 20 // 32 MB

type Proxy interface {
	Param(key string) string
}

var _ Context = (*wrapper)(nil)

// Context is an HTTP Context.
type Context interface {
	context.Context
	//Vars() url.Values
	Query() url.Values
	Form() url.Values

	Param(key string) string

	//ClientIp []string{"X-Forwarded-For"}
	ClientIp(remoteIPHeaders []string) string

	FormFile(name string) (*multipart.FileHeader, error)
	FormFiles(name string) ([]*multipart.FileHeader, error)

	Header() http.Header
	Request() *http.Request
	Response() http.ResponseWriter

	Bind(interface{}) error
	BindQuery(interface{}) error
	BindJson(interface{}) error
	BindForm(interface{}) error
	BindMultipartForm(interface{}) error

	//Returns(interface{}, error) error
	//Result(int, interface{}) error
	JSON(int, interface{}) func() error
	XML(int, interface{}) func() error
	String(int, string) func() error
	Blob(int, string, []byte) func() error
	Stream(int, string, io.Reader) func() error

	//Reset(http.ResponseWriter, *http.Request)
	ErrorEncoder(err error, isProdEnv bool) error
}

func NewHttpContext(res http.ResponseWriter, req *http.Request, proxy Proxy) Context {
	w := &wrapper{}
	w.Reset(res, req, proxy)
	return w
}

type responseWriter struct {
	code int
	w    http.ResponseWriter
}

func (w *responseWriter) rest(res http.ResponseWriter) {
	w.w = res
	w.code = http.StatusOK
}
func (w *responseWriter) Header() http.Header        { return w.w.Header() }
func (w *responseWriter) WriteHeader(statusCode int) { w.code = statusCode }
func (w *responseWriter) Write(data []byte) (int, error) {
	w.w.WriteHeader(w.code)
	return w.w.Write(data)
}

type wrapper struct {
	proxy Proxy
	req   *http.Request
	res   http.ResponseWriter
	w     responseWriter
}

func (c *wrapper) Header() http.Header {
	return c.req.Header
}

func (c *wrapper) Form() url.Values {
	if err := c.req.ParseForm(); err != nil {
		return url.Values{}
	}
	return c.req.Form
}

func (c *wrapper) Query() url.Values {
	return c.req.URL.Query()
}

func (c *wrapper) Request() *http.Request        { return c.req }
func (c *wrapper) Response() http.ResponseWriter { return c.res }

func (c *wrapper) JSON(code int, v interface{}) func() error {
	return func() error {
		c.res.Header().Set("Content-Type", "application/json")
		c.res.WriteHeader(code)
		return json.NewEncoder(c.res).Encode(v)
	}
}

func (c *wrapper) XML(code int, v interface{}) func() error {
	return func() error {
		c.res.Header().Set("Content-Type", "application/xml")
		c.res.WriteHeader(code)
		return xml.NewEncoder(c.res).Encode(v)
	}
}

func (c *wrapper) String(code int, text string) func() error {
	return func() error {
		c.res.Header().Set("Content-Type", "text/plain")
		c.res.WriteHeader(code)
		_, err := c.res.Write([]byte(text))
		return err
	}
}

func (c *wrapper) Blob(code int, contentType string, data []byte) func() error {
	return func() error {
		c.res.Header().Set("Content-Type", contentType)
		c.res.WriteHeader(code)
		_, err := c.res.Write(data)
		return err
	}

}

func (c *wrapper) Stream(code int, contentType string, rd io.Reader) func() error {
	return func() error {
		c.res.Header().Set("Content-Type", contentType)
		c.res.WriteHeader(code)
		_, err := io.Copy(c.res, rd)
		return err
	}
}

func (c *wrapper) Reset(res http.ResponseWriter, req *http.Request, proxy Proxy) {
	c.w.rest(res)
	c.res = res
	c.req = req
	c.proxy = proxy
}

func (c *wrapper) Deadline() (time.Time, bool) {
	if c.req == nil {
		return time.Time{}, false
	}
	return c.req.Context().Deadline()
}

func (c *wrapper) Done() <-chan struct{} {
	if c.req == nil {
		return nil
	}
	return c.req.Context().Done()
}

func (c *wrapper) Err() error {
	if c.req == nil {
		return context.Canceled
	}
	return c.req.Context().Err()
}

func (c *wrapper) Value(key interface{}) interface{} {
	if c.req == nil {
		return nil
	}
	return c.req.Context().Value(key)
}

func (c *wrapper) ErrorEncoder(err error, isProdEnv bool) error {
	thinkErr := think.FromError(err)
	var data any
	if isProdEnv && think.IsErrSystemSpace(thinkErr) {
		data = thinkErr.Response.Data
		thinkErr.Response.Data = thinkErr.Code.ToString()
	}

	_ = c.JSON(thinkErr.Code.HttpCode(), thinkErr.Response)()

	if isProdEnv && think.IsErrSystemSpace(thinkErr) {
		thinkErr.Response.Data = data
	}
	return nil
}

func (c *wrapper) filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func (c *wrapper) ContentType() string {
	return c.filterFlags(c.req.Header.Get("Content-Type"))
}

// Bind 支持 x-www-form-urlencoded / form-data / url?user=xx&p=123 / json
func (c *wrapper) Bind(v interface{}) error {
	b := binding.Default(c.req.Method, c.ContentType())
	return b.Bind(c.req, v)
}

func (c *wrapper) BindJson(v interface{}) error {
	return binding.JSON.Bind(c.req, v)
}

// BindQuery 支持url?user=xx&p=123
func (c *wrapper) BindQuery(v interface{}) error {
	return binding.Query.Bind(c.req, v)
}

// BindForm 支持 x-www-form-urlencoded / form-data / url?user=xx&p=123
func (c *wrapper) BindForm(v interface{}) error {
	return binding.Form.Bind(c.req, v)
}

func (c *wrapper) BindMultipartForm(v interface{}) error {
	return binding.FormMultipart.Bind(c.req, v)
}

func (c *wrapper) FormFile(name string) (*multipart.FileHeader, error) {
	if c.req.MultipartForm == nil {
		if err := c.req.ParseMultipartForm(DefaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.req.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//src, err := fh.Open()
	//if err != nil {
	//}
	//defer src.Close()
	//
	//bs, err := io.ReadAll(src)
	return fh, err
}

// FormFiles is the parsed multipart form, including file uploads.
// key : "upload[]"
func (c *wrapper) FormFiles(key string) ([]*multipart.FileHeader, error) {
	err := c.req.ParseMultipartForm(DefaultMultipartMemory)
	if err != nil {
		return nil, err
	}
	files, ok := c.req.MultipartForm.File[key]
	if ok {
		return files, err
	}
	return nil, think.RecordNotFound
}
func (c *wrapper) Param(key string) string {
	return c.proxy.Param(key)
}

func (c *wrapper) ClientIp(remoteIPHeaders []string) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.req.RemoteAddr))
	if err != nil {
		return ""
	}

	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		return ""
	}

	for _, headerName := range remoteIPHeaders {
		ip, valid := host.ValidateHeaderIp(c.Header().Get(headerName))
		if valid {
			return ip
		}
	}
	return remoteIP.String()
}
