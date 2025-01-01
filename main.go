/*
Diary Writer Application
Author: Rafi Rasheed TC
Repo: https://github.com/rafitc/diary-writer.repo

This application gets messages from Telegram and writes them into SQLite3.
It also has a scheduled process to upload the data to GitHub and deploy the new app on a daily basis.

Two tasks running on startup:
1. Telegram Bot
2. Cron job to read data and publish it on the web
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

	log.Info("Starting Diary-Writer Application")

	// starting CronJob to Write and Publish the data
	writer.StartCronJob()

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

	// start the telegram bot
	bot.StartTelegramBot(ctx, dataChannel)
}
