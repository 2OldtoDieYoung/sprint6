package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task тип данных
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTask эндпоинт запроса для всех задач
func getTask(res http.ResponseWriter, req *http.Request) {
	respon, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(respon)
}

// postTask эндпоинт для создания задачи
func postTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if _, check := tasks[task.ID]; check {
		http.Error(res, "Такая задача уже существует", http.StatusBadRequest)
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

}

// getTaskId эндпоинт для получения Task по ID
func getTaskId(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(res, "ID не найден", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(resp)
}

// deleteTask эндпоинт для удаления задачи
func deleteTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(res, "Такого ID нет, или он удален", http.StatusNoContent)
		return
	}

	delete(tasks, id)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// регистрация обработчиков
	r.Get("/tasks", getTask)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskId)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
