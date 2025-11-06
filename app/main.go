package main

import (
	"fmt"
	"net/http"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var nextID = 1

func main() {
	loadTodos()
	http.HandleFunc("/todo", todoHandler)
	http.HandleFunc("/todo/", todoIDHandler) // for /todo/1
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}