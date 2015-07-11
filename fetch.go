package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		writeError(w, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	var fr fetchRequest
	err = json.Unmarshal(body, &fr)

	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	count := len(fr.URLs)
	results := make(chan fetchResult, count)

	for id, url := range fr.URLs {
		go func(id int, url string) {
			fetch(id, url, results)
		}(id, url)
	}

	response := fetchResponse{
		Results: make([]fetchResult, count),
	}

	var received int
	for result := range results {
		response.Results[result.ID] = result

		received++

		if received == count {
			break
		}
	}

	close(results)

	answer, err := json.Marshal(response)

	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	w.Write(answer)
}

func fetch(id int, url string, results chan<- fetchResult) {

	r := fetchResult{
		ID:      id,
		URL:     url,
		Success: false,
	}

	resp, err := http.Get(url)

	if err != nil {
		results <- r
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		results <- r
		return
	}

	r.Success = true
	r.Body = string(body)

	results <- r

}

func writeError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(
		w,
		`{ "error" : "%s" }`,
		http.StatusText(statusCode),
	)
}

type fetchRequest struct {
	URLs []string `json:"urls"`
}

type fetchResult struct {
	ID      int    `json:"-"`
	URL     string `json:"url"`
	Body    string `json:"body"`
	Success bool   `json:"success"`
}

type fetchResponse struct {
	Results []fetchResult `json:"results"`
}
