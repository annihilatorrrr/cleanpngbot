package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/anaskhan96/soup"
)

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = ctx.EffectiveMessage.Reply(b, "I'm alive, send me a word to search in cleanpng.com!\nBy @Memers_Gallery!\nSource code: https://github.com/annihilatorrrr/cleanpngbot",
		&gotgbot.SendMessageOpts{DisableWebPagePreview: true})
	return ext.EndGroups
}

func procequery(rquery string) string {
	query := strings.ToLower(rquery)
	if strings.Contains(query, " ") {
		query = strings.Join(strings.Split(query, " "), "-")
	}
	raw, err := soup.Get(fmt.Sprintf("https://www.cleanpng.com/free/%s.html", query))
	if err != nil {
		return err.Error()
	}
	datas := soup.HTMLParse(raw).FindAll("article")
	aa := false
	txt := fmt.Sprintf("<b>Here's the search results for %s with thier resolution and disk size:</b>\n\n", query)
	for _, rdata := range datas {
		pd := rdata.FindAll("p")
		if len(pd) < 3 {
			continue
		}
		aa = true
		txt += fmt.Sprintf(`<b>> <a href="https://www.cleanpng.com%s">%s</a> - %s - %s</b>`+"\n",
			rdata.Find("a").Attrs()["href"],
			pd[0].Find("a").Text(),
			pd[1].Find("span").Text(),
			pd[2].Find("span").Text(),
		)
	}
	txt += "\n<b>@Memers_Gallery</b>"
	if !aa {
		txt = "No data Found!\n<b>@Memers_Gallery</b>"
	}
	return txt
}

func sendres(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.Text == "" {
		return ext.EndGroups
	}
	if len(msg.Text) > 50 {
		_, _ = msg.Reply(b, "Query too big!", nil)
		return ext.EndGroups
	}
	em, err := msg.Reply(b, "Finding ...", nil)
	if err != nil {
		return ext.EndGroups
	}
	txt := procequery(msg.Text)
	_, _, err = em.EditText(b, txt, &gotgbot.EditMessageTextOpts{DisableWebPagePreview: true, ParseMode: "html"})
	if err != nil {
		_, _, _ = em.EditText(b, err.Error(), nil)
	}
	return ext.EndGroups
}

func sendinline(b *gotgbot.Bot, ctx *ext.Context) error {
	q := ctx.InlineQuery
	if q.Query == "" {
		_, err := q.Answer(b, []gotgbot.InlineQueryResult{
			gotgbot.InlineQueryResultArticle{Title: "Error:", Description: "Write some query!", InputMessageContent: gotgbot.InputTextMessageContent{MessageText: "Provide some query!"}},
		}, nil)
		if err != nil {
			log.Println(err.Error())
		}
		return ext.EndGroups
	}
	txt := procequery(q.Query)
	_, err := q.Answer(b, []gotgbot.InlineQueryResult{
		gotgbot.InlineQueryResultArticle{Title: "Results!", InputMessageContent: gotgbot.InputTextMessageContent{MessageText: txt, ParseMode: "html", DisableWebPagePreview: true}},
	}, nil)
	if err != nil {
		log.Println(err.Error())
	}
	return ext.EndGroups
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("No bot token was found!")
	}
	url := os.Getenv("URL")
	if url == "" {
		panic("No webhook url was found!")
	}
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		panic("No port was found to bind!")
	}
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic(err.Error())
	}
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Printf("An error occurred while handling update:\n%s", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: 20,
		},
	})
	_, err = b.SetWebhook(url, &gotgbot.SetWebhookOpts{
		DropPendingUpdates: true,
		AllowedUpdates:     []string{"message"},
		MaxConnections:     20,
		SecretToken:        "xyzzz",
	})
	if err != nil {
		panic(err.Error())
	}
	err = updater.StartWebhook(b,
		ext.WebhookOpts{
			Port: port,
		})
	if err != nil {
		panic(err.Error())
	}
	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewMessage(message.ChatType("private"), sendres))
	dispatcher.AddHandler(handlers.NewInlineQuery(nil, sendinline))
	log.Printf("%s has been started!\n", b.User.Username)
	runtime.GC()
	updater.Idle()
	_ = updater.Stop()
}
