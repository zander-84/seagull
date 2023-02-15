package app

import (
	"context"
	"sort"
	"sync"
	"time"
)

// int  从小到大排序

type Bootstrap struct {
	defaultTimeout time.Duration

	initEvents  map[int][]Event
	initTimeout time.Duration

	beforeStartEvents  map[int][]Event
	beforeStartTimeout time.Duration

	afterStartEvents  map[int][]Event
	afterStartTimeout time.Duration

	beforeStopEvents  map[int][]Event
	beforeStopTimeout time.Duration

	afterStopEvents  map[int][]Event
	afterStopTimeout time.Duration

	finalEvents  map[int][]Event
	finalTimeout time.Duration
}

func NewBootstrap() *Bootstrap {
	bs := new(Bootstrap)
	bs.defaultTimeout = time.Minute
	bs.initEvents = make(map[int][]Event, 0)
	bs.beforeStartEvents = make(map[int][]Event, 0)
	bs.afterStartEvents = make(map[int][]Event, 0)
	bs.beforeStopEvents = make(map[int][]Event, 0)
	bs.afterStopEvents = make(map[int][]Event, 0)
	bs.finalEvents = make(map[int][]Event, 0)
	return bs
}
func (b *Bootstrap) Init() {
	keys := b.getAscKey(b.initEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		b.doEvents(b.initEvents[v], b.getTimeout(b.initTimeout, b.defaultTimeout))
	}
}
func (b *Bootstrap) SetDefaultTimeOut(timeout time.Duration) {
	b.defaultTimeout = timeout
}

func (b *Bootstrap) RegisterInitEvents(key int, events ...Event) {
	es, ok := b.initEvents[key]
	if !ok {
		es = make([]Event, 0)
	}
	b.initEvents[key] = append(es, events...)
}
func (b *Bootstrap) SetInitTimeOut(timeout time.Duration) {
	b.initTimeout = timeout
}

func (b *Bootstrap) RegisterBeforeStartEvents(key int, events ...Event) {
	es, ok := b.beforeStartEvents[key]
	if !ok {
		es = make([]Event, 0)
	}
	b.beforeStartEvents[key] = append(es, events...)
}

func (b *Bootstrap) SetBeforeStartTimeout(timeout time.Duration) {
	b.beforeStartTimeout = timeout
}

func (b *Bootstrap) RegisterAfterStartEvents(key int, events ...Event) {
	es, ok := b.afterStartEvents[key]
	if !ok {
		es = make([]Event, 0)
	}
	b.afterStartEvents[key] = append(es, events...)
}

func (b *Bootstrap) SetAfterStartTimeout(timeout time.Duration) {
	b.afterStartTimeout = timeout
}

func (b *Bootstrap) RegisterBeforeStopEvents(key int, events ...Event) {
	es, ok := b.beforeStopEvents[key]
	if !ok {
		es = make([]Event, 0)
	}
	b.beforeStopEvents[key] = append(es, events...)
}

func (b *Bootstrap) SetBeforeStopTimeout(timeout time.Duration) {
	b.beforeStopTimeout = timeout
}

func (b *Bootstrap) RegisterAfterStopEvents(key int, events ...Event) {
	es, ok := b.afterStopEvents[key]
	if !ok {
		es = make([]Event, 0)
	}
	b.afterStopEvents[key] = append(es, events...)
}

func (b *Bootstrap) SetAfterStopTimeout(timeout time.Duration) {
	b.afterStopTimeout = timeout
}

func (b *Bootstrap) RegisterFinalEvents(key int, events ...Event) {
	es, ok := b.finalEvents[key]
	if !ok {
		es = make([]Event, 0)
	}
	b.finalEvents[key] = append(es, events...)
}
func (b *Bootstrap) SetFinalTimeout(timeout time.Duration) {
	b.finalTimeout = timeout
}

func (b *Bootstrap) execBeforeStartEvents() {
	keys := b.getAscKey(b.beforeStartEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		b.doEvents(b.beforeStartEvents[v], b.getTimeout(b.beforeStartTimeout, b.defaultTimeout))
	}
}
func (b *Bootstrap) execAfterStartEvents() {
	keys := b.getAscKey(b.afterStartEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		b.doEvents(b.afterStartEvents[v], b.getTimeout(b.afterStartTimeout, b.defaultTimeout))
	}
}

func (b *Bootstrap) execBeforeStopEvents() {
	keys := b.getAscKey(b.beforeStopEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		b.doEvents(b.beforeStopEvents[v], b.getTimeout(b.beforeStopTimeout, b.defaultTimeout))
	}
}
func (b *Bootstrap) execAfterStopEvents() {
	keys := b.getAscKey(b.afterStopEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		b.doEvents(b.afterStopEvents[v], b.getTimeout(b.afterStopTimeout, b.defaultTimeout))
	}
}

func (b *Bootstrap) execFinalEvents() {
	keys := b.getAscKey(b.finalEvents)
	if len(keys) < 1 {
		return
	}
	for _, v := range keys {
		b.doEvents(b.finalEvents[v], b.getTimeout(b.finalTimeout, b.defaultTimeout))
	}
}

func (b *Bootstrap) getTimeout(timeout time.Duration, defaultTimeout time.Duration) time.Duration {
	if timeout < 1 {
		return defaultTimeout
	}
	return timeout
}

func (b *Bootstrap) doEvents(events []Event, eventsTimeOut time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), eventsTimeOut)
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

func (b *Bootstrap) getAscKey(in map[int][]Event) []int {
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
