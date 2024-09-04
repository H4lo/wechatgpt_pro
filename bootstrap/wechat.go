package bootstrap

import (
	"fmt"
	"os"

	// "time"
	"wechatbot/handler/wechat"

	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	// "github.com/go-co-op/gocron"
)

func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}

func StartWebChat() {
	log.Info("Start WebChat Bot")
	bot := openwechat.DefaultBot(openwechat.Desktop)

	// 处理消息的回调
	bot.MessageHandler = wechat.Handler

	// 输入网址登录
	// bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 显示二维码登录
	bot.UUIDCallback = ConsoleQrCode

	reloadStorage := openwechat.NewJsonFileHotReloadStorage("token.json")
	err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
	if err != nil {
		err := os.Remove("token.json")
		if err != nil {
			return
		}

		reloadStorage = openwechat.NewJsonFileHotReloadStorage("token.json")
		err = bot.HotLogin(reloadStorage)
		if err != nil {
			return
		}
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Fatal(err)
		return
	}
	// log.Println();

	log.Printf("好友列表: ")
	friends, err := self.Friends()
	for i, friend := range friends {
		log.Println(i, friend)
	}

	log.Printf("好友数量: %d", friends.Count())

	// Storage = bot
	wechat.SetGlobalBot(bot)
	wechat.SetSchedule()

	wechat.SetDbConnection()

	groups, err := self.Groups()
	for i, group := range groups {
		// 好友名称
		log.Println(i, group.NickName)
	}

	log.Printf("公众号列表: ")
	mps, err := self.Mps()
	for i, mp := range mps {
		// 公众号名称
		log.Println(i, mp.NickName)
	}

	err = bot.Block()
	if err != nil {
		// log.Fatal(err)
		// return
		StartWebChat()
	}
}
