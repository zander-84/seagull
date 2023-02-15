package memory

import (
	"github.com/patrickmn/go-cache"
	"github.com/zander-84/seagull/think"
	"sync"
	"sync/atomic"
	"time"
)

type Memory struct {
	engine *cache.Cache
	conf   Conf
	once   int64
	err    error
	lock   sync.Mutex
}

func NewMemory(conf Conf) *Memory {
	this := &Memory{}
	this.init(conf)
	return this
}

func (m *Memory) init(conf Conf) {
	m.conf = conf.SetDefault()
	m.err = think.UnImpl
	atomic.StoreInt64(&m.once, 0)
	m.engine = nil
}

func (m *Memory) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if atomic.CompareAndSwapInt64(&m.once, 0, 1) {
		m.engine = cache.New(time.Duration(m.conf.Expiration)*time.Minute, time.Duration(m.conf.CleanupInterval)*time.Minute)
		m.err = nil
	}

	return m.err
}

func (m *Memory) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.engine = nil
	atomic.StoreInt64(&m.once, 0)
	m.err = think.UnImpl
	return nil
}

func (m *Memory) Restart(conf Conf) error {
	m.Stop()
	m.init(conf)
	return m.Start()
}

func (m *Memory) Engine() *cache.Cache {
	return m.engine
}
