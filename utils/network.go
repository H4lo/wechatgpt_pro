package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"log"
	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Iid     int    `json:"iid"`
	Title   string `json:"title"`
	AddDate int64  `json:"add_date"`
	More 	string `json:"more"`
}

type Sub struct {
	Items []Item `json:"items"`
}

type Site struct {
	Subs []Sub `json:"subs"`
}

type Data struct {
	Site Site `json:"site"`
}

type Response struct {
	Data Data `json:"data"`
}

type ResponseData struct {
	Data struct {
		Klines []string `json:"klines"`
	} `json:"data"`
}

func Weibo(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	count := 1
	res := ""
	for _, sub := range response.Data.Site.Subs {
		for _, item := range sub.Items {
			if(count > 10){
				break
			}
			t := time.Unix(item.AddDate, 0)
			res += fmt.Sprintf("%d. %s [%s]   热度: %s \nhttps://www.anyknew.com/go/%d\n\n", count, item.Title, t.Format("2006-01-02 15:04:05"), item.More, item.Iid)
			count++
		}
	}

	return res
}

func GetMpContentByUrl(url string) string{

	resp, err := http.Get(url)
	if err != nil {
		// log.Fatal(err)
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// log.Fatal(err)
		return ""
	}

	// 微信公众号链接的摘要方法
	if(strings.Contains(url, "mp.weixin")){
		content := doc.Find(".rich_media_content.js_underline_content").Text()
		// fmt.Println(content)
		return content
	}else{
		// 遍历html中所有p标签中的内容，但是如果是动态页面就无法获取到
		
		var texts []string
		doc.Find("p").Each(func(i int, s *goquery.Selection){
			texts = append(texts, s.Text())
		})

		allTexts := strings.Join(texts, " ")
		// fmt.Println(allTexts)
		return allTexts
	}
}

func GetStock() string{
	resp, err := http.Get("https://push2his.eastmoney.com/api/qt/stock/kline/get?cb=al&secid=1.688023&fields1=f1&fields2=f51,f53&klt=101&fqt=1&end=20500101&lmt=1")
	if err != nil {
		fmt.Printf("http.Get -> err : %v\n", err)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll -> err : %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	// Remove the 'al(' at the beginning and the ');' at the end
	bodyStr := string(body)
	jsonStr := bodyStr[3 : len(bodyStr)-2]

	var data ResponseData
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Printf("json.Unmarshal -> err : %v\n", err)
		return ""
	}

	// Split the kline string on the comma and get the second part
	klineParts := strings.Split(data.Data.Klines[0], ",")
	numberStr := klineParts[1]

	// fmt.Println("The number is:", numberStr)

	return numberStr
}