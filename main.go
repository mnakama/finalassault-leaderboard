package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type playerData struct {
	Rank        uint
	DisplayName string `json:"displayName"`
	PlayerID    uint   `json:"playerid"`
	Platform    string `json:"platform"`
	Info        *playerInfo
}

type leaderboardData struct {
	Result  string       `json:"result"`
	Players []playerData `json:"players"`
	ShowAll bool
}

type pageData struct {
	Players *[]playerData
	ShowAll bool
}

var t_leaderboard *template.Template

func main() {
	t_leaderboard = template.Must(template.New("leaderboard.html").
		ParseFiles("leaderboard.html"))

	/*lbData, err := getLeaderboardData()
	if err != nil {
		log.Panic(err)
	}*/

	//printLeaderboardTable(&lbData)
	//printLeaderboardHTML(&lbData)

	router := mux.NewRouter()
	router.HandleFunc("/", serveLeaderboard)
	router.HandleFunc("/leaderboard", serveLeaderboard)
	router.HandleFunc("/finalassault/leaderboard", serveLeaderboard)

	http.ListenAndServe(":8000", router)
}

func serveLeaderboard(w http.ResponseWriter, r *http.Request) {
	var pData pageData

	err := t_leaderboard.ExecuteTemplate(w, "head", pData)
	if err != nil {
		log.Println(err)
	}
	w.(http.Flusher).Flush()

	data, err := getLeaderboardData()
	if err != nil {
		log.Println(err)
		return
	}

	rows, err := parseLeaderboardData(data)
	if err != nil {
		log.Println(err)
		return
	}

	pData.Players = &rows

	err = t_leaderboard.ExecuteTemplate(w, "body", pData)
	if err != nil {
		log.Println(err)
	}
}

func getLeaderboardData() ([]byte, error) {
	resp, err := http.Get("https://phasermm.com/api/dashboards/publicLeaderboard/retail/0")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//data, err := ioutil.ReadFile("/tmp/leaderboard.json")
	return body, err
}

func parseLeaderboardData(data []byte) ([]playerData, error) {
	var lbData leaderboardData
	var cooked []playerData

	if err := json.Unmarshal(data, &lbData); err != nil {
		return cooked, err
	}

	if lbData.Result != "great success" {
		return cooked, fmt.Errorf("server returned failure: %s", lbData.Result)
	}

	for i := range lbData.Players {
		lbData.Players[i].DisplayName = nameToUnicode(lbData.Players[i].DisplayName)
		lbData.Players[i].Rank = uint(i) + 1
		info, ok := playerLookup[lbData.Players[i].PlayerID]
		if ok {
			lbData.Players[i].Info = &info
		}
	}

	return lbData.Players, nil
}

// lookup table for converting the broken data back into its proper utf8 form
var latinMap = map[rune]byte{
	0x152:  0x8c, // steam name: YPA达瓦里希
	0x153:  0x9c, // steam name: ✠ Balloneta ✠
	0x160:  0xa6,
	0x178:  0x9f, // steam name: Rob888 🔥
	0x192:  0x83, // steam name: Larsenik - ラルセン
	0x2c6:  0x88, // steam name: 装甲驱逐舰レ级
	0x2013: 0x96, // steam name: 千ㄖㄒ卄卂几 丂山丨千ㄒ
	0x2014: 0x97, // steam name: 时空之门VR-03
	0x2018: 0x91, // steam name: 骑着蜗牛奔宝马
	0x2019: 0x92, // steam name: 千ㄖㄒ卄卂几 丂山丨千ㄒ
	0x201a: 0x82, // steam name: Larsenik - ラルセン
	0x201c: 0x93, // steam name: [鴻哥] [Hong]
	0x201d: 0x94, // steam name: Rob888 🔥
	0x201e: 0x84, // steam name: [0xae] BensEye™
	0x2021: 0x87, // steam name: DD-class大凤
	0x2026: 0x85, // steam name lookup on: elacc波兰妹
	0x2030: 0x89, // steam name: 骑着蜗牛奔宝马
	0x203a: 0x9b, // steam name: 骑着蜗牛奔宝马
	0x2039: 0x8b, // steam name: 时空之门VR-03
	0x20ac: 0x80, // steam name: 装甲驱逐舰レ级
	0x2122: 0x99, // steam name (guessed): 不慌 不忙
}

func nameToUnicode(name string) string {
	// Final Assault's backend reads the database in latin-1 mode, but the
	// data itself is stored in unicode. Here, we convert the unicode runes
	// back into bytes, then return the bytes as a unicode (utf-8) string.
	//var foundSpecial bool

	runes := []rune(name)
	out := make([]byte, len(runes))
	for i, rune := range runes {
		if rune < 0x100 {
			// fits in one byte; convert it back
			out[i] = byte(rune)
		} else if b, ok := latinMap[rune]; ok {
			// multi-byte rune. Use lookup table to return it to its original value.
			out[i] = b
		} else {
			// not in the table? write a 0 byte. This makes an invalid utf-8 character,
			// but it's better than showing the corrupted version
			//foundSpecial = true
			out[i] = 0
		}

	}

	/*if foundSpecial {
		var runeStr string
		for _, rune := range runes {
			runeStr += fmt.Sprintf("%x ", rune)
		}
		log.Printf("Runes: %s\nin: %s out: %s\n", runeStr, name, out)
	}*/

	return string(out)
}

func (p *playerData) GetSteam() string {
	if p.Info != nil && p.Info.Steam != "" {
		return fmt.Sprintf("https://steamcommunity.com/%s", p.Info.Steam)
	}

	return fmt.Sprintf("https://steamcommunity.com/search/users/#text=%s", p.DisplayName)
}
