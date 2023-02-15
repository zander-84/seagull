package goredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang/groupcache/singleflight"
	"github.com/zander-84/seagull/think"
	"sync"
	"sync/atomic"
	"time"
)

type Rdb struct {
	engine          *redis.Client
	conf            Conf
	once            int64
	err             error
	lock            sync.Mutex
	context         context.Context
	funSingleFlight singleflight.Group
	getSingleFlight singleflight.Group
	hashCtl         *HashCtl
}

func NewRdb(conf Conf) *Rdb {
	this := new(Rdb)
	this.init(conf)
	return this
}

func (r *Rdb) init(conf Conf) {
	r.conf = conf.SetDefault()
	r.err = think.UnImpl
	r.context = context.Background()
	r.engine = nil

	atomic.StoreInt64(&r.once, 0)
}

func (r *Rdb) Start() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if atomic.CompareAndSwapInt64(&r.once, 0, 1) {
		r.engine = redis.NewClient(&redis.Options{
			Addr:         r.conf.Addr,
			Password:     r.conf.Password,
			DB:           r.conf.Db,
			PoolSize:     r.conf.PoolSize,
			PoolTimeout:  time.Duration(r.conf.PoolTimeout) * time.Second,
			MinIdleConns: r.conf.MinIdle,
		})

		r.err = r.engine.Ping(context.Background()).Err()
		if r.err == nil {
			r.hashCtl, r.err = newHashCtl(r.engine)
		}
	}

	return r.err
}

func (r *Rdb) Stop() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.engine != nil {
		_ = r.engine.Close()
	}
	r.engine = nil
	atomic.StoreInt64(&r.once, 0)
	r.err = think.UnImpl
	return nil
}

func (r *Rdb) Restart(conf Conf) error {
	r.Stop()
	r.init(conf)
	return r.Start()
}

func (r *Rdb) Engine() *redis.Client {
	return r.engine
}
