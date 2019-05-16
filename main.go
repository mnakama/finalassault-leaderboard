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
		lbData.Players[i].Rank = uint(i) + 1
	}
	return lbData.Players, nil
}
