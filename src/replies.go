package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	welcomeMessage = "Hi there!👋 My name is FeedThemBot and I am here to help you eat regularly.\nI know, it's difficult! 😫" + 
					" \nI could remind you when you should probably grab a bite. 😋" + 
					" \nJust press a button and set up a meal notification. You can also use me to remind you about feeding your pets 🐈"
	timezoneMessage = "First, I need you to select your timezone in order to deliver messages at right time ⏳"
	patienceMessage1 = "Hmmm... I'm not sure if you're using me the right way"
	patienceMessage2 = "Could you please stop it?"
	dunnoMessage = "I don't know what to do. Sorry 😢"
)

const (
	feedMeWelcomeMessage = "OK, I shall remind you whenever you ought to eat. How often you’d like to eat per day?"
	feedPetWelcomeMessage = "I don't know how to feed someone else. Sorry 😢"
)

const (
	myFirstMealMessage = "When you’d like to have 1st meal?"
	mySetOtherMealsMessage = "When you’d like to have (set up time manually or set period):"
)

var timezoneReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT-12","GMT-12"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-10","GMT-10"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-8","GMT-08"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-7","GMT-07"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT-1","GMT-01"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-3","GMT-03"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-4","GMT-04"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-5","GMT-05"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT","GMT+00"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT+4","GMT+04"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+3","GMT+03"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+2","GMT+02"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+1","GMT+01"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT+5","GMT+05"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+7","GMT+07"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+9","GMT+09"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+10","GMT+10"),
	),
)

var setupReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Me!","Feed Me!"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Someone Else!","Feed Someone Else!"),
	),
)

var dayFreqReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1","1t"),
		tgbotapi.NewInlineKeyboardButtonData("2","2t"),
		tgbotapi.NewInlineKeyboardButtonData("3","3t"),
		tgbotapi.NewInlineKeyboardButtonData("4","4t"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5","5t"),
		tgbotapi.NewInlineKeyboardButtonData("6","6t"),
		tgbotapi.NewInlineKeyboardButtonData("7","7t"),
	    tgbotapi.NewInlineKeyboardButtonData("8","8t"),
	),
)

var firstMealReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("07:00","07:00"),
		tgbotapi.NewInlineKeyboardButtonData("08:00","08:00"),
		tgbotapi.NewInlineKeyboardButtonData("09:00","09:00"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("10:00","10:00"),
		tgbotapi.NewInlineKeyboardButtonData("11:00","11:00"),
		tgbotapi.NewInlineKeyboardButtonData("12:00","12:00"),
	),
)

var agreeDisagreeReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Submit","Submit"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cancel","Cancel"),
	),
)

func printEditMealsReplies(mealsNumber int) tgbotapi.InlineKeyboardMarkup {
	replies := [8]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("2nd meal", "2nd meal"),
		tgbotapi.NewInlineKeyboardButtonData("3rd meal", "3rd meal"),
		tgbotapi.NewInlineKeyboardButtonData("4th meal", "4th meal"),
		tgbotapi.NewInlineKeyboardButtonData("5th meal", "5th meal"),
		tgbotapi.NewInlineKeyboardButtonData("6th meal", "6th meal"),
		tgbotapi.NewInlineKeyboardButtonData("7th meal", "7th meal"),
		tgbotapi.NewInlineKeyboardButtonData("8th meal", "8th meal"),
		tgbotapi.NewInlineKeyboardButtonData("periodically", "periodically")}
	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	var row3 []tgbotapi.InlineKeyboardButton
	var firstMealReplies tgbotapi.InlineKeyboardMarkup

	half := mealsNumber / 2
	row1 = append(row1, replies[:half]...)
	row2 = append(row2, replies[half:mealsNumber-1]...)
	row3 = append(row3, replies[len(replies)-1])
	if mealsNumber == 2 {
		firstMealReplies = tgbotapi.NewInlineKeyboardMarkup(
			row1,
			row3,
		)
		return firstMealReplies
	}
	firstMealReplies = tgbotapi.NewInlineKeyboardMarkup(
		row1,
		row2,
		row3,
	)
	return firstMealReplies
}