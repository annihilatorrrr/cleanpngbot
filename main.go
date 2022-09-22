package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/anaskhan96/soup"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, _ = ctx.EffectiveMessage.Reply(b, "I'm alive, send me a word to search in cleanpng site!\nBy @Memers_Gallery!", nil)
	return ext.EndGroups
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
	query := strings.ToLower(msg.Text)
	raw, err := soup.Get(fmt.Sprintf("https://www.cleanpng.com/free/%s.html", query))
	if err != nil {
		_, _, _ = em.EditText(b, err.Error(), nil)
		return ext.EndGroups
	}
	data := soup.HTMLParse(raw)
	log.Println(data.Attrs())
	_, _, _ = em.EditText(b, "Check logs! (Devs!)", nil)
	return ext.EndGroups
}

func main() {
	token := "5793391546:AAFaACVFk0lpKn8jWLTn78Hpl2pnNs8sEvs" // os.Getenv("TOKEN")
	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic(err.Error())
	}
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("An error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: 20,
		},
	})
	_, err = b.SetWebhook(os.Getenv("URL"), &gotgbot.SetWebhookOpts{
		DropPendingUpdates: true,
		AllowedUpdates:     []string{"message"},
		MaxConnections:     20,
		SecretToken:        "xyzzz",
	})
	if err != nil {
		panic(err.Error())
	}
	port, _ := strconv.Atoi(os.Getenv("PORT"))
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
	log.Printf("%s has been started!\n", b.User.Username)
	runtime.GC()
	updater.Idle()
	_ = updater.Stop()
}
