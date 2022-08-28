package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Eliad-S/Permutation_web_service/db"
	"github.com/Eliad-S/Permutation_web_service/statistics"
)

type Similar_words struct {
	Similar []string `jsob:"similar`
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)
	if first == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, "This is my website!\n")
}

func GetStats(w http.ResponseWriter, r *http.Request) {

	stats := statistics.GetStats()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetSimilar(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ctx := r.Context()

	hasWord := r.URL.Query().Has("word")
	word := r.URL.Query().Get("word")

	fmt.Printf("%s: got / request. word(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasWord, word)
	if word == "" {
		w.Header().Set("x-missing-field", "word")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	similar_words, err := db.Get_similar_words(word)
	if err != nil {
		panic(err.Error())
	}
	similar := Similar_words{Similar: similar_words}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(similar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	duration := time.Since(start)

	statistics.Inc_requests(duration.Nanoseconds())
}
