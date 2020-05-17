package main

import (
	"log"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"golang.org/x/net/proxy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"regexp"
)

var BotToken = os.Getenv("TOKEN")
var URL = os.Getenv("URL")
var USER = os.Getenv("USER")
var PASS = os.Getenv("PASS")

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

func isSetMeal(s string) bool {
	value, _ := regexp.MatchString(`\d[nrt][dh] meal`, s)
	return value;
}

func feedThemBot() {
	if len(BotToken) == 0 {
		log.Printf("TOKEN environment variable is missing: you can request one by creating a bot in BotFather")
    	return
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
			fmt.Fprintln(os.Stderr, "an error occured while connecting to server", err)
		}

		tr := &http.Transport{Dial: dialer.Dial}

		myClient := &http.Client{
			Transport: tr,
		}
		
		bot, err = tgbotapi.NewBotAPIWithClient(BotToken, "https://api.telegram.org/bot%s/%s", myClient)
		if err != nil {
			log.Panic(err)
		}
	} else {
		var err error
		bot, err = tgbotapi.NewBotAPI(BotToken)
		if err != nil {
			log.Panic(err)
		}
	}
	
	bot.Debug = true

	log.Printf("Successfully authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

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
				if (dayFrequency > 1) { 
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n" + 
					mySetOtherMealsMessage)
					msg.ReplyMarkup = printEditMealsReplies(dayFrequency)
				} else {
					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n") 
				}
			}
			if isSetMeal(callbackData) {
			    msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "When you’d like to have " + callbackData + "?")
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

func initDB() {
	if err := createUsageTable(); err != nil {
		panic(err)
	}
	if err := createDailyUsersTable(); err != nil {
		panic(err)
	}
}

func main() {
	initDB()
	feedThemBot()
}