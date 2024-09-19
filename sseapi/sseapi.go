package sseapi

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"io"
	"sseBot/config"
	"sseBot/utils"
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

func GetPosts(postChannel chan Post, config *config.BotConfig) {
	client := &http.Client{}
	loginReq, err := utils.LoginSSEReq(config)
	if err != nil {
		log.Println(err)
		return
	}
	req, err := utils.GetPostsReq(config)
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
	for _, post := range posts {
		postChannel <- post
	}
}

func GetPostContent(id int, config *config.BotConfig) (Post, error) {
	client := &http.Client{}
	loginReq, err := utils.LoginSSEReq(config)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}
	req, err := utils.GetPostContentReq(id, config)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}

	_, err = client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}
	var post Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}
	return post, nil
}
