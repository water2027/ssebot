package variable

import (
	"sync"
	"github.com/eatmoreapple/openwechat"
)

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

var postId int = 481  //从这个id开始获取新的post
var mu sync.Mutex
var targetGroup1 *openwechat.Group //目标群组
var targetGroup2 *openwechat.Group //目标群组

func GroupInit(group1 *openwechat.Group, group2 *openwechat.Group) {
	mu.Lock()
	defer mu.Unlock()
	targetGroup1 = group1
	targetGroup2 = group2
}

func GetPostId() *int {
	mu.Lock()
	defer mu.Unlock()
	return &postId
}

func GetGroup1() *openwechat.Group {
	mu.Lock()
	defer mu.Unlock()
	return targetGroup1
}

func GetGroup2() *openwechat.Group {
	mu.Lock()
	defer mu.Unlock()
	return targetGroup2
}