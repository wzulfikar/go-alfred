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
	retryCount := 0
	parseMode := "Markdown"
	for {
		if retryCount > 1 {
			break
		}

		inlineConf := createInlineQueryResult(inlineQueryID, results, parseMode)
		if _, err := bot.AnswerInlineQuery(*inlineConf); err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), "Bad Request: can't parse entities") {
				log.Printf("failed to parse result in markdown for query %s. retrying without markdown..\n", inlineQueryID)
				parseMode = ""
				retryCount++
				continue
			}
			bot.AnswerInlineQuery(inlineQueryErrorConf(inlineQueryID))
		}

		if retryCount > 0 {
			log.Println("retry completed for query", inlineQueryID)
		}
		break
	}
}

func inlineQueryErrorConf(inlineQueryID string) tgbotapi.InlineConfig {
	article := tgbotapi.InlineQueryResultArticle{
		Type:        "article",
		ID:          inlineQueryID,
		Title:       "Internal error",
		Description: "Whoops! Something went wrong. Please try again :)",
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text:                  "`luxtagbot internal error: please try again.`",
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

func createInlineQueryResult(inlineQueryID string, results *[]contracts.Result, parseMode string) *tgbotapi.InlineConfig {
	articles := make([]interface{}, len(*results))
	for i, result := range *results {
		text := result.Text
		if text == "" {
			text = fmt.Sprintf("*%s*\n%s\n\n––\nOpen in browser:\n%s",
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
				ParseMode:             parseMode,
				DisableWebPagePreview: true,
			},
		}
	}

	inlineConf := &tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		IsPersonal:    true,
		CacheTime:     15,
		Results:       articles,
	}

	return inlineConf
}
