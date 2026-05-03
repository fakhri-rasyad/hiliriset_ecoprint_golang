# Hiliriset Ecoprint API

Backend API for monitoring the ecoprint fabric boiling process in real time.

## Tech Stack

- **Go** + Fiber v3
- **PostgreSQL** + GORM
- **MQTT** (Eclipse Paho) — ESP device communication
- **WebSocket** — live session monitoring
- **JWT** — authentication
- **Swagger** — API documentation

---

## Architecture

```
┌─────────────────────┐
│   Client (Android)  │
└────────┬────────────┘
         │ REST API (JWT)
         │ WebSocket (live telemetry)
         ▼
┌─────────────────────┐
│  Backend (Go+Fiber) │
└──────┬──────────────┘
       │
       ├──── PostgreSQL (GORM)
       │
       └──── MQTT Broker (Mosquitto)
                    │
             ┌──────┴──────┐
             │ ESP Device  │
             │ → telemetry │
             │ ← start/stop│
             └─────────────┘
```

---

## Project Structure

```
hiliriset_ecoprint_golang/
├── config/             # DB and env config
├── controllers/        # HTTP handlers
├── services/           # Business logic
├── repositories/       # Database queries
├── models/             # GORM, Base, Request, Response structs
├── routes/             # Route registration
├── mqtt_package/       # MQTT client and handler
├── websocket_utils/    # WebSocket hub and handler
├── utils/              # JWT, response helpers, password hashing
├── migrations/         # SQL migration files
├── docs/               # Swagger generated files
└── main.go
```

---

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Mosquitto MQTT broker (or any MQTT broker)
- golang-migrate CLI
- swag CLI (for generating Swagger docs)

---

### Installing golang-migrate

**Windows (scoop)**
scoop install migrate

**Mac (homebrew)**
brew install golang-migrate

**Linux**
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

Or via Go install:
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

For other installation options see: https://github.com/golang-migrate/migrate

---

### Installing swag

go install github.com/swaggo/swag/cmd/swag@latest

---

## Environment Variables

Create a `.env` file in the root directory:

```
APP_PORT=3000
JWT_SECRET=your_jwt_secret

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecoprint_golang

MQTT_HOST=localhost
MQTT_PORT=1883
MQTT_USERNAME=
MQTT_PASSWORD=
```

---

## Running Locally

1. Clone the repository
   git clone https://github.com/fakhri-rasyad/hiliriset_ecoprint_golang.git
   cd hiliriset_ecoprint_golang

2. Install dependencies
   go mod tidy

3. Generate Swagger docs
   swag init

4. Set up the database
   createdb ecoprint_golang

5. Run migrations
   migrate -path database/migrations -database "postgres://user:password@localhost:5432/ecoprint_golang?sslmode=disable" up

6. Start the MQTT broker
   mosquitto -v

7. Run the server
   go run main.go

8. Open Swagger
   http://localhost:3000/swagger/index.html

---

## API Overview

### Auth

| Method | Endpoint          | Description       |
| ------ | ----------------- | ----------------- |
| POST   | /v1/auth/register | Register new user |
| POST   | /v1/auth/login    | Login, get token  |

### ESP Devices

| Method | Endpoint                | Description    |
| ------ | ----------------------- | -------------- |
| POST   | /api/v1/esps            | Register ESP   |
| GET    | /api/v1/esps            | List user ESPs |
| GET    | /api/v1/esps/:public_id | ESP detail     |
| DELETE | /api/v1/esps/:public_id | Delete ESP     |

### Kompors

| Method | Endpoint                   | Description       |
| ------ | -------------------------- | ----------------- |
| POST   | /api/v1/kompors            | Add kompor        |
| GET    | /api/v1/kompors            | List user kompors |
| GET    | /api/v1/kompors/:public_id | Kompor detail     |
| DELETE | /api/v1/kompors/:public_id | Delete kompor     |

### Boiling Sessions

| Method | Endpoint                             | Description                |
| ------ | ------------------------------------ | -------------------------- |
| POST   | /api/v1/sessions                     | Create session (2hr timer) |
| GET    | /api/v1/sessions                     | List user sessions         |
| GET    | /api/v1/sessions/:public_id          | Session detail             |
| GET    | /api/v1/sessions/:session_id/records | Session records            |

### WebSocket

| Endpoint                         | Description           |
| -------------------------------- | --------------------- |
| ws://host/api/v1/sessions/:id/ws | Live telemetry stream |

#### WebSocket message format

// Telemetry (every 2 seconds)
{"air_temp": 36.5, "water_temp": 98.2, "humidity": 72.1}

// Session finished
{"event": "finished"}

---

## MQTT Topics

| Topic                     | Direction     | Payload                                         |
| ------------------------- | ------------- | ----------------------------------------------- |
| esp/{public_id}/cmd       | Backend → ESP | "start" or "stop"                               |
| esp/{public_id}/telemetry | ESP → Backend | {"air_temp": 0, "water_temp": 0, "humidity": 0} |

---

## Session Flow

1. Client creates a session → backend sets finished_at to now + 2 hours (For demonstration purposes, this project only uses 5 minutes rather than 2 hours)
2. Backend publishes "start" to esp/{public_id}/cmd
3. ESP sends telemetry every 2 seconds to esp/{public_id}/telemetry
4. Backend saves each record and broadcasts to connected WebSocket clients
5. When finished_at is reached, backend publishes "stop" to ESP
6. WebSocket clients receive {"event": "finished"} and close the connection

---

## License

MIT
