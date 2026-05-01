package main

import (
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	cfg := NewConfig()

	server := &http.Server{
		Addr: cfg.Listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleProxy(w, r, cfg)
		}),
	}

	cfg.Logger.Info("chaos proxy listening", "port", cfg.Listen)
	cfg.Logger.Error("server error", "error", server.ListenAndServe())
}

func handleProxy(w http.ResponseWriter, r *http.Request, cfg Config) {
	if shouldDrop(cfg) {
		cfg.Logger.Info("dropping", "method", r.Method, "host", r.Host)
		http.Error(w, "chaos proxy dropped request", http.StatusServiceUnavailable)
		return
	}

	injectDelay(cfg)

	if r.Method == http.MethodConnect {
		handleConnect(w, r, cfg)
		return
	}

	handleHTTP(w, r, cfg)
}

func handleHTTP(w http.ResponseWriter, r *http.Request, cfg Config) {
	cfg.Logger.Info("http info", "method", r.Method, "url", r.URL.String())

	r.RequestURI = ""

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	_, _ = io.Copy(w, resp.Body)
}

func handleConnect(w http.ResponseWriter, r *http.Request, cfg Config) {
	cfg.Logger.Info("connect info", "host", r.Host)

	upstream, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer upstream.Close()

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "connection hijacking unsupported", http.StatusInternalServerError)
		return
	}

	client, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer client.Close()

	_, _ = client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	done := make(chan struct{}, 2)

	go func() {
		_, _ = io.Copy(upstream, client)
		done <- struct{}{}
	}()

	go func() {
		_, _ = io.Copy(client, upstream)
		done <- struct{}{}
	}()

	<-done
}
