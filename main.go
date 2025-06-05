package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"github.com/bgalek/container-log-sanitizer/redactor"
)

// Plugin activation handshake
func activate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1.1+json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Implements": ["docker.logdriver"]}`))
}

// Capabilities endpoint
func capabilities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"Capabilities": {"ReadLogs": false}}`))
}

// StartLogging endpoint
func startLogging(w http.ResponseWriter, r *http.Request) {
	var req struct {
		File string                 `json:"File"`
		Info map[string]interface{} `json:"Info"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	red, err := redactor.NewRedactor()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go func(file string) {
		f, err := os.OpenFile(file, os.O_RDWR, 0600)
		if err != nil {
			log.Printf("Failed to open log pipe: %v", err)
			return
		}
		defer f.Close()
		r := io.Reader(f)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			redacted := red.RedactLine(line)
			os.Stdout.Write([]byte(redacted + "\n"))
		}
	}(req.File)
	w.WriteHeader(http.StatusOK)
}

// StopLogging endpoint
func stopLogging(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// No-op for now
}

func main() {
	log.Println("container-log-sanitizer plugin starting...")

	http.HandleFunc("/Plugin.Activate", activate)
	http.HandleFunc("/LogDriver.Capabilities", capabilities)
	http.HandleFunc("/LogDriver.StartLogging", startLogging)
	http.HandleFunc("/LogDriver.StopLogging", stopLogging)

	// Docker plugins listen on Unix socket, default to /run/docker/plugins/<plugin>.sock
	sock := "/run/docker/plugins/container-log-sanitizer.sock"
	os.Remove(sock) // Remove if exists
	listener, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatalf("Failed to listen on unix socket: %v", err)
	}
	defer listener.Close()
	os.Chmod(sock, 0660)
	log.Printf("Listening on unix socket: %s", sock)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
