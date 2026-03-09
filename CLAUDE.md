# CLAUDE.md — AI Assistant Guide for go-f1

## Project Overview

**go-f1** is a live Formula 1 data streaming application written in Go. It connects to the [OpenF1](https://openf1.org) MQTT broker and captures real-time racing telemetry, logging all incoming data to both the console and a file.

**Module:** `otterbyte.co.uk/go-f1`
**Go version:** 1.24+

---

## Repository Structure

```
go-f1/
├── main.go               # Entry point: MQTT client setup, subscriptions, signal handling
├── openf1_structs.go     # Go structs for all 13 OpenF1 data types
├── go.mod                # Module definition and direct dependencies
├── go.sum                # Dependency checksums
├── .gitignore            # Standard Go gitignore
├── example.json          # Sample MQTT payload (for reference/testing)
└── f1_live_data.log      # Append-only log of received MQTT messages (git-ignored)
```

All production code lives in the `main` package — no subdirectories or internal packages yet.

---

## Key Dependencies

| Package | Purpose |
|---|---|
| `github.com/eclipse/paho.mqtt.golang` | MQTT client (core dependency) |
| `github.com/joho/godotenv` | Load environment config from `.env` |
| `github.com/gin-gonic/gin` | HTTP framework (imported but not yet used — planned) |
| `github.com/gorilla/websocket` | WebSocket transport for MQTT |

---

## Configuration

The application reads from environment variables. Create a `.env` file in the repo root (it is git-ignored):

```env
MQTT_BROKER=<broker URL>
MQTT_USER=<username>
OPENF1_TOKEN=<authentication token>
```

The `.env` file is optional — the application logs a warning if absent and falls back to whatever environment variables are already set.

---

## Build and Run

```bash
# Build
go build -o go-f1 .

# Run (requires .env or exported env vars)
./go-f1

# Build and run in one step
go run .
```

Stop the application with `Ctrl+C` (SIGINT) or `SIGTERM` — both trigger a graceful MQTT disconnect.

---

## MQTT Topics

The application subscribes to 13 OpenF1 topics at QoS 0:

| Topic | Data |
|---|---|
| `v1/car_data` | Speed, RPM, throttle, brake, gear, DRS |
| `v1/position` | Driver race position |
| `v1/location` | 3D track position (X, Y, Z) |
| `v1/weather` | Temperature, humidity, rainfall, wind |
| `v1/session_status` | Session state changes |
| `v1/laps` | Lap timing and sector splits |
| `v1/intervals` | Live gaps between drivers |
| `v1/pit` | Pit stop events |
| `v1/race_control` | Flags, safety car messages |
| `v1/stints` | Tyre compound and age |
| `v1/team_radio` | Radio message URLs |
| `v1/drivers` | Driver metadata |
| `v1/overtakes` | Overtake events |

---

## Data Structures (`openf1_structs.go`)

Each struct maps 1:1 to an OpenF1 topic payload using `json` struct tags:

- `CarData` — telemetry snapshot per driver
- `Driver` — driver number, name, team, colours, headshot URL
- `Interval` — live gap/interval between drivers
- `Lap` — full lap record with sector times and speed traps
- `Location` — X/Y/Z car coordinates on track
- `Meeting` — race weekend metadata
- `Pit` — pit stop event with lap and duration
- `Position` — current race position
- `RaceControl` — flag status and messages
- `Session` — session type, circuit, UTC offset
- `Stint` — tyre compound, lap counts, fresh/used state
- `TeamRadio` — URL to radio recording
- `Weather` — environmental conditions

All timestamp fields use `time.Time` with UTC-based JSON unmarshaling.

---

## Code Conventions

### Error Handling
- **`panic()`** — used for unrecoverable startup failures (e.g., cannot connect to broker)
- **`log.Fatal()`** — used for file I/O failures at startup
- **`log.Printf()`** — used for runtime errors that should be logged but not terminate the process
- JSON unmarshal errors in the message handler are currently logged silently (known gap)

### Naming
- Constants: `PascalCase` (e.g., `ClientID`, `AllTopics`)
- Package-level variables: `PascalCase` (e.g., `Broker`, `Token`, `Username`)
- Functions: `camelCase` (e.g., `connectHandler`, `connectLostHandler`)
- Structs and exported types: `PascalCase`
- JSON field tags: `snake_case` to match the OpenF1 API format

### Formatting
- All code must be `gofmt`-formatted. Run `gofmt -w .` before committing.

---

## Testing

There are currently no tests in the repository. When adding tests:

```bash
go test ./...          # run all tests
go test -v ./...       # verbose output
go test -cover ./...   # with coverage
```

The mocking framework `go.uber.org/mock` is available as an indirect dependency and should be used for any MQTT client or handler mocks.

---

## Known Issues / Gaps

1. **Unused import:** `gin-gonic/gin` is declared in `go.mod` but not imported anywhere — it is reserved for a planned HTTP API layer.
2. **No reconnect logic:** If the MQTT broker drops the connection, the application does not attempt to reconnect automatically.
3. **Silent JSON errors:** Malformed MQTT payloads log an error but are otherwise dropped with no further handling.
4. **No topic-specific dispatch:** All 13 topics share a single message handler; there is no per-topic processing or routing yet.

---

## Planned / Future Work

- HTTP API (REST/WebSocket) using Gin to expose live F1 data to clients
- Per-topic message handlers and data routing
- Automatic broker reconnection
- Unit and integration tests

---

## Git Workflow

- The `master` branch is the main branch.
- Feature and task branches follow the convention `claude/<description>-<id>`.
- Commit messages are short and imperative (e.g., `Add MQTT reconnection logic`).
