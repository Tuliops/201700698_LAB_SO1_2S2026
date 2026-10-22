package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/recibir", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		// Reenviar al cliente gRPC en localhost:8083
		resp, err := http.Post("http://localhost:8083/procesar", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("[D1-API] Error contactando al cliente gRPC: %v\n", err)
			http.Error(w, "Error interno", 500)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reporte procesado y enviado a Go D1"))
	})

	log.Println("[D1-API] Escuchando en :8081")
	http.ListenAndServe(":8081", nil)
}
