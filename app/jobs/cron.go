package jobs

import (
	"github.com/jasonlvhit/gocron"
	"github.com/lbryio/lighthouse/app/jobs/blocked"
	"github.com/lbryio/lighthouse/app/jobs/chainquery"
	"github.com/lbryio/lighthouse/app/jobs/internalapis"
	"github.com/sirupsen/logrus"
)

var cronRunning chan bool
var scheduler *gocron.Scheduler

// Start starts the jobs that run in the background after initialization
func Start() {
	scheduler = gocron.NewScheduler()
	var channels *string
	scheduler.Every(15).Minutes().Do(chainquery.Sync, channels)
	scheduler.Every(6).Hours().Do(internalapis.Sync)
	scheduler.Every(1).Minutes().Do(blocked.ProcessBlockedList)
	scheduler.Every(1).Minutes().Do(blocked.ProcessFilteredList)

	cronRunning = scheduler.Start()
}

// Shutdown is used to shutdown the background jobs.
func Shutdown() {
	logrus.Debug("Shutting down cron jobs...")
	scheduler.Clear()
	close(cronRunning)
}
