package bot

import (
	"context"
	"main/core"
	"main/logger"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var log = logger.Logger

type MessageData struct {
	Message string
	Asset   string
	Date    string
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
			messageData.Date = time.Unix(int64(update.Message.Date), 0).Format(time.RFC3339)
		}

		if len(update.Message.Photo) > 0 {
			if len(update.Message.Caption) != 0 {
				messageData.Message = update.Message.Caption
			}

			// Get the file path
			// Process all photos
			for _, photo := range update.Message.Photo {
				fileID := photo.FileID
				// Save the file ID
				messageData.Asset = fileID
			}
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
