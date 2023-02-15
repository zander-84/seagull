package project

var PkgCorsTpl = `package pkg

import (
	"github.com/zander-84/seagull/endpoint/middleware/cors"
	"github.com/zander-84/seagull/think"
	"time"
)

func GinCors(mode think.Mode) cors.Config {
	allowOrigins := []string{"*"}
	if mode.IsProd() {
		allowOrigins = []string{"https://xxxxx.52haoka.com"}
	} else if mode.IsDev() {
		//allowOrigins = []string{"https://xxxxxx.52haoka.com"}
	}

	return cors.Config{
		AllowOrigins:  allowOrigins,
		AllowMethods:  []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Vas-Authorization", "Trace-Id"},
		ExposeHeaders: []string{"Trace-Id", "Content-Disposition", "Vas-Authorization"},

		//AllowOriginFunc: func(origin string) bool {
		//	return true
		//},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}
}
`

var PkgEtcdTpl = `package pkg

import (
	"context"
	"errors"
	"github.com/zander-84/seagull/drive/etcd"
	"time"
)

func newEtcdInstance(conf etcd.Conf) (*etcd.Etcd, error) {
	instance := etcd.NewEtcd(conf)
	if err := instance.Start(); err != nil {
		return nil, errors.New("etcd start err:" + err.Error())
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := instance.Engine().Get(ctx, "test")
	if err != nil {
		return nil, errors.New("etcd start err:" + err.Error())
	}

	return instance, nil
}
`

var PkgMongoTpl = `package pkg

import (
	"errors"
	"github.com/zander-84/seagull/drive/mongo"
)

func newMongoInstance(conf mongo.Conf) (*mongo.Mongo, error) {
	instance := mongo.NewMongo(conf)
	if err := instance.Start(); err != nil {
		return nil, errors.New("Mongo start err:" + err.Error())
	}

	return instance, nil
}

`

var PkgMysqlTpl = `package pkg

import (
	"errors"
	"github.com/zander-84/seagull/drive/gormmysql"
)

func newMysqlInstance(conf gormmysql.Conf) (*gormmysql.Gdb, error) {
	instance := gormmysql.NewGdb(conf)
	if err := instance.Start(); err != nil {
		return nil, errors.New("mysql start err:" + err.Error())
	}

	return instance, nil
}
`

var PkgRedisTpl = `package pkg

import (
	"errors"
	"github.com/zander-84/seagull/drive/goredis"
)

func newRedisInstance(conf goredis.Conf) (*goredis.Rdb, error) {
	instance := goredis.NewRdb(conf)
	if err := instance.Start(); err != nil {
		return nil, errors.New("redis start err:" + err.Error())
	}
	return instance, nil
}
`

var PkgInitTpl = `package pkg

import (
	"github.com/ulule/limiter/v3"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/drive/etcd"
	"github.com/zander-84/seagull/drive/goredis"
	"github.com/zander-84/seagull/drive/gormmysql"
	"github.com/zander-84/seagull/drive/mongo"
)

var (
	err             error
	Processor       contract.Processor
	ProcessorCancel func() error

	Mysql *gormmysql.Gdb

	Mongo *mongo.Mongo

	Redis *goredis.Rdb

	Cache   contract.Cache
	Limiter *limiter.Limiter

	Etcd *etcd.Etcd
)

func InitLib() {
	//Processor = worker.NewProcessor()
	//ProcessorCancel = func() error { return Processor.Wait(time.Minute * 2) }
	//
	//if Mysql, err = newMysqlInstance(config.Data.Mysql); err != nil {
	//	log.Fatal(err.Error())
	//}
	//if Mongo, err = newMongoInstance(config.Data.Mongo); err != nil {
	//	log.Fatal(err.Error())
	//}
	//
	//if Etcd, err = newEtcdInstance(config.Data.Etcd); err != nil {
	//	log.Fatal(err.Error())
	//}
	//
	//if Redis, err = newRedisInstance(config.Data.Redis); err != nil {
	//	log.Fatal(err.Error())
	//}
	//
	//if Limiter, err = middleware.NewRedisLimiter("1-S", Redis.Engine()); err != nil {
	//	log.Fatal(err.Error())
	//}
	//Cache = cache.NewRedisCache(Redis.Engine(), codec.GetCodec(codec.Json), Processor, 3)

}

`
