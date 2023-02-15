package crontab

import (
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
)

var _ contract.Crontab = (*Crontab)(nil)

func (c *Crontab) AddJob(cronJob contract.CronJob) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.jobs[cronJob.ID]; ok {
		return think.RecordExist
	}
	// 添加job
	_job := newJob(cronJob)
	id, err := c.engine.AddJob(cronJob.Spec, _job)
	if err == nil {
		_job.id = c.engine.Entry(id).ID
		c.jobs[cronJob.ID] = _job
	}
	return err
}

// RemoveJob 移除
func (c *Crontab) RemoveJob(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _job, ok := c.jobs[id]; ok {
		c.engine.Remove(_job.id)
		delete(c.jobs, id)
	}

	return c.err
}

func (c *Crontab) StatusJobs() map[string]contract.CronJob {
	c.lock.Lock()
	defer c.lock.Unlock()
	out := make(map[string]contract.CronJob)
	for _, _job := range c.jobs {
		entry := c.engine.Entry(_job.id)
		cronJob := _job.cronJob
		cronJob.Next = entry.Next
		cronJob.Prev = entry.Prev
		out[_job.cronJob.ID] = cronJob
	}
	return out
}

func (c *Crontab) StartJobs() error {
	c.engine.Start()
	return c.err
}

func (c *Crontab) StopJobs() error {
	c.engine.Stop()
	return c.err

}
