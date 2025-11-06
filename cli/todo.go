package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var nextID int = 1

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

func addTodo(title string) {
	newTodo := Todo{
		ID:        nextID,
		Title:     title,
		Completed: false,
	}
	nextID++
	todos = append(todos, newTodo)
	saveTodos()
	fmt.Println("Todo added:", title)
}

func listTodos() {
	if len(todos) == 0 {
		fmt.Println("No todos yet.")
		return
	}

	fmt.Println("\n--- Pending Todos ---")
	for _, t := range todos {
		if !t.Completed {
			fmt.Printf("[%d] %s \n", t.ID, t.Title)
		}
	}

	fmt.Println("\n--- Completed Todos ---")
	for _, t := range todos {
		if t.Completed {
			fmt.Printf("[%d] %s \n", t.ID, t.Title)
		}
	}
}

func markCompleted(id int) {
	for i, t := range todos {
		if t.ID == id {
			todos[i].Completed = true
			saveTodos()
			fmt.Println("Todo marked completed:", t.Title)
			return
		}
	}
	fmt.Println("Todo not found.")
}

func deleteTodo(id int) {
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos()
			fmt.Println("Todo deleted:", t.Title)
			return
		}
	}
	fmt.Println("Todo not found.")
}
