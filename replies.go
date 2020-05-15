package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	EmojiCat = "\xF0\x9F\x90\x88"
	EmojiHard = "\xF0\x9F\x98\xAB"
	EmojiPity = "\xF0\x9F\x98\x96"
	EmojiYumm = "\xF0\x9F\x98\x8B"
	EmojiWave = "\xF0\x9F\x91\x8B"
)

const (
	welcomeMessage = 	"Hi there!" + EmojiWave + " My name is FeedThemBot and I am here to help you eat regularly.\nI know, it's difficult! " + EmojiHard + 
					" \nI could remind you when you should probably grab a bite. " + EmojiYumm + 
					" \nJust press a button and set up a meal notification. You can also use me to remind you about feeding your pets " + EmojiCat
	patienceMessage1 = "Hmmm... I'm not sure if you're using me the right way"
	patienceMessage2 = "Could you please stop it?"
	dunnoMessage = "I don't know what to do. Sorry" + EmojiPity
)

const (
	feedMeWelcomeMessage = "OK, I shall remind you whenever you ought to eat. How often you’d like to eat per day?"
	feedPetWelcomeMessage = "I don't know how to feed someone else. Sorry" + EmojiPity
)

const (
	myFirstMealMessage = "When you’d like to have 1st meal?"
	mySetOtherMealsMessage = "When you’d like to have (set up time manually or set period):"
)

var startReplies = tgbotapi.NewInlineKeyboardMarkup(
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
