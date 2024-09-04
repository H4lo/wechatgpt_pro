package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	ChatGpt ChatGptConfig `json:"chatgpt" mapstructure:"chatgpt" yaml:"chatgpt"`
	GroupChat GroupChatConfig `json:"wechat_group_onfig" mapstructure:"wechat_group_onfig" yaml:"wechat_group_onfig"`
}

type ChatGptConfig struct {
	Token         string  `json:"token,omitempty"  mapstructure:"token,omitempty"  yaml:"token,omitempty"`
	Wechat        *string `json:"wechat,omitempty" mapstructure:"wechat,omitempty" yaml:"wechat,omitempty"`
	WechatKeyword *string `json:"wechat_keyword"   mapstructure:"wechat_keyword"   yaml:"wechat_keyword"`
	Telegram      *string `json:"telegram"         mapstructure:"telegram"         yaml:"telegram"`
	TgWhitelist   *string `json:"tg_whitelist"     mapstructure:"tg_whitelist"     yaml:"tg_whitelist"`
	TgKeyword     *string `json:"tg_keyword"       mapstructure:"tg_keyword"       yaml:"tg_keyword"`
	OpenAiUrl	  *string `json:"openai_url"       mapstructure:"openai_url"       yaml:"openai_url"`
	OpenAiModel	  *string `json:"openai_model"       mapstructure:"openai_model"       yaml:"openai_model"`
	Prompt		  *string `json:"prompt"       mapstructure:"prompt"       yaml:"prompt"`
}


type GroupChatConfig struct {
	DaliyGroupName string `json:"daliy_group_name"       mapstructure:"daliy_group_name"       yaml:"daliy_group_name"`
	SelfName string `json:"self_name"       mapstructure:"self_name"       yaml:"self_name"`
	RobotName string `json:"robot_name"       mapstructure:"robot_name"       yaml:"robot_name"`
}


func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./local")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

func GetDaliyGroupName() string {
	return config.GroupChat.DaliyGroupName
}

func GetSelfName() string {
	return config.GroupChat.SelfName
}

func GetRobotName() string {
	return config.GroupChat.RobotName
}

func GetWechat() *string {
	wechat := getEnv("wechat")
	if wechat != nil {
		return wechat
	}

	if config == nil {
		return nil
	}

	if wechat == nil {
		wechat = config.ChatGpt.Wechat
	}
	return wechat
}

func GetWechatKeyword() *string {
	keyword := getEnv("wechat_keyword")

	if keyword != nil {
		return keyword
	}

	if config == nil {
		return nil
	}

	if keyword == nil {
		keyword = config.ChatGpt.WechatKeyword
	}
	return keyword
}

func GetTelegram() *string {
	tg := getEnv("telegram")
	if tg != nil {
		return tg
	}

	if config == nil {
		return nil
	}

	if tg == nil {
		tg = config.ChatGpt.Telegram
	}
	return tg
}

func GetTelegramKeyword() *string {
	tgKeyword := getEnv("tg_keyword")

	if tgKeyword != nil {
		return tgKeyword
	}

	if config == nil {
		return nil
	}

	if tgKeyword == nil {
		tgKeyword = config.ChatGpt.TgKeyword
	}
	return tgKeyword
}

func GetTelegramWhitelist() *string {
	tgWhitelist := getEnv("tg_whitelist")

	if tgWhitelist != nil {
		return tgWhitelist
	}

	if config == nil {
		return nil
	}

	if tgWhitelist == nil {
		tgWhitelist = config.ChatGpt.TgWhitelist
	}
	return tgWhitelist
}

func GetOpenAiApiKey() *string {
	apiKey := getEnv("api_key")
	if apiKey != nil {
		return apiKey
	}

	if config == nil {
		return nil
	}

	if apiKey == nil {
		apiKey = &config.ChatGpt.Token
	}
	return apiKey
}

func GetOpenAiUrl() *string {
	openAiUrl := getEnv("openAiUrl")
	if openAiUrl != nil {
		return openAiUrl
	}

	if openAiUrl == nil {
		openAiUrl = config.ChatGpt.OpenAiUrl
	}
	return openAiUrl
}

func GetOpenAiModel() *string {
	OpenAiModel := getEnv("OpenAiModel")
	if OpenAiModel != nil {
		return OpenAiModel
	}

	if OpenAiModel == nil {
		OpenAiModel = config.ChatGpt.OpenAiModel
	}
	return OpenAiModel
}

func GetOpenAiPrompt() *string {
	prompt := getEnv("prompt")
	if prompt != nil {
		return prompt
	}

	if prompt == nil {
		prompt = config.ChatGpt.Prompt
	}
	return prompt
}

func getEnv(key string) *string {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = os.Getenv(strings.ToUpper(key))
	}

	if len(value) > 0 {
		return &value
	}

	if config == nil {
		return nil
	}

	if len(value) > 0 {
		return &value
	}

	if config.ChatGpt.WechatKeyword != nil {
		value = *config.ChatGpt.WechatKeyword
	}
	return nil
}
