package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubTodoStore struct {
	todos map[int]Todo
}

func (s *stubTodoStore) GetTodos() map[int]Todo {
	return s.todos
}

func (s *stubTodoStore) CreateTodo(todo Todo) {
	s.todos[len(s.todos)] = todo
}

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(rr, req)

	expected := `{"alive": true}`

	assertStatus(t, rr.Code, http.StatusOK)
	assertResponseBody(t, rr.Body.String(), expected)
}

func TestGetTodos(t *testing.T) {
	todos := map[int]Todo{
		0: {Title: "test"},
		1: {Title: "test2"},
	}
	store := &stubTodoStore{todos}
	server := &TodoServer{store}
	t.Run("return all todos", func(t *testing.T) {

		request := getTodosRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()

		want, _ := json.Marshal(todos)

		assertResponseBody(t, got, string(want))
	})
}

func TestCreateTodo(t *testing.T) {
	store := &stubTodoStore{map[int]Todo{}}
	server := &TodoServer{store}
	todo := Todo{Title: "to-do app ready!"}
	t.Run("create a todo", func(t *testing.T) {
		request := createTodoRequest(todo)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()

		want, _ := json.Marshal(map[int]Todo{
			0: {Title: "to-do app ready!"},
		})

		assertStatus(t, response.Code, http.StatusCreated)

		assertResponseBody(t, got, string(want))
	})

	t.Run("should return 400 when empty passed", func(t *testing.T) {
		newTodo := Todo{Title: ""}

		request := createTodoRequest(newTodo)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()

		want := "{\"code\":400,\"reason\":\"Invalid Title\"}\n"

		assertStatus(t, response.Code, http.StatusBadRequest)

		assertResponseBody(t, got, string(want))

	})
}

func getTodosRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/todos"), nil)
	return req
}

func createTodoRequest(todo Todo) *http.Request {
	body := fmt.Sprintf(`{"title":"%s"}`, todo.Title)
	req, _ := http.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
	return req
}

func assertResponseBody(t testing.TB, got, want string) {

	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status is wrong, got %d want %d", uint64(got), uint64(want))
	}
}
