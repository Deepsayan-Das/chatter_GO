# chatter_GO

A real-time chat backend built with Go, featuring WebSocket-based messaging, JWT authentication, and PostgreSQL persistence.

## Overview

chatter_GO provides a REST API for user management and room administration, combined with a persistent WebSocket connection for real-time room messaging and direct messaging between users.

## Technology Stack

| Component        | Technology                  |
|------------------|-----------------------------|
| Language         | Go 1.22                     |
| HTTP Framework   | Gin                         |
| WebSocket        | gorilla/websocket           |
| Database         | PostgreSQL 16               |
| Auth             | JWT (HS256)                 |
| Password Hashing | bcrypt                      |
| DB Driver        | pgx/v5                      |
| Containerization | Docker / Docker Compose     |

## Project Structure

```
chatter_GO/
├── cmd/
│   └── server/
│       └── main.go             # Entry point, route registration
├── internal/
│   ├── db/                     # Database connection
│   ├── handlers/               # HTTP request handlers
│   │   ├── auth.handler.go
│   │   ├── room.handler.go
│   │   └── message.handler.go
│   ├── middleware/
│   │   └── auth.middleware.go  # JWT authentication middleware
│   ├── models/                 # Shared data models
│   ├── services/               # Business logic layer
│   │   ├── auth.services.go
│   │   ├── room.service.go
│   │   ├── message.service.go
│   │   └── dm.service.go
│   ├── utils/
│   │   ├── hash.go             # bcrypt helpers
│   │   ├── jwt.go              # JWT generation and validation
│   │   └── utils_test.go       # Unit tests
│   └── websocket/
│       ├── hub.go              # Central client registry
│       ├── manager.go          # Room-scoped connection management
│       ├── client.go           # Per-connection read/write pumps
│       ├── websocket.handler.go
│       ├── websockets.models.go
│       └── websocket_test.go   # Unit tests
├── migrations/                 # SQL migration files
├── dockerfile
├── docker-compose.yml
└── .env
```

## Prerequisites

- Go 1.22 or later
- PostgreSQL 16 (or Docker)

## Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/chatapp
JWT_SECRET=your_secret_key
PORT=8080
```

## Running Locally

**Without Docker:**

```bash
go mod tidy
go run ./cmd/server/main.go
```

**With Docker Compose:**

```bash
docker-compose up --build
```

This starts both the PostgreSQL database and the backend server. Apply database migrations manually after the containers are running.

> **Note:** The Dockerfile uses a multi-stage build. The binary is compiled in a `golang:1.22-alpine` builder stage and copied into a minimal `alpine:3.19` runtime image, producing a final image of approximately 10 MB.

## Testing

The following packages have unit test coverage:

| Package              | Coverage                                           |
|----------------------|----------------------------------------------------|
| `internal/utils`     | JWT generation, token validation, password hashing |
| `internal/websocket` | Manager join/leave/broadcast, Hub initialization   |

Run all tests:

```bash
go test ./...
```

Run with verbose output:

```bash
go test ./internal/utils/... ./internal/websocket/... -v
```

## REST API Reference

All protected routes require an `Authorization: Bearer <token>` header.

### Authentication

| Method | Route       | Auth | Description      |
|--------|-------------|------|------------------|
| POST   | /register   | No   | Create account   |
| POST   | /login      | No   | Obtain JWT token |

**Register request body:**
```json
{ "username": "alice", "email": "alice@example.com", "password": "secret" }
```

**Login response:**
```json
{ "token": "<jwt>" }
```

### Rooms

| Method | Route          | Auth | Description                        |
|--------|----------------|------|------------------------------------|
| POST   | /rooms/create  | Yes  | Create a new room                  |
| POST   | /rooms/join    | Yes  | Join an existing room              |
| POST   | /rooms/leave   | Yes  | Leave a room                       |
| GET    | /rooms/search  | Yes  | Search rooms by name (`?name=...`) |
| GET    | /rooms/my      | Yes  | List rooms the user has joined     |

### Messages (REST)

| Method | Route                  | Auth | Description                   |
|--------|------------------------|------|-------------------------------|
| POST   | /messages              | Yes  | Send a message (persisted)    |
| GET    | /rooms/:id/messages    | Yes  | Fetch last 50 messages        |
| PATCH  | /messages/:id          | Yes  | Edit own message              |
| DELETE | /messages/:id          | Yes  | Soft-delete own message       |

## WebSocket Protocol

Connect to the WebSocket endpoint with a valid JWT passed as a query parameter:

```
ws://localhost:8080/ws?token=<jwt>
```

Once connected, send JSON frames to the server:

### Message Types

| Type value | Action       | Required fields              |
|------------|--------------|------------------------------|
| `0`        | join_room    | `room_id`                    |
| `1`        | leave_room   | `room_id`                    |
| `2`        | send_message | `room_id`, `content`         |
| `3`        | send_dm      | `receiver_id`, `content`     |

**Example — join a room:**
```json
{ "type": 0, "room_id": 1 }
```

**Example — send a room message:**
```json
{ "type": 2, "room_id": 1, "content": "Hello, room!" }
```

**Example — send a direct message:**
```json
{ "type": 3, "receiver_id": 42, "content": "Hey!" }
```

### Server-pushed Events

**New room message:**
```json
{
  "type": "new_message",
  "message_id": 101,
  "room_id": 1,
  "user_id": 7,
  "content": "Hello, room!"
}
```

**New direct message:**
```json
{
  "type": "new_dm",
  "message_id": 55,
  "sender_id": 7,
  "receiver_id": 42,
  "content": "Hey!"
}
```

### Connection Behavior

- The server sends a WebSocket ping every 54 seconds; the client must respond with a pong.
- The connection is closed if no pong is received within 60 seconds.
- Write operations time out after 10 seconds.
- Maximum inbound frame size is 4096 bytes.

## Architecture Notes

- **Hub**: Maintains the global registry of connected clients, keyed both by pointer and by user ID for direct DM delivery.
- **Manager**: Tracks per-room client subscriptions. All map access is protected by a `sync.RWMutex`.
- **ReadPump / WritePump**: Each client runs two goroutines. ReadPump handles inbound frames and dispatches to Manager or Hub. WritePump drains the client's `Send` channel and handles ping/pong keepalive.
- **Persistence**: Room messages are persisted to PostgreSQL via `services.CreateMessage` before being broadcast. DMs are likewise persisted via `services.CreateDM`.