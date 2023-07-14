package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type Response struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

func main() {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: err.Error(),
			})
			return
		}

		if message.Receiver == "" || message.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "receiver or message is empty",
			})
			return
		}

		fmt.Printf("Receiver: %s\nMessage: %s\n", message.Receiver, message.Message)

		json.NewEncoder(w).Encode(Response{
			StatusCode:   http.StatusOK,
			ErrorMessage: "",
		})
	})

	log.Fatal(http.ListenAndServe(":3002", nil))
}
