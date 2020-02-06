package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"os"
)

var stats map[tgbotapi.User]UserStat

func main() {
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
		processMessage(update.Message, bot)
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

func processMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	user := message.From
	if _, ok := stats[*user]; ok {} else {
		stats[*user] = UserStat{
			User:             *user,
			RelationshipRate: 1,
			Penalties:		  0,
		}
	}

	stat := stats[*user]

	if len(message.Text) > 0 {
		if len(stat.LastMessages) == 4 && message.Date - stat.LastMessages[0].Date < 30 {
			stat.registerPenalty()
			msg := tgbotapi.NewMessage(message.Chat.ID, stat.prepareMessage())
			bot.Send(msg)
		}
		stack := addMessageToStack(stat.LastMessages, *message)
		stat.LastMessages = stack
	}

}

func addMessageToStack(stack []tgbotapi.Message, message tgbotapi.Message) []tgbotapi.Message{
	if len(stack) > 4 {
		stack = stack[1:]
	}
	stack = append(stack, message)
	return stack
}