package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/inlinequery"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/anaskhan96/soup"
	"github.com/google/uuid"
)

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	search := ""
	_, _ = ctx.EffectiveMessage.Reply(b, "I'm alive, send me a word or try me inline by just writing my username in text box to search in cleanpng.com!\nBy @Memers_Gallery!\nSource code: https://github.com/annihilatorrrr/cleanpngbot",
		&gotgbot.SendMessageOpts{
			DisableWebPagePreview: true,
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:              "Try me inline!",
					SwitchInlineQuery: &search,
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
	if page != "" || page != "0" {
		srchstr = "https://www.cleanpng.com/free/%s" + fmt.Sprintf(",%s", page) + ".html"
	}
	raw, err := soup.Get(fmt.Sprintf(srchstr, query))
	if err != nil {
		return err.Error()
	}
	datas := soup.HTMLParse(raw).FindAll("article")
	aa := false
	txt := fmt.Sprintf("<b>Here's the search results for %s with thier resolutions and disk sizes:</b>\n\n", query)
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

func callbackhand(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	splited := strings.Split(query.Data, "=")
	data := splited[1]
	page := splited[2]
	txt := procequery(data, page)
	if strings.Contains(txt, "No data Found!") || strings.Contains(txt, "error") {
		_, _, _ = query.Message.EditText(b, txt, nil)
		return ext.EndGroups
	}
	intpage, _ := strconv.Atoi(page)
	backint := intpage - 1
	if backint < 0 {
		backint = 0
	}
	var err error
	if backint == 0 || backint < 0 {
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
	return ext.EndGroups
}

func sendres(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.Text == "" || msg.ViaBot != nil {
		return ext.EndGroups
	}
	if len(msg.Text) > 50 {
		_, _ = msg.Reply(b, "Query is too big to search!", nil)
		return ext.EndGroups
	}
	em, err := msg.Reply(b, "Finding ...", nil)
	if err != nil {
		return ext.EndGroups
	}
	txt := procequery(msg.Text, "")
	_, _, err = em.EditText(b, txt, &gotgbot.EditMessageTextOpts{
		DisableWebPagePreview: true,
		ParseMode:             "html",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{{
				Text:         "Next Page >",
				CallbackData: fmt.Sprintf("call=%s=1", msg.Text),
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
	txt := procequery(q.Query, "")
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
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:         "Next Page >",
					CallbackData: fmt.Sprintf("call=%s=1", q.Query),
				}},
			}},
		},
	}, nil)
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

	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewInlineQuery(inlinequery.All, sendinline))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("call="), callbackhand))
	dispatcher.AddHandler(handlers.NewMessage(message.ChatType("private"), sendres))

	_, err = b.SetWebhook(url, &gotgbot.SetWebhookOpts{
		DropPendingUpdates: true,
		AllowedUpdates:     []string{"message", "inline_query", "chosen_inline_result", "callback_query"},
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
	log.Printf("%s has been started!\n", b.User.Username)
	updater.Idle()
}
