// To-Do List Server — Go (tanpa library eksternal)
// Build  : go build -o server server.go
// Jalankan: ./server   (atau: go run server.go)
// Buka   : http://localhost:969

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	port     = 969
	todoFile = "todo.txt"
	htmlFile = "index.html"
)

// ── Struct ────────────────────────────────────────────────────────────────────

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
	Note string `json:"note"`
}

// ── File R/W ──────────────────────────────────────────────────────────────────

func readTodos() []Todo {
	data, err := os.ReadFile(todoFile)
	if err != nil {
		return []Todo{}
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	todos := []Todo{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 3)
		if len(parts) < 2 {
			continue
		}
		note := ""
		if len(parts) == 3 {
			note = strings.TrimSpace(parts[2])
		}
		todos = append(todos, Todo{
			ID:   i,
			Task: strings.TrimSpace(parts[0]),
			Done: strings.TrimSpace(strings.ToLower(parts[1])) == "true",
			Note: note,
		})
	}
	return todos
}

func writeTodos(todos []Todo) error {
	var sb strings.Builder
	for _, t := range todos {
		task := strings.ReplaceAll(t.Task, ",", ";")
		note := strings.ReplaceAll(t.Note, ",", ";")
		done := "false"
		if t.Done {
			done = "true"
		}
		sb.WriteString(fmt.Sprintf("%s,%s,%s\n", task, done, note))
	}
	return os.WriteFile(todoFile, []byte(sb.String()), 0644)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func sendJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// ── Routes ────────────────────────────────────────────────────────────────────

var (
	reTodosBase   = regexp.MustCompile(`^/api/todos$`)
	reTodosItem   = regexp.MustCompile(`^/api/todos/(\d+)$`)
	reTodosToggle = regexp.MustCompile(`^/api/todos/(\d+)/toggle$`)
)

func handler(w http.ResponseWriter, r *http.Request) {
	cors(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := r.URL.Path
	method := r.Method

	// Serve HTML
	if method == http.MethodGet && (path == "/" || path == "/index.html") {
		abs, _ := filepath.Abs(htmlFile)
		http.ServeFile(w, r, abs)
		return
	}

	// GET /api/todos
	if method == http.MethodGet && reTodosBase.MatchString(path) {
		sendJSON(w, readTodos(), http.StatusOK)
		return
	}

	// POST /api/todos — tambah
	if method == http.MethodPost && reTodosBase.MatchString(path) {
		var body struct {
			Task string `json:"task"`
			Note string `json:"note"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, map[string]string{"error": "Request tidak valid"}, 400)
			return
		}
		task := strings.TrimSpace(body.Task)
		if task == "" {
			sendJSON(w, map[string]string{"error": "Task tidak boleh kosong"}, 400)
			return
		}
		todos := readTodos()
		newTodo := Todo{ID: len(todos), Task: task, Done: false, Note: strings.TrimSpace(body.Note)}
		todos = append(todos, newTodo)
		writeTodos(todos)
		sendJSON(w, map[string]any{"success": true, "todo": newTodo}, http.StatusOK)
		return
	}

	// POST /api/todos/:id/toggle
	if method == http.MethodPost && reTodosToggle.MatchString(path) {
		m := reTodosToggle.FindStringSubmatch(path)
		idx, _ := strconv.Atoi(m[1])
		todos := readTodos()
		if idx < 0 || idx >= len(todos) {
			sendJSON(w, map[string]string{"error": "Tidak ditemukan"}, 404)
			return
		}
		todos[idx].Done = !todos[idx].Done
		writeTodos(todos)
		sendJSON(w, map[string]any{"success": true, "todo": todos[idx]}, http.StatusOK)
		return
	}

	// PUT /api/todos/:id — edit
	if method == http.MethodPut && reTodosItem.MatchString(path) {
		m := reTodosItem.FindStringSubmatch(path)
		idx, _ := strconv.Atoi(m[1])
		todos := readTodos()
		if idx < 0 || idx >= len(todos) {
			sendJSON(w, map[string]string{"error": "Tidak ditemukan"}, 404)
			return
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if v, ok := body["task"].(string); ok { todos[idx].Task = strings.TrimSpace(v) }
		if v, ok := body["note"].(string); ok { todos[idx].Note = strings.TrimSpace(v) }
		if v, ok := body["done"].(bool);   ok { todos[idx].Done = v }
		writeTodos(todos)
		sendJSON(w, map[string]any{"success": true, "todo": todos[idx]}, http.StatusOK)
		return
	}

	// DELETE /api/todos/:id
	if method == http.MethodDelete && reTodosItem.MatchString(path) {
		m := reTodosItem.FindStringSubmatch(path)
		idx, _ := strconv.Atoi(m[1])
		todos := readTodos()
		if idx < 0 || idx >= len(todos) {
			sendJSON(w, map[string]string{"error": "Tidak ditemukan"}, 404)
			return
		}
		removed := todos[idx]
		todos = append(todos[:idx], todos[idx+1:]...)
		writeTodos(todos)
		sendJSON(w, map[string]any{"success": true, "removed": removed}, http.StatusOK)
		return
	}

	http.NotFound(w, r)
}

func main() {
	// Buat todo.txt jika belum ada
	if _, err := os.Stat(todoFile); os.IsNotExist(err) {
		os.WriteFile(todoFile, []byte(""), 0644)
	}

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("✅ Server berjalan di http://localhost:%d\n", port)
	fmt.Printf("📄 Membaca/menulis ke: %s\n", todoFile)
	fmt.Println("   Tekan Ctrl+C untuk berhenti.\n")

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}