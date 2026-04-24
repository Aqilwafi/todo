<?php
/**
 * To-Do List Server — PHP (tanpa framework)
 * Jalankan: php -S localhost:969 server.php
 * Buka: http://localhost:969
 */

define('PORT',      969);
define('TODO_FILE', __DIR__ . '/todo.txt');
define('HTML_FILE', __DIR__ . '/index.html');

// ── CORS ──────────────────────────────────────────────────────────────────────
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

// ── File R/W ──────────────────────────────────────────────────────────────────

function read_todos(): array {
    if (!file_exists(TODO_FILE)) return [];
    $lines = file(TODO_FILE, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
    $todos = [];
    foreach ($lines as $i => $line) {
        $parts = explode(',', trim($line), 3);
        if (count($parts) < 2) continue;
        $todos[] = [
            'id'   => $i,
            'task' => trim($parts[0]),
            'done' => strtolower(trim($parts[1])) === 'true',
            'note' => isset($parts[2]) ? trim($parts[2]) : '',
        ];
    }
    return $todos;
}

function write_todos(array $todos): void {
    $lines = array_map(function ($t) {
        $task = str_replace(',', ';', $t['task']);
        $note = str_replace(',', ';', $t['note'] ?? '');
        $done = $t['done'] ? 'true' : 'false';
        return "{$task},{$done},{$note}";
    }, $todos);
    file_put_contents(TODO_FILE, implode("\n", $lines) . (count($lines) ? "\n" : ''));
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function send_json($data, int $status = 200): void {
    http_response_code($status);
    header('Content-Type: application/json; charset=utf-8');
    echo json_encode($data, JSON_UNESCAPED_UNICODE);
    exit;
}

function body(): array {
    $raw = file_get_contents('php://input');
    return json_decode($raw ?: '{}', true) ?? [];
}

function path_segments(): array {
    $path = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
    return array_values(array_filter(explode('/', $path)));
}

// ── Router ────────────────────────────────────────────────────────────────────

$method   = $_SERVER['REQUEST_METHOD'];
$segments = path_segments();   // e.g. ['api','todos','2','toggle']
$path     = '/' . implode('/', $segments);

// GET / atau /index.html
if ($method === 'GET' && ($path === '/' || $path === '' || $path === '/index.html')) {
    if (!file_exists(HTML_FILE)) { http_response_code(404); echo 'index.html tidak ditemukan'; exit; }
    header('Content-Type: text/html; charset=utf-8');
    readfile(HTML_FILE);
    exit;
}

// GET /api/todos
if ($method === 'GET' && $path === '/api/todos') {
    send_json(read_todos());
}

// POST /api/todos — tambah
if ($method === 'POST' && $path === '/api/todos') {
    $data = body();
    $task = trim($data['task'] ?? '');
    if (!$task) send_json(['error' => 'Task tidak boleh kosong'], 400);
    $todos   = read_todos();
    $new     = ['id' => count($todos), 'task' => $task, 'done' => false, 'note' => trim($data['note'] ?? '')];
    $todos[] = $new;
    write_todos($todos);
    send_json(['success' => true, 'todo' => $new]);
}

// POST /api/todos/:id/toggle
if ($method === 'POST' && count($segments) === 4 && $segments[3] === 'toggle') {
    $idx   = (int) $segments[2];
    $todos = read_todos();
    if ($idx < 0 || $idx >= count($todos)) send_json(['error' => 'Tidak ditemukan'], 404);
    $todos[$idx]['done'] = !$todos[$idx]['done'];
    write_todos($todos);
    send_json(['success' => true, 'todo' => $todos[$idx]]);
}

// PUT /api/todos/:id — edit
if ($method === 'PUT' && count($segments) === 3) {
    $idx   = (int) $segments[2];
    $data  = body();
    $todos = read_todos();
    if ($idx < 0 || $idx >= count($todos)) send_json(['error' => 'Tidak ditemukan'], 404);
    if (isset($data['task'])) $todos[$idx]['task'] = trim($data['task']);
    if (isset($data['note'])) $todos[$idx]['note'] = trim($data['note']);
    if (isset($data['done'])) $todos[$idx]['done'] = (bool) $data['done'];
    write_todos($todos);
    send_json(['success' => true, 'todo' => $todos[$idx]]);
}

// DELETE /api/todos/:id
if ($method === 'DELETE' && count($segments) === 3) {
    $idx     = (int) $segments[2];
    $todos   = read_todos();
    if ($idx < 0 || $idx >= count($todos)) send_json(['error' => 'Tidak ditemukan'], 404);
    $removed = array_splice($todos, $idx, 1)[0];
    write_todos($todos);
    send_json(['success' => true, 'removed' => $removed]);
}

http_response_code(404);
echo 'Tidak ditemukan';