package alfred

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/wzulfikar/alfred/contracts"
	"github.com/wzulfikar/alfred/util"
)

func NewTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false
	log.Println("bot initialized:", bot.Self.UserName)
	return bot, nil
}

func SendMsg(bot *tgbotapi.BotAPI, chatId int, text string) {
	msg := tgbotapi.NewMessage(int64(chatId), text)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "markdown"

	bot.Send(msg)
}

func AnswerInlineQuery(bot *tgbotapi.BotAPI, inlineQueryID string, results *[]contracts.Result) {
	articles := make([]interface{}, len(*results))
	for i, result := range *results {
		text := result.Text
		if text == "" {
			text = fmt.Sprintf("*%s*\n%s\n\n––\nLink:\n%s",
				result.Title,
				util.Truncate(util.EscapeMarkdown(result.Description), "...\\[redacted]"),
				result.URL)
		}

		articles[i] = tgbotapi.InlineQueryResultArticle{
			Type:        "article",
			ID:          result.ID,
			Title:       util.EscapeMarkdown(result.Title),
			URL:         result.URL,
			ThumbURL:    result.ThumbURL,
			Description: util.EscapeMarkdown(util.Truncate(result.Description, "")),
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text:                  text,
				ParseMode:             "Markdown",
				DisableWebPagePreview: true,
			},
		}
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		IsPersonal:    true,
		CacheTime:     15,
		Results:       articles,
	}

	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "Bad Request: can't parse entities") {
			log.Printf("results: %v\n", inlineConf.Results)
		}
		bot.AnswerInlineQuery(inlineQueryErrorConf(inlineQueryID))
	}
}

func inlineQueryErrorConf(inlineQueryID string) tgbotapi.InlineConfig {
	article := tgbotapi.InlineQueryResultArticle{
		Type:        "article",
		ID:          inlineQueryID,
		Title:       "Internal error",
		Description: "Whoops! Something went wrong. Please try again :)",
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text:                  "Internal error: please try again.",
			ParseMode:             "Markdown",
			DisableWebPagePreview: true,
		},
	}
	return tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		IsPersonal:    true,
		CacheTime:     100,
		Results:       []interface{}{article},
	}
}
