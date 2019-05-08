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
}

type leaderboardData struct {
	Result  string       `json:"result"`
	Players []playerData `json:"players"`
}

var t_leaderboard *template.Template

func add1(n int) int {
	return n + 1
}

var funcMap = template.FuncMap{
	"add1": add1,
}

func main() {
	t_leaderboard = template.Must(template.New("leaderboard.html").
		Funcs(funcMap).
		ParseFiles("leaderboard.html"))

	/*lbData, err := getLeaderboardData()
	if err != nil {
		log.Panic(err)
	}*/

	//printLeaderboardTable(&lbData)
	//printLeaderboardHTML(&lbData)

	router := mux.NewRouter()
	router.HandleFunc("/", serveLeaderboard)
	router.HandleFunc("/finalassault/leaderboard", serveLeaderboard)

	http.ListenAndServe(":8000", router)
}

func serveLeaderboard(w http.ResponseWriter, r *http.Request) {
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

	err = t_leaderboard.Execute(w, rows)
	if err != nil {
		log.Println(err)
	}
}

func getLeaderboardData() ([]byte, error) {
	data, err := ioutil.ReadFile("/tmp/leaderboard.json")
	return data, err
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

	cooked = make([]playerData, 0, len(lbData.Players))
	var nullCount uint
	var zmode int
	for i, row := range lbData.Players {
		row.Rank = uint(i) + 1

		if row.DisplayName == "" {
			nullCount++
			if nullCount >= 5 && zmode >= 0 {
				zmode = 1
			}
			continue
		}

		nullCount = 0

		if zmode == 1 {
			if row.DisplayName[0] == 'Z' || row.DisplayName[0] == 'z' {
				zmode = 2
			}
			continue
		}

		if zmode == 2 {
			if row.DisplayName[0] != 'Z' && row.DisplayName[0] != 'z' {
				zmode = -1
			} else {
				continue
			}
		}
		cooked = append(cooked, row)
	}

	return cooked, nil
}
