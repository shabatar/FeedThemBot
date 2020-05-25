package main

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"os"
)

var BotToken = os.Getenv("TOKEN")
var URL = os.Getenv("URL")
var USER = os.Getenv("USER")
var PASS = os.Getenv("PASS")
var SEND = os.Getenv("SEND")

func authorizeBot() (*tgbotapi.BotAPI, error) {
	if len(BotToken) == 0 {
		log.Printf("TOKEN environment variable is missing: you can request one by creating a bot in BotFather")
		return nil, errors.New("error on authorization phase, perhaps a TOKEN is missing")
	}
	var bot *tgbotapi.BotAPI
	if len(URL) != 0 {
		dialer, err := proxy.SOCKS5(
			"tcp",
			URL,
			&proxy.Auth{
				User:     USER,
				Password: PASS,
			},
			proxy.Direct)
		if err != nil {
			return nil, errors.New("an error occurred while connecting to server")
		}

		tr := &http.Transport{Dial: dialer.Dial}

		myClient := &http.Client{
			Transport: tr,
		}

		bot, err = tgbotapi.NewBotAPIWithClient(BotToken, "https://api.telegram.org/bot%s/%s", myClient)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		bot, err = tgbotapi.NewBotAPI(BotToken)
		if err != nil {
			return nil, err
		}
	}
	return bot, nil
}
