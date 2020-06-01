package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"time"
)

func processCallback(userName string, callbackData string) (msg string, replyMarkup *tgbotapi.InlineKeyboardMarkup, payload string) {
	msg = dunnoMessage
	payload = "dunno"
	usr, _ := getUserData(userName)
	if callbackData == "Feed Me!" {
		if usr.userTimezone == -100 {
			msg = timezoneMessage
			payload = "timezone"
			replyMarkup = &timezoneReplies
		} else {
			msg = explanationMessage
			payload = "editFirst"
			usr, _ = getUserData(userName)
			replyMarkup = printMarkedMealReplies(usr.userMealsUTC)
		}
	} else if callbackData == "Submit" {
		msg = "You are all set! Wait for the notifications"
		payload = "submit"
		syncTimezone(userName)
		migrateDailyUser(userName)
		updateDailySchedule()
	} else if callbackData == "Cancel" {
		msg = explanationMessage
		payload = "editNext"
		createUserDailyMeals(userName)
		usr, _ = getUserData(userName)
		replyMarkup = printMarkedMealReplies(usr.userMealsUTC)
	} else if callbackData == "Skip meal" {
		msg = "Okay, skipping meal."
		payload = "skip"
		setDailyUserSkipLunch(userName)
	} else if callbackData == "Stop reminder" {
		msg = "Okay, stopping reminder."
		payload = "stop"
		stopDailySchedule(userName)
		updateDailySchedule()
		removeStoppedFromDailySchedule()
	} else if isDayTime(callbackData) {
		msg = explanationMessage
		payload = "editNext"
		mealEditIndex := usr.userMealEditIndex
		if mealEditIndex == -100 {
			setUserMealEditIndex(userName, "1")
		} else {
			setUserMealEditIndex(userName, strconv.Itoa(usr.userMealEditIndex+1))
		}
		updateUserDailyMeal(userName, callbackData)
		usr, _ = getUserData(userName)
		replyMarkup = printMarkedMealReplies(usr.userMealsUTC)
	} else if isTimezone(callbackData) {
		setUserTimezone(userName, callbackData[len(callbackData)-3:])
		msg = explanationMessage
		payload = "editFirst"
		usr, _ = getUserData(userName)
		replyMarkup = printMarkedMealReplies(usr.userMealsUTC)
	} else if callbackData == "/start" {
		_ = clearInsertUser(userName)
		_ = removeFromDailySchedule(userName)
		msg = welcomeMessage
		replyMarkup = &setupReplies
		payload = "welcome"
	} else if callbackData == "/restart" {
		_ = clearInsertUser(userName)
		_ = removeFromDailySchedule(userName)
		msg = welcomeMessage
		replyMarkup = &setupReplies
		payload = "welcome"
	} else if callbackData == "/stop" {
		msg = "Okay, stopping reminder."
		stopDailySchedule(userName)
		updateDailySchedule()
		removeStoppedFromDailySchedule()
		payload = "stop"
	}
	return msg, replyMarkup, payload
}

func processUpdate(bot *tgbotapi.BotAPI, userName string, chatID int64, messageID int, messageText string, callbackID string) {
	answerText, replyMarkup, payload := processCallback(userName, messageText)
	msg := tgbotapi.NewMessage(chatID, answerText)
	if replyMarkup != nil {
		msg.ReplyMarkup = *replyMarkup
	}
	if payload == "editNext" {
		editMessageID, _ := getMessageIDFromUsage(userName, "editFirst")
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      msg.BaseChat.ChatID,
				MessageID:   editMessageID,
				ReplyMarkup: replyMarkup,
			},
			Text:      msg.Text,
			ParseMode: "HTML",
		}
		answer, _ := bot.Send(edit)
		insertUsageStats(userName, chatID, messageID, messageText, answer.MessageID, answer.Text, payload)
		return
	}
	_ = setUserPatience(userName, 2)
	_ = insertUser(userName, chatID)
	answer, _ := bot.Send(msg)
	if callbackID != "" {
		_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, msg.Text))
	}
	insertUsageStats(userName, chatID, messageID, messageText, answer.MessageID, answer.Text, payload)
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
		ticker := time.NewTicker(30 * time.Second)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				updateDailySchedule()
				removeStoppedFromDailySchedule()
				users, _ := getClosestDailyUsers()
				for _, userName := range users {
					chatID, _ := getUserChatID(userName)
					msg := tgbotapi.NewMessage(chatID, "Knock-knock. Who's there? Your stomach. \nFeed me! \n🍳🧀🥪🌮🥧🍦")
					msg.ReplyMarkup = skipLunchReplies
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
			insertUsageStats(userName, update.Message.Chat.ID, update.Message.MessageID, "impatient", answer.MessageID, answer.Text, "impatient1")
		}
		if patience < 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, patienceMessage2)
			answer, _ := bot.Send(msg)
			setUserPatience(userName, 1)
			insertUsageStats(userName, update.Message.Chat.ID, update.Message.MessageID, "impatient", answer.MessageID, answer.Text, "impatient2")
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
	if err := createDailySchedule(); err != nil {
		panic(err)
	}
	if err := createDiffFunction(); err != nil {
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
