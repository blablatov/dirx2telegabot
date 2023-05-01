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

	os.Setenv("telega_botoken", "5853322065:AAHwqJwOEVOrLMcpKf-vOW5rOYp4eByFevs")

	// TLS connect. Подключение по протоколу TLS
	//mux := http.NewServeMux()
	//mux.HandleFunc("/Заявка", handler)
	//mux.HandleFunc("/Задание", handler)
	http.HandleFunc("/Заявка", handler)
	//log.Fatal(http.ListenAndServeTLS("localhost:8077", crtFile, keyFile, nil))
	log.Fatal(http.ListenAndServe("localhost:8077", nil))
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
		//if update.Message != nil || r.URL.Path != "" { // If we got a message. Если получаем сообщение
		if r.URL.Path != "" { // If we got a message. Если получаем сообщение
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			tmsg := r.URL.Path
			msg.Text = (tmsg)

			bot.Send(msg)
			//r.URL.Path = ""
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
