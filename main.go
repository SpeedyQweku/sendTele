package main

import (
	"flag"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string, wg *sync.WaitGroup) {
	defer wg.Done()

	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func main() {
	log.Println("Send messages to Telegram...")
	var (
		btoken string
		chatid int64
		mgs    string
	)

	flag.StringVar(&btoken, "t", "", "Bot token")
	flag.Int64Var(&chatid, "c", 0, "Chat ID")
	flag.StringVar(&mgs, "m", "", "Message")
	flag.Parse()

	bot, err := tgbotapi.NewBotAPI(btoken)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go sendMessage(bot, chatid, mgs, &wg)

	wg.Wait()
}
