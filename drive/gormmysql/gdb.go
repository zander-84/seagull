package gormmysql

import (
	"database/sql"
	"fmt"
	"github.com/zander-84/seagull/think"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Gdb struct {
	engine *gorm.DB
	sqlDB  *sql.DB
	conf   Conf
	once   int64
	err    error
	lock   sync.Mutex
}

func (g *Gdb) init(conf Conf) {
	g.conf = conf.SetDefault()
	g.err = think.UnImpl
	atomic.StoreInt64(&g.once, 0)
	g.engine = nil
	g.sqlDB = nil
}

func NewGdb(conf Conf) *Gdb {
	out := new(Gdb)
	out.init(conf)
	return out
}

func (g *Gdb) Start() error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if atomic.CompareAndSwapInt64(&g.once, 0, 1) {

		// 时间配置

		// 配置文件
		gormCnf := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: true,
			},
			AllowGlobalUpdate: false,
			NowFunc: func() time.Time {
				return time.Now()
			},
		}
		// debug
		var LogLevel logger.LogLevel
		if g.conf.Debug {
			LogLevel = logger.Info
		} else {
			LogLevel = logger.Silent
		}
		gormCnf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      LogLevel,    // Log level
				Colorful:      true,        // 禁用彩色打印
			},
		)

		// mysql conf
		mysqlCnf := mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s", g.conf.User, g.conf.Pwd, g.conf.Host, g.conf.Port, g.conf.Database, g.conf.Charset, url.QueryEscape(g.conf.TimeZone)), // DSN data source name
			//DefaultStringSize:         256,   // string 类型字段的默认长度
			//DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
		})

		// 开始初始化
		g.engine, g.err = gorm.Open(mysqlCnf, gormCnf)
		if g.err != nil {
			return g.err
		}
		g.sqlDB, g.err = g.engine.DB()
		if g.err != nil {
			return g.err
		}

		g.sqlDB.SetMaxIdleConns(g.conf.MaxIdleconns)
		g.sqlDB.SetMaxOpenConns(g.conf.MaxOpenconns)
		g.sqlDB.SetConnMaxLifetime(time.Duration(g.conf.ConnMaxLifetime) * time.Second)

		//if this.conf.RemoveSomeCallbacks {
		//	_ = this.engine.Callback().Create().Remove("gorm:save_before_associations")
		//	_ = this.engine.Callback().Create().Remove("gorm:force_reload_after_create")
		//	_ = this.engine.Callback().Create().Remove("gorm:save_after_associations")
		//	_ = this.engine.Callback().Update().Remove("gorm:save_before_associations")
		//	_ = this.engine.Callback().Update().Remove("gorm:save_after_associations")
		//}
	}
	return g.err
}

func (g *Gdb) Stop() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.sqlDB != nil {
		_ = g.sqlDB.Close()
	}
	atomic.StoreInt64(&g.once, 0)
	g.err = think.UnImpl
	g.engine = nil
	g.sqlDB = nil
	return nil
}

func (g *Gdb) Restart(conf Conf) error {
	g.Stop()
	g.init(conf)
	return g.Start()
}

func (g *Gdb) Engine() *gorm.DB {
	return g.engine
}

func (g *Gdb) SqlDB() *sql.DB {
	return g.sqlDB
}
