package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func getTodos(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func postTodo(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var newTodo Todo
	json.Unmarshal(body, &newTodo)

	newTodo.ID = nextID
	nextID++
	newTodo.Completed = false

	todos = append(todos, newTodo)
	saveTodos()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request, id int) {
	body, _ := io.ReadAll(r.Body)
	var updated Todo
	json.Unmarshal(body, &updated)

	for i, t := range todos {
		if t.ID == id {
			if updated.Title != "" {
				todos[i].Title = updated.Title
			}
			todos[i].Completed = updated.Completed
			saveTodos()
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	http.NotFound(w, r)
}

func deleteTodo(w http.ResponseWriter, r *http.Request, id int) {
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}