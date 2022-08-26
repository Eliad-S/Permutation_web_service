package api

import (
	"fmt"
	"io"
	"net/http"
)

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")
}

func getStats(w http.ResponseWriter, r *http.Request) {
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

func getSimilar(w http.ResponseWriter, r *http.Request) {
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
