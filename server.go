package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Todo struct {
	Title string `json:"title"`
}

type TodoStore interface {
	GetTodos() map[int]Todo
	CreateTodo(Todo)
}

type TodoServer struct {
	store TodoStore
}

func getTodosJson(ts TodoStore) ([]byte, error) {

	return json.Marshal(ts.GetTodos())
}

func (t *TodoServer) GetAllTodos(w http.ResponseWriter) {
	data, err := getTodosJson(t.store)
	if err != nil {
		error := httpError{Code: http.StatusInternalServerError, Reason: err.Error()}
		http.Error(w, string(error.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (t *TodoServer) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		error := httpError{Code: http.StatusInternalServerError, Reason: err.Error()}
		http.Error(w, string(error.Error()), http.StatusInternalServerError)
		return
	}

	arr := []byte(string(jsonData))

	err = json.Unmarshal(arr, &todo)
	if err != nil {
		error := httpError{Code: http.StatusInternalServerError, Reason: err.Error()}
		http.Error(w, string(error.Error()), http.StatusInternalServerError)
		return
	}

	requestStatus := http.StatusBadRequest

	if todo.Title == "" {
		error := httpError{Code: http.StatusBadRequest, Reason: "Invalid Title"}
		http.Error(w, string(error.Error()), http.StatusBadRequest)
		return
	}

	t.store.CreateTodo(todo)
	requestStatus = http.StatusCreated
	data, err := getTodosJson(t.store)

	if err != nil {
		error := httpError{Code: http.StatusBadRequest, Reason: err.Error()}
		http.Error(w, string(error.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(requestStatus)
	w.Write(data)
}

func (t *TodoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	switch r.Method {
	case http.MethodGet:
		t.GetAllTodos(w)
	case http.MethodPost:
		t.CreateTodo(w, r)
	}
}
