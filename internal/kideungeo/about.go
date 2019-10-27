package kideungeo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// AboutCommand is a command for about
type AboutCommand struct {
	Command
}

// Handle is a generic callback
func (command *AboutCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle About")
	return nil
}

// HandleMessageCreate is used when a message is created
func (command *AboutCommand) HandleMessageCreate(bot *KideungeoBot, m *discordgo.MessageCreate) error {
	bot.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("이름: 키등어봇, 버전: %v", bot.BotVersion))
	return nil
}
