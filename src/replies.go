package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

const (
	welcomeMessage = "Hi there!👋 \nMy name is FeedThemBot and I am here to help you eat regularly. I know, it's hard! 😫" +
		" \nI could remind you when you should probably grab a bite 😋" +
		" \nJust press a button below and set up a meal notification.\n" +
		"I will notify you with cheerful message once time's up 🍏"
	timezoneMessage  = "First, I need you to select your timezone in order to deliver messages at right time ⏳"
	patienceMessage1 = "Hmmm... I'm not sure if you're using me the right way 🤔"
	patienceMessage2 = "Whatever you are doing, could you please stop it? 🙏"
	dunnoMessage     = "I don't know what to do. Sorry 😢"
)

const (
	feedMeWelcomeMessage  = "OK, I shall remind you whenever you ought to eat. \nHow many meal reminders you'd like to have each day?"
	feedPetWelcomeMessage = "I don't know how to feed someone else. Sorry 😢"
)

const (
	myFirstMealMessage     = "When you’d like to have breakfast?🍳"
	mySetOtherMealsMessage = "When you’d like to have (set up time manually or set period):"
)

var mealEmojis = []string{"☕", "🥐", "🍏", "🧀", "🍌", "🥨", "🍓", "🍻"}

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
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Feed Someone Else! (tbd)", "Feed Someone Else!"),
	),
)

var dayFreqReplies = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1"+mealEmojis[0], "1f"),
		tgbotapi.NewInlineKeyboardButtonData("2"+mealEmojis[1], "2f"),
		tgbotapi.NewInlineKeyboardButtonData("3"+mealEmojis[2], "3f"),
		tgbotapi.NewInlineKeyboardButtonData("4"+mealEmojis[3], "4f"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5"+mealEmojis[4], "5f"),
		tgbotapi.NewInlineKeyboardButtonData("6"+mealEmojis[5], "6f"),
		tgbotapi.NewInlineKeyboardButtonData("7"+mealEmojis[6], "7f"),
		tgbotapi.NewInlineKeyboardButtonData("8"+mealEmojis[7], "8f"),
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

func printMarkedMealReplies(meals []string) tgbotapi.InlineKeyboardMarkup {
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
	return resultMealReplies
}

func printEditMealsReplies(mealsNumber int) tgbotapi.InlineKeyboardMarkup {
	replies := [9]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("1st meal"+mealEmojis[0], "1st meal"),
		tgbotapi.NewInlineKeyboardButtonData("2nd meal"+mealEmojis[1], "2nd meal"),
		tgbotapi.NewInlineKeyboardButtonData("3rd meal"+mealEmojis[2], "3rd meal"),
		tgbotapi.NewInlineKeyboardButtonData("4th meal"+mealEmojis[3], "4th meal"),
		tgbotapi.NewInlineKeyboardButtonData("5th meal"+mealEmojis[4], "5th meal"),
		tgbotapi.NewInlineKeyboardButtonData("6th meal"+mealEmojis[5], "6th meal"),
		tgbotapi.NewInlineKeyboardButtonData("7th meal"+mealEmojis[6], "7th meal"),
		tgbotapi.NewInlineKeyboardButtonData("8th meal"+mealEmojis[7], "8th meal"),
		tgbotapi.NewInlineKeyboardButtonData("periodically", "periodically")}
	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	var row3 []tgbotapi.InlineKeyboardButton
	var firstMealReplies tgbotapi.InlineKeyboardMarkup

	half := mealsNumber / 2
	row1 = append(row1, replies[:half]...)
	row2 = append(row2, replies[half:mealsNumber]...)
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
