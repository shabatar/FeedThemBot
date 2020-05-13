package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"golang.org/x/net/proxy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("testlink.com","http://testlink.com"),
		tgbotapi.NewInlineKeyboardButtonData("Feed Me","Feed Me"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Someone","Feed Someone"),
		tgbotapi.NewInlineKeyboardButtonData("Help","Help"),
	),
)

const (
	BotToken               = ""
	ProxyURL               = ""
	ProxyUser              = ""
	ProxyPass              = ""
)

func isText(update tgbotapi.Update) bool {
	return reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != ""
}

func main() {
	proxyAuth := proxy.Auth{
	    User:     ProxyUser,
	    Password: ProxyPass,
	}
	dialer, err := proxy.SOCKS5(
		        "tcp", 
		        ProxyURL,
		        &proxyAuth, 
	            proxy.Direct)
	if err != nil {
	    fmt.Fprintln(os.Stderr, "an error occured while connecting to proxy", err)
	}

	tr := &http.Transport{Dial: dialer.Dial}
	myClient := &http.Client{
	    Transport: tr,
	}
	bot, err := tgbotapi.NewBotAPIWithClient(BotToken, "https://api.telegram.org/bot%s/%s", myClient)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

    patience := 2;

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		if (isText(update)) {
			switch update.Message.Text {
			    case "/start", "/restart":
			    	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "RRR!")
			    	// msg.ReplyMarkup = numericKeyboard
			    	bot.Send(msg)
			    default: 
			        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You " + update.Message.Text + ", " + update.Message.From.FirstName + "!")
			    	bot.Send(msg) 
			}
		} else {
			patience -= 1
		}
	    if patience == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hmmm... I'm not sure if you're using me the right way")
		    bot.Send(msg)
		}
	    if patience < 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Could you please stop it?")
		    bot.Send(msg)
		    patience = 1
		}

	}
}