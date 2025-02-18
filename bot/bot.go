package bot

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"ssebot/post"
)

type bot struct {
	webhook      string
	postHandler  post.PostHandler
	templateStr string
}

type BotHandler interface {
	ReceiveMessage(ctx context.Context) error
	SendMessage(resp string) error
	Run() error
}

func NewBot(webhook string, postHandler post.PostHandler, templateStr string) *bot {
	return &bot{
		webhook: webhook,
		postHandler: postHandler,
		templateStr: templateStr,
	}
}

func (b *bot) ReceiveMessage(ctx context.Context) error {
	if err := b.postHandler.Get(ctx); err != nil {
		return err
	}
	return nil
}

func (b *bot) SendMessage(resp string) error {
	data := fmt.Sprintf(`{"msgtype":"markdown","markdown":{"content":"%s"}}`, resp)
	req, err := http.NewRequest("POST", b.webhook, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (b *bot) Run() error {
	ctx := context.Background()
	go b.ReceiveMessage(ctx)
	postChan := b.postHandler.GetChan()
    for post := range *postChan {
        msg := fmt.Sprintf(b.templateStr, post.Title, post.PostID,)
		log.Println(msg)
        err := b.SendMessage(msg)
		if err != nil {
			return err
		}
    }
	return nil
}
