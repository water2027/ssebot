package sseapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"sseBot/config"
	"sseBot/utils"
	"sseBot/variable"
	"strings"
)

type loginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func GetPosts(postChannel chan variable.Post, config *config.BotConfig) {
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

	var loginResponse loginResponse

	body, _ := io.ReadAll(loginResp.Body)
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		log.Println(err)
		return
	}

	if loginResponse.Code == 200 {
		// 将token添加到第二个请求的header中
		req.Header.Add("Authorization", "Bearer "+loginResponse.Data.Token)
	}

	defer loginResp.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var posts []variable.Post
	body, err = io.ReadAll(resp.Body)
	log.Println(string(body))
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(body, &posts)
	//让posts按照PostID升序排列
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PostID < posts[j].PostID
	})
	id := variable.GetPostId()
	for _, post := range posts {
		//如果post.Title以test开头，就不放入postChannel
		if post.PostID > *id && !strings.HasPrefix(post.Title, "test") {
			postChannel <- post
			*id = post.PostID
		}
	}
}

func GetPostContent(id int, config *config.BotConfig) (variable.Post, error) {
	client := &http.Client{}
	loginReq, err := utils.LoginSSEReq(config)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}
	req, err := utils.GetPostContentReq(id, config)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}

	var loginResponse loginResponse

	body, _ := io.ReadAll(loginResp.Body)
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}

	if loginResponse.Code == 200 {
		// 将token添加到第二个请求的header中
		req.Header.Add("Authorization", "Bearer "+loginResponse.Data.Token)
	}

	defer loginResp.Body.Close()


	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}
	var post variable.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}
	return post, nil
}

func GetHeatPosts(postChannel chan variable.Post, config *config.BotConfig) {
	client := &http.Client{}
	loginReq, err := utils.LoginSSEReq(config)
	if err != nil {
		log.Println(err)
		return
	}
	req, err := utils.GetHeatPostsReq(config)
	if err != nil {
		log.Println(err)
		return
	}
	loginResp,_ :=client.Do(loginReq)
	var loginResponse loginResponse

	body, _ := io.ReadAll(loginResp.Body)
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		log.Println(err)
		return
	}

	if loginResponse.Code == 200 {
		// 将token添加到第二个请求的header中
		req.Header.Add("Authorization", "Bearer "+loginResponse.Data.Token)
	}

	defer loginResp.Body.Close()
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	var posts []variable.Post
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(body, &posts)
	for _, post := range posts {
		postChannel <- post
	}
}
