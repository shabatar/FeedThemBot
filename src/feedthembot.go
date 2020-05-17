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

func isTimezone(s string) bool {
	value, _ := regexp.MatchString(`GMT[+-]\d\d?`, s)
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

	for update := range updates {
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, dunnoMessage)
			callbackData := update.CallbackQuery.Data
			userName := update.CallbackQuery.From.UserName
			log.Printf("[%s] %s", userName, update.CallbackQuery.Data)
			mealsSet, _ := userMealsUTCSet(userName)
			dayFrequency, _ := getUserSelectedFrequency(userName)
			if mealsSet {
				msg = tgbotapi.NewMessage(chatID, "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?")
				msg.ReplyMarkup = agreeDisagreeReplies
				bot.Send(msg)
				continue
			}
			switch callbackData {
				case "Feed Me!":
					userTimezone, _ := getUserTimezone(userName)
					if userTimezone == -100 {
						msg = tgbotapi.NewMessage(chatID, timezoneMessage)
						msg.ReplyMarkup = timezoneReplies
						bot.Send(msg)
						continue
					}
					msg = tgbotapi.NewMessage(chatID, feedMeWelcomeMessage)
					msg.ReplyMarkup = dayFreqReplies
				case "Feed Someone Else!":
					msg = tgbotapi.NewMessage(chatID, feedPetWelcomeMessage)
				case "Submit":
					msg = tgbotapi.NewMessage(chatID, "You are all set! Wait for the notifications")
				case "Cancel":
					msg = tgbotapi.NewMessage(chatID, "Okay...😢")
					createUserDailyMeals(userName)
			}
			if isDayFrequency(callbackData) {
				msg = tgbotapi.NewMessage(chatID, "Shall you eat " + 
					callbackData[:1] + " times per day!\n" + myFirstMealMessage)
				setUserSelectedFrequency(userName, callbackData[:1])
				createUserDailyMeals(userName)
				msg.ReplyMarkup = firstMealReplies
			}
			if isDayTime(callbackData) {
				// userName := update.CallbackQuery.From.UserName
				dayFrequency, _ := getUserSelectedFrequency(userName)
				// Update first meal
				setUserMealEditIndex(userName, "1")
				updateUserDailyMeal(userName, callbackData)
				if (dayFrequency > 1) { 
					msg = tgbotapi.NewMessage(chatID, "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n" + 
					mySetOtherMealsMessage)
					msg.ReplyMarkup = printEditMealsReplies(dayFrequency)
				} else {
					msg = tgbotapi.NewMessage(chatID, "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n") 
					mealsSet, _ := userMealsUTCSet(userName)
					if mealsSet {
						msg = tgbotapi.NewMessage(chatID, "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?")
						msg.ReplyMarkup = agreeDisagreeReplies
					}
				}
			}
			if isSetMeal(callbackData) {
				userName := update.CallbackQuery.From.UserName
			    msg = tgbotapi.NewMessage(chatID, "When you’d like to have " + callbackData + "? \n(enter HH:MM:SS or HH:MM)")
			    setUserMealEditIndex(userName, callbackData[:1])
			}
			if isTimezone(callbackData) {
			    setUserTimezone(userName, callbackData[len(callbackData)-3:])
			    msg = tgbotapi.NewMessage(chatID, "Timezone has been set")
		    	msg = tgbotapi.NewMessage(chatID, feedMeWelcomeMessage)
				msg.ReplyMarkup = dayFreqReplies
			}
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
			bot.Send(msg)
		}
		if update.Message == nil {
			continue
		}
		userName := update.Message.From.UserName
		log.Printf("[%s] %s", userName, update.Message.Text)
		if (isText(update)) {
			insertUser(userName) 
			setUserPatience(userName, 2)
			chatID := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, dunnoMessage)
			msgText := update.Message.Text
			mealsSet, _ := userMealsUTCSet(userName)
			log.Printf("boooooolean %t", mealsSet)
			dayFrequency, _ := getUserSelectedFrequency(userName)
			if mealsSet {
				msg = tgbotapi.NewMessage(chatID, "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?")
				msg.ReplyMarkup = agreeDisagreeReplies
				bot.Send(msg)
				continue
			}
			if isDayTime(msgText) {
				// Update meal #K
			    log.Printf("setting userMeal to =  : " + msgText)
				updateUserDailyMeal(userName, msgText)
				if (dayFrequency > 1) { 
					msg = tgbotapi.NewMessage(chatID, "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n" + 
					mySetOtherMealsMessage)
					msg.ReplyMarkup = printEditMealsReplies(dayFrequency)
					bot.Send(msg)
				} else {
					msg = tgbotapi.NewMessage(chatID, "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n") 
				}
				continue
			} else {
				switch msgText {
					case "/start", "/restart":
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
						msg.ReplyMarkup = setupReplies
				}
			}
			bot.Send(msg)
		} else {
			insertUser(userName)
			addUserPatience(userName, -1)
		}
		patience, err := getUserPatience(userName)
		if err != nil {
			continue
		}
		if patience == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage1)
			bot.Send(msg)
		}
		if patience < 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage2)
			bot.Send(msg)
			setUserPatience(userName, 1)
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
	if err := createUsersTable(); err != nil {
		panic(err)
	}
}

func main() {
	initDB()
	feedThemBot()
}