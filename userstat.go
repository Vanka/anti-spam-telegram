package main

import tgbotapi "github.com/Syfaro/telegram-bot-api"

type UserStat struct {
	User tgbotapi.User
	RelationshipRate float32
	LastMessages []tgbotapi.Message
	Penalties int
}

func (stat UserStat) registerPenalty() {
	stat.Penalties += 1
	stat.RelationshipRate *= 0.8
}

func (stat UserStat) increaseReputation() {
	stat.RelationshipRate *= 1.2
}

func (stat UserStat) lastMessage() *tgbotapi.Message {
	messageSize := len(stat.LastMessages)
	if messageSize > 0 {
		return &stat.LastMessages[messageSize - 1]
	}
	return nil
}

func (stat UserStat) prepareMessage() string {
	var ask string
	if stat.RelationshipRate > 0.8 {
		ask = "пожалуйста, солнышко, прекрати так себя вести"
	} else if stat.RelationshipRate >= 0.5 {
		ask = "ну хватит себя вести как сволота!"
	} else if stat.RelationshipRate >= 0.2 {
		ask = "ну ты реально ведешь себя как сучка какая-то, прекрати, а!"
	} else {
		ask = "заебал, пиздуй-ка ты нахуй со своим спамом!"
	}

	return getUserName(stat.User) + ", " + ask;
}

func getUserName(user tgbotapi.User) string {
	var result string

	if len(user.UserName) > 0 {
		result += "@" + user.UserName
	} else {
		result += user.FirstName + " " + user.LastName
	}
	return result
}