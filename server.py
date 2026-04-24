#!/usr/bin/env python3
"""
Simple To-Do List Server
Format file: to-do,boolean,note
Jalankan: python server.py
Buka: http://localhost:969
"""

import http.server
import json
import os
import urllib.parse
from pathlib import Path

PORT = 969
TODO_FILE = "todo.txt"
HTML_FILE = "index.html"

def read_todos():
    """Baca semua todo dari file."""
    todos = []
    if not os.path.exists(TODO_FILE):
        return todos
    with open(TODO_FILE, "r", encoding="utf-8") as f:
        for i, line in enumerate(f):
            line = line.strip()
            if not line:
                continue
            parts = line.split(",", 2)
            if len(parts) < 2:
                continue
            task = parts[0].strip()
            done = parts[1].strip().lower() == "true"
            note = parts[2].strip() if len(parts) > 2 else ""
            todos.append({"id": i, "task": task, "done": done, "note": note})
    return todos

def write_todos(todos):
    """Tulis semua todo ke file."""
    with open(TODO_FILE, "w", encoding="utf-8") as f:
        for todo in todos:
            done_str = "true" if todo["done"] else "false"
            note = todo.get("note", "").replace(",", ";")  # hindari konflik delimiter
            task = todo["task"].replace(",", ";")
            f.write(f"{task},{done_str},{note}\n")


class TodoHandler(http.server.BaseHTTPRequestHandler):

    def log_message(self, format, *args):
        print(f"[{self.address_string()}] {format % args}")

    def send_json(self, data, status=200):
        body = json.dumps(data, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", len(body))
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()
        self.wfile.write(body)

    def serve_file(self, path, content_type):
        if not os.path.exists(path):
            self.send_error(404, f"File tidak ditemukan: {path}")
            return
        with open(path, "rb") as f:
            content = f.read()
        self.send_response(200)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", len(content))
        self.end_headers()
        self.wfile.write(content)

    def do_OPTIONS(self):
        self.send_response(200)
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type")
        self.end_headers()

    def do_GET(self):
        parsed = urllib.parse.urlparse(self.path)
        path = parsed.path

        if path == "/" or path == "/index.html":
            self.serve_file(HTML_FILE, "text/html; charset=utf-8")
        elif path == "/api/todos":
            self.send_json(read_todos())
        else:
            self.send_error(404, "Tidak ditemukan")

    def do_POST(self):
        parsed = urllib.parse.urlparse(self.path)
        path = parsed.path

        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)

        try:
            data = json.loads(body)
        except json.JSONDecodeError:
            self.send_json({"error": "JSON tidak valid"}, 400)
            return

        if path == "/api/todos":
            # Tambah todo baru
            task = data.get("task", "").strip()
            if not task:
                self.send_json({"error": "Task tidak boleh kosong"}, 400)
                return
            todos = read_todos()
            new_todo = {
                "id": len(todos),
                "task": task,
                "done": False,
                "note": data.get("note", "").strip()
            }
            todos.append(new_todo)
            write_todos(todos)
            self.send_json({"success": True, "todo": new_todo})

        elif path.startswith("/api/todos/") and path.endswith("/toggle"):
            # Toggle done/undone
            try:
                idx = int(path.split("/")[3])
            except (IndexError, ValueError):
                self.send_json({"error": "ID tidak valid"}, 400)
                return
            todos = read_todos()
            if idx < 0 or idx >= len(todos):
                self.send_json({"error": "Todo tidak ditemukan"}, 404)
                return
            todos[idx]["done"] = not todos[idx]["done"]
            write_todos(todos)
            self.send_json({"success": True, "todo": todos[idx]})

        else:
            self.send_error(404, "Endpoint tidak ditemukan")

    def do_PUT(self):
        parsed = urllib.parse.urlparse(self.path)
        path = parsed.path

        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)

        try:
            data = json.loads(body)
        except json.JSONDecodeError:
            self.send_json({"error": "JSON tidak valid"}, 400)
            return

        if path.startswith("/api/todos/"):
            try:
                idx = int(path.split("/")[3])
            except (IndexError, ValueError):
                self.send_json({"error": "ID tidak valid"}, 400)
                return
            todos = read_todos()
            if idx < 0 or idx >= len(todos):
                self.send_json({"error": "Todo tidak ditemukan"}, 404)
                return
            if "task" in data:
                todos[idx]["task"] = data["task"].strip()
            if "note" in data:
                todos[idx]["note"] = data["note"].strip()
            if "done" in data:
                todos[idx]["done"] = bool(data["done"])
            write_todos(todos)
            self.send_json({"success": True, "todo": todos[idx]})
        else:
            self.send_error(404, "Endpoint tidak ditemukan")

    def do_DELETE(self):
        parsed = urllib.parse.urlparse(self.path)
        path = parsed.path

        if path.startswith("/api/todos/"):
            try:
                idx = int(path.split("/")[3])
            except (IndexError, ValueError):
                self.send_json({"error": "ID tidak valid"}, 400)
                return
            todos = read_todos()
            if idx < 0 or idx >= len(todos):
                self.send_json({"error": "Todo tidak ditemukan"}, 404)
                return
            removed = todos.pop(idx)
            write_todos(todos)
            self.send_json({"success": True, "removed": removed})
        else:
            self.send_error(404, "Endpoint tidak ditemukan")


if __name__ == "__main__":
    server = http.server.HTTPServer(("", PORT), TodoHandler)
    print(f"✅ Server berjalan di http://localhost:{PORT}")
    print(f"📄 Membaca/menulis ke: {os.path.abspath(TODO_FILE)}")
    print("   Tekan Ctrl+C untuk berhenti.\n")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n🛑 Server dihentikan.")
        server.server_close()