package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var (
	tasks  = []Task{}
	nextId = 1
	mu     sync.Mutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", tasksHandler)
	mux.HandleFunc("/tasks/", taskHandler)

	handler := PanicMiddleware(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(tasks)
	case http.MethodPost:
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		task.ID = nextId
		nextId++
		tasks = append(tasks, task)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(task)  
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, task := range tasks {
		if task.ID == id {
			switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(task)
				return
			case http.MethodPut:
				var updatedTask Task
				if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				tasks[i] = updatedTask
				tasks[i].ID = task.ID
				json.NewEncoder(w).Encode(tasks[i])
				return
			case http.MethodDelete:
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
	}
	http.NotFound(w, r)
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic: %v\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
