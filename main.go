package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RequestBody structure
type RequestBody struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ResponseBody structure
type ResponseBody struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CurrentTime string `json:"current_time"`
	Say         string `json:"say"`
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "Server address")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var requestBody RequestBody

		requestBodyRaw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(requestBodyRaw, &requestBody)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		responseBody := ResponseBody{
			ID:          requestBody.ID,
			FirstName:   requestBody.FirstName + md5Hash(requestBody.FirstName),
			LastName:    requestBody.LastName + md5Hash(requestBody.LastName),
			CurrentTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Say:         "Go is best",
		}

		responseBodyRaw, err := json.Marshal(responseBody)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBodyRaw)
	})

	log.Println("Server listen on " + *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
