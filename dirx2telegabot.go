package main

import (
	"bufio"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// Getting token from config. Получение токена из конфига.
	chkey := make(chan string)
	go func() {
		chkey <- ReadBotKey()
	}()
	sbot := <-chkey

	bot, err := tgbotapi.NewBotAPI(sbot)
	if err != nil {
		//log.Panic(err)
		log.Fatalf("Error of bot-api: %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message. Если получаем сообщение
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

// Func reads token from file the ./botoken.conf. Метод получения токена из конфига.
func ReadBotKey() string {
	var botkey string
	bk, err := os.Open("botoken.conf")
	if err != nil {
		log.Fatalf("Error open config: %v", err)
	}
	defer bk.Close()
	input := bufio.NewScanner(bk)
	for input.Scan() {
		botkey = input.Text()
	}
	return botkey
}
