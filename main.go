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

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

var config botConfig
var postId int = 0
var targetGroup *openwechat.Group
var mu sync.Mutex

type botConfig struct {
	TargetGroupName string `json:"targetGroupName"`
	TimeInterval    int    `json:"timeInterval"`
	Telephone       string `json:"telephone"`
	Email           string `json:"email"`
	Password        string `json:"password"` //用go没复现出来，只能手动复制了
	Str 		   string `json:"str"`
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
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	fmt.Println("startBot")
	bot := openwechat.DefaultBot(openwechat.Desktop)
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	bot.UUIDCallback = consoleQrcode

	if err = bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println("loginErr", err)
		return
	}

	// if err = bot.Login(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() {
			fmt.Println(msg.Content)
			msg.ReplyText("你好")
		}
	}

	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println("getCurrentUserErr", err)
		return
	}
	groups, err := self.Groups()
	if err != nil {
		fmt.Println("getGroupsErr", err)
		return
	}
	targetGroup = groups.GetByNickName(config.TargetGroupName)
	if targetGroup == nil {
		fmt.Println("groupNotFound")
		return
	}
	_, err = targetGroup.SendText("hello")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = bot.Block(); err != nil {
		fmt.Println("logout", err)
		return
	}
}

func getPosts() {
	str := config.Str
	//login
	loginData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, config.Email, config.Password)
	loginReq, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/login", bytes.NewBuffer([]byte(loginData)))
	if err != nil {
		fmt.Println(err)
		return
	}
	loginReq.Header.Set("Content-Type", "application/json")

	//get posts
	getPostsData := fmt.Sprintf(`{"limit":5,"offset":0,"partition":"主页","searchsort":"home","userTelephone":"%s"}`, config.Telephone)
	req, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/browse", bytes.NewBuffer([]byte(getPostsData)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	//login
	loginResp, err := client.Do(loginReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer loginResp.Body.Close()
	//get posts
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var posts []Post
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.Unmarshal(body, &posts)
	//让posts按照PostID升序排列
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PostID < posts[j].PostID
	})

	for _, post := range posts {
		fmt.Println(post.PostID)
		fmt.Println(post.UserName)
		fmt.Println(post.Title)
		fmt.Println(post.Tag)
		id := getPostId()
		if post.PostID > *id {
			*id = post.PostID
			go func(post Post) {
				target := GetGroup()
				urlpc := "https://ssemarket.cn/pc/postDetails?id=" + fmt.Sprint(post.PostID)
				urlmb := "https://ssemarket.cn/mb/postDetails?id=" + fmt.Sprint(post.PostID)
				msg := fmt.Sprintf(str, post.UserName, post.Title, post.Tag, urlpc, urlmb)
				_, err := target.SendText(msg)
				if err != nil {
					fmt.Println(err)
				}
			}(post)
		}
	}
}

func main() {
	go initBot()
	fmt.Println("请在30秒内扫码登录")
	time.Sleep(30 * time.Second)
	ticker := time.NewTicker(time.Duration(config.TimeInterval) * time.Minute)
	for range ticker.C {
		getPosts()
	}
}
