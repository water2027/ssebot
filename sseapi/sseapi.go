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
	defer loginResp.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var posts []variable.Post
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
	id := variable.GetPostId()
	for _, post := range posts {
		//如果post.Title以test开头，就不放入postChannel
		if post.PostID > *id&&!strings.HasPrefix(post.Title, "test") {
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

	_, err = client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return variable.Post{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
	client.Do(loginReq)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	var posts []variable.Post
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(body, &posts)
	for _, post := range posts {
		postChannel <- post
	}
}
