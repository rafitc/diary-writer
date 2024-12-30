package bot

import (
	"context"
	"main/core"
	"main/logger"
	"path"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var log = logger.Logger

type MessageData struct {
	Message             string
	Asset               string
	Asset_extension     string
	Asset_download_link string
	Date                time.Time
}

func StartTelegramBot(ctx context.Context, ch chan MessageData) {
	handlerWithChannel := createHandler(ch)

	opts := []bot.Option{
		bot.WithDefaultHandler(handlerWithChannel),
	}

	b, err := bot.New(core.Config.BOT.TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func createHandler(ch chan MessageData) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// Initialize MessageData
		messageData := MessageData{}
		if update.Message != nil {

			messageData.Message = update.Message.Text
			messageData.Date = time.Unix(int64(update.Message.Date), 0)
		}

		if len(update.Message.Photo) > 0 {
			if len(update.Message.Caption) != 0 {
				messageData.Message = update.Message.Caption
			}

			// Get the last photo (highest resolution)
			lastPhoto := update.Message.Photo[len(update.Message.Photo)-1] //TODO need to handle multiple photos case
			// Retrieve the file info
			file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: lastPhoto.FileID})
			if err != nil {
				log.Errorf("Failed to get file info: %v", err)
				return
			}

			// Extract the file extension from the file path
			filePath := file.FilePath
			fileName := path.Base(filePath)                      // Get the base name of the file
			fileExtension := strings.ToLower(path.Ext(fileName)) // Extract the file extension
			// store file-name with extension

			// get download link
			downloadLink := b.FileDownloadLink(file)

			log.Infof("File Name: %s, File Extension: %s", fileName, fileExtension)

			messageData.Asset = lastPhoto.FileID
			messageData.Asset_extension = fileExtension
			messageData.Asset_download_link = downloadLink
		}

		log.Debugf("Passing content to the channel: %v", messageData)
		// Send message data to the channel
		ch <- messageData

		// Respond to the user
		responseText := "Thank you for your message!"
		if messageData.Asset != "" {
			responseText += " Photo path successfully."
		}

		log.Debugf("Sending response: %s", responseText)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   responseText,
		})
	}
}
