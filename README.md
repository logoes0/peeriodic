# 📄 Collaborative Real-Time Editor

A full-stack real-time collaborative editor that enables users to create, share, and co-edit text documents in isolated rooms. Built with **Go** for the backend, **React** for the frontend, and **WebSockets** for live collaboration.

## 🚀 Features

- 🔗 Room-based document collaboration with shareable links
- 📡 Real-time sync using WebSockets (Gorilla WebSocket + React)
- 💾 Manual and autosave to PostgreSQL for persistence across restarts
- 🧠 Rename, delete, and revisit previous rooms (stored in localStorage)
- ✅ Clean UI with real-time updates, saving indicators, and clipboard copy
- 🔒 Planned support for Firebase authentication (per-user room ownership)

## 🛠️ Tech Stack

| Layer     | Technology            |
|-----------|------------------------|
| Frontend  | React.js, JavaScript   |
| Backend   | Go (net/http, Gorilla) |
| Database  | PostgreSQL             |
| Realtime  | WebSocket (Gorilla)    |
| Hosting   | Local                  |

## ⚙️ Getting Started

### 🧱 Backend (Go)

```bash
make run-be
```

> Runs backend on http://localhost:5000

### 🌐 Frontend (React)

```bash
make run-fe
```

> Runs frontend on http://localhost:3000

## 🔄 API Endpoints

- `GET  /api/rooms?id={roomId}` – fetch room content
- `POST /api/save?room={roomId}` – save content to DB
- `GET  /ws?room={roomId}` – open WebSocket connection

## 🖼️ Screenshots

| Homepage | Editor |
|----------|--------|
| Room creation, rename, delete | Live collaboration, autosave |
