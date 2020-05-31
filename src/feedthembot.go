package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"time"
)

func processCallback(userName string, callbackData string) (string, *tgbotapi.InlineKeyboardMarkup) {
	mealsSet, _ := userMealsUTCSet(userName)
	dayFrequency, _ := getUserSelectedFrequency(userName)
	var replyMarkup *tgbotapi.InlineKeyboardMarkup
	msg := dunnoMessage
	usr, _ := getUserData(userName)
	if callbackData == "Feed Me!" {
		if usr.userTimezone == -100 {
			msg = timezoneMessage
			replyMarkup = &timezoneReplies
		} else {
			msg = feedMeWelcomeMessage
			replyMarkup = &dayFreqReplies
		}
	} else if callbackData == "Feed Someone Else!" {
		msg = feedPetWelcomeMessage
	} else if callbackData == "Submit" {
		msg = "You are all set! Wait for the notifications"
		syncTimezone(userName)
		migrateDailyUser(userName)
	} else if callbackData == "Cancel" {
		msg = "Okay...😢"
		createUserDailyMeals(userName)
	} else if isDayFrequency(callbackData) {
		msg = "Shall you eat " + callbackData[:1] + " times per day!\n" + myFirstMealMessage
		setUserSelectedFrequency(userName, callbackData[:1])
		createUserDailyMeals(userName)
		replyMarkup = printMarkedMealReplies(usr.userMealsUTC)
	} else if isSetMeal(callbackData) {
		msg = "When you’d like to have " + callbackData + "? \n"
		setUserMealEditIndex(userName, callbackData[:1])
		replyMarkup = printMarkedMealReplies(usr.userMealsUTC)
	} else if isDayTime(callbackData) {
		dayFrequency := usr.selectedFrequency
		// Update first meal
		mealEditIndex := usr.userMealEditIndex
		if mealEditIndex == -100 {
			setUserMealEditIndex(userName, "1")
		}
		updateUserDailyMeal(userName, callbackData)
		mealsSet, _ = userMealsUTCSet(userName)
		if dayFrequency > 0 && mealsSet {
			msg = "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?"
			replyMarkup = &agreeDisagreeReplies
		} else if dayFrequency > 1 {
			msg = "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n" + mySetOtherMealsMessage
			replyMarkup = printEditMealsReplies(dayFrequency)
		} else {
			msg = "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n"
		}
	} else if isTimezone(callbackData) {
		setUserTimezone(userName, callbackData[len(callbackData)-3:])
		msg = feedMeWelcomeMessage
		replyMarkup = &dayFreqReplies
	} else if callbackData == "/start" {
		_ = clearInsertUser(userName)
		msg = welcomeMessage
		replyMarkup = &setupReplies
	} else if callbackData == "/restart" {
		_ = clearInsertUser(userName)
		msg = welcomeMessage
		replyMarkup = &setupReplies
	} else if dayFrequency > 0 && mealsSet {
		msg = "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?"
		replyMarkup = &agreeDisagreeReplies
	}
	return msg, replyMarkup
}

func processUpdate(bot *tgbotapi.BotAPI, userName string, chatID int64, messageID int, messageText string, callbackID string) {
	answerText, replyMarkup := processCallback(userName, messageText)
	msg := tgbotapi.NewMessage(chatID, answerText)
	if replyMarkup != nil {
		msg.ReplyMarkup = *replyMarkup
	}
	_ = setUserPatience(userName, 2)
	_ = insertUser(userName, chatID)
	answer, _ := bot.Send(msg)
	if callbackID != "" {
		_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, msg.Text))
	}
	insertUsageStats(userName, chatID, messageID, messageText, answer.MessageID, answer.Text)
}

func feedThemBot() {
	bot, err := authorizeBot()
	if err != nil {
		log.Panic(err)
		return
	}
	bot.Debug = true
	log.Printf("Successfully authorized on account %s", bot.Self.UserName)

	if SEND == "SEND" {
		ticker := time.NewTicker(840 * time.Second)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				users, _ := getClosestDailyUsers()
				for _, userName := range users {
					chatID, _ := getUserChatID(userName)
					msg := tgbotapi.NewMessage(chatID, "Knock-knock. Who's there? Your stomach. \nFeed me! \n🍳🧀🥪🌮🥧🍦")
					_, _ = bot.Send(msg)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			userName := update.CallbackQuery.From.UserName
			log.Printf("[%s] %s", userName, update.CallbackQuery.Data)
			processUpdate(bot, userName, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, callbackData, update.CallbackQuery.ID)
			continue
		}
		if update.Message == nil {
			continue
		}
		userName := update.Message.From.UserName
		log.Printf("[%s] %s", userName, update.Message.Text)
		if isText(update) {
			processUpdate(bot, userName, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text, "")
		} else {
			insertUser(userName, 1)
			addUserPatience(userName, -1)
		}
		patience, err := getUserPatience(userName)
		if err != nil {
			continue
		}
		if patience == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage1)
			answer, _ := bot.Send(msg)
			insertUsageStats(userName, update.Message.Chat.ID, update.Message.MessageID, "impatient", answer.MessageID, answer.Text)
		}
		if patience < 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage2)
			answer, _ := bot.Send(msg)
			setUserPatience(userName, 1)
			insertUsageStats(userName, update.Message.Chat.ID, update.Message.MessageID, "impatient", answer.MessageID, answer.Text)
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
	if SEND == "SEND" {
		feedThemBot()
	} else {
		initDB()
		feedThemBot()
	}
}
