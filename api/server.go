package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func InitilizeServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/similar", GetSimilar)
	mux.HandleFunc("/api/v1/stats", GetStats)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}

}
