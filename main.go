package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"os"
)

func main() {
	var lastMessage *tgbotapi.Message
	bot, err := tgbotapi.NewBotAPI("954724330:AAH7XJVLIOUveij2XTNr6IJgnAuvviZg49c")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	log.Printf("Http Listener switched on port %s", os.Getenv("PORT"))
	
	updates := fetchUpdates(bot)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s][%d] %s", update.Message.From.UserName, update.Message.Date, update.Message.Text)

		if lastMessage != nil && update.Message.Date - lastMessage.Date < 3 && lastMessage.From.UserName == update.Message.From.UserName {
			messageText := "@" + getUserName(update.Message.From) + ", сука, задолбал спамить!"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
		lastMessage = update.Message

	}

}

func fetchUpdates(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	_, err := bot.SetWebhook(tgbotapi.NewWebhook("https://anti-peedrila.herokuapp.com/" + bot.Token))
	if err != nil {
		log.Panic(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	return updates
}

func getUserName(user *tgbotapi.User) string {
	var result string

	if len(user.UserName) > 0 {
		result += user.UserName
	} else {
		result += user.FirstName + " " + user.LastName
	}
	return result
}