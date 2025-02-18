package sseapi

import (
	"net/http"
	"io"
	"log"
	"encoding/json"

	"ssebot/config"
)

type loginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Post struct {
	PostID        int    `json:"PostID"`
	UserID        int    `json:"UserID"`
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

func GetPosts() []Post {
	client := &http.Client{}
	loginReq, err := loginSSEReq(&config.BotConfig)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	req, err := getPostsReq(&config.BotConfig)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	var loginResponse loginResponse

	body, _ := io.ReadAll(loginResp.Body)
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	if loginResponse.Code == 200 {
		// 将token添加到第二个请求的header中
		req.Header.Add("Authorization", "Bearer "+loginResponse.Data.Token)
	}

	defer loginResp.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	defer resp.Body.Close()

	var posts []Post
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	json.Unmarshal(body, &posts)
	return posts
}