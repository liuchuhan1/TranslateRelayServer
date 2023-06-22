package main

import (
	"TranslateRelayServer/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/gtrans", handler.GoogleHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to Start Server,err:%s", err.Error())
	}
}
