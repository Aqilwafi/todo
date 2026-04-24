# 📝 To-Do List — localhost:969

Aplikasi to-do list sederhana yang berjalan di **localhost** dengan penyimpanan berbasis file teks (`todo.txt`). Tidak memerlukan database, tidak memerlukan framework — cukup satu file server + satu HTML.

---

## 📁 Struktur File

```
.
├── index.html          # Antarmuka web (sama untuk semua server)
├── todo.txt            # Penyimpanan data (dibuat otomatis)
├── server.py           # Server Python
├── server.js           # Server Node.js
├── server.php          # Server PHP
├── server.go           # Server Go
├── Dockerfile          # Docker image
├── docker-compose.yml  # Docker Compose
└── README.md
```

> `index.html` **tidak perlu diubah** untuk server manapun — semua server mengimplementasikan API yang sama.

---

## 🚀 Cara Menjalankan

Pilih salah satu runtime yang tersedia di komputer kamu:

### <img src="https://skillicons.dev/icons?i=python" width="30" height="30" valign="middle" /> Python
```bash
python server.py
```
> Memerlukan Python 3.x (sudah termasuk di macOS & kebanyakan Linux)

### <img src="https://skillicons.dev/icons?i=nodejs" width="30" height="30" valign="middle" /> Node.js
```bash
node server.js
```
> Memerlukan [Node.js](https://nodejs.org) versi berapa pun. Tidak butuh `npm install`.

### <img src="https://skillicons.dev/icons?i=php" width="30" height="30" valign="middle" /> PHP
```bash
php -S localhost:969 server.php
```
> Memerlukan PHP 7.4+. Sudah tersedia di macOS, banyak Linux, dan hosting shared.

### <img src="https://skillicons.dev/icons?i=go" width="30" height="30" valign="middle" /> Go
```bash
# Langsung jalankan (tanpa build):
go run server.go

# Atau build dulu jadi binary:
go build -o server server.go
./server
```
> Memerlukan [Go](https://go.dev) 1.18+. Tanpa dependensi eksternal.

### <img src="https://skillicons.dev/icons?i=docker" width="30" height="30" valign="middle" /> Docker
```bash
# Build & jalankan (background):
docker compose up -d

# Lihat log:
docker compose logs -f

# Hentikan:
docker compose down
```
> Data `todo.txt` tetap tersimpan di host meski container dihapus.

Setelah server berjalan, buka browser ke:

```
http://localhost:969
```

---

## 📄 Format `todo.txt`

Setiap baris mewakili satu tugas dengan format:

```
<tugas>,<selesai>,<catatan>
```

| Kolom | Tipe | Keterangan |
|---|---|---|
| `tugas` | string | Nama tugas |
| `selesai` | `true` / `false` | Status penyelesaian |
| `catatan` | string | Catatan opsional |

**Contoh:**

```
Beli sayuran,false,Di pasar Minggu pagi
Kerjakan laporan,true,
Meeting mingguan,false,Zoom jam 10
```

> **Catatan:** Karakter koma (`,`) dalam teks tugas/catatan akan otomatis diganti titik koma (`;`) saat disimpan agar tidak konflik dengan delimiter.

---

## 🔌 API Reference

Semua server mengimplementasikan endpoint yang sama di `http://localhost:969`:

| Method | Endpoint | Fungsi |
|---|---|---|
| `GET` | `/` | Serve `index.html` |
| `GET` | `/api/todos` | Ambil semua tugas |
| `POST` | `/api/todos` | Tambah tugas baru |
| `POST` | `/api/todos/:id/toggle` | Toggle selesai/belum |
| `PUT` | `/api/todos/:id` | Edit tugas |
| `DELETE` | `/api/todos/:id` | Hapus tugas |

### Contoh request

**Tambah tugas baru:**
```bash
curl -X POST http://localhost:969/api/todos \
  -H "Content-Type: application/json" \
  -d '{"task": "Beli kopi", "note": "Arabica"}'
```

**Toggle selesai:**
```bash
curl -X POST http://localhost:969/api/todos/0/toggle \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Edit tugas:**
```bash
curl -X PUT http://localhost:969/api/todos/0 \
  -H "Content-Type: application/json" \
  -d '{"task": "Beli kopi decaf", "note": "Toko sebelah"}'
```

**Hapus tugas:**
```bash
curl -X DELETE http://localhost:969/api/todos/0
```

---

## ✨ Fitur

- ✅ Tambah tugas dengan catatan opsional
- ✅ Centang / batalkan tugas (toggle)
- ✅ Edit tugas secara inline
- ✅ Hapus tugas
- ✅ Filter: Semua / Belum / Selesai
- ✅ Progress bar & statistik
- ✅ Semua perubahan langsung tersimpan ke `todo.txt`
- ✅ Tidak butuh database atau framework

---

## <img src="https://skillicons.dev/icons?i=docker" width="30" height="30" valign="middle" /> Kustomisasi Docker

`Dockerfile` dan `docker-compose.yml` secara default menggunakan **Python**. Untuk mengganti ke runtime lain, buka file tersebut dan ikuti komentar di dalamnya — tinggal uncomment blok yang diinginkan.

**Ganti ke Node.js di `docker-compose.yml`:**
```yaml
services:
  todo:
    image: node:22-alpine
    working_dir: /app
    command: node server.js
    ports:
      - "969:969"
    volumes:
      - .:/app
```

---

## 📋 Persyaratan Sistem

| Runtime | Versi Minimum |
|---|---|
| Python | 3.x |
| Node.js | 12+ |
| PHP | 7.4+ |
| Go | 1.18+ |
| Docker | 20+ |

Semua server menggunakan **library standar bawaan** — tidak ada `npm install`, `pip install`, `composer install`, atau `go get` yang diperlukan.

---

## 📜 Lisensi

Bebas digunakan dan dimodifikasi untuk keperluan apapun.