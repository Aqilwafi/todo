# ── To-Do List — Dockerfile ───────────────────────────────────────────────────
# Menggunakan Python (paling ringan, tanpa build step).
# Ganti ke node/php/go di bawah jika diinginkan.
#
# Build : docker build -t todo-app .
# Run   : docker run -p 969:969 -v $(pwd)/todo.txt:/app/todo.txt todo-app
# Buka  : http://localhost:969
# ─────────────────────────────────────────────────────────────────────────────

FROM python:3.12-alpine

WORKDIR /app

# Salin file aplikasi
COPY server.py  ./
COPY index.html ./

# Buat todo.txt kosong jika belum ada (akan di-mount oleh volume)
RUN touch todo.txt

# Expose port
EXPOSE 969

# Jalankan server
CMD ["python", "server.py"]


# ══════════════════════════════════════════════════════════════════════════════
# ALTERNATIF: uncomment salah satu blok di bawah jika ingin pakai runtime lain
# (hapus/comment blok Python di atas terlebih dahulu)
# ══════════════════════════════════════════════════════════════════════════════

# ── Node.js ───────────────────────────────────────────────────────────────────
# FROM node:22-alpine
# WORKDIR /app
# COPY server.js  ./
# COPY index.html ./
# RUN touch todo.txt
# EXPOSE 969
# CMD ["node", "server.js"]

# ── PHP ───────────────────────────────────────────────────────────────────────
# FROM php:8.3-cli-alpine
# WORKDIR /app
# COPY server.php ./
# COPY index.html ./
# RUN touch todo.txt
# EXPOSE 969
# CMD ["php", "-S", "0.0.0.0:969", "server.php"]

# ── Go (multi-stage, binary kecil) ───────────────────────────────────────────
# FROM golang:1.22-alpine AS builder
# WORKDIR /build
# COPY server.go ./
# RUN go build -o server server.go
#
# FROM alpine:3.19
# WORKDIR /app
# COPY --from=builder /build/server ./
# COPY index.html ./
# RUN touch todo.txt
# EXPOSE 969
# CMD ["./server"]