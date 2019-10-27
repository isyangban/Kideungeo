package kideungeo

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/chromedp"
)

const (
	HOUKAI_URL = "https://houkai3rd.epictools.net/home"
)

type HoukaiCommand struct {
	Command
}

// SubCommand: name, args
func handleMessageCreateSubCommand(
	bot *KideungeoBot, m *discordgo.MessageCreate, name string, args []string,
) {
	fmt.Println("handle sub command")
	switch name {
	case "쿠폰":
		couponNumber, _ := getCouponNumber(bot.ChromeContext)
		bot.Session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("이번주 쿠폰 번호야: %v", couponNumber))
	}
}

func getCouponNumber(pCtx context.Context) (string, error) {
	// Force Timeout
	ctx, cancel := context.WithTimeout(
		pCtx, 15*time.Second,
	)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	couponPath := "et-weekly-info > dl > dd:nth-child(2)"
	var couponNumber string
	err := chromedp.Run(ctx,
		chromedp.Navigate(HOUKAI_URL),
		chromedp.Text(couponPath, &couponNumber),
	)

	if err != nil {
		fmt.Println(err)
		return couponNumber, err
	}
	return couponNumber, nil
}

func (command *HoukaiCommand) HanndleMessageCreate(
	bot *KideungeoBot,
	m *discordgo.MessageCreate,
) error {
	fmt.Println("Houkai Command", *command)
	handleMessageCreateSubCommand(bot, m, command.arguments[0], command.arguments[1:])
	return nil
}
