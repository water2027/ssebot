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

var postId int = 425  //从这个id开始获取新的post
var mu sync.Mutex
var targetGroup *openwechat.Group //目标群组

func GroupInit(group *openwechat.Group) {
	mu.Lock()
	defer mu.Unlock()
	targetGroup = group
}

func GetPostId() *int {
	mu.Lock()
	defer mu.Unlock()
	return &postId
}

func GetGroup() *openwechat.Group {
	mu.Lock()
	defer mu.Unlock()
	return targetGroup
}