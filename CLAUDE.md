# CLAUDE.md — AI Assistant Guide for go-f1

## Project Overview

**go-f1** is a live Formula 1 data streaming application written in Go. It connects to the official F1 live timing service (`livetiming.formula1.com`) via the SignalR WebSocket protocol, subscribes to all available timing and telemetry topics, and logs all incoming data to both the console and a file.

**Module:** `otterbyte.co.uk/go-f1`
**Go version:** 1.24+

---

## Repository Structure

```
go-f1/
├── main.go               # Entry point: SignalR connection, subscriptions, message loop
├── openf1_structs.go     # Go structs for F1 live timing data types
├── go.mod                # Module definition and direct dependencies
├── go.sum                # Dependency checksums
├── .gitignore            # Standard Go gitignore
├── example.json          # Sample payload (for reference/testing)
└── f1_live_data.log      # Append-only log of received messages (git-ignored)
```

All production code lives in the `main` package — no subdirectories or internal packages yet.

---

## Key Dependencies

| Package | Purpose |
|---|---|
| `github.com/gorilla/websocket` | WebSocket client for the SignalR connection |
| `github.com/joho/godotenv` | Load environment config from `.env` |

---

## Configuration

The application reads from environment variables. Create a `.env` file in the repo root (it is git-ignored):

```env
F1_AUTH_TOKEN=<optional F1 TV Pro bearer token>
```

`F1_AUTH_TOKEN` is optional. Without it the app connects anonymously and receives all public timing data. Providing a valid F1 TV Pro token may unlock additional streams (e.g. detailed telemetry).

---

## Build and Run

```bash
# Build
go build -o go-f1 .

# Run
./go-f1

# Build and run in one step
go run .
```

Stop the application with `Ctrl+C` (SIGINT) or `SIGTERM` — both trigger a graceful WebSocket close.

---

## Connection Flow

The official F1 live timing service uses the **old ASP.NET SignalR protocol** (not SignalR Core) over WebSockets. The connection sequence is:

1. **Negotiate** — `GET https://livetiming.formula1.com/signalr/negotiate?connectionData=...&clientProtocol=1.5`
   Returns a `ConnectionToken` and sets session cookies.
2. **Connect** — `wss://livetiming.formula1.com/signalr/connect?...&connectionToken=<TOKEN>`
   Must include headers `User-Agent: BestHTTP` and `Accept-Encoding: gzip,identity` (case-sensitive).
3. **Subscribe** — Send JSON over the WebSocket:
   ```json
   {"H":"Streaming","M":"Subscribe","A":[["Heartbeat","TimingData",...]],  "I":1}
   ```
4. **Receive** — Messages arrive as `{"C":"...","M":[{"H":"Streaming","M":"feed","A":["TopicName",{...},"timestamp"]}]}`

On connection loss the app waits 5 seconds and re-negotiates automatically.

---

## Topics

| Topic | Data |
|---|---|
| `Heartbeat` | Server keep-alive with UTC timestamp |
| `CarData.z` | Per-driver telemetry snapshot (compressed) |
| `Position.z` | Per-driver X/Y/Z track position (compressed) |
| `ExtrapolatedClock` | Session clock with extrapolation flag |
| `TopThree` | Current top-3 leaderboard |
| `RcmSeries` | Race control message series |
| `TimingStats` | Sector and speed trap bests |
| `TimingAppData` | Tyre and stint timing overlays |
| `WeatherData` | Temperature, humidity, rainfall, wind |
| `TrackStatus` | Track condition (AllClear, Yellow, SC, Red, VSC) |
| `DriverList` | Driver metadata map keyed by racing number |
| `RaceControlMessages` | Flags, safety car, penalty messages |
| `SessionInfo` | Meeting, circuit, session type and times |
| `SessionData` | Sector and qualifying session state |
| `LapCount` | Current and total lap count |
| `TimingData` | Live per-driver lap/sector timing (incremental patch) |
| `TeamRadio` | Radio recording URLs |
| `AudioStreams` | Audio stream metadata |
| `ContentStreams` | Video/content stream metadata |

### Compressed topics

`CarData.z` and `Position.z` payloads are **zlib-compressed and base64-encoded**. The `decompressZlib()` function in `main.go` handles decoding before the JSON is logged or parsed.

---

## Data Structures (`openf1_structs.go`)

Structs for the well-known topics:

- `Heartbeat` — UTC keep-alive timestamp
- `TrackStatus` — status code + message string
- `WeatherData` — all values as strings (as sent by the feed)
- `LapCount` — current and total laps
- `ExtrapolatedClock` — remaining time and extrapolation flag
- `DriverInfo` / `DriverList` — per-driver metadata, keyed map
- `RaceControlMessage` / `RaceControlMessages` — flag and SC events
- `TeamRadioCapture` / `TeamRadio` — radio recording paths
- `SessionInfo` — meeting, circuit, session type
- `CarDataPayload` — decompressed telemetry (RPM, speed, gear, throttle, brake, DRS)
- `PositionPayload` — decompressed X/Y/Z car positions

**Important:** `TimingData` and `TimingAppData` arrive as incremental JSON patches (not full state), so they do not have simple struct representations. They are currently logged as raw JSON.

---

## Code Conventions

### Error Handling
- **`log.Fatalf()`** — used for unrecoverable startup failures (e.g., cannot open log file)
- **`log.Printf()`** — used for recoverable runtime errors (parse failures, write errors)
- Connection errors trigger a reconnect loop, not a crash

### Naming
- Constants: `camelCase` at package level (e.g., `negotiateURL`, `hubName`)
- Package-level variables: `PascalCase` (e.g., `AllTopics`)
- Functions: `camelCase` (e.g., `negotiate`, `connectWS`, `handleMessage`)
- Structs and exported types: `PascalCase`
- JSON field tags: match the F1 live timing feed's `PascalCase` keys exactly

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

---

## Known Issues / Gaps

1. **`TimingData` is not parsed** — it arrives as an incremental JSON patch; no struct or merge logic exists yet.
2. **No subscription reuse across reconnects** — each reconnect re-negotiates from scratch (necessary for the old SignalR protocol).
3. **No per-topic dispatch** — all topics share a single handler; data is logged but not routed or processed per type.

---

## Planned / Future Work

- HTTP API (REST/WebSocket) using Gin to expose live F1 data to clients
- Per-topic message handlers and data routing
- JSON patch merging for `TimingData` to maintain full driver state
- Unit and integration tests

---

## Git Workflow

- The `master` branch is the main branch.
- Feature and task branches follow the convention `claude/<description>-<id>`.
- Commit messages are short and imperative (e.g., `Add TimingData patch merging`).
