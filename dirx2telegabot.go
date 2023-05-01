package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	crtFile = filepath.Join(".", "certs", "server.crt")
	keyFile = filepath.Join(".", "certs", "server.key")
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	//os.Setenv("telega_botoken", "5853322065:AAHwqJwOEVOrLMcpKf-vOW5rOYp4eByFevs")

	// TLS connect. Подключение по протоколу TLS
	mux := http.NewServeMux()
	mux.HandleFunc("/Приказ", http.HandlerFunc(handler))
	mux.HandleFunc("/Заявка", http.HandlerFunc(handler))
	mux.HandleFunc("/Задание", http.HandlerFunc(handler))
	//http.HandleFunc("/Заявка", handler)
	//log.Fatal(http.ListenAndServeTLS("localhost:8077", crtFile, keyFile, nil))
	log.Fatal(http.ListenAndServe("localhost:8077", mux))
}

// This handler is returning component path of URL. Обработчик возвращает путь к компоненту URL
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)

	// Getting token from config. Получение токена из конфига.
	chkey := make(chan string)
	go func() {
		chkey <- readToken()
	}()

	//bot, err := tgbotapi.NewBotAPI(os.Getenv("telega_botoken"))
	bot, err := tgbotapi.NewBotAPI(<-chkey)
	if err != nil {
		log.Fatalf("Error of bot-api: %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Text {
			case "open":
				msg.ReplyMarkup = numericKeyboard

			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}

// Func reads token from file the ./botoken.conf. Метод получения токена из конфига.
func readToken() string {
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
