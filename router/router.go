package router

import (
	"go-learn/controller"

	"github.com/gorilla/mux"
)

// InitializeRouter sets up all routes
func InitializeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/todo/get-todo-list", controller.GetTodos).Methods("GET")
	router.HandleFunc("/api/todo/get-todo-by-id/{id}", controller.GetTodo).Methods("GET")
	router.HandleFunc("/api/todo/add-todo", controller.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todo/update-todo/{id}", controller.UpdateTodo).Methods("PUT")
	router.HandleFunc("/api/todo/delete-todo/{id}", controller.DeleteTodo).Methods("DELETE")

	return router
}
