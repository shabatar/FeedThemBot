package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"golang.org/x/net/proxy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"encoding/json"
)

var startReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Me!","Feed Me!"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Someone Else!","Feed Someone Else!"),
	),
)

const (
	EmojiCat = "\xF0\x9F\x90\x88"
	EmojiHard = "\xF0\x9F\x98\xAB"
	EmojiPity = "\xF0\x9F\x98\x96"
	EmojiYumm = "\xF0\x9F\x98\x8B"
	EmojiWave = "\xF0\x9F\x91\x8B"
)

type Proxy struct {
	ProxyURL  string `json:"URL"`
	ProxyUser string `json:"user"`
	ProxyPass string `json:"pass"`
} 

type Settings struct {
	BotToken string `json:"token"`
	Proxy Proxy `json:"proxy"` 
}


func isText(update tgbotapi.Update) bool {
	return reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != ""
}

func main() {
	configFile, err := os.Open("botsettings.json")
	var settings Settings
	if err != nil {
		fmt.Fprintln(os.Stderr, "an error occured while opening config file: ", err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
		fmt.Fprintln(os.Stderr, "an error occured while parsing config file: ", err)
	}

	var Proxy Proxy = settings.Proxy
	proxyAuth := proxy.Auth{
		User:     Proxy.ProxyUser,
		Password: Proxy.ProxyPass,
	}
	dialer, err := proxy.SOCKS5(
				"tcp", 
				Proxy.ProxyURL,
				&proxyAuth, 
				proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "an error occured while connecting to proxy", err)
	}

	tr := &http.Transport{Dial: dialer.Dial}

	myClient := &http.Client{
		Transport: tr,
	}
	
	bot, err := tgbotapi.NewBotAPIWithClient(settings.BotToken, "https://api.telegram.org/bot%s/%s", myClient)
	
	if err != nil {
		log.Panic(err)
	}
	
	bot.Debug = true

	log.Printf("Successfully authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	patience := 2;

	for update := range updates {
		if update.CallbackQuery != nil {
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			switch update.CallbackQuery.Data {
				case "Feed Me!":
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
						"I cannot feed you right now. Sorry" + EmojiPity))
				case "Feed Someone Else!":
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
						"I cannot feed them right now. Sorry" + EmojiPity))
				default:
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, 
						"I don't know what to do. Sorry" + EmojiPity))
			}
		}
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if (isText(update)) {
			patience = 2
			switch update.Message.Text {
				case "/start", "/restart":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, 
						"Hi there!" + EmojiWave + " My name is FeedThemBot and I am here to help you eat regularly.\nI know, it's difficult! " + EmojiHard + 
						" \nI could remind you when you should probably grab a bite. " + EmojiYumm + 
						" \nJust press a button and set up a meal notification. You can also use me to remind you about feeding your pets " + EmojiCat)
					msg.ReplyMarkup = startReplies
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