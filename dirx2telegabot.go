// Demo packet exchange of rest-data from Directum RX to chatbot Telegram.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	crtFile = filepath.Join(".", "certs", "server.crt")
	keyFile = filepath.Join(".", "certs", "server.key")
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Общие папки", "http://directum-server.ru"),
		tgbotapi.NewInlineKeyboardButtonData("Входящие", "Переход в Входящие СЭД Directum RX"),
		tgbotapi.NewInlineKeyboardButtonData("Исходящие", "Переход в Исходящие СЭД Directum RX"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Недавние документы", "Переход в Н.Документы СЭД Directum RX"),
		tgbotapi.NewInlineKeyboardButtonURL("Избранное", "http://directum-server.ru/Избранное"),
		tgbotapi.NewInlineKeyboardButtonURL("Общие папки", "http://directum-server.ru/Общие папки"),
	),
)

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	//os.Setenv("telega_botoken", "telega-token_qqqqqqqqqwwwwwwEEEEEe")

	// TLS or simple connect. Подключение по протоколу TLS или базовое
	mux := http.NewServeMux()
	mux.HandleFunc("/Документ на исполнение", http.HandlerFunc(handler))
	mux.HandleFunc("/Документ на доработку", http.HandlerFunc(handler))
	mux.HandleFunc("/Уведомление", http.HandlerFunc(handler))
	//log.Fatal(http.ListenAndServeTLS("localhost:8077", crtFile, keyFile, nil))
	log.Fatal(http.ListenAndServe("localhost:8077", mux))
}

// This handler is returning component path of URL
// Обработчик возвращает путь к компоненту URL
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nRequest of client:\n")

	fmt.Fprintf(w, "Method = %s\nURL = %s\nProto = %s\n", r.Method, r.URL, r.Proto)
	fmt.Printf("Method = %s\nURL = %s\nProto = %s\n", r.Method, r.URL, r.Proto)

	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		fmt.Printf("Header[%q] = %q\n", k, v)
	}

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Printf("Host = %q\n", r.Host)

	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	fmt.Printf("RemoteAddr = %q\n", r.RemoteAddr)

	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
		fmt.Printf("Form[%q] = %q\n", k, v)
	}
	// TODO: дескрипторы полей структуры HTTP-обработчика

	surl := strings.TrimPrefix(r.URL.Path, "/")

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

	// Loop update of data. Циклическое обновление данных
	for update := range updates {

		// Ignore any not system Message, without http data
		// Фильтрация несистемных сообщений, при получении http-данных
		if surl == "" && update.Message.Text != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.Text = "Очередь Directum RX пуста"
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			} else {
				continue
			}
		}

		// Checks if we've gotten a new messages. Проверка новых сообщений
		if update.Message != nil && surl != "" {
			// Construct a new message with chat ID and containing rest-data that we received
			// Создание нового сообщения с ID чата и полученными rest-данными
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, surl)

			// If the message got, select needs link of keyboard
			// Если сообщение поступило, пройдите по нужной ссылке на клаве
			switch update.Message.Text {
			case "dirx":
				msg.ReplyMarkup = numericKeyboard
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.Text = "Введите: dirx"
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
				continue
			}

			// Send the message. Отправка сообщений
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, Telegram show the user a message with the data received
			// Отвечая на запрос, Telegram показывает пользователю сообщение с полученными данными
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// Sends a message containing the data received
			// Отправляет сообщение, содержащее полученные данные
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
		surl = ""
	}
}

// Func reads token from file the ./botoken.conf.
// Метод получения токена из конфига
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
