package main

import (
	"core"
	"fmt"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	l := len(os.Args)
	if l < 4 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: kursNotifier ini_file_name chat_ids_filename result_folder [-q]")
		return
	}

	config, err := core.LoadConfiguration(os.Args[1])
	if err != nil {
		panic(err)
	}

	enabledChatIDs, err := core.LoadEnabledChatIDs(os.Args[2])
	if err != nil {
		panic(err)
	}

	quiet := l == 4 && os.Args[4] == "-q"
	doNotSend := l == 4 && os.Args[4] == "-n"

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := bot.GetUpdates(updateConfig)
	if err != nil {
		panic(err)
	}
	for _, update := range updates {
		if update.Message != nil {
			if !quiet {
				fmt.Printf("Message from %v: %v\n", update.Message.Chat.ID, update.Message.Text)
			}
			for _, chatID := range config.ChatIDs {
				if update.Message.Chat.ID == chatID {
					switch update.Message.Text {
					case "/start":
						enabledChatIDs, _ = addToEnabledChatIDs(os.Args[2], enabledChatIDs, chatID)
					case "/stop":
						enabledChatIDs, _ = removeFromEnabledChatIDs(os.Args[2], enabledChatIDs, chatID)
					}
				}
			}
		}
	}

	message := core.BuildKursMessage(os.Args[3], quiet, config.Sources)
	if message != "" {
		if !quiet {
			fmt.Print(message)
		}
		if !doNotSend {
			for _, chatID := range enabledChatIDs {
				msg := tgbotapi.NewMessage(chatID, message)
				_, _ = bot.Send(msg)
			}
		}
	}
}

func addToEnabledChatIDs(fileName string, enabledChatIDs []int64, chatID int64) ([]int64, error) {
	for _, id := range enabledChatIDs {
		if id == chatID {
			return enabledChatIDs, nil
		}
	}
	enabledChatIDs = append(enabledChatIDs, chatID)
	return enabledChatIDs, core.SaveEnabledChatIDs(fileName, enabledChatIDs)
}

func removeFromEnabledChatIDs(fileName string, enabledChatIDs []int64, chatID int64) ([]int64, error) {
	var result []int64
	for _, id := range enabledChatIDs {
		if id != chatID {
			result = append(result, id)
		}
	}
	return result, core.SaveEnabledChatIDs(fileName, result)
}
