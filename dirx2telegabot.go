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

	for update := range updates {
		if r.URL.Path == "" { //
			continue
		}

		/*if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}*/

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch r.URL.Path {
		case "/Приказ":
			update.Message.Command()
			msg.Text = "Ознакомтесь с Приказом."
			r.URL.Path = ""
		case "/Заявка":
			update.Message.Command()
			msg.Text = "Вам поступила Заявка."
			r.URL.Path = ""
		case "/Задание":
			update.Message.Command()
			msg.Text = "Вам поступило Задание."
			r.URL.Path = ""
		default:
			update.Message.Command()
			msg.Text = "Поступил документ."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
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
