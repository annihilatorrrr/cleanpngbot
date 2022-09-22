package main

import (
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		token = "111:3333kkkk"
	}
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
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 5,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 5,
			},
		},
	})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("%s has been started!\n", b.User.Username)
	runtime.GC()
	updater.Idle()
	_ = updater.Stop()
}
