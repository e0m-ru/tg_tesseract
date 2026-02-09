package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_TOKEN"))
	if err != nil {
		log.Fatalf("Auth error: %s", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}
		if update.Message.Command() == "/report" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "http://80.73.24.45")
			msg.ReplyToMessageID = update.Message.MessageID
		}

		if update.Message.Photo != nil {
			photos := update.Message.Photo
			photo_url, _ := bot.GetFileDirectURL(photos[len(photos)-1].FileID)
			pu, _ := url.Parse(photo_url)
			fn := path.Base(pu.Path)
			response, err := http.Get(photo_url)
			if err != nil {
				log.Fatalf("GetPhoto error: %s", err)
			}

			fp := path.Join("tmp/", fn)
			file, err := os.Create(fp)

			if err != nil {
				log.Fatalf("Create file error: %s", err)
			}
			_, err = io.Copy(file, response.Body)
			if err != nil {
				log.Fatalf("Copy to file error: %s", err)
			}

			response.Body.Close()

			cmd := exec.Command("tesseract", fp, "stdout", "-l", "rus")
			b, err := cmd.Output()
			if err != nil {
				log.Fatalf("Tesseract error: %s", err)
			}

			file.Close()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(b))
			msg.ReplyToMessageID = update.Message.MessageID

			_, err = bot.Send(msg)

			if err != nil {
				log.Fatalf("Message error: %s", err)
			}
		}
	}
}
