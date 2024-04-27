package logger

import (
	"context"
	"github.com/filecoin-project/pubsub"
	"github.com/gnasnik/titan-quest/core/dao"
	model2 "github.com/gnasnik/titan-quest/core/generated/model"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("logger")

const (
	loggerLoginTopic     string = "login"
	loggerOperationTopic        = "operation"
)

var g *logger

func init() {
	g = &logger{
		logger: pubsub.New(50),
	}
}

type logger struct {
	logger *pubsub.PubSub
}

func (l *logger) pub(topic string, v interface{}) {
	l.logger.Pub(v, topic)
}

func (l *logger) sub(ctx context.Context) {
	login := l.logger.Sub(loggerLoginTopic)
	operator := l.logger.Sub(loggerOperationTopic)
	go func() {
		defer l.logger.Unsub(login)
		defer l.logger.Unsub(operator)

		for {
			select {
			case msg := <-login:
				err := dao.AddLoginLog(ctx, msg.(*model2.LoginLog))
				if err != nil {
					log.Errorf("add login log: %v", err)
				}
			case msg := <-operator:
				err := dao.AddOperationLog(ctx, msg.(*model2.OperationLog))
				if err != nil {
					log.Errorf("add operation log: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func AddLoginLog(v interface{}) {
	g.pub(loggerLoginTopic, v)
}

func AddOperationLog(v interface{}) {
	g.pub(loggerOperationTopic, v)
}

func Subscribe(ctx context.Context) {
	g.sub(ctx)
}
