package main

import (
	"auiapp/http_api"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/index", http_api.IndexHandler)
	http.HandleFunc("/save", http_api.SaveHandler)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server: http://localhost:8080/index")
	http.ListenAndServe(":8080", nil)
}
