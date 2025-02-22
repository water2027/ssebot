package post

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"ssebot/config"
	"ssebot/sseapi"
)

type PostChan struct {
	channel chan sseapi.Post
}

type PostHandler interface {
	Get(ctx context.Context) error
	GetChan() *chan sseapi.Post
}

func NewPostChan(channel chan sseapi.Post) *PostChan {
    return &PostChan{
        channel: channel,
    }
}

func (pc *PostChan) Get(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(config.BotConfig.TimeInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			close(pc.channel)
			return errors.New("stop")
        case <-ticker.C:
            posts := sseapi.GetPosts()
            slices.Reverse(posts)
            for _, post := range posts {
                if post.PostID > config.NowNum && !strings.HasPrefix(post.Title, "test") {
                    pc.channel <- post
                    config.NowNum = post.PostID
                }
            }
		}
	}
}

func (pc *PostChan) GetChan() *chan sseapi.Post {
	return &pc.channel
}
