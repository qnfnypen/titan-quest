package telegram

import (
	"context"
	"fmt"
	"github.com/gnasnik/titan-quest/core/dao"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

var (
	telegramMembersKey = "gm::telegram::members"
)

func RunTelegramBot(token string, groupId int64) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	ctx := context.Background()

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("create telegram cmd: %v", err)
		return
	}

	fmt.Println("telegram cmd id:", b.Me.ID)

	b.Handle(tele.OnUserJoined, func(c tele.Context) error {
		if c.Chat().ID != groupId {
			return nil
		}

		fmt.Println(c.Message().OriginalSender, c.Message().Sender)

		log.Printf("User %d join", c.Sender().ID)

		_, err := dao.RedisCache.SAdd(ctx, telegramMembersKey, c.Sender().ID).Result()
		if err != nil {
			log.Printf("redis sadd: %s %v\n", telegramMembersKey, err)
			return err
		}

		return nil
	})

	b.Handle(tele.OnUserLeft, func(c tele.Context) error {
		if c.Chat().ID != groupId {
			return nil
		}

		fmt.Println(c.Message().OriginalSender, c.Message().Sender)

		log.Printf("User %d left", c.Sender().ID)

		_, err := dao.RedisCache.SRem(ctx, telegramMembersKey, c.Sender().ID).Result()
		if err != nil {
			log.Printf("redis srem: %s %v\n", telegramMembersKey, err)
			return err
		}

		return nil
	})

	b.Start()
}
