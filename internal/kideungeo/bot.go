package kideungeo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type KideungeoBot struct {
	SelfName   string
	AppID      uint64
	AppName    string
	Token      string
	Session    *discordgo.Session
	BotVersion string
	BotAuthor  string
}

const BotVersion = "0.0.1"
const BotAuthor = "choijeongwoon"

func New(token string) *KideungeoBot {
	bot := &KideungeoBot{
		Token:      token,
		SelfName:   "Kideungeo",
		BotVersion: BotVersion,
		BotAuthor:  BotAuthor,
	}
	return bot
}

func (kb *KideungeoBot) Start() error {
	dg, err := discordgo.New("Bot " + kb.Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return nil
	}
	kb.Session = dg
	kb.Session.AddHandler(newMessageCreateHandler(kb))
	dg.Open()

	return nil
}

func (kb *KideungeoBot) Close() {
	kb.Session.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func newMessageCreateHandler(bot *KideungeoBot) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Parse Content => Command or Discussion
		// Coammnd Syntax: ![키등어] commandName commandArg1 commandArg2 commandArg3
		// If it is a relevent discussion -> Do a reaction
		parsedContent := parseContent(m.Content)
		switch v := parsedContent.(type) {
		case Commander:
			v.HandleMessageCreate(bot, m)
		}
	}
}

// Parses Discord Message Conent and Returns Parsed Command or Discussion
func parseContent(content string) interface{} {
	// Start with ![Alphanumberal] -> Its a command
	r, _ := regexp.Compile("^![a-zA-Z0-9가-힣]*")
	command := r.FindString(content)
	if len(command) != 0 {
		commandPrefix := string([]rune(command)[1:])
		commandRest := content[len(command):]
		command, ok := SupportedCommands[commandPrefix]
		if ok {
			// Parse arguments and add arguments
			argsR := regexp.MustCompile("[^\\s\"']+|\"([^\"]*)\"|'([^']*)'")
			args := argsR.FindAllString(commandRest, len(strings.Fields(commandRest)))
			for i, arg := range args {
				args[i] = strings.Trim(arg, "'\"")
			}
			command.SetArguments(args)
			return command
		} else {
			return content
		}
	} else {
		return content
	}
}
