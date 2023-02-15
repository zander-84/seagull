package endpoint

import (
	"context"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool"
)

type Body interface {
	Get(key any) (any, bool)
	Set(key any, value any)
}

// Header is the storage medium used by a Header.
type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
	Foreach(f func(k, v string) error)
}
type header map[string][]string

func NewHeader(in map[string][]string) Header {
	if in == nil {
		in = make(map[string][]string, 0)
	}
	return header(in)
}
func (h header) Foreach(f func(k, v string) error) {
	for k, _ := range h {
		if err := f(k, h.Get(k)); err != nil {
			return
		}
	}
}
func (h header) Get(key string) string {
	v := h[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (h header) Set(key string, value string) {
	h[key] = []string{value}
}

func (h header) Keys() []string {
	out := make([]string, 0, len(h))
	for k, _ := range h {
		out = append(out, k)
	}
	return out
}

// Transporter is transport context value interface.
type Transporter interface {
	// Kind transporter
	// grpc
	// http
	Kind() Kind

	Method() Method
	// Endpoint return server or client endpoint
	// Server Transport: grpc://127.0.0.1:9000
	// Client Transport: discovery:///provider-demo
	//Endpoint() string
	// Operation Service full method selector generated by protobuf
	// example: /helloworld.Greeter/SayHello
	Operation() string
	// RequestHeader return transport request header
	// http: http.Header
	// grpc: metadata.MD
	RequestHeader() Header
	// ReplyHeader return transport reply/response header
	// only valid for server transport
	// http: http.Header
	// grpc: metadata.MD
	ReplyHeader() Header

	Body() Body

	SetCode(code think.Code)
	SetBizCode(bizCode string)

	Code() think.Code
	BizCode() string
}

type transporter struct {
	body          *tool.ConcurrentMap
	kind          Kind
	requestHeader Header
	replyHeader   Header
	operation     string
	method        Method
	code          think.Code
	bizCode       string
}

type endpointKey struct{}

func GetTransporter(ctx context.Context) Transporter {
	v, ok := ctx.Value(endpointKey{}).(Transporter)
	if !ok {
		panic("err ctx")
	}
	return v
}

func NewTransporter(kind Kind, inHeader Header, operation string, method Method) Transporter {
	t := new(transporter)
	t.body = tool.NewConcurrentMap()
	t.kind = kind
	t.requestHeader = inHeader
	t.operation = operation
	t.replyHeader = make(header, 0)
	t.method = method
	t.code = think.Code_Success
	return t
}

func WithContext(ctx context.Context, endpointCtx Transporter) context.Context {
	return context.WithValue(ctx, endpointKey{}, endpointCtx)
}

func (t *transporter) Kind() Kind {
	return t.kind
}

func (t *transporter) Operation() string {
	return t.operation
}
func (t *transporter) Method() Method {
	return t.method
}
func (t *transporter) RequestHeader() Header {
	return t.requestHeader
}

func (t *transporter) ReplyHeader() Header {
	return t.replyHeader
}

func (t *transporter) Body() Body {
	return t.body
}
func (t *transporter) Code() think.Code {
	return t.code
}
func (t *transporter) BizCode() string {
	return t.bizCode
}
func (t *transporter) SetCode(code think.Code) {
	t.code = code
}
func (t *transporter) SetBizCode(bizCode string) {
	t.bizCode = bizCode
}
