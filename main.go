package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
)

func main() {
	var lastMessage *tgbotapi.Message
	bot, err := tgbotapi.NewBotAPIWithClient("954724330:AAH7XJVLIOUveij2XTNr6IJgnAuvviZg49c", &http.Client{})
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 6000

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s][%d] %s", update.Message.From.UserName, update.Message.Date, update.Message.Text)

		if lastMessage != nil && update.Message.Date - lastMessage.Date < 3 && lastMessage.From.UserName == update.Message.From.UserName {
			messageText := "@" + update.Message.From.UserName + ", сука, задолбал спамить!"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
		lastMessage = update.Message

	}
}