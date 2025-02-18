package sseapi

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"ssebot/config"
)

func loginSSEReq(config *config.Config) (*http.Request, error) {
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

func getPostsReq(config *config.Config) (*http.Request, error) {
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
