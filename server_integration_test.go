package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStoreTodoAndRetrievingThem(t *testing.T) {
	store := NewInMemoryTodoStore()
	server := TodoServer{store}

	todo := Todo{}
	todo.Title = "integrationTest"

	server.ServeHTTP(httptest.NewRecorder(), createTodoRequest(todo))
	server.ServeHTTP(httptest.NewRecorder(), createTodoRequest(todo))
	server.ServeHTTP(httptest.NewRecorder(), createTodoRequest(todo))

	response := httptest.NewRecorder()

	server.ServeHTTP(response, getTodosRequest())
	assertStatus(t, response.Code, http.StatusOK)

	todos := map[int]Todo{
		0: {Title: "integrationTest"},
		1: {Title: "integrationTest"},
		2: {Title: "integrationTest"},
	}
	want, _ := json.Marshal(todos)
	assertResponseBody(t, response.Body.String(), string(want))

}
