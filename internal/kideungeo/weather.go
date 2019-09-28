package kideungeo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type WeatherCommand struct {
	Command
}

func (command *WeatherCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle About")
	return nil
}

func (command *WeatherCommand) HandleMessageCreate(bot *KideungeoBot, m *discordgo.MessageCreate) error {
	fmt.Println("Handle About")
	return nil
}
