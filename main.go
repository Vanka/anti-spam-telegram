package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"os"
)

var stats map[tgbotapi.User]UserStat

func main() {
	appPort := os.Getenv("PORT")
	log.Printf("App Port %s", appPort)
	botToken := os.Getenv("BOTTOKEN")
	log.Printf("Bot Token %s", botToken)
	mongoURI := os.Getenv("MONGODB_URI")
	log.Printf("Mongo Uri %s", mongoURI)

	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Panic("Error while creating a client: %s", err)
	//}

	stats = make(map[tgbotapi.User]UserStat)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Error while making: %s", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go http.ListenAndServe(":" + appPort, nil)
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
	_, err := bot.SetWebhook(tgbotapi.NewWebhook("https://anti-peedrila-765923146deb.herokuapp.com/" + bot.Token))
	if err != nil {
		log.Panic(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	return updates
}

func processMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	user := *message.From
	if _, ok := stats[user]; ok {} else {
		stats[user] = UserStat{
			User:             user,
			RelationshipRate: 1,
			Penalties:		  0,
			LastMessages:     []tgbotapi.Message{},
		}
	}

	stat := stats[user]
	log.Printf("Stats for %s before processing:", getUserName(stat.User))
	log.Print(stat)
	log.Print(message.Voice)

	if message.Voice != nil {
		
	}

	if message.ForwardFrom == nil {
		if len(stat.LastMessages) > 4 && message.Date - stat.LastMessages[0].Date < 30 || (stat.lastMessage() != nil && message.Date - stat.lastMessage().Date <= 5) {
			stat.Penalties += 1
			stat.RelationshipRate *= 0.5
			msg := tgbotapi.NewMessage(message.Chat.ID, stat.prepareMessage())
			// msg := tgbotapi.NewMessage(message.Chat.ID, "Вова голубоглазый пидр")
			bot.Send(msg)
		} else {
			stat.RelationshipRate *= 1.05
		}
		stack := addMessageToStack(stat.LastMessages, *message)
		stat.LastMessages = stack
	}
	log.Printf("Stats for %s after processing:", getUserName(stat.User))
	log.Printf("LastMessage stack length: %d. First message: [%d] %s. Last message: [%d] %s", len(stat.LastMessages), stat.LastMessages[0].Date, stat.LastMessages[0].Text, message.Date, message.Text)
	stats[user] = stat
}

func addMessageToStack(stack []tgbotapi.Message, message tgbotapi.Message) []tgbotapi.Message{
	if len(stack) > 4 {
		stack = stack[1:]
	}
	stack = append(stack, message)
	return stack
}