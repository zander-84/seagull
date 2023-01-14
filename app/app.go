package app

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/transport"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Info is application context value.
type Info interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

// App is an application components lifecycle manager.
type App struct {
	opts     options
	ctx      context.Context
	cancel   func()
	mu       sync.Mutex
	instance *contract.ServiceInstance
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	o := options{
		ctx:               context.Background(),
		sigs:              []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout:  10 * time.Second,
		stopTimeout:       10 * time.Second,
		beforeStartEvents: make(map[int][]Event, 0),
		afterStartEvents:  make(map[int][]Event, 0),
		beforeStopEvents:  make(map[int][]Event, 0),
		afterStopEvents:   make(map[int][]Event, 0),
		finalEvents:       make(map[int][]Event, 0),
		eventsTimeOut:     time.Minute,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.
func (a *App) Endpoint() []string {
	if a.instance != nil {
		return a.instance.Endpoints
	}
	return nil
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()
	eg, ctx := errgroup.WithContext(NewContext(a.ctx, a))

	a.beforeStart()

	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(NewContext(a.opts.ctx, a), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(NewContext(a.opts.ctx, a))
		})
	}
	wg.Wait()
	if a.opts.registrar != nil {
		rctx, rcancel := context.WithTimeout(ctx, a.opts.registrarTimeout)
		defer rcancel()
		if err := a.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}

	a.afterStart()

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			a.beforeStop()
			return a.Stop()
		}
	})

	err = eg.Wait()

	a.afterStop()
	a.finalStop()
	return err
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(a.ctx, a), a.opts.registrarTimeout)
		defer cancel()
		if err := a.opts.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *App) buildInstance() (*contract.ServiceInstance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	return &contract.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s Info) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s Info, ok bool) {
	s, ok = ctx.Value(appKey{}).(Info)
	return
}
func (a *App) IsStop() bool {
	select {
	case <-a.ctx.Done():
		return true
	default:
		return false
	}
}

func (a *App) Context() context.Context {
	return a.ctx
}

// beforeStart 服务停止前事件
func (a *App) beforeStart() {
	keys := getAscKey(a.opts.beforeStartEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		a.doEvents(a.opts.beforeStartEvents[v])
	}

}

// afterStart 服务启动前事件
func (a *App) afterStart() {
	keys := getAscKey(a.opts.afterStartEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		a.doEvents(a.opts.afterStartEvents[v])
	}
}

// beforeStop 服务停止后事件
func (a *App) beforeStop() {
	keys := getAscKey(a.opts.beforeStopEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		a.doEvents(a.opts.beforeStopEvents[v])
	}
}

// afterStop 服务停止后事件
func (a *App) afterStop() {
	keys := getAscKey(a.opts.afterStopEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		a.doEvents(a.opts.afterStopEvents[v])
	}
}

// finalStop 服务停止后事件
func (a *App) finalStop() {
	keys := getAscKey(a.opts.finalEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		a.doEvents(a.opts.finalEvents[v])
	}
}

func (a *App) doEvents(events []Event) {
	ctx, cancel := context.WithTimeout(context.Background(), a.opts.eventsTimeOut)
	defer cancel()
	fin := make(chan struct{}, 1)
	wg := sync.WaitGroup{}

	for _, event := range events {
		wg.Add(1)
		go func(e func() error) {
			defer wg.Done()
			defer func() {
				if rerr := recover(); rerr != nil {
				}
			}()
			_ = e()
		}(event.handler)
	}

	go func() {
		wg.Wait()
		fin <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-fin:
		return
	}
}

func getAscKey(in map[int][]Event) []int {
	if len(in) < 1 {
		return []int{}
	}
	keys := make([]int, 0)
	for k, _ := range in {
		keys = append(keys, k)
	}

	sort.Sort(sort.IntSlice(keys))
	return keys
}
