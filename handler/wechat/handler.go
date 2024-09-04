package wechat

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"wechatbot/config"
	"wechatbot/utils"

	"github.com/eatmoreapple/openwechat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type Type string

const (
	GroupHandler = "group"
)

var handlers map[Type]MessageHandlerInterface

var gBot *openwechat.Bot

var db *sql.DB

func init() {
	handlers = make(map[Type]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
}

func Handler(msg *openwechat.Message) {
	err := handlers[GroupHandler].handle(msg)
	if err != nil {
		log.Errorf("handle error: %s\n", err.Error())
		return
	}
}

func SetGlobalBot(bot *openwechat.Bot) {
	gBot = bot
}

func GetGlobalBot() *openwechat.Bot {
	return gBot
}

func timeTips(tips string) {
	log.Println("Trigger gocron tasks...")
	self, _ := GetGlobalBot().GetCurrentUser()

	// bootstrap.gBot.Self
	var sendToGroupName = config.GetDaliyGroupName()
	if sendToGroupName == "" {
		log.Println("Error happened in get group name.")
		return
	}

	// 获取到group对象
	groups, err := self.Groups()
	if err != nil {
		log.Fatal(err)
	}

	for i, group := range groups {
		log.Println(i, group)
		if group.NickName == sendToGroupName {
			img, _ := os.Open("/home/test/dev/wechatgpt/imgs/time.jpg")
			defer img.Close()
			var sendMsg, _ = self.SendImageToGroup(group, img)

			log.Println(sendMsg)
			// return nil
		}
	}

	var g = groups.SearchByNickName(20, sendToGroupName)

	log.Println(g)

	var ret = g.SendText(tips, 500)
	log.Println(ret)
}

func hoTips() {
	self, _ := GetGlobalBot().GetCurrentUser()

	// bootstrap.gBot.Self
	var sendTo = config.GetDaliyGroupName()

	// 获取到group对象
	groups, err := self.Groups()
	if err != nil {
		log.Fatal(err)
	}

	wb_top := utils.Weibo("http://www.anyknew.com/api/v1/sites/weibo")

	var g = groups.SearchByNickName(20, sendTo)
	log.Println(g)

	g.SendText("群友们早上好！这是今天的微博热搜：", 500)
	g.SendText(wb_top, 500)

	type StockChange struct {
		StockPrice float64
		ChatCounts int
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfDay = startOfDay.Add(-24 * time.Hour)

	// // Get the start of the next day
	startOfNextDay := startOfDay.Add(time.Hour * 24)

	row := db.QueryRow("SELECT StockPrice, ChatCounts FROM StockChange WHERE Time >= ? AND Time < ? ORDER BY Time DESC LIMIT 1", startOfDay, startOfNextDay)

	var sc StockChange
	err = row.Scan(&sc.StockPrice, &sc.ChatCounts)
	if err != nil {
		log.Fatal(err)
	}
	g.SendText(fmt.Sprintf("昨日688023股价：%.2f\n昨日本群聊天摸鱼总次数：%d", sc.StockPrice, sc.ChatCounts), 500)

}

func SetDbConnection() {
	var err error
	db, err = sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/wechat?parseTime=true&loc=Asia%2FShanghai")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	//     log.Fatal(err)
	// }

	fmt.Println("Successfully connected to the database")
}

func GetDBObj() *sql.DB {
	return db
}

// 工作日提醒
func doTipsCron(s *gocron.Scheduler, time string, tipStr string) {
	s.Every(1).Monday().At(time).Do(timeTips, tipStr)
	s.Every(1).Tuesday().At(time).Do(timeTips, tipStr)
	s.Every(1).Wednesday().At(time).Do(timeTips, tipStr)
	s.Every(1).Thursday().At(time).Do(timeTips, tipStr)
	s.Every(1).Friday().At(time).Do(timeTips, tipStr)

}

func SetSchedule() {
	s := gocron.NewScheduler()

	// 工作日打卡
	doTipsCron(s, "09:40", "Smile, 该上班打卡了")
	doTipsCron(s, "11:20", "Smile, 该吃饭了，没有什么工作比吃饭更重要！")
	doTipsCron(s, "12:20", "Smile, 该午休了，中间睡两个时才是对自己拼命工作最大的奖励！")
	doTipsCron(s, "17:20", "Smile, 该下班了，工作是老板的，命是自己的！")

	// 每日提醒
	s.Every(1).Day().At("09:00").Do(hoTips)
	s.Every(1).Day().At("18:30").Do(SaveStockInfo)

	s.Start()

	// select {}
}
