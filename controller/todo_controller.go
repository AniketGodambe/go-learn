package controller

import (
	"encoding/json"
	"fmt"
	"go-learn/model"
	"go-learn/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetTodos returns all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	utils.SetHeaders(w, http.MethodGet)
	fmt.Println("METHOD GET", r.Method)
	fmt.Println("METHOD NEEDED", http.MethodGet)

	utils.ValidateMethod(r.Method, http.MethodGet, w)

	todos := model.GetAllTodos()

	if todos == nil {
		utils.SendJSON(w, http.StatusOK, "Todo's retrieved successfully", true, []model.Todo{})
		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo's retrieved successfully", true, model.GetAllTodos())

}

// GetTodo returns a single todo
func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	todo, err := model.GetTodoByID(id)
	if err != nil || todo == nil {

		utils.SendJSON(w, http.StatusNotFound, "Todo not found", false, nil)

		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo retrieved successfully", true, todo)

}

// CreateTodo adds a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todo model.Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)

	newTodo := model.CreateTodo(todo)
	utils.SendJSON(w, http.StatusOK, "Todo added successfully", true, newTodo)

}

// UpdateTodo modifies an existing todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedTodo model.Todo
	_ = json.NewDecoder(r.Body).Decode(&updatedTodo)

	todo, err := model.UpdateTodo(id, updatedTodo)
	if err != nil || todo == nil {
		utils.SendJSON(w, http.StatusNotFound, "Todo not found", false, todo)
		return
	}

	utils.SendJSON(w, http.StatusOK, "Todo updated successfully", true, todo)

}

// DeleteTodo removes a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	todo, err := model.GetTodoByID(id)
	if err != nil || todo == nil {

		utils.SendJSON(w, http.StatusNotFound, "Todo not found", false, nil)

		return
	}
	model.DeleteTodo(id)

	utils.SendJSON(w, http.StatusOK, "Todo deleted successfully", true, todo)

}
