// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	fmt.Println("Server is listening on :8000")
	http.ListenAndServe(":8000", nil)
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, false)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, true)
}

func handleRequest(w http.ResponseWriter, r *http.Request, concurrent bool) {
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	var sortedArrays [][]int
	if concurrent {
		sortedArrays = sortConcurrently(req.ToSort)
	} else {
		sortedArrays = sortSequentially(req.ToSort)
	}

	elapsedTime := time.Since(startTime).Nanoseconds()

	response := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       elapsedTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sortSequentially(arrays [][]int) [][]int {
	for i := range arrays {
		sort.Ints(arrays[i])
	}
	return arrays
}

func sortConcurrently(arrays [][]int) [][]int {
	var wg sync.WaitGroup
	wg.Add(len(arrays))

	for i := range arrays {
		go func(i int) {
			defer wg.Done()
			sort.Ints(arrays[i])
		}(i)
	}

	wg.Wait()

	return arrays
}
