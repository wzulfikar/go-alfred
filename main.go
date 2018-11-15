package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// const token = "607984507:AAF1rT7hxQbMxMNnNys9ReFXFz-JTIL_JQQ" // luxtag_bot
const token = "286275707:AAFtt5GmfNK6fHhaBUD-8wSU3m0q29rKy9A" // wzulfikar_bot

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	for {
		log.Println("fetching update..")
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("updates found:", len(updates))
		for update := range updates {
			if update.InlineQuery.Query == "" {
				continue
			}

			queryText := update.InlineQuery.Query
			ytQuery := fmt.Sprintf("query=%s", url.QueryEscape(queryText))
			issues, err := FetchIssue(ytQuery)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(issues)

			articles := []interface{}{}
			for i := 0; i < len(*issues); i++ {
				issue := (*issues)[i]

				article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, issue.Summary, shorten(issue.Description))
				article.Description = fmt.Sprintf("*%s*\n%s\n[→ Read more](%s)",
					issue.Summary,
					shorten(issue.Description),
					GetLink(issue.ID))
				articles = append(articles, article)
			}

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       []interface{}{articles},
			}

			if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
				log.Println(err)
			}
		}

		log.Println("waiting for next update")
		time.Sleep(time.Second * 5)
	}
}
