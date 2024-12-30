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
	"context"
	"main/bot"
	"main/logger"
	"main/writer"
	"os"
	"os/signal"
)

var log = logger.Logger

func main() {
	// Start a thread to check for data in sqlite server
	// If data is present, Then compare it with current time, if the data is old then push into github and trigger build

	log.Info("Starting the writer")
	writer.StartCronJob()

	// Wait for the Cron job to run
	// time.Sleep(5 * time.Minute)

	// Start the telegram bot
	log.Info("Starting the telegram bot")
	// create a context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Create a channel to fetch the data from the bot
	dataChannel := make(chan bot.MessageData)

	go func() {
		for dataChannel != nil {
			select {
			case data := <-dataChannel:
				log.Info("Writing data into the database")
				writer.InsertDataIntoDB(data.Message, data.Asset, data.Asset_extension, data.Asset_download_link, data.Date)
				// writer.WriteData(data)
			case <-ctx.Done():
				log.Info("Exiting the bot")
				return
			}
		}
	}()

	bot.StartTelegramBot(ctx, dataChannel)

}
