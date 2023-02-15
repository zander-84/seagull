package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/zander-84/seagull/think"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"sync/atomic"
	"time"
)

type Etcd struct {
	engine  *clientv3.Client
	conf    Conf
	once    int64
	err     error
	lock    sync.Mutex
	context context.Context
}

func (e *Etcd) init(conf Conf) {
	e.conf = conf.SetDefault()
	e.err = think.UnImpl
	e.context = context.Background()
	e.engine = nil
	atomic.StoreInt64(&e.once, 0)
}
func NewEtcd(conf Conf) *Etcd {
	this := new(Etcd)
	this.init(conf)
	return this
}

func (e *Etcd) Start() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if atomic.CompareAndSwapInt64(&e.once, 0, 1) {
		var tlsconf *tls.Config
		if e.conf.TlsCa != "" && e.conf.TlsKey != "" && e.conf.TlsPem != "" {
			ce, err := tls.X509KeyPair([]byte(e.conf.TlsPem), []byte(e.conf.TlsKey))
			if err != nil {
				e.err = err
				return e.err
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM([]byte(e.conf.TlsCa))
			tlsconf = &tls.Config{
				Certificates: []tls.Certificate{ce},
				RootCAs:      pool,
			}
		}
		e.engine, e.err = clientv3.New(clientv3.Config{
			Endpoints:            e.conf.Endpoints,
			AutoSyncInterval:     0,
			DialTimeout:          20 * time.Second,
			DialKeepAliveTime:    0,
			DialKeepAliveTimeout: 0,
			MaxCallSendMsgSize:   0,
			MaxCallRecvMsgSize:   0,
			TLS:                  tlsconf,
			Username:             "",
			Password:             "",
			RejectOldCluster:     false,
			DialOptions:          nil,
			Context:              nil,
			Logger:               nil,
			LogConfig:            nil,
			PermitWithoutStream:  false,
		})
		if e.err != nil {
			return e.err
		}
		e.err = e.Ping()

	}
	return e.err
}
func (e *Etcd) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := e.engine.Get(ctx, "/ping")
	return err
}
func (e *Etcd) Stop() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.engine != nil {
		_ = e.engine.Close()
	}
	e.engine = nil
	atomic.StoreInt64(&e.once, 0)
	e.err = think.UnImpl
	return nil
}

func (e *Etcd) Restart(conf Conf) error {
	e.Stop()
	e.init(conf)
	return e.Start()
}
func (e *Etcd) Engine() *clientv3.Client {
	return e.engine
}
