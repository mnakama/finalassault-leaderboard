package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/json-iterator/go"
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
}

type pageData struct {
	Players *[]playerData
}

var (
	t_leaderboard *template.Template
	client        = http.Client{}
)

func main() {
	if os.Getenv("PROFILE") != "" {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)

		f, err := os.Create("/tmp/cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()

		go (func() {
			<-c
			log.Println("Caught interrupt")
			pprof.StopCPUProfile()
			f.Close()
			os.Exit(0)
		})()
	}

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
	router.PathPrefix("/res/").Handler(
		http.StripPrefix("/res/",
			http.FileServer(http.Dir("./res"))))
	router.PathPrefix("/style/").Handler(
		http.StripPrefix("/style/",
			http.FileServer(http.Dir("./style"))))

	const path string = "/tmp/finalassault-leaderboard"
	os.Remove(path)

	var listener net.Listener
	var err error
	port := os.Getenv("PORT")
	if port != "" {
		listener, err = net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		listener, err = net.Listen("unix", path)
		if err != nil {
			log.Panic(err)
		}
		defer os.Remove(path)

		if err := os.Chmod(path, 0666); err != nil {
			log.Panic(err)
		}
	}

	log.Fatal(http.Serve(listener, router))
}

func serveLeaderboard(w http.ResponseWriter, r *http.Request) {
	var pData pageData

	start := time.Now()
	err := t_leaderboard.ExecuteTemplate(w, "head", pData)
	if err != nil {
		log.Println(err)
	}
	w.(http.Flusher).Flush()
	head := time.Now()

	rows, err := getLeaderboardData()
	if err != nil {
		log.Println(err)
		return
	}
	fetch := time.Now()

	pData.Players = &rows

	// skip the template; it's too slow. This saves about 120ms out of 250ms on page load time.
	for _, player := range rows {
		safeName := template.HTMLEscapeString(player.DisplayName)

		fmt.Fprintf(w,
			"<tr><td>%d"+
				"<td><a class=anchor id=\"%s\"></a><a href=\"#%s\">%s</a>"+
				"<td>",
			player.Rank,
			safeName,
			safeName,
			safeName,
		)

		lcPlatform := strings.ToLower(player.Platform)
		if lcPlatform == "steam" {
			fmt.Fprintf(w, "<a href=\"%s\"><img width=30 alt=Steam title=Steam src=/res/steam-logo.svg></a>", player.GetSteam())
		} else if lcPlatform == "oculus" {
			fmt.Fprint(w, "<img width=23 alt=Oculus title=Oculus src=/res/oculus-logo.svg>")
		} else if lcPlatform == "vive" {
			fmt.Fprint(w, "<img width=23 alt=Vive title=Vive src=/res/vive-logo.svg>")
		} else if lcPlatform == "playstation" {
			fmt.Fprint(w, "<img width=23 alt=PlayStation title=PlayStation src=/res/playstation-logo.svg>")
		} else {
			fmt.Fprint(w, player.Platform)
		}

		fmt.Fprintf(w, "<td>")
		if player.Info != nil {
			if player.Info.Twitch != "" {
				fmt.Fprintf(w, "<a title=Twitch href=\"https://www.twitch.tv/%s\"><img width=18 alt=Twitch src=/res/twitch-logo.svg></a> ", player.Info.Twitch)
			}

			if player.Info.Youtube != "" {
				fmt.Fprintf(w, "<a title=Youtube href=\"https://www.youtube.com/%s\"><img width=25 alt=Youtube src=/res/youtube-logo.svg></a> ", player.Info.Youtube)
			}
		}

		// uncomment to see playerID in an HTML comment
		//fmt.Fprintf(w, "<!-- %d -->", player.PlayerID)
	}

	end := time.Now()

	log.Printf("fetch: %s render: %s\nTotal: %v\n",
		fetch.Sub(head), end.Sub(fetch), end.Sub(start))
}

// used to create an array of the correct size
var playerCount = 7000

func getLeaderboardData() ([]playerData, error) {
	var lbData leaderboardData
	start := time.Now()

	// fetch json

	resp, err := client.Get("https://phasermm.com/api/dashboards/publicLeaderboard/retail/0")

	fb := time.Now()

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := jsoniter.NewDecoder(resp.Body)

	// parse json

	lbData.Players = make([]playerData, 0, playerCount)

	if err := dec.Decode(&lbData); err != nil {
		return lbData.Players, err
	}

	if lbData.Result != "great success" {
		return lbData.Players, fmt.Errorf("server returned failure: %s", lbData.Result)
	}

	playerCount = len(lbData.Players)

	for i := range lbData.Players {
		lbData.Players[i].DisplayName = nameToUnicode(lbData.Players[i].DisplayName)
		lbData.Players[i].Rank = uint(i) + 1
		info, ok := playerLookup[lbData.Players[i].PlayerID]
		if ok {
			lbData.Players[i].Info = &info
		}
	}

	end := time.Now()
	log.Printf("ttfb: %v load+parse: %v\n", fb.Sub(start), end.Sub(fb))

	return lbData.Players, nil
}

// lookup table for converting the broken data back into its proper utf8 form.
//
// To add to this map, find the character's utf-8 encoding on https://unicode-table.com/en
// and replace the garbled runes with what they should be. Most utf-8 encodings are 3 messed up
// runes long.
//
// Example (in hex): e2 2dc 153 => e2 98 9c
// results in unicode character: ☜

var latinMap = map[rune]byte{
	0x152:  0x8c, // steam name: YPA达瓦里希
	0x153:  0x9c, // steam name: ✠ Balloneta ✠
	0x160:  0xa6,
	0x178:  0x9f, // steam name: Rob888 🔥
	0x17d:  0x8e, // 虎 steam name: Morty 虎龍
	0x17e:  0x9e, // ☞
	0x192:  0x83, // steam name: Larsenik - ラルセン
	0x2c6:  0x88, // steam name: 装甲驱逐舰レ级
	0x2dc:  0x98, // ☜
	0x2013: 0x96, // steam name: 千ㄖㄒ卄卂几 丂山丨千ㄒ
	0x2014: 0x97, // steam name: 时空之门VR-03
	0x2018: 0x91, // steam name: 骑着蜗牛奔宝马
	0x2019: 0x92, // steam name: 千ㄖㄒ卄卂几 丂山丨千ㄒ
	0x201a: 0x82, // steam name: Larsenik - ラルセン
	0x201c: 0x93, // steam name: [鴻哥] [Hong]
	0x201d: 0x94, // steam name: Rob888 🔥
	0x201e: 0x84, // steam name: [0xae] BensEye™
	0x2020: 0x86, // 冬 steam name: ❄冬❄雪❄
	0x2021: 0x87, // steam name: DD-class大凤
	0x2022: 0x95, // ʕ steam name: Martin ʕ´•ᴥ•`ʔ (GER)
	0x2026: 0x85, // steam name lookup on: elacc波兰妹
	0x2030: 0x89, // steam name: 骑着蜗牛奔宝马
	0x2039: 0x8b, // steam name: 时空之门VR-03
	0x203a: 0x9b, // steam name: 骑着蜗牛奔宝马
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
