package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /agent", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Agent information retrieved successfully"}`))
	})
	return mux
}