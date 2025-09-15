package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type CronJob struct {
	cron *cron.Cron
}

func NewCron() *CronJob {
	c := CronJob{
		cron: cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow))),
	}
	return &c
}

func (c *CronJob) AddNewJob(schedule string) {
	_, err := c.cron.AddFunc(schedule, func() {
		fmt.Println("Job executed at:", time.Now().Format("15:04:05"))
	})
	if err != nil {
		panic(err)
	}
}
