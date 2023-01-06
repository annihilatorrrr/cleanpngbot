package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
)

const startMsg = `
I'm alive, send me a word or try me inline by just writing my username in text box or send /search command followed by the query to search in cleanpng.com!
Send /download cleanpng_link to send that PNG as photo in telegram or send just send the link to download.

By @Memers_Gallery!
Source code: https://github.com/annihilatorrrr/cleanpngbot`

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ""
	_, _ = ctx.EffectiveMessage.Reply(b, startMsg,
		&gotgbot.SendMessageOpts{
			DisableWebPagePreview: true,
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:              "Try me inline!",
					SwitchInlineQuery: &query,
				}},
			}},
		})
	return ext.EndGroups
}

func procequery(rquery, page string) string {
	query := strings.ToLower(rquery)
	if strings.Contains(query, " ") {
		query = strings.Join(strings.Split(query, " "), "-")
	}
	srchstr := "https://www.cleanpng.com/free/%s.html"
	txt := fmt.Sprintf("<b>Here's the search results for %s with thier resolutions and disk sizes:</b>", query)
	if page != "0" {
		srchstr = "https://www.cleanpng.com/free/%s" + fmt.Sprintf(",%s", page) + ".html"
		txt += fmt.Sprintf("\n<b>Page: %s</b>\n\n", page)
	} else {
		txt += "\n\n"
	}
	raw, err := soup.Get(fmt.Sprintf(srchstr, query))
	if err != nil {
		return "<b>No data Found!<b>\n" + err.Error()
	}
	datas := soup.HTMLParse(raw).FindAll("article")
	aa := false
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
		txt = "<b>No data Found!\n@Memers_Gallery</b>"
	}
	return txt
}

func callbackhand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	splited := strings.Split(query.Data, "=")
	data := splited[1]
	page := splited[2]
	txt := procequery(data, page)
	if strings.Contains(txt, "<b>No data Found!") {
		_, _, _ = query.Message.EditText(b, txt, &gotgbot.EditMessageTextOpts{ParseMode: "html"})
		return ext.EndGroups
	}
	intpage, _ := strconv.Atoi(page)
	backint := intpage - 1
	if backint < 1 {
		backint = 0
	}
	var err error
	if backint == 0 {
		_, _, err = query.Message.EditText(b, txt, &gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
			ParseMode:             "html",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:         "Next Page >",
					CallbackData: fmt.Sprintf("call=%s=%d", data, intpage+1),
				}},
			}},
		})
	} else {
		_, _, err = query.Message.EditText(b, txt, &gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
			ParseMode:             "html",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text:         "< Previous Page",
						CallbackData: fmt.Sprintf("call=%s=%d", data, backint),
					},
					{
						Text:         "Next Page >",
						CallbackData: fmt.Sprintf("call=%s=%d", data, intpage+1),
					},
				},
			}},
		})
	}
	if err != nil {
		_, _, _ = query.Message.EditText(b, err.Error(), nil)
	}
	_, _ = query.Answer(b, nil)
	return ext.EndGroups
}

func sendres(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.Text == "" || msg.ViaBot != nil {
		return ext.EndGroups
	}
	todwnld := strings.Contains(msg.Text, "https://www.cleanpng.com/png-")
	if len(msg.Text) > 50 && !todwnld {
		_, _ = msg.Reply(b, "Query is too big to search!", nil)
		return ext.EndGroups
	}
	em, err := msg.Reply(b, "Processing ...", nil)
	if err != nil {
		return ext.EndGroups
	}
	if todwnld {
		link := downloader(msg.Text)
		if link == "" {
			_, _, _ = em.EditText(b, "Download link not found!", nil)
			return ext.EndGroups
		}
		_, _ = em.Delete(b, nil)
		_, _ = b.SendPhoto(msg.Chat.Id, link, &gotgbot.SendPhotoOpts{ReplyToMessageId: msg.MessageId})
		return ext.EndGroups
	}
	txt := procequery(msg.Text, "0")
	if strings.Contains(txt, "<b>No data Found!") {
		_, _, _ = em.EditText(b, txt, &gotgbot.EditMessageTextOpts{ParseMode: "html"})
		return ext.EndGroups
	}
	_, _, err = em.EditText(b, txt, &gotgbot.EditMessageTextOpts{
		DisableWebPagePreview: true,
		ParseMode:             "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{
				Text:         "Next Page >",
				CallbackData: fmt.Sprintf("call=%s=2", msg.Text),
			}},
		}},
	})
	if err != nil {
		_, _, _ = em.EditText(b, err.Error(), nil)
	}
	return ext.EndGroups
}

func sendinline(b *gotgbot.Bot, ctx *ext.Context) error {
	q := ctx.InlineQuery
	if q.Query == "" {
		_, _ = q.Answer(b, []gotgbot.InlineQueryResult{
			gotgbot.InlineQueryResultArticle{
				Id:                  uuid.NewString(),
				Title:               "Error:",
				Description:         "Write some query!",
				InputMessageContent: gotgbot.InputTextMessageContent{MessageText: "Provide some query!"},
			},
		}, nil)
		return ext.EndGroups
	}
	if len(q.Query) > 50 {
		_, _ = q.Answer(b, []gotgbot.InlineQueryResult{
			gotgbot.InlineQueryResultArticle{
				Id:                  uuid.NewString(),
				Title:               "Error:",
				Description:         "Query too big!",
				InputMessageContent: gotgbot.InputTextMessageContent{MessageText: "Query is too big to search!"},
			},
		}, nil)
		return ext.EndGroups
	}
	txt := procequery(q.Query, "0")
	_, _ = q.Answer(b, []gotgbot.InlineQueryResult{
		gotgbot.InlineQueryResultArticle{
			Id:          uuid.NewString(),
			Title:       "Results!",
			Description: "Found something!",
			InputMessageContent: gotgbot.InputTextMessageContent{
				MessageText:           txt,
				ParseMode:             "html",
				DisableWebPagePreview: true,
			},
		},
	}, nil)
	return ext.EndGroups
}

func search(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	args := ctx.Args()[1:]
	if len(args) == 0 {
		_, _ = msg.Reply(b, "I need a query to search!", nil)
		return ext.EndGroups
	}
	query := args[0]
	if len(query) > 50 {
		_, _ = msg.Reply(b, "Query is too big to search!", nil)
		return ext.EndGroups
	}
	em, err := msg.Reply(b, "Finding ...", nil)
	if err != nil {
		return ext.EndGroups
	}
	txt := procequery(query, "0")
	if strings.Contains(txt, "<b>No data Found!") {
		_, _, _ = em.EditText(b, txt, &gotgbot.EditMessageTextOpts{ParseMode: "html"})
		return ext.EndGroups
	}
	_, _, err = em.EditText(b, txt, &gotgbot.EditMessageTextOpts{
		DisableWebPagePreview: true,
		ParseMode:             "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{
				Text:         "Next Page >",
				CallbackData: fmt.Sprintf("call=%s=2", query),
			}},
		}},
	})
	if err != nil {
		_, _, _ = em.EditText(b, err.Error(), nil)
	}
	return ext.EndGroups
}

func downloader(url string) string {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return fmt.Sprintf("%sdownload-png.html", url)
}

func download(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	args := ctx.Args()[1:]
	if len(args) == 0 {
		_, _ = msg.Reply(b, "I need a link!", nil)
		return ext.EndGroups
	}
	if !strings.Contains(args[0], "https://www.cleanpng.com/png-") {
		_, _ = msg.Reply(b, "Not a valid link!", nil)
		return ext.EndGroups
	}
	nm, err := msg.Reply(b, "Processing ...", nil)
	if err != nil {
		return ext.EndGroups
	}
	link := downloader(args[0])
	if link == "" {
		_, _, _ = nm.EditText(b, "Download link not found!", nil)
		return ext.EndGroups
	}
	_, _ = nm.Delete(b, nil)
	_, _ = b.SendPhoto(msg.Chat.Id, link, &gotgbot.SendPhotoOpts{
		ReplyToMessageId: msg.MessageId,
		Caption:          "<b>By @CleanPNGRoBot from @Memers_Gallery</b>",
		ParseMode:        "html",
		MessageThreadId:  msg.MessageThreadId,
	})
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
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			UnhandledErrFunc: func(err error) {
				log.Printf("An error occurred while handling update:\n%s", err.Error())
			},
			MaxRoutines: -1,
		}),
	})

	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("search", search))
	dispatcher.AddHandler(handlers.NewCommand("download", download))
	dispatcher.AddHandler(handlers.NewInlineQuery(inlinequery.All, sendinline))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("call="), callbackhand))
	dispatcher.AddHandler(handlers.NewMessage(message.ChatType("private"), sendres))

	if _, err = b.SetWebhook(url, &gotgbot.SetWebhookOpts{
		DropPendingUpdates: true,
		AllowedUpdates:     []string{"message", "inline_query", "chosen_inline_result", "callback_query"},
		SecretToken:        "xyzzz",
	}); err != nil {
		panic(err.Error())
	}
	if err = updater.StartWebhook(b,
		token,
		ext.WebhookOpts{
			Port:        port,
			SecretToken: "xyzzz",
		}); err != nil {
		panic(err.Error())
	}
	log.Printf("%s has been started!\n", b.User.Username)
	updater.Idle()
}
