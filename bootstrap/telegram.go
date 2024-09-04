package bootstrap

import (
	"strings"
	"time"

	"wechatbot/config"
	"wechatbot/handler/telegram"
	"wechatbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func StartTelegramBot() {
	log.Info("Start Telegram Bot")
	telegramKey := config.GetTelegram()
	if telegramKey == nil {
		log.Info("未找到tg token,不启动tg bot")
		return
	}

	bot, err := tgbotapi.NewBotAPI(*telegramKey)
	if err != nil {
		log.Error("tg bot 启动失败：", err.Error())
		return
	}

	bot.Debug = false
	log.Info("Authorized on account: ", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	for len(updates) != 0 {
		<-updates
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		chatID := update.Message.Chat.ID
		chatUserName := update.Message.Chat.UserName

		tgUserNameStr := config.GetTelegramWhitelist()
		if tgUserNameStr != nil {
			tgUserNames := strings.Split(*tgUserNameStr, ",")
			if len(tgUserNames) > 0 {
				found := false
				for _, name := range tgUserNames {
					if name == chatUserName {
						found = true
						break
					}
				}

				if !found {
					log.Error("用户设置了私人私用，白名单以外的人不生效: ", chatUserName)
					continue
				}
			}
		}

		tgKeyWord := config.GetTelegramKeyword()
		var reply *string
		// 如果设置了关键字就以关键字为准，没设置就所有消息都监听
		if tgKeyWord != nil {
			content, key := utils.ContainsI(text, *tgKeyWord)
			if len(key) == 0 {
				continue
			}

			splitItems := strings.Split(content, key)
			if len(splitItems) < 2 {
				continue
			}

			requestText := strings.TrimSpace(splitItems[1])
			log.Info("问题：", requestText)
			reply = telegram.Handle(requestText)
		} else {
			log.Info("问题：", text)
			reply = telegram.Handle(text)
		}

		if reply == nil {
			continue
		}

		msg := tgbotapi.NewMessage(chatID, *reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Errorf("发送消息出错:%s", err.Error())
			continue
		}

		log.Info("回答：", *reply)
	}

	select {}
}
