package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"strings"
	"time"
)

func ParseGetRequest(line string) string {
	serviceName := line[4:]
	serviceName = strings.Replace(serviceName, " ", "", -1)
	return serviceName
}

func ParseDelRequest(line string) string {
	serviceName := line[4:]
	serviceName = strings.Replace(serviceName, " ", "", -1)
	return serviceName
}

func ParseSetRequest(line string) (string, string, string, error) {
	splittedLine := strings.Split(line[4:], ",")
	if len(splittedLine) != 3 {
		return "", "", "", errors.New("Incorrect format of input data")
	}
	serviceName := strings.Replace(splittedLine[0], " ", "", -1)
	login := strings.Replace(splittedLine[1], " ", "", -1)
	password := strings.Replace(splittedLine[2], " ", "", -1)
	return serviceName, login, password, nil
}

var service *Server

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text[:4] == "/get" {
		serviceName := ParseGetRequest(update.Message.Text)
		login, password, err := service.Get(update.Message.Chat.ID, serviceName)

		var textAnswer string

		if err != nil {
			textAnswer = "There was an error with get method"
		} else {
			textAnswer = fmt.Sprintf("Login : %s \n Password : %s", login, password)
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   textAnswer,
		})
	}
	if update.Message.Text[:4] == "/set" {
		serviceName, login, password, errRead := ParseSetRequest(update.Message.Text)

		err := service.Set(update.Message.Chat.ID, serviceName, login, password)
		var textAnswer string

		if errRead != nil {
			textAnswer = "There was an error with reading your data"
		} else if err != nil {
			textAnswer = fmt.Sprintf("There was an error with set method : ", err.Error())
		} else {
			textAnswer = "Done"
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   textAnswer,
		})
	}
	if update.Message.Text[:4] == "/del" {
		serviceName := ParseDelRequest(update.Message.Text)
		service.Delete(update.Message.Chat.ID, serviceName)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Succesfully deleted",
		})
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New("6013761542:AAEqQLoLXwjX6Go3Gzf31vFlVvyIEllJFT8", opts...)
	if err != nil {
		panic(err)
	}
	service = NewServer()

	go func() {
		currentKeysSet := make(map[string]bool)
		for {
			time.Sleep(time.Minute)
			newKeys := service.GetKeys()
			newKeysSet := make(map[string]bool)
			for _, key := range newKeys {
				newKeysSet[key] = true
				if _, ok := currentKeysSet[key]; ok {
					service.Delete(DecodeService(key))
					fmt.Println("Deleted key ", key)
				}
			}
			currentKeysSet = newKeysSet
		}
	}()
	b.Start(ctx)
}
