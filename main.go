/*
Daiary Writer
Author :- Rafi Rasheed TC
Repo :- https://github.com/rafitc

Script to get message from Telegram and write into SQLite3. then a scheduled process to upload the data into git and deploy the new app on daily basis.

Two tasks running on startup.
1. Telegram Bot
3. Cron to read data and publish in the web

Telegram Bot
*/

package main

import (
	"main/logger"
	"main/writer"
	"time"
)

var log = logger.Logger

func main() {
	// Start a thread to check for data in sqlite server
	// If data is present, Then compare it with current time, if the data is old then push into github and trigger build

	log.Info("Starting the writer")
	writer := writer.Writer{}
	writer.StartCronJob()

	// Wait for the Cron job to run
	time.Sleep(5 * time.Minute)

}
