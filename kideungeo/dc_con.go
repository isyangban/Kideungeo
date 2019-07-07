package kideungeo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

const (
	DCCON_HOME_URL    = "https://dccon.dcinside.com/"
	DCCON_SEARCH_URL  = "https://dccon.dcinside.com/hot/1/title/"
	DCCON_DETAILS_URL = "https://dccon.dcinside.com/index/package_detail"
	DCCON_IMG_URL     = "http://dcimg5.dcinside.com/dccon.php"
)

type DCConCommand struct {
	Command
}

type DCConPackage struct {
	packageIdx string
	name       string
	seller     string
}

type DCConDetail struct {
	Detail []DCCon
}

type DCCon struct {
	Ext        string `json:"ext"`
	Idx        string `json:"idx"`
	PackageIdx string `json:"package_idx"`
	Path       string `json:"path"`
	Sort       string `json:"sort"`
	Title      string `json:"title"`
}

func (command *DCConCommand) Handle(bot *KideungeoBot) error {
	fmt.Println("Handle About")
	return nil
}

func (command *DCConCommand) HandleMessageCreate(
	bot *KideungeoBot, m *discordgo.MessageCreate,
) error {
	fmt.Println("DCCon Command", *command)
	args := command.arguments
	if len(args) == 1 { // !콘 붕괴
		fmt.Println("args 1")
		dcCons := searchDCCons(args[0])
		sendDCCons(bot, m, dcCons, args[0])
	} else if len(args) == 2 { // !콘 붕괴 Blah
		fmt.Println("args 2")
		dcCons := searchDCCon(args[0])
		if num, err := strconv.Atoi(args[1]); err == nil {
			sendDCCon(bot, m, dcCons, "", num)
		} else {
			sendDCCon(bot, m, dcCons, args[1], -1)
		}
	} else {
		// Show help
	}
	return nil
}

func sendDCCons(
	bot *KideungeoBot, m *discordgo.MessageCreate, dcCons []DCConPackage, packageName string,
) error {
	dcConsCountLimit := 10
	if len(dcCons) < dcConsCountLimit {
		dcConsCountLimit = len(dcCons)
	}

	reply := &discordgo.MessageEmbed{
		Color:       0x0099ff,
		Title:       "디시콘 검색",
		Description: fmt.Sprintf("오옷 내 차례인건가? 결과는 %v개", dcConsCountLimit),
		URL:         searchDCConsURL(packageName),
	}

	for _, dcCon := range dcCons[:dcConsCountLimit] {
		reply.Fields = append(reply.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("(%v): ", dcCon.packageIdx) + dcCon.name,
			Value:  dcCon.seller,
			Inline: true,
		})
	}

	bot.Session.ChannelMessageSendEmbed(m.ChannelID, reply)
	// bot.Session.UserChannelCreate()
	return nil
}

func searchDCConsURL(packageName string) string {
	searchUrl, _ := url.Parse(DCCON_SEARCH_URL)
	searchUrl.Path = path.Join(searchUrl.Path, packageName)
	return searchUrl.String()
}

func searchDCCons(packageName string) []DCConPackage {
	searchUrl, _ := url.Parse(DCCON_SEARCH_URL)
	searchUrl.Path = path.Join(searchUrl.Path, packageName)
	res, err := http.Get(searchUrl.String())
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var dcCons []DCConPackage = make([]DCConPackage, 0)
	listItems := doc.Find("ul.dccon_shop_list > li")
	listItems.Each(func(i int, s *goquery.Selection) {
		packageIdx, _ := s.Find("a.link_product").Attr("href")
		name := s.Find("strong.dcon_name").Text()
		seller := s.Find("span.dcon_seller").Text()
		dcCons = append(dcCons, DCConPackage{
			packageIdx, name, seller,
		})
	})
	return dcCons
}

func searchDCCon(packageName string) []DCCon {
	fmt.Println("searching dc con", packageName)
	var dcCons []DCCon = make([]DCCon, 0)
	detail := DCConDetail{Detail: dcCons}
	var packageIdx string
	if strings.HasPrefix(packageName, "#") {
		packageIdx = packageName[1:]
	} else {
		packages := searchDCCons(packageName)
		if len(packages) == 0 {
			fmt.Println("No result")
			return dcCons
		} else {
			packageIdx = packages[0].packageIdx[1:]
		}
	}

	data := url.Values{}
	data.Set("package_idx", packageIdx)
	data.Set("code", "")
	req, _ := http.NewRequest("POST", DCCON_DETAILS_URL, strings.NewReader(data.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://dccon.dcinside.com")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &detail)
	return detail.Detail
}

func getDCConImageURL(path string) string {
	return DCCON_IMG_URL + "?no=" + path
}

func sendDCCon(
	bot *KideungeoBot, m *discordgo.MessageCreate,
	dcCons []DCCon, title string, index int,
) {
	if len(dcCons) == 0 {
		return
	}
	sort.SliceStable(dcCons, func(i, j int) bool {
		in, _ := strconv.Atoi(dcCons[i].Sort)
		jn, _ := strconv.Atoi(dcCons[j].Sort)
		return in < jn
	})
	var dcCon DCCon
	if len(title) != 0 {
		for _, v := range dcCons {
			if v.Title == title {
				dcCon = v
				break
			}
		}
	} else if index > 0 && index < len(dcCons) {
		dcCon = dcCons[index]
	} else {
		dcCon = dcCons[0]
	}
	// Get DCONImageURL
	imageURL := getDCConImageURL(dcCon.Path)
	fmt.Println("image url ", imageURL)
	req, _ := http.NewRequest("GET", imageURL, nil)

	client := &http.Client{}
	req.Header.Set("Referer", "https://dccon.dcinside.com")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Code: ", res.Status)
	defer res.Body.Close()
	// Send file
	bot.Session.ChannelFileSend(m.ChannelID, dcCon.Title+"."+dcCon.Ext, res.Body)
	return
}
