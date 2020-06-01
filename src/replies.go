package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

const (
	welcomeMessage = "Hi there!👋 \nMy name is FeedThemBot and I am here to help you eat regularly. \nI know, it's hard! 😫" +
		" \nI could remind you when you should grab a bite 😋" +
		" \nJust press a button below and set up a daily meal notification.\n" +
		"I will notify you with cheerful message shortly before time's up 🍏"
	timezoneMessage    = "First, I need you to select your timezone in order to deliver messages at right time ⏳"
	explanationMessage = "Now, let's set daily meal notifications.\nWhat time would you like to be notified each day?\n\nTap on buttons below and select one or more meals.\nOnce done, tap Submit to confirm."
	patienceMessage1   = "Hmmm... I'm not sure if you're using me the right way 🤔"
	patienceMessage2   = "Whatever you are doing, could you please stop it? 🙏"
	dunnoMessage       = "I don't know what to do. Sorry 😢"
)

var mealEmojis = []string{"☕", "🥐", "🍏", "🧀", "🍌", "🥨", "🍓", "🍻", "🍞", "🥕", "🥦", "🌽", "🍕", "🍩", "🍪", "🍳", "🥚", "🍆", "🍰"}

var timezoneReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT-12", "GMT-12"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-10", "GMT-10"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-8", "GMT-08"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-7", "GMT-07"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT-1", "GMT-01"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-3", "GMT-03"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-4", "GMT-04"),
		tgbotapi.NewInlineKeyboardButtonData("GMT-5", "GMT-05"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT", "GMT+00"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT+4", "GMT+04"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+3", "GMT+03"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+2", "GMT+02"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+1", "GMT+01"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("GMT+5", "GMT+05"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+7", "GMT+07"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+9", "GMT+09"),
		tgbotapi.NewInlineKeyboardButtonData("GMT+10", "GMT+10"),
	),
)

var setupReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Me!", "Feed Me!"),
	),
)

var agreeDisagreeReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Submit", "Submit"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cancel", "Cancel"),
	),
)

var skipLunchReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("On my way!😋", "Skip meal"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Stop reminder", "Stop reminder"),
	),
)

func printMarkedMealReplies(meals []string) *tgbotapi.InlineKeyboardMarkup {
	var resultMealReplies = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("06:00", "06:00"),
			tgbotapi.NewInlineKeyboardButtonData("07:00", "07:00"),
			tgbotapi.NewInlineKeyboardButtonData("08:00", "08:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("09:00", "09:00"),
			tgbotapi.NewInlineKeyboardButtonData("10:00", "10:00"),
			tgbotapi.NewInlineKeyboardButtonData("11:00", "11:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("12:00", "12:00"),
			tgbotapi.NewInlineKeyboardButtonData("13:00", "13:00"),
			tgbotapi.NewInlineKeyboardButtonData("14:00", "14:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("15:00", "15:00"),
			tgbotapi.NewInlineKeyboardButtonData("16:00", "16:00"),
			tgbotapi.NewInlineKeyboardButtonData("17:00", "17:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("18:00", "18:00"),
			tgbotapi.NewInlineKeyboardButtonData("19:00", "19:00"),
			tgbotapi.NewInlineKeyboardButtonData("20:00", "20:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("21:00", "21:00"),
			tgbotapi.NewInlineKeyboardButtonData("22:00", "22:00"),
			tgbotapi.NewInlineKeyboardButtonData("23:00", "23:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Clear", "Cancel"),
			tgbotapi.NewInlineKeyboardButtonData("Submit", "Submit"),
		),
	)
	log.Printf("[" + strings.Join(meals, ",") + "]")
	for r, row := range resultMealReplies.InlineKeyboard {
		for c, _ := range row {
			for j, meal := range meals {
				if resultMealReplies.InlineKeyboard[r][c].Text == meal {
					resultMealReplies.InlineKeyboard[r][c].Text = resultMealReplies.InlineKeyboard[r][c].Text + mealEmojis[j]
				}
			}
		}
	}
	return &resultMealReplies
}
