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
	"strconv"
	"regexp"
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

func isInt(s string) bool {
	_, err := strconv.Atoi(s);
	return err == nil
}

func isDayFrequency(s string) bool {
	value, _ := regexp.MatchString(`\dt`, s)
	return value;
}

func isDayTime(s string) bool {
	value, _ := regexp.MatchString(`\d\d:\d\d`, s)
	return value;
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

	dayFrequency := 0;

	for update := range updates {
		if update.CallbackQuery != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, dunnoMessage)
			callbackData := update.CallbackQuery.Data
			switch callbackData {
				case "Feed Me!":
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, feedMeWelcomeMessage)
					msg.ReplyMarkup = dayFreqReplies
				case "Feed Someone Else!":
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, feedPetWelcomeMessage)
			}
			if isDayFrequency(callbackData) {
				dayFrequency, _ = strconv.Atoi(callbackData[:1]);
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Shall you eat " + 
					callbackData[:1] + " times per day!\n" + myFirstMealMessage)
				msg.ReplyMarkup = firstMealReplies
			}
			if isDayTime(callbackData) {
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, strconv.Itoa(dayFrequency) + " meals\n" + 
				mySetOtherMealsMessage)
			}
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
			bot.Send(msg)
		}
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if (isText(update)) {
			patience = 2
			switch update.Message.Text {
				case "/start", "/restart":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
					msg.ReplyMarkup = startReplies
					bot.Send(msg)
				default: 
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, dunnoMessage))
			}
		} else {
			patience -= 1
		}
		if patience == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage1)
			bot.Send(msg)
		}
		if patience < 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage2)
			bot.Send(msg)
			patience = 1
		}

	}
}