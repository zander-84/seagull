package endpoint

import (
	"context"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/transport"
	"log"
	"runtime"
	"strings"
)

// Rmc Resource Management Center
type Rmc interface {
	Group(prefix string, options ...Options) Rmc
	Use(options ...Options) Rmc
	Endpoint(method transport.Method, path string, e HandlerFunc, codec Codecs, options ...Options)
	Proxy(proxy ProxyEndpoint, kind transport.Kind)
	GetEndpoint(p transport.Kind, method transport.Method, path string) (HandlerFunc, error)
	MustGetEndpoint(p transport.Kind, method transport.Method, path string) HandlerFunc

	GetCodec(p transport.Kind, method transport.Method, path string) (Codec, error)
	MustGetCodec(p transport.Kind, method transport.Method, path string) Codec
}

type Conf struct {
	Path     string
	Method   transport.Method
	FullPath Path

	//ps         []Protocol
	Middleware      Middleware
	InnerMiddleware Middleware

	HandlerFunc HandlerFunc

	Codecs Codecs

	ErrorEncoder ErrorEncoder
}

func newRmcConf() Conf {
	out := Conf{
		//DecFunc: func(ctx context.Context, protocol Protocol, request interface{}) (response interface{}, err error) {
		//	return
		//},
	}
	return out
}

type Options func(*Conf)

func OptMW(m ...Middleware) Options {
	return func(rmc *Conf) {
		rmc.Middleware = ChainMerge(rmc.Middleware, m...)
	}
}

func OptInnerMW(m ...Middleware) Options {
	return func(rmc *Conf) {
		rmc.InnerMiddleware = ChainMerge(rmc.InnerMiddleware, m...)
	}
}

func OptCodec(codecs Codecs) Options {
	return func(rmc *Conf) {
		rmc.Codecs = codecs
	}
}

func OptErrorEncoder(ee ErrorEncoder) Options {
	return func(rmc *Conf) {
		if ee != nil {
			rmc.ErrorEncoder = ee
		}
	}
}

// rmc Resource Management Center
type rmc struct {
	conf      Conf
	endpoints map[string]Conf
}

func NewRmc() Rmc {
	return &rmc{
		conf:      newRmcConf(),
		endpoints: make(map[string]Conf, 0),
	}
}

func (r *rmc) copy() *rmc {
	nr := new(rmc)
	nr.conf = r.conf
	nr.endpoints = r.endpoints
	return nr
}

func (r *rmc) Group(prefix string, options ...Options) Rmc {
	nr := r.copy()
	nr.conf.Path += prefix
	for _, v := range options {
		v(&nr.conf)
	}
	return nr
}

func (r *rmc) Use(options ...Options) Rmc {
	nr := r.copy()
	for _, v := range options {
		v(&nr.conf)
	}
	return nr
}

func (r *rmc) MustGetCodec(p transport.Kind, method transport.Method, path string) Codec {
	out, err := r.GetCodec(p, method, path)
	if err != nil {
		panic("miss codec method: 【" + string(method) + "】 path: 【" + path + "】")
	}
	return out
}

func (r *rmc) GetCodec(p transport.Kind, method transport.Method, path string) (Codec, error) {
	conf, err := r.getConfig(method, path)
	if err != nil {
		return Codec{}, err
	}

	if conf.Codecs == nil {
		return Codec{}, think.ErrRecordNotFound("endpoint 404")
	}
	codec, ok := conf.Codecs[p]
	if !ok {
		return Codec{}, think.ErrRecordNotFound("endpoint 404")
	}
	return codec, nil
}

func (r *rmc) MustGetEndpoint(p transport.Kind, method transport.Method, path string) HandlerFunc {
	out, err := r.GetEndpoint(p, method, path)
	if err != nil {
		panic("miss endpoint method: 【" + string(method) + "】 path: 【" + path + "】")
	}
	return out
}

func (r *rmc) GetEndpoint(p transport.Kind, method transport.Method, path string) (HandlerFunc, error) {
	conf, err := r.getConfig(method, path)
	if err != nil {
		return nil, err
	}

	if conf.Codecs == nil {
		return nil, think.ErrRecordNotFound("endpoint 404")
	}
	codec, ok := conf.Codecs[p]
	if !ok {
		return nil, think.ErrRecordNotFound("endpoint 404")
	}
	return r._endpoint(conf.HandlerFunc, codec.Dec, codec.Enc, conf.Middleware, conf.InnerMiddleware, conf.ErrorEncoder), nil
}

func (r *rmc) Endpoint(method transport.Method, path string, hf HandlerFunc, codecs Codecs, options ...Options) {
	key := Key(method, path)
	if _, ok := r.endpoints[key]; ok {
		log.Panicf("路径已经注册 %s", key)
	}

	nr := r.copy()

	for _, v := range options {
		v(&nr.conf)
	}
	nr.conf.Codecs = codecs
	nr.conf.Method = method
	nr.conf.Path += path
	nr.conf.FullPath = NewPath(nr.conf.Path, method)
	nr.conf.HandlerFunc = hf

	if nr.conf.Middleware == nil {
		nr.conf.Middleware = func(handlerFunc HandlerFunc) HandlerFunc {
			return handlerFunc
		}
	}
	nr.endpoints[key] = nr.conf
}

func (r *rmc) getConfig(method transport.Method, path string) (*Conf, error) {
	key := Key(method, path)
	conf, ok := r.endpoints[key]
	if !ok {
		return nil, think.ErrRecordNotFound("endpoint 404")
	}
	return &conf, nil

}

func Recover(ctx context.Context) error {
	if rErr := recover(); rErr != nil {
		buf := make([]byte, 64<<10)
		n := runtime.Stack(buf, false)
		buf = buf[:n]
		//log.Printf("Printf err: %v \n", rErr)
		//log.Println(string(buf))

	}
	return nil
}

func (r *rmc) _endpoint(hf HandlerFunc, dec Dec, enc Enc, middleware Middleware, innerMiddleware Middleware, errorEncoder ErrorEncoder) HandlerFunc {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		kind := transport.GetTransporter(ctx).Kind()

		resp, err = middleware(func(hf HandlerFunc) HandlerFunc {
			return func(ctx context.Context, data interface{}) (interface{}, error) {
				var err error

				if dec != nil {
					if data, err = dec(ctx, data); err != nil {
						return nil, err
					}
				}
				if innerMiddleware == nil {
					if data, err = hf(ctx, data); err != nil {
						return nil, err
					}
				} else {
					if data, err = innerMiddleware(hf)(ctx, data); err != nil {
						return nil, err
					}
				}

				if enc != nil {
					if data, err = enc(ctx, data); err != nil {
						return nil, err
					}
				}
				return data, nil
			}
		}(hf))(ctx, request)

		if err != nil {
			if errorEncoder != nil {
				err = errorEncoder(ctx, kind, err)
			}
		} else {
			if sender, ok := resp.(func() (any, error)); ok {
				if resp, err = sender(); err != nil {
					if errorEncoder != nil {
						err = errorEncoder(ctx, kind, err)
					}
				}
			}
		}

		return resp, err
	}
}

func (r *rmc) Proxy(proxy ProxyEndpoint, kind transport.Kind) {
	for _, v := range r.endpoints {
		conf, err := r.getConfig(v.Method, v.Path)
		if err != nil {
			return
		}

		if codec, ok := conf.Codecs[kind]; ok {
			proxy(kind, v.FullPath, r._endpoint(conf.HandlerFunc, codec.Dec, codec.Enc, conf.Middleware, conf.InnerMiddleware, v.ErrorEncoder))
		}
	}
}

func Key(method transport.Method, path string) string {
	return string(method) + ":" + path
}

func parseKey(key string) (transport.Method, string) {
	data := strings.Split(key, ":")
	if len(data) < 2 {
		return transport.Method(data[0]), ""
	}

	return transport.Method(data[0]), strings.Join(data[1:], ":")
}
