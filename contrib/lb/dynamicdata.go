package lb

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var _ Listener = (*DynamicData)(nil)

type DynamicData struct {
	name       string
	data       map[any]int
	version    uint64
	lock       sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	err        error
	errLock    sync.RWMutex
}

// NewListener name 在etcd中是prefix key
func NewListener(name string) *DynamicData {
	dd := new(DynamicData)
	dd.version = 0
	dd.data = make(map[any]int)
	dd.name = name

	dd.ctx, dd.cancelFunc = context.WithCancel(context.Background())
	return dd
}
func (dd *DynamicData) Name() string {
	return dd.name
}

func (dd *DynamicData) Close() {
	dd.lock.Lock()
	defer dd.lock.Unlock()
	dd.cancelFunc()
	atomic.StoreUint64(&dd.version, 0)
	dd.data = make(map[any]int)

}

func (dd *DynamicData) Context() context.Context {
	return dd.ctx
}

func (dd *DynamicData) Exist(addr any) bool {
	dd.lock.RLock()
	defer dd.lock.RUnlock()
	_, ok := dd.data[addr]
	return ok
}

// Set  全量设置
func (dd *DynamicData) Set(data map[any]int) error {
	dd.lock.Lock()
	defer dd.lock.Unlock()
	dd.data = data
	atomic.AddUint64(&dd.version, 1)
	return nil
}

// Add 增量添加
func (dd *DynamicData) Add(addr any) error {
	return dd.AddWithWeight(addr, 1)
}

func (dd *DynamicData) AddWithWeight(addr any, weight int) error {
	dd.lock.Lock()
	defer dd.lock.Unlock()
	if weight == 0 {
		weight = 1
	}
	dd.data[addr] = weight
	atomic.AddUint64(&dd.version, 1)
	return nil
}

func (dd *DynamicData) Remove(addr any) error {
	dd.lock.Lock()
	defer dd.lock.Unlock()
	delete(dd.data, addr)
	atomic.AddUint64(&dd.version, 1)
	return nil
}

func (dd *DynamicData) Version() uint64 {
	return atomic.LoadUint64(&dd.version)
}

func (dd *DynamicData) Get() (map[any]int, []any, uint64) {
	dd.lock.RLock()
	defer dd.lock.RUnlock()

	var dataMap = make(map[any]int, 0)
	var dataSlice = make([]any, 0)
	for k, v := range dd.data {
		dataMap[k] = v
		dataSlice = append(dataSlice, k)
	}

	return dataMap, dataSlice, atomic.LoadUint64(&dd.version)
}

func (dd *DynamicData) SetErr(err error) {
	dd.errLock.Lock()
	defer dd.errLock.Unlock()
	dd.err = err
}

func (dd *DynamicData) Err() error {
	dd.errLock.RLock()
	defer dd.errLock.RUnlock()
	return dd.err
}

func (dd *DynamicData) Println() {
	addr, _, version := dd.Get()
	for k, v := range addr {
		fmt.Printf("地址:%s  权重:%d  版本:%d\n", k, v, version)
	}
	fmt.Println()
}
