package kideungeo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type AboutCommand struct {
	Command
}

func (command *AboutCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle About")
	return nil
}

func (command *AboutCommand) HandleMessageCreate(bot *KideungeoBot, m *discordgo.MessageCreate) error {
	bot.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("이름: 키등어봇, 버전: %v", bot.BotVersion))
	return nil
}
