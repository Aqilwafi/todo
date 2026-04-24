#!/usr/bin/env node
/**
 * To-Do List Server — Node.js (tanpa library)
 * Jalankan: node server.js
 * Buka: http://localhost:969
 */

const http = require("http");
const fs = require("fs");
const path = require("path");
const url = require("url");

const PORT = 969;
const TODO_FILE = path.join(__dirname, "todo.txt");
const HTML_FILE = path.join(__dirname, "index.html");

// ── File R/W ─────────────────────────────────────────────────────────────────

function readTodos() {
  if (!fs.existsSync(TODO_FILE)) return [];
  const lines = fs.readFileSync(TODO_FILE, "utf-8").split("\n");
  const todos = [];
  lines.forEach((line, i) => {
    line = line.trim();
    if (!line) return;
    const parts = line.split(",");
    if (parts.length < 2) return;
    const task = parts[0].trim();
    const done = parts[1].trim().toLowerCase() === "true";
    const note = parts.slice(2).join(",").trim();
    todos.push({ id: i, task, done, note });
  });
  return todos;
}

function writeTodos(todos) {
  const lines = todos.map((t) => {
    const task = t.task.replace(/,/g, ";");
    const note = (t.note || "").replace(/,/g, ";");
    return `${task},${t.done ? "true" : "false"},${note}`;
  });
  fs.writeFileSync(TODO_FILE, lines.join("\n") + (lines.length ? "\n" : ""), "utf-8");
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function sendJSON(res, data, status = 200) {
  const body = JSON.stringify(data);
  res.writeHead(status, {
    "Content-Type": "application/json; charset=utf-8",
    "Content-Length": Buffer.byteLength(body),
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
    "Access-Control-Allow-Headers": "Content-Type",
  });
  res.end(body);
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let body = "";
    req.on("data", (chunk) => (body += chunk));
    req.on("end", () => {
      try { resolve(JSON.parse(body || "{}")); }
      catch { reject(new Error("JSON tidak valid")); }
    });
    req.on("error", reject);
  });
}

// ── Server ────────────────────────────────────────────────────────────────────

const server = http.createServer(async (req, res) => {
  const parsed = url.parse(req.path || req.url);
  const pathname = parsed.pathname;
  const method = req.method.toUpperCase();

  // CORS preflight
  if (method === "OPTIONS") {
    res.writeHead(204, {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
      "Access-Control-Allow-Headers": "Content-Type",
    });
    return res.end();
  }

  // Serve HTML
  if (method === "GET" && (pathname === "/" || pathname === "/index.html")) {
    if (!fs.existsSync(HTML_FILE)) {
      res.writeHead(404); return res.end("index.html tidak ditemukan");
    }
    const html = fs.readFileSync(HTML_FILE);
    res.writeHead(200, { "Content-Type": "text/html; charset=utf-8" });
    return res.end(html);
  }

  // GET /api/todos
  if (method === "GET" && pathname === "/api/todos") {
    return sendJSON(res, readTodos());
  }

  // POST /api/todos — tambah
  if (method === "POST" && pathname === "/api/todos") {
    try {
      const data = await readBody(req);
      const task = (data.task || "").trim();
      if (!task) return sendJSON(res, { error: "Task tidak boleh kosong" }, 400);
      const todos = readTodos();
      const newTodo = { id: todos.length, task, done: false, note: (data.note || "").trim() };
      todos.push(newTodo);
      writeTodos(todos);
      return sendJSON(res, { success: true, todo: newTodo });
    } catch { return sendJSON(res, { error: "Request tidak valid" }, 400); }
  }

  // POST /api/todos/:id/toggle
  const toggleMatch = pathname.match(/^\/api\/todos\/(\d+)\/toggle$/);
  if (method === "POST" && toggleMatch) {
    const idx = parseInt(toggleMatch[1]);
    const todos = readTodos();
    if (idx < 0 || idx >= todos.length) return sendJSON(res, { error: "Tidak ditemukan" }, 404);
    todos[idx].done = !todos[idx].done;
    writeTodos(todos);
    return sendJSON(res, { success: true, todo: todos[idx] });
  }

  // PUT /api/todos/:id — edit
  const itemMatch = pathname.match(/^\/api\/todos\/(\d+)$/);
  if (method === "PUT" && itemMatch) {
    try {
      const idx = parseInt(itemMatch[1]);
      const data = await readBody(req);
      const todos = readTodos();
      if (idx < 0 || idx >= todos.length) return sendJSON(res, { error: "Tidak ditemukan" }, 404);
      if (data.task !== undefined) todos[idx].task = data.task.trim();
      if (data.note !== undefined) todos[idx].note = data.note.trim();
      if (data.done !== undefined) todos[idx].done = Boolean(data.done);
      writeTodos(todos);
      return sendJSON(res, { success: true, todo: todos[idx] });
    } catch { return sendJSON(res, { error: "Request tidak valid" }, 400); }
  }

  // DELETE /api/todos/:id
  if (method === "DELETE" && itemMatch) {
    const idx = parseInt(itemMatch[1]);
    const todos = readTodos();
    if (idx < 0 || idx >= todos.length) return sendJSON(res, { error: "Tidak ditemukan" }, 404);
    const [removed] = todos.splice(idx, 1);
    writeTodos(todos);
    return sendJSON(res, { success: true, removed });
  }

  res.writeHead(404); res.end("Tidak ditemukan");
});

server.listen(PORT, () => {
  console.log(`✅ Server berjalan di http://localhost:${PORT}`);
  console.log(`📄 Membaca/menulis ke: ${TODO_FILE}`);
  console.log("   Tekan Ctrl+C untuk berhenti.\n");
});