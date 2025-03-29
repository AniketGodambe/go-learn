package model

import (
	"time"
)

// Todo represents a todo item
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// In-memory database simulation
var todos []Todo
var currentID = 1

// GetAllTodos returns all todos
func GetAllTodos() []Todo {
	return todos
}

// GetTodoByID returns a single todo by ID
func GetTodoByID(id int) (*Todo, error) {
	for _, todo := range todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, nil
}

// CreateTodo adds a new todo
func CreateTodo(todo Todo) Todo {
	todo.ID = currentID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	currentID++
	todos = append(todos, todo)
	return todo
}

// UpdateTodo modifies an existing todo
func UpdateTodo(id int, updatedTodo Todo) (*Todo, error) {
	for i, todo := range todos {
		if todo.ID == id {
			updatedTodo.ID = id
			updatedTodo.CreatedAt = todo.CreatedAt
			updatedTodo.UpdatedAt = time.Now()
			todos[i] = updatedTodo
			return &todos[i], nil
		}
	}
	return nil, nil
}

// DeleteTodo removes a todo
func DeleteTodo(id int) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return nil
}
