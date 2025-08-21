package main

import (
	"log"
	"net/http"
	handler "study/handler"
)

func main() {
	http.HandleFunc("/send", handler.SendHandler)
	http.HandleFunc("/verify/", handler.VerifyHandler)

	// Добавьте эту строку для запуска сервера
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
