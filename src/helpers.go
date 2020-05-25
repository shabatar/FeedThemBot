package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"reflect"
	"regexp"
	"strconv"
)

func isText(update tgbotapi.Update) bool {
	return reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != ""
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isDayFrequency(s string) bool {
	value, _ := regexp.MatchString(`\df`, s)
	return value
}

func isDayTime(s string) bool {
	value, _ := regexp.MatchString(`\d\d:\d\d`, s)
	return value
}

func isSetMeal(s string) bool {
	value, _ := regexp.MatchString(`\d[nrt][dh] meal`, s)
	return value
}

func isTimezone(s string) bool {
	value, _ := regexp.MatchString(`GMT[+-]\d\d?`, s)
	return value
}
