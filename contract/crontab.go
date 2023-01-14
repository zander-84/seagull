package contract

import "time"

/*
CronJob 描述
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/
type CronJob struct {
	ID   string
	Desc string
	Spec string
	Next time.Time
	Prev time.Time
	Cmd  Job
}

type Crontab interface {
	AddJob(jobInfo CronJob) error

	RemoveJob(id string) error

	StatusJobs() map[string]CronJob

	StartJobs() error

	StopJobs() error

	CheckParse(spec string) error
}
