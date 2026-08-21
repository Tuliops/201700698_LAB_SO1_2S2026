package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const Carnet = "201700698"
const VMNum = 1

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        int    `json:"VM"`
	Carnet    string `json:"carnet"`
}

type CallResponse struct {
	Apiname    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

// Retorna el estado de  API2
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := HealthResponse{
		Status:    "UP",
		Message:   "API2 is Ready",
		Timestamp: time.Now().Format(time.RFC3339),
		VM:        VMNum,
		Carnet:    Carnet,
	}
	json.NewEncoder(w).Encode(resp)
}

// Consulta el endpoint /health de la API destino y evalua su estado
func callAPI(targetAPI string, targetURL string, targetVM string) CallResponse {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(targetURL)

	if err != nil || resp.StatusCode != http.StatusOK {
		return CallResponse{
			Apiname:    targetAPI,
			Message:    fmt.Sprintf("ERROR: The %s located on the %s is not working", targetAPI, targetVM),
			Connection: false,
			Carnet:     Carnet,
		}
	}
	defer resp.Body.Close()

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil || health.Status != "UP" {
		return CallResponse{
			Apiname:    targetAPI,
			Message:    fmt.Sprintf("ERROR: The %s located on the %s is not working", targetAPI, targetVM),
			Connection: false,
			Carnet:     Carnet,
		}
	}

	return CallResponse{
		Apiname:    targetAPI,
		Message:    fmt.Sprintf("The %s located on the %s is working", targetAPI, targetVM),
		Connection: true,
		Carnet:     Carnet,
	}
}

// Registra los endpoints e inicia el servidor en el puerto 8082
func main() {
	http.HandleFunc("/health", healthHandler)

	http.HandleFunc(fmt.Sprintf("/api2/%s/call-api1", Carnet), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		url := os.Getenv("API1_URL")
		if url == "" {
			url = "http://localhost:8081/health"
		}
		json.NewEncoder(w).Encode(callAPI("API1", url, "VM1"))
	})

	http.HandleFunc(fmt.Sprintf("/api2/%s/call-api3", Carnet), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		url := os.Getenv("API3_URL")
		if url == "" {
			url = "http://192.168.122.12:8083/health"
		}
		json.NewEncoder(w).Encode(callAPI("API3", url, "VM2"))
	})

	fmt.Println("API2 escuchando en el puerto 8082 (VM1)...")
	http.ListenAndServe(":8082", nil)
}

