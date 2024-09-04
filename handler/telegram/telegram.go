package telegram

import (
	"strings"

	"wechatbot/openai"

	log "github.com/sirupsen/logrus"
)

func Handle(msg string) *string {
	requestText := strings.TrimSpace(msg)
	reply, err := openai.Completions(requestText, 0)
	if err != nil {
		log.Error(err)
	}
	return reply
}
