package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"os"
	"sync"

	"log"
	"strconv"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

var config botConfig              //配置文件
var postId int = 394              //从这个id开始获取新的post
var targetGroup *openwechat.Group //目标群组
var mu sync.Mutex
var intChannel chan int //用于退出程序
var wg sync.WaitGroup

type botConfig struct {
	TargetGroupName string `json:"targetGroupName"`
	TimeInterval    int    `json:"timeInterval"`
	Telephone       string `json:"telephone"`
	Email           string `json:"email"`
	Password        string `json:"password"` //用go没复现出来，只能手动复制了
	Str             string `json:"str"`
}

type Post struct {
	PostID        int    `json:"PostID"`
	UserName      string `json:"UserName"`
	UserScore     int    `json:"UserScore"`
	UserTelephone string `json:"UserTelephone"`
	UserAvatar    string `json:"UserAvatar"`
	UserIdentity  string `json:"UserIdentity"`
	Title         string `json:"Title"`
	Content       string `json:"Content"`
	Like          int    `json:"Like"`
	Comment       int    `json:"Comment"`
	Browse        int    `json:"Browse"`
	Heat          int    `json:"Heat"`
	PostTime      string `json:"PostTime"`
	IsSaved       bool   `json:"IsSaved"`
	IsLiked       bool   `json:"IsLiked"`
	Photos        string `json:"Photos"`
	Tag           string `json:"Tag"`
}

func initLog() {
	logFile, err := os.OpenFile("ssebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
		return
	}
	log.SetOutput(logFile)
}

func getPostId() *int {
	mu.Lock()
	defer mu.Unlock()
	return &postId
}

func GetGroup() *openwechat.Group {
	mu.Lock()
	defer mu.Unlock()
	return targetGroup
}

func consoleQrcode(uuid string) {
	q, err := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(q.ToSmallString(true))
}

func initBot() {
	var err error
	fmt.Println("startBot")
	bot := openwechat.DefaultBot(openwechat.Desktop)
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	bot.UUIDCallback = consoleQrcode
	bot.SyncCheckCallback = nil

	// if err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
	// 	fmt.Println("loginErr", err)
	// 	return
	// }

	if err = bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		log.Println("loginErr", err)
		intChannel <- 1
		return
	}

	// if err = bot.Login(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	bot.MessageHandler = func(msg *openwechat.Message) {
		log.Println(msg, msg.IsSendByGroup())
		if msg.IsSendByGroup() {
			content := msg.Content
			if strings.HasPrefix(content, "@机器人") {
				log.Println(content)
				trimmedMessage := strings.TrimPrefix(content, "@机器人")
				IDRecieved := strings.Split(trimmedMessage, " ")[0]
				cleanedInput := strings.ReplaceAll(IDRecieved, "\u2005", "")
				ID, err := strconv.Atoi(cleanedInput)
				log.Println(ID)
				if err != nil {
					log.Println(err)
					return
				}
				client := &http.Client{}
				loginReq, err := loginSSEReq()
				if err != nil {
					log.Println(err)
					return
				}
				req, err := getPostContent(ID)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = client.Do(loginReq)
				if err != nil {
					log.Println(err)
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					return
				}
				var post Post
				err = json.Unmarshal(body, &post)
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("post", post)
				log.Println("post.Content", post.Content)
				msg.ReplyText(post.Content)
			}

		}
	}

	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Println("getCurrentUserErr", err)
		intChannel <- 1
		return
	}
	groups, err := self.Groups()
	if err != nil {
		log.Println("getGroupsErr", err)
		intChannel <- 1
		return
	}
	targetGroup = groups.GetByNickName(config.TargetGroupName)
	if targetGroup == nil {
		log.Println("groupNotFound")
		intChannel <- 1
		return
	}
	_, err = targetGroup.SendText("hello")
	if err != nil {
		log.Println(err)
		intChannel <- 1
		return
	}
	wg.Done()

	if err = bot.Block(); err != nil {
		log.Println("logout", err)
		intChannel <- 1
		return
	}
}

func loginSSEReq() (*http.Request, error) {
	//login
	loginData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, config.Email, config.Password)
	loginReq, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/login", bytes.NewBuffer([]byte(loginData)))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	loginReq.Header.Set("Content-Type", "application/json")
	return loginReq, nil
}

func getPostsReq() (*http.Request, error) {
	//get posts
	getPostsData := fmt.Sprintf(`{"limit":5,"offset":0,"partition":"主页","searchsort":"home","userTelephone":"%s"}`, config.Telephone)
	req, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/browse", bytes.NewBuffer([]byte(getPostsData)))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func getPostContent(id int) (*http.Request, error) {
	//get post content
	phonenum,err := strconv.Atoi(config.Telephone)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	str := fmt.Sprintf(`{"userTelephone": "%d","postID": %d}`, phonenum, id)
	req, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/showDetails", bytes.NewBuffer([]byte(str)))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	return req, nil
}

func doGetPosts(postChannel chan Post) {
	client := &http.Client{}

	loginReq, err := loginSSEReq()
	if err != nil {
		log.Println(err)
		return
	}
	req, err := getPostsReq()
	if err != nil {
		log.Println(err)
		return
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return
	}
	defer loginResp.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var posts []Post
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(body, &posts)
	//让posts按照PostID升序排列
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PostID < posts[j].PostID
	})
	id := getPostId()
	for _, post := range posts {
		if post.PostID > *id {
			postChannel <- post
			*id = post.PostID
		}
	}
}

func sendPost(postChannel chan Post) {
	str := config.Str
	for post := range postChannel {
		target := GetGroup()
		urlmb := fmt.Sprintf("https://ssemarket.cn/mb/#/postDetails?id=%d", post.PostID)
		msg := fmt.Sprintf(str, post.Title, urlmb)
		log.Println(msg)
		_, err := target.SendText(msg)
		if err != nil {
			log.Println(err)
		}
	}

}

func initConfig() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
}

func main() {
	initConfig()
	initLog()
	wg.Add(1)
	go initBot()
	intChannel = make(chan int)
	postChannel := make(chan Post)
	wg.Wait()
	go sendPost(postChannel)
	ticker := time.NewTicker(time.Duration(config.TimeInterval) * time.Minute)
	for {
		select {
		case <-intChannel:
			return
		case <-ticker.C:
			doGetPosts(postChannel)
		}
	}

}
