package kideungeo

import "github.com/bwmarrin/discordgo"

// Specifies the commands
type Commander interface {
	Handle(*KideungeoBot) error
	HandleMessageCreate(*KideungeoBot, *discordgo.MessageCreate) error
	Help() string
	CommandPrefix() string
	SetArguments([]string)
}

type Command struct {
	name          string
	commandPrefix string
	help          string
	arguments     []string
}

func (command *Command) Help() string {
	return command.help
}

func (command *Command) CommandPrefix() string {
	return command.commandPrefix
}

func (command *Command) SetArguments(arguments []string) {
	command.arguments = arguments
	return
}

var SupportedCommands = map[string]Commander{
	"날씨": &WeatherCommand{
		Command: Command{name: "Weather", commandPrefix: "날씨", help: "!날씨 [도시 이름]"},
	},
	"야구": &BaseballCommand{
		Command: Command{name: "Baseball", commandPrefix: "야구", help: "!야구 [?팀명]"},
	},
	"콘": &DCConCommand{
		Command: Command{name: "DCCon", commandPrefix: "콘", help: "!콘 [디시콘 이름] [디시콘 종류]"},
	},
	"키등어": &AboutCommand{
		Command: Command{name: "About", commandPrefix: "키등어", help: "!키등어"},
	},
	"도움": &HelpCommand{
		Command: Command{name: "Help", commandPrefix: "도움", help: "!도움"},
	},
}
