package main

import (
	"bufio"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	crtFile = filepath.Join(".", "certs", "client.crt")
	keyFile = filepath.Join(".", "certs", "client.key")
)

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// TLS connect. Подключение по протоколу TLS
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS("localhost:8077", crtFile, keyFile, nil))
}

// This handler is returning component path of URL. Обработчик возвращает путь к компоненту URL
func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)

	////////////////////////////////////////////////////////////////////
	// Parsing rows of data from a http-request
	// Парсинг строк данных http-запроса
	//sd := &DataType{}
	counts := make(map[string]int)
	//Slice for save all keys from mapping. Срез для хранения всех ключей мапы
	datakeys := make([]string, 0, len(counts))
	for _, line := range strings.Split(string(pars), ":") {
		counts[line]++
		log.Println(line) // Checks data of mapping. Проверка данных мапы
	}
	// Sorts keys to list and to sets values in order
	// Сортировка ключей для перечисления и присваивания значений по порядку
	for countkeys := range counts {
		datakeys = append(datakeys, countkeys)
	}
	sort.Strings(datakeys)
	for _, countkeys := range datakeys {
		fmt.Printf("\nCountkeys: %v\nCounts: %v\n", countkeys, counts[countkeys])
		if countkeys != "" {
			if sp.weight == "" {
				sp.mu.Lock()
				sp.weight = countkeys
				sp.mu.Unlock()
			} else {
				if sp.plan_default == "" {
					sp.mu.Lock()
					sp.plan_default = countkeys
					sp.mu.Unlock()
				} else {
					sp.mu.Lock()
					sp.value = countkeys
					sp.mu.Unlock()
				}
			}
		}
	}
	// Output for test. Тестовый вывод данных.
	log.Println("Weight: ", sp.weight)
	log.Println("Plan_default: ", sp.plan_default)
	log.Println("Value : ", sp.value)

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
