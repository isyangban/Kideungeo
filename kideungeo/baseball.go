package kideungeo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type BaseballCommand struct {
	Command
}

func (command *BaseballCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle About")
	return nil
}

func (command *BaseballCommand) HandleMessageCreate(bot *KideungeoBot, m *discordgo.MessageCreate) error {
	fmt.Println("Handle About")
	return nil
}
