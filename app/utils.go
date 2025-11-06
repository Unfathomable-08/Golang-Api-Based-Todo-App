package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func todoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getTodos(w)
	case http.MethodPost:
		postTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func todoIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todo/")
	id, _ := strconv.Atoi(idStr)

	switch r.Method {
	case http.MethodPut:
		updateTodo(w, r, id)
	case http.MethodDelete:
		deleteTodo(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func loadTodos() {
	data, err := os.ReadFile("todo.json")
	if err != nil {
		todos = []Todo{}
		return
	}
	json.Unmarshal(data, &todos)
	if len(todos) > 0 {
		nextID = todos[len(todos)-1].ID + 1
	}
}

func saveTodos() {
	data, _ := json.MarshalIndent(todos, "", "  ")
	os.WriteFile("todo.json", data, 0644)
}