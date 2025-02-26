package reporter

import (
	"fmt"
	"github.com/go-lark/lark"
)

type Provider string

const (
	FeiShu Provider = "feishu"
)

func Report(provider Provider, url, report string) error {
	switch provider {
	case FeiShu:
		return sendFeiShu(url, report)
	default:
		return fmt.Errorf("provider %s is not supported", provider)
	}
}

func sendFeiShu(url string, report string) (err error) {
	bot := lark.NewNotificationBot(url)
	_, err = bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(report).Build())
	if err != nil {
		err = fmt.Errorf("send feishu meesage error: %s", err)
		return
	}
	return
}
