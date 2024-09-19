package bot

import (
	"fmt"
	"log"
	"sseBot/config"
	"strconv"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"sseBot/sseapi"
	"sseBot/variable"
)

func keepAlive(bot *openwechat.Self) {
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	for range ticker.C {
		heartBeat(bot)
	}
}

func heartBeat(bot *openwechat.Self) {
	// 生成要发送的消息
	outMessage := fmt.Sprintf("防微信自动退出登录[%d]", time.Now().Unix())
	bot.SendTextToFriend(openwechat.NewFriendHelper(bot), outMessage)
}

func consoleQrcode(uuid string) {
	q, err := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(q.ToSmallString(true))
}

func InitBot(config *config.BotConfig, intChannel chan int, postChannel chan variable.Post) {
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
				log.Println(cleanedInput == "热点")
				if cleanedInput == "热点" {
					sseapi.GetHeatPosts(postChannel, config)
					return
				}
				ID, err := strconv.Atoi(cleanedInput)
				log.Println(ID)
				if err != nil {
					log.Println(err)
					return
				}
				post, err := sseapi.GetPostContent(ID, config)
				if err != nil {
					log.Println(err)
					return
				}
				msg.ReplyText(fmt.Sprintf("标题：%s\n内容：%s", post.Title, post.Content))
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
	target := groups.GetByNickName(config.TargetGroupName)
	variable.GroupInit(target)
	targetGroup := variable.GetGroup()
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
	go keepAlive(self)
	go sendPost(postChannel, config)

	if err = bot.Block(); err != nil {
		log.Println("logout", err)
		intChannel <- 1
		return
	}
}

func sendPost(postChannel chan variable.Post, config *config.BotConfig) {
	str := config.Str
	for post := range postChannel {
		target := variable.GetGroup()
		urlmb := fmt.Sprintf("https://ssemarket.cn/mb/#/postDetails?id=%d", post.PostID)
		msg := fmt.Sprintf(str, post.Title, urlmb)
		log.Println(msg)
		_, err := target.SendText(msg)
		if err != nil {
			log.Println(err)
		}

	}
}
