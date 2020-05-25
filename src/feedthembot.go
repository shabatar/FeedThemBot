package main

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"time"
)

func processCallback(userName string, callbackData string) (string, tgbotapi.InlineKeyboardMarkup, error) {
	mealsSet, _ := userMealsUTCSet(userName)
	dayFrequency, _ := getUserSelectedFrequency(userName)
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	msg := dunnoMessage
	// TODO: detect blank replyMarkup
	var err = errors.New("error processing callback: could not create keyboard buttons")
	if callbackData == "Feed Me!" {
		userTimezone, _ := getUserTimezone(userName)
		if userTimezone == -100 {
			msg = timezoneMessage
			replyMarkup = timezoneReplies
			err = nil
		} else {
			msg = feedMeWelcomeMessage
			replyMarkup = dayFreqReplies
			err = nil
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
		replyMarkup = mealReplies
		err = nil
	} else if isSetMeal(callbackData) {
		msg = "When you’d like to have " + callbackData + "? \n"
		setUserMealEditIndex(userName, callbackData[:1])
		replyMarkup = mealReplies
		err = nil
	} else if isDayTime(callbackData) {
		dayFrequency, _ := getUserSelectedFrequency(userName)
		// Update first meal
		mealEditIndex, _ := getUserMealEditIndex(userName)
		if mealEditIndex == -100 {
			setUserMealEditIndex(userName, "1")
		}
		updateUserDailyMeal(userName, callbackData)
		mealsSet, _ = userMealsUTCSet(userName)
		if dayFrequency > 0 && mealsSet {
			msg = "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?"
			replyMarkup = agreeDisagreeReplies
			err = nil
		} else if dayFrequency > 1 {
			msg = "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n" + mySetOtherMealsMessage
			replyMarkup = printEditMealsReplies(dayFrequency)
			err = nil
		} else {
			msg = "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n"
		}
	} else if isTimezone(callbackData) {
		setUserTimezone(userName, callbackData[len(callbackData)-3:])
		msg = feedMeWelcomeMessage
		replyMarkup = dayFreqReplies
		err = nil
	} else if callbackData == "/start" {
		msg = welcomeMessage
		replyMarkup = setupReplies
		err = nil
	} else if callbackData == "/restart" {
		_ = clearInsertUser(userName)
		msg = welcomeMessage
		replyMarkup = setupReplies
		err = nil
	} else if dayFrequency > 0 && mealsSet {
		msg = "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?"
		replyMarkup = agreeDisagreeReplies
		err = nil
	}
	return msg, replyMarkup, err
}

func processText(userName string, messageText string) (string, tgbotapi.InlineKeyboardMarkup, error) {
	mealsSet, _ := userMealsUTCSet(userName)
	dayFrequency, _ := getUserSelectedFrequency(userName)
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	msg := dunnoMessage
	// TODO: detect blank replyMarkup
	var err = errors.New("error processing callback: could not create keyboard buttons")
	if messageText == "/start" {
		msg = welcomeMessage
	} else if messageText == "/restart" {
		_ = clearInsertUser(userName)
		msg = welcomeMessage
		replyMarkup = setupReplies
		err = nil
	} else if isDayTime(messageText) {
		// Update meal #K
		_ = updateUserDailyMeal(userName, messageText)
		mealsSet, _ = userMealsUTCSet(userName)
		if dayFrequency > 0 && mealsSet {
			msg = "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?"
			replyMarkup = agreeDisagreeReplies
			err = nil
		} else if dayFrequency > 1 {
			msg = "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n" + mySetOtherMealsMessage
			replyMarkup = printEditMealsReplies(dayFrequency)
			err = nil
		} else {
			msg = "You eat " + strconv.Itoa(dayFrequency) + " meals/day\n"
		}
	} else if dayFrequency > 0 && mealsSet {
		msg = "You have successfully set " + strconv.Itoa(dayFrequency) + " meals/day\n" + "\n Agree?"
		replyMarkup = agreeDisagreeReplies
		err = nil
	}
	return msg, replyMarkup, err
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

			msgText, replyMarkup, err := processCallback(userName, callbackData)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
			if err == nil {
				msg.ReplyMarkup = replyMarkup
			}
			_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, msg.Text))
			_, _ = bot.Send(msg)
			continue
		}
		if update.Message == nil {
			continue
		}
		userName := update.Message.From.UserName
		log.Printf("[%s] %s", userName, update.Message.Text)
		if isText(update) {
			chatID := update.Message.Chat.ID
			msgText := update.Message.Text
			answerText, replyMarkup, err := processCallback(userName, msgText)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, answerText)
			if err == nil {
				msg.ReplyMarkup = replyMarkup
			}
			_ = setUserPatience(userName, 2)
			_ = insertUser(userName, chatID)
			_, _ = bot.Send(msg)
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
	if SEND == "SEND" {
		feedThemBot()
	} else {
		initDB()
		feedThemBot()
	}
}
