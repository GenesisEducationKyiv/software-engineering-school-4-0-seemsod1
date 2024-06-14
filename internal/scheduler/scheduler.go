package scheduler

import (
	"log"
	"strconv"

	"github.com/go-co-op/gocron/v2"
)

// GoCron is a struct that embeds the gocron.Scheduler
type GoCron struct {
	Scheduler gocron.Scheduler
}

func NewGoCronScheduler() *GoCron {
	sch, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	return &GoCron{Scheduler: sch}
}

func (gs *GoCron) Start() {
	gs.Scheduler.Start()
}

func (gs *GoCron) AddEverydayJob(task func(), minute int) error {
	if _, err := gs.Scheduler.NewJob(
		gocron.CronJob(
			strconv.Itoa(minute)+" * * * *", // every hour at minute 1
			false,
		),
		gocron.NewTask(task),
	); err != nil {
		return err
	}
	return nil
}
