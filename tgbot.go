package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const botToken = "607984507:AAF1rT7hxQbMxMNnNys9ReFXFz-JTIL_JQQ" // wzulfikar_bot
const chatId = -1001489244273                                    // Luxtag::Knowledge

var bot *tgbotapi.BotAPI

func init() {
	_bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot = _bot

	bot.Debug = false
	log.Println("Bot initialized:", bot.Self.UserName)
}

func SendMsg(text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "markdown"

	bot.Send(msg)
}
