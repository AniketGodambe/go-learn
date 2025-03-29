package controller

import (
	"encoding/json"
	"go-learn/model"
	"go-learn/utils"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const dataFile = "todos.json"

var (
	todos     []model.Todo
	currentID int
	mu        sync.Mutex
)

// Initialize loads todos from file
func Initialize() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, start with empty todos
			todos = []model.Todo{}
			currentID = 1
			return nil
		}
		return err
	}

	if err := json.Unmarshal(file, &todos); err != nil {
		return err
	}

	// Set currentID to the highest ID + 1
	if len(todos) > 0 {
		currentID = todos[len(todos)-1].ID + 1
	} else {
		currentID = 1
	}

	return nil
}

func saveTodos() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, data, 0644)
}

// GetTodos returns all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	utils.SetHeaders(w, http.MethodGet)

	utils.ValidateMethod(r.Method, http.MethodGet, w)

	todos := GetAllTodos()

	if todos == nil {
		utils.SendJSON(w, http.StatusOK, "Todo's retrieved successfully", true, []model.Todo{})
		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo's retrieved successfully", true, GetAllTodos())

}

func GetAllTodos() []model.Todo {
	mu.Lock()
	defer mu.Unlock()
	return todos
}

// GetTodo returns a single todo
func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	todo, err := GetTodoByID(id)
	if err != nil || todo == nil {

		utils.SendJSON(w, http.StatusNotFound, "Todo not found", false, nil)

		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo retrieved successfully", true, todo)

}

func GetTodoByID(id int) (*model.Todo, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, todo := range todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, nil
}

// CreateTodo adds a new todo

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo model.Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)

	newTodo, err := CreateTodoFn(todo)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "Failed to create todo", false, nil)
		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo added successfully", true, newTodo)
}

func CreateTodoFn(todo model.Todo) (model.Todo, error) {
	mu.Lock()
	todo.ID = currentID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	currentID++
	todos = append(todos, todo)
	mu.Unlock()

	if err := saveTodos(); err != nil {
		return model.Todo{}, err
	}
	return todo, nil
}

// UpdateTodo modifies an existing todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedTodo model.Todo
	_ = json.NewDecoder(r.Body).Decode(&updatedTodo)

	todo, err := UpdateTodoFn(id, updatedTodo)
	if err != nil || todo == nil {
		utils.SendJSON(w, http.StatusNotFound, "Todo not found", false, todo)
		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo updated successfully", true, todo)

}

func UpdateTodoFn(id int, updatedTodo model.Todo) (*model.Todo, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			updatedTodo.ID = id
			updatedTodo.CreatedAt = todo.CreatedAt
			updatedTodo.UpdatedAt = time.Now()
			todos[i] = updatedTodo

			if err := saveTodos(); err != nil {
				return nil, err
			}
			return &todos[i], nil
		}
	}
	return nil, nil
}

// DeleteTodo removes a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	todo, err := GetTodoByID(id)
	if err != nil || todo == nil {
		utils.SendJSON(w, http.StatusNotFound, "Todo not found", false, nil)
		return
	}

	if err := DeleteTodoFn(id); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, "Failed to delete todo", false, nil)
		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo deleted successfully", true, todo)
}

func DeleteTodoFn(id int) error {
	// First find the index with minimal locking
	index := -1
	mu.Lock()
	for i, todo := range todos {
		if todo.ID == id {
			index = i
			break
		}
	}
	mu.Unlock()

	if index == -1 {
		return nil // Not found
	}

	// Perform the deletion and save with proper locking
	mu.Lock()
	todos = append(todos[:index], todos[index+1:]...)
	mu.Unlock()

	// Save to file outside the lock
	return saveTodos()
}
