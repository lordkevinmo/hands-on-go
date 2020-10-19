package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lordkevinmo/hands-on-go/chapter7/urlShortner/helper"
	"github.com/lordkevinmo/hands-on-go/chapter7/urlShortner/utils"
)

// DBClient holds the link to the database
type DBClient struct {
	db *sql.DB
}

// Record holds the data necessary for generate the shortned url
type Record struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func main() {
	db, err := helper.InitDB()
	if err != nil {
		panic(err)
	}
	dbClient := &DBClient{db: db}
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Create new Router
	r := mux.NewRouter()
	// Attach an elegant path with handler
	r.HandleFunc("v1/short/{encoded_string:[a-zA-Z0-9]*}", dbClient.GetOriginalURL).Methods(http.MethodGet)
	r.HandleFunc("v1/short", dbClient.GenerateShortURL).Methods(http.MethodPost)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// GenerateShortURL adds URL to DB and gives back shortened string
func (driver *DBClient) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	var id int
	var record Record
	postBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(postBody, &record)
	err = driver.db.QueryRow("INSERT INTO web_url(url) VALUES() RETURNING id", record.URL).Scan(&id)
	responseMap := map[string]string{"encoded_string": utils.ToBase62(id)}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(responseMap)
		w.Write(response)
	}
}

// GetOriginalURL fetches the original URL for the given encoded (short) string
func (driver *DBClient) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	var url string
	vars := mux.Vars(r)
	// Get ID from base62
	id := utils.ToBase10(vars["encoded_string"])
	err := driver.db.QueryRow("SELECT url FROM web_url WHERE id = $1", id).Scan(&url)
	// Handle Response details
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		responseMap := map[string]interface{}{"url": url}
		response, _ := json.Marshal(responseMap)
		w.Write(response)
	}
}
