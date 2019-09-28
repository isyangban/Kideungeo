package kideungeo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type HelpCommand struct {
	Command
}

func (command *HelpCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle About")
	return nil
}

func (command *HelpCommand) HandleMessageCreate(bot *KideungeoBot, m *discordgo.MessageCreate) error {
	fmt.Println("Handle About")
	// var fields discordgo.MessageEmbedField = []

	reply := &discordgo.MessageEmbed{
		Color: 0x0099ff,
		Title: "도움",
		Description: `키등어는 멍청해서 아이짱이 대신 대답해줄게!
		 지원되는 명령어는 다음과 같아~`,
	}

	for _, supportedCommand := range SupportedCommands {
		reply.Fields = append(reply.Fields, &discordgo.MessageEmbedField{
			Name:   supportedCommand.CommandPrefix(),
			Value:  supportedCommand.Help(),
			Inline: true,
		})
	}

	bot.Session.ChannelMessageSendEmbed(m.ChannelID, reply)
	return nil
}
