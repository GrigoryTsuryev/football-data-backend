package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type Data struct {
	Competitions []Competitions `json:"competitions"`
}

type Competitions struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Area Area   `json:"area"`
}

type Area struct {
	Country string `json:"name"`
}

type StandingsData struct {
	Standings []Standings `json:"standings"`
}

type Standings struct {
	Type  string  `json:"type"`
	Table []Table `json:"table"`
}

type Table struct {
	Position int  `json:"position"`
	Team     Team `json:"team"`
	Won      int  `json:"won"`
	Draw     int  `json:"draw"`
	Lost     int  `json:"lost"`
	Points   int  `json:"points"`
}

type Team struct {
	Name string `json:"name"`
}

type Error struct {
	Massage   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

var url = "http://api.football-data.org/v2/"

// var API_KEY = "dee538714c524f9d9deaa2b6202d20a7"

var API_KEY = "94cd3029a48845e895049063a19072bc"

// var API_KEY = "e39ebec8d17146238570c00bad724c99"

var allowedProviders = []int{2000, 2001, 2002, 2003, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2021}
var clientR = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func main() {

	pong, err := clientR.Ping().Result()
	fmt.Println(pong, err)

	r := mux.NewRouter()
	r.HandleFunc("/api/competitions", fetchAllCompetitions).Methods("GET")
	r.HandleFunc("/api/competitions/{id}/standings/", fetchStandingsById).Methods("GET")
	http.ListenAndServe(":8000", r)
}

func fetchAllCompetitions(w http.ResponseWriter, r *http.Request) {

	var competitions = "competitions"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	val, err := clientR.Get(competitions).Result()
	if err == redis.Nil {
		fmt.Println("Fetching from the web")
	} else if err != nil {
		panic(err)
	}
	// var val string
	var data Data
	var filteredData Data
	if len(val) == 0 {

		resp, err := http.Get(url + "competitions")
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(body, &data)
		for i := range data.Competitions {
			var el = data.Competitions[i].Id
			if allowedFreeAPIProvider(el, allowedProviders) {
				filteredData.Competitions = append(filteredData.Competitions, data.Competitions[i])
			}
		}

		filteredBody, _ := json.Marshal(filteredData)
		error := clientR.SetNX(competitions, filteredBody, 360*time.Second).Err()
		if error != nil {
			panic(error)
		}

		json.NewEncoder(w).Encode(filteredData)

	} else {
		fmt.Println("Fetching from the cache")
		json.Unmarshal([]byte(val), &data)
		json.NewEncoder(w).Encode(data)

	}

}
func allowedFreeAPIProvider(el int, list []int) bool {
	for _, b := range list {
		if b == el {
			return true
		}
	}
	return false
}

func fetchStandingsById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	var id = vars["id"]

	var key = "competitions" + id + "standinds"

	val, err := clientR.Get(key).Result()
	if err == redis.Nil {
		fmt.Println("Fetching from the web")
	} else if err != nil {
		panic(err)
	}

	var data StandingsData
	var totalData StandingsData

	if len(val) == 0 {
		req, err := http.NewRequest("GET", url+"competitions/"+id+"/standings", nil)
		req.Header.Add("X-Auth-Token", API_KEY)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &data)

		for i, el := range data.Standings {

			if el.Type == "TOTAL" {
				totalData.Standings = append(totalData.Standings, data.Standings[i])
			}
		}
		totalBody, _ := json.Marshal(totalData)
		error := clientR.SetNX(key, totalBody, 360*time.Second).Err()
		if error != nil {
			panic(error)
		}
		json.NewEncoder(w).Encode(totalData)

	} else {
		fmt.Println("Fetching from the cache")
		json.Unmarshal([]byte(val), &data)
		json.NewEncoder(w).Encode(data)
	}

}
