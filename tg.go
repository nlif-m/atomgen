package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nlif-m/atomgen/ytdlp"
)

func TgBot(ag Atomgen, atomfileUpdateChan chan bool) {
	bot, err := tgbotapi.NewBotAPI(ag.cfg.TelegramBotToken)
	if err != nil {
		log.Panic(err)
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
		log.Printf("%d\t%s\n", update.Message.From.ID, ag.cfg.TelegramAdminId)
		if fmt.Sprint(update.Message.From.ID) != ag.cfg.TelegramAdminId {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access Denied")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		txt := update.Message.Text
		ytType, url, down := ag.ytdlp.IsDownloadable(txt)
		if !down {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%q is not downloadable", txt))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		switch ytType {
		case ytdlp.YoutubeVideoType, ytdlp.VkVideoType:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%q is downloadable and is %q, start downloading %q", txt, ytType, url))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

			go func(ctx tgbotapi.Update) {
				err := ag.DownloadURL(url, true)
				if err != nil {
					msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, fmt.Sprintf("Sorry but failed to download %q", url))
					msg.ReplyToMessageID = ctx.Message.MessageID
					bot.Send(msg)
					return
				}

				atomfileUpdateChan <- true
				msg := tgbotapi.NewMessage(ctx.Message.Chat.ID, fmt.Sprintf("Successfully downloaded %q", url))
				msg.ReplyToMessageID = ctx.Message.MessageID
				bot.Send(msg)
			}(update)
		case ytdlp.YoutubePlaylistType:
			// TODO: Implement adding playlists to config url list
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%q is youtube playlist and in new version it will be supported", url))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}
