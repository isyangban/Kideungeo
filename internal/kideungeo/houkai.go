package kideungeo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/chromedp"
	"github.com/schollz/closestmatch"
)

const (
	HOUKAI_HOME     = "https://houkai3rd.epictools.net/home"
	HOUKAI_VALKYRIE = "https://houkai3rd.epictools.net/valkyrie"
	HOUKAI_BASE     = "https://houkai3rd.epictools.net/"
)

type HoukaiCommand struct {
	Command
}

type ValkyrieSimple struct {
	Name string
	Suit string
}

type StigmaSimple struct {
	Name     string
	ImageURL string
}

type WeaponSimple struct {
	Name     string
	ImageURL string
}

type Equipment struct {
	Weapon       WeaponSimple
	StigmaTop    StigmaSimple
	StigmaMiddle StigmaSimple
	StigmaBottom StigmaSimple
}

type ValkyrieDetail struct {
	Name                  string
	Suit                  string
	Nickname              string
	QTECondition          string
	QTETrigger            string
	ImageURL              string
	PartySkill            string
	LeaderSkill           string
	RecommendedEquipments []Equipment
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
	case "발키리":
		valkyries, _ := getValkyrieList(bot.ChromeContext)
		if len(args) > 0 {
			valkyrie, _ := getValkyrieInfo(bot.ChromeContext, getClosestMatchingValkyre(valkyries, args[0]))
			sendValkyrieInfo(bot, m, valkyrie)
		} else {
			sendValkyrieList(bot, m, valkyries)
		}
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
		chromedp.Navigate(HOUKAI_HOME),
		chromedp.Text(couponPath, &couponNumber),
	)

	if err != nil {
		fmt.Println(err)
		return couponNumber, err
	}
	return couponNumber, nil
}

func getValkyrieList(pCtx context.Context) ([]ValkyrieSimple, error) {
	ctx, cancel := context.WithTimeout(
		pCtx, 15*time.Second,
	)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(HOUKAI_VALKYRIE),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("body", &html, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
	}

	vLists := doc.Find(".row.align-items-start").Children()
	var valkyreName string
	var valkyreList []ValkyrieSimple
	for i := range vLists.Nodes {
		node := vLists.Eq(i)
		if node.HasClass("main-heading") {
			// Name Text has main-heading
			valkyreName = node.Text()
		} else if node.HasClass("ng-star-inserted") {
			// Suit Text has ng-star-insterted
			valkyreSuit := strings.TrimSpace(node.Text())
			valkyreList = append(valkyreList, ValkyrieSimple{
				Name: valkyreName, Suit: valkyreSuit,
			})
		}
	}

	return valkyreList, nil
}

func getValkyrieInfo(pCtx context.Context, suit string) (ValkyrieDetail, error) {
	var valkyre ValkyrieDetail
	ctx, cancel := context.WithTimeout(
		pCtx, 15*time.Second,
	)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	var html string
	url, err := JoinURL(HOUKAI_VALKYRIE, suit)
	if err != nil {
		return valkyre, err
	}
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("body", &html, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println(err)
		return valkyre, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return valkyre, err
	}
	fmt.Println("path")
	fmt.Println(html)

	image := doc.Find("img.full-image")
	imageUrl, ok := image.Attr("src")
	if ok {
		fullImageURL, _ := JoinURL(HOUKAI_BASE, imageUrl)
		valkyre.ImageURL = fullImageURL
	}
	table := doc.Find(".table")
	for i := range table.Nodes {
		fmt.Println("node: ", i)
		node := table.Eq(i)
		row := node.Find("tr")
		row.Each(func(i int, s *goquery.Selection) {
			header := strings.TrimSpace(s.Find("th").Text())
			data := strings.TrimSpace(s.Find("td").Text())
			switch header {
			case "발키리 슈트":
				valkyre.Suit = data
			case "인물":
				valkyre.Name = data
			case "별명":
				valkyre.Nickname = data
			case "QTE 조건":
				valkyre.QTECondition = data
			case "QTE 트리거":
				valkyre.QTETrigger = data
			}
		})
	}

	valkyre.PartySkill = doc.Find("et-valkyrie-detail dl dd:nth-child(2)").Text()
	valkyre.LeaderSkill = doc.Find("et-valkyrie-detail dl dd:nth-child(4)").Text()

	return valkyre, nil
}

func sendValkyrieList(
	bot *KideungeoBot, m *discordgo.MessageCreate, valkyries []ValkyrieSimple,
) error {
	reply := &discordgo.MessageEmbed{
		Color:       0x0099ff,
		Title:       "발키리 검색",
		Description: fmt.Sprintf("붕괴에는 총 %v개의 발키리가 있어!", len(valkyries)),
		URL:         HOUKAI_VALKYRIE,
	}

	valkyrieMap := make(map[string][]string)

	for _, valkyrie := range valkyries {
		valkyrieMap[valkyrie.Name] = append(valkyrieMap[valkyrie.Name], valkyrie.Suit)
	}

	for name, suits := range valkyrieMap {
		reply.Fields = append(reply.Fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  strings.Join(suits, ", "),
			Inline: true,
		})
	}

	bot.Session.ChannelMessageSendEmbed(m.ChannelID, reply)
	// bot.Session.UserChannelCreate()
	return nil
}

func sendValkyrieInfo(
	bot *KideungeoBot, m *discordgo.MessageCreate, valkyre ValkyrieDetail,
) error {
	reply := &discordgo.MessageEmbed{
		Color:       0x0099ff,
		Title:       valkyre.Suit,
		Description: fmt.Sprintf("%v에 대한 정보야", valkyre.Nickname),
		URL:         HOUKAI_VALKYRIE,
		Image:       &discordgo.MessageEmbedImage{URL: valkyre.ImageURL},
	}
	reply.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "발키리 슈트",
			Value:  valkyre.Suit,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "인물",
			Value:  valkyre.Name,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "별명",
			Value:  valkyre.Nickname,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "QTE 조건",
			Value:  valkyre.QTECondition,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "QTE 트리거",
			Value:  valkyre.QTETrigger,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "파티 스킬",
			Value:  valkyre.PartySkill,
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "리터 스킬",
			Value:  valkyre.LeaderSkill,
			Inline: true,
		},
	}
	bot.Session.ChannelMessageSendEmbed(m.ChannelID, reply)
	// bot.Session.UserChannelCreate()
	return nil
}

func getClosestMatchingValkyre(valkyres []ValkyrieSimple, input string) string {
	var wordsToTest []string
	for _, valkyre := range valkyres {
		wordsToTest = append(wordsToTest, valkyre.Suit)
	}
	bagSizes := []int{4}
	cm := closestmatch.New(wordsToTest, bagSizes)
	return cm.Closest(input)
}

func (command *HoukaiCommand) HandleMessageCreate(
	bot *KideungeoBot,
	m *discordgo.MessageCreate,
) error {
	fmt.Println("Houkai Command", *command)
	handleMessageCreateSubCommand(bot, m, command.arguments[0], command.arguments[1:])
	return nil
}

func (command *HoukaiCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle Houkai")
	return nil
}
