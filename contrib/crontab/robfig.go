package crontab

import (
	"github.com/robfig/cron/v3"
	"github.com/zander-84/seagull/think"
	"sync"
	"sync/atomic"
	"time"
)

type Crontab struct {
	engine       *cron.Cron
	engineParser cron.Parser
	conf         Conf
	jobs         map[string]*job
	err          error
	lock         sync.Mutex
	once         int64
}

func NewCrontab(conf Conf) *Crontab {
	c := new(Crontab)
	c.init(conf)
	return c
}

func (c *Crontab) init(conf Conf) {
	c.conf = conf.SetDefault()
	c.err = think.UnImpl
	c.jobs = make(map[string]*job)
	atomic.StoreInt64(&c.once, 0)
}

func (c *Crontab) Start() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if atomic.CompareAndSwapInt64(&c.once, 0, 1) {
		c.engine = cron.New(cron.WithLocation(time.Now().Location()), cron.WithSeconds())
		c.engineParser = cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		)
		c.err = nil
	}
	return c.err
}
func (c *Crontab) CheckParse(spec string) error {
	_, err := c.engineParser.Parse(spec)
	return err
}
func (c *Crontab) Stop() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.engine != nil {
		c.engine.Stop()
	}

	c.engine = nil
	atomic.StoreInt64(&c.once, 0)
	c.err = think.UnImpl
	c.jobs = make(map[string]*job)
	return nil
}

func (c *Crontab) Restart(conf Conf) error {
	c.Stop()
	c.init(conf)
	return c.Start()
}

func (c *Crontab) Engine() *cron.Cron {
	return c.engine
}
