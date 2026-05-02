package main

import "net/http"

func main() {
	cfg := NewConfig()

	server := &http.Server{
		Addr: cfg.Listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleProxy(w, r, cfg)
		}),
	}

	cfg.Logger.Info("turbulence proxy listening", "port", cfg.Listen)
	cfg.Logger.Error("server error", "error", server.ListenAndServe())
}
