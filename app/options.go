package app

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/transport"
	"net/url"
	"os"
	"time"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	registrar        contract.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []transport.Server

	// int  从小到大排序
	eventsTimeOut time.Duration

	beforeStartEvents map[int][]Event
	afterStartEvents  map[int][]Event
	beforeStopEvents  map[int][]Event
	afterStopEvents   map[int][]Event
	finalEvents       map[int][]Event
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Endpoint with service endpoint.
func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Server with transport servers.
func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Registrar with service registry.
func Registrar(r contract.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

// RegistrarTimeout with registrar timeout.
func RegistrarTimeout(t time.Duration) Option {
	return func(o *options) { o.registrarTimeout = t }
}

func EventsTimeOut(eventsTimeOut time.Duration) Option {
	return func(o *options) {
		o.eventsTimeOut = eventsTimeOut
	}
}

// StopTimeout with app stop timeout.
func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

func AppendBeforeStartEvents(key int, events ...Event) Option {
	return func(o *options) {
		if o.beforeStartEvents == nil {
			o.beforeStartEvents = make(map[int][]Event, 0)
		}
		es, ok := o.beforeStartEvents[key]
		if !ok {
			es = make([]Event, 0)
		}
		o.beforeStartEvents[key] = append(es, events...)
	}
}

func AppendAfterStartEvents(key int, events ...Event) Option {
	return func(o *options) {
		if o.afterStartEvents == nil {
			o.afterStartEvents = make(map[int][]Event, 0)
		}
		es, ok := o.afterStartEvents[key]
		if !ok {
			es = make([]Event, 0)
		}
		o.afterStartEvents[key] = append(es, events...)
	}
}

func AppendBeforeStopEvents(key int, events ...Event) Option {
	return func(o *options) {
		if o.beforeStopEvents == nil {
			o.beforeStopEvents = make(map[int][]Event, 0)
		}
		es, ok := o.beforeStopEvents[key]
		if !ok {
			es = make([]Event, 0)
		}
		o.beforeStopEvents[key] = append(es, events...)
	}
}

func AppendAfterStopEvents(key int, events ...Event) Option {
	return func(o *options) {
		if o.afterStopEvents == nil {
			o.afterStopEvents = make(map[int][]Event, 0)
		}
		es, ok := o.afterStopEvents[key]
		if !ok {
			es = make([]Event, 0)
		}
		o.afterStopEvents[key] = append(es, events...)
	}
}

func AppendFinalEvents(key int, events ...Event) Option {
	return func(o *options) {
		if o.finalEvents == nil {
			o.finalEvents = make(map[int][]Event, 0)
		}
		es, ok := o.finalEvents[key]
		if !ok {
			es = make([]Event, 0)
		}
		o.finalEvents[key] = append(es, events...)
	}
}
