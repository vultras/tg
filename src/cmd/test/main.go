package main

import (
    "log"
    "os"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "boteval/src/cmd/behx"
)

var startScreen = behx.NewScreen(
	behx.NewKeyboard(
		behx.NewButtonRow(
			behx.NewButton(
				"PRESS ME",
				behx.NewCustomAction(func(c *behx.Context){
					log.Println("pressed the button!")
				}),
			),
			behx.NewButton("PRESS ME 2", behx.NewCustomAction(func(c *behx.Context){
				log.Println("pressed another button!")
			})),
		),
		behx.NewButtonRow(
			behx.NewButton("PRESS ME 3", behx.NewCustomAction(func(c *behx.Context){
				log.Println("pressed third button!")
			})),
		),
	),
)

var behaviour = behx.Behaviour{
	StartAction: behx.NewScreenChange("start"),
	Screens: behx.ScreenMap{
		"start": startScreen,
	},
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("1"),
        tgbotapi.NewKeyboardButton("2"),
        tgbotapi.NewKeyboardButton("3"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("4"),
        tgbotapi.NewKeyboardButton("5"),
        tgbotapi.NewKeyboardButton("6"),
    ),
)

var otherKeyboard = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("a"),
        tgbotapi.NewKeyboardButton("b"),
        tgbotapi.NewKeyboardButton("c"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("d"),
        tgbotapi.NewKeyboardButton("e"),
        tgbotapi.NewKeyboardButton("f"),
    ),
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore non-Message updates
            continue
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

        switch update.Message.Text {
        case "open":
            msg.ReplyMarkup = numericKeyboard
        case "letters" :
            msg.ReplyMarkup = otherKeyboard
        case "close":
            msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
        }

        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}

