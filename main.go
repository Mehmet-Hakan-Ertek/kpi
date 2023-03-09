package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	Name     string
	Complete bool
}

type Todos []Todo

var todoList Todos

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Todo{})

func main() {

	router := mux.NewRouter()

	

	router.HandleFunc("/getTodos", getTodos).Methods("GET")
	router.HandleFunc("/postTodos", postTodos).Methods("POST")
	router.HandleFunc("/updateTodo/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/removeTodo/{id}", removeTodo).Methods("DELETE")

	fmt.Println("Listenin at 8080")
	http.ListenAndServe(":8080", router)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	todos := todoList

	db.First(&todo, 1)
	respondWithJSON(w, http.StatusOK, todos)
}

func postTodos(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todo); err != nil {
		respondWithError(w, http.StatusOK, "yok")
	}

	db.Create(&todo)

	respondWithJSON(w, http.StatusOK, todo)
	todoList = append(todoList, todo)

}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todo Todo

	s, err := strconv.ParseUint(vars["id"], 10, 32)

	db.Model(&todo).Update(todo)
	todoList[s] = todo
	fmt.Println(err)

	todos := todoList

	respondWithJSON(w, http.StatusOK, todos)
}

func removeTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	s, err := strconv.ParseUint(vars["id"], 10, 32)

	fmt.Println(err)

	todoList = append(todoList[:s], todoList[s+1:]...)

	db.Delete(&todo, s)
	respondWithJSON(w, http.StatusOK, todoList)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
