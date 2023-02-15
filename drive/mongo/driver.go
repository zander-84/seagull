package mongo

import (
	"context"
	"fmt"
	"github.com/tidwall/pretty"
	"github.com/zander-84/seagull/think"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"sync/atomic"
)

type Mongo struct {
	engine  *mongo.Client
	conf    Conf
	once    int64
	err     error
	lock    sync.Mutex
	context context.Context
}

func (m *Mongo) init(conf Conf) {
	m.conf = conf.SetDefault()
	m.err = think.UnImpl
	atomic.StoreInt64(&m.once, 0)
	m.engine = nil
	m.context = context.Background()
}

func NewMongo(conf Conf) *Mongo {
	m := new(Mongo)
	m.init(conf)
	return m
}

func (m *Mongo) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if atomic.CompareAndSwapInt64(&m.once, 0, 1) {

		dns := fmt.Sprintf("mongodb://%s:%s", m.conf.Host, m.conf.Port)

		mongoOptions := new(options.ClientOptions)
		mongoOptions.ApplyURI(dns)

		monitor := &event.CommandMonitor{}
		if m.conf.Debug {
			monitor.Started = func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				data := pretty.Pretty([]byte(startedEvent.Command.String()))
				log.Println(fmt.Sprintf("mongo请求：Id:%d, 内容：%s", startedEvent.RequestID, string(data)))
			}
		}

		if m.conf.DebugReply {
			monitor.Succeeded = func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
				data := pretty.Pretty([]byte(succeededEvent.Reply.String()))
				log.Println(fmt.Sprintf("mongo响应：Id:%d, 内容：%s", succeededEvent.RequestID, string(data)))
			}
		}

		if m.conf.Debug || m.conf.DebugReply {
			mongoOptions.SetMonitor(monitor)
		}

		if m.conf.User != "" && m.conf.Pwd != "" {
			mongoOptions.SetAuth(options.Credential{
				AuthMechanism:           "",
				AuthMechanismProperties: nil,
				AuthSource:              m.conf.Database,
				Username:                m.conf.User,
				Password:                m.conf.Pwd,
				PasswordSet:             false,
			})
		}
		MaxPoolSize := m.conf.MaxPoolSize
		MinPoolSize := m.conf.MinPoolSize
		mongoOptions.MaxPoolSize = &MaxPoolSize
		mongoOptions.MinPoolSize = &MinPoolSize

		m.engine, m.err = mongo.Connect(context.Background(), mongoOptions)
		if m.err != nil {
			return m.err
		}

		if m.err = m.engine.Ping(context.Background(), nil); m.err != nil {
			return m.err
		}
	}
	return m.err
}

func (m *Mongo) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.engine != nil {
		_ = m.engine.Disconnect(m.context)
	}
	m.engine = nil
	atomic.StoreInt64(&m.once, 0)
	m.err = think.UnImpl
	return nil
}

func (m *Mongo) Restart(conf Conf) error {
	m.Stop()
	m.init(conf)
	return m.Start()
}

func (m *Mongo) Engine() *mongo.Client {
	return m.engine
}

func (m *Mongo) DB() *mongo.Database {
	return m.engine.Database(m.conf.Database)
}

func (m *Mongo) Collection(collection string) *mongo.Collection {
	return m.engine.Database(m.conf.Database).Collection(collection)
}

func (m *Mongo) GetDB(dbname string) *mongo.Database {
	return m.engine.Database(dbname)
}
func (m *Mongo) GetCollection(dbname string, collection string) *mongo.Collection {
	return m.GetDB(dbname).Collection(collection)
}
