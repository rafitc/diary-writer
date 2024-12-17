package writer

import (
	"main/logger"

	"github.com/robfig/cron"
)

type Writer struct {
	Name string
}

var log = logger.Logger

func (writer *Writer) StartCronJob() {
	log.Infof("starting %v job", writer.Name)
	// Starting a cronJob to check is there anything in db
	c := cron.New()
	c.AddFunc("0/5 * * * * *", func() { dailyDiaryDataChecker(writer.Name) })
	log.Infof("Started %s cronJob", writer.Name)
	c.Start()

}

func dailyDiaryDataChecker(name string) {
	log.Infof("Starting %s cron", name)
}
