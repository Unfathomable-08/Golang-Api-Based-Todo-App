package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Println(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func printOptions() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1 - Add Todo")
		fmt.Println("2 - List Todos")
		fmt.Println("3 - Mark Completed")
		fmt.Println("4 - Delete Todo")
		fmt.Println("5 - Exit")

		opt, _ := getInput("Enter option number:", reader)

		switch opt {
		case "1":
			title, _ := getInput("Enter todo title:", reader)
			addTodo(title)

		case "2":
			listTodos()

		case "3":
			idStr, _ := getInput("Enter todo ID to mark completed:", reader)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID")
				continue
			}
			markCompleted(id)

		case "4":
			idStr, _ := getInput("Enter todo ID to delete:", reader)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID")
				continue
			}
			deleteTodo(id)

		case "5":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid option")
		}
	}
}

func main() {
	loadTodos()
	printOptions()
}
