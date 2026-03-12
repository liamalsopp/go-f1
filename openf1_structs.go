package main

import "time"

// ---------------------------------------------------------------------
// F1 Live Timing Data Structures
// Source: https://livetiming.formula1.com/signalr (reverse-engineered)
// Reference: https://github.com/theOehrly/Fast-F1
// ---------------------------------------------------------------------

// Heartbeat is a periodic keep-alive message from the server.
// Topic: Heartbeat
type Heartbeat struct {
	Utc time.Time `json:"Utc"`
}

// TrackStatus reports the current track condition.
// Topic: TrackStatus
// Status codes: 1=AllClear, 2=Yellow, 4=SafetyCar, 5=Red, 6=VSC, 7=SCEnding
type TrackStatus struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

// WeatherData reports environmental conditions at the circuit.
// Topic: WeatherData
// Note: All numeric values arrive as strings from the live timing feed.
type WeatherData struct {
	AirTemp       string `json:"AirTemp"`
	Humidity      string `json:"Humidity"`
	Pressure      string `json:"Pressure"`
	Rainfall      string `json:"Rainfall"`
	TrackTemp     string `json:"TrackTemp"`
	WindDirection string `json:"WindDirection"`
	WindSpeed     string `json:"WindSpeed"`
}

// LapCount reports the current and total lap count.
// Topic: LapCount
type LapCount struct {
	CurrentLap int `json:"CurrentLap"`
	TotalLaps  int `json:"TotalLaps"`
}

// ExtrapolatedClock reports the session clock with extrapolation flag.
// Topic: ExtrapolatedClock
type ExtrapolatedClock struct {
	Utc           string `json:"Utc"`
	Remaining     string `json:"Remaining"`
	Extrapolating bool   `json:"Extrapolating"`
}

// DriverInfo holds metadata for a single driver.
// Used within DriverList (keyed by racing number string).
type DriverInfo struct {
	RacingNumber  string `json:"RacingNumber"`
	BroadcastName string `json:"BroadcastName"`
	FullName      string `json:"FullName"`
	Tla           string `json:"Tla"`          // Three-letter abbreviation, e.g. "VER"
	Line          int    `json:"Line"`         // Current race position line
	TeamName      string `json:"TeamName"`
	TeamColour    string `json:"TeamColour"`   // Hex colour, e.g. "3671C6"
	FirstName     string `json:"FirstName"`
	LastName      string `json:"LastName"`
	HeadshotUrl   string `json:"HeadshotUrl"`
	CountryCode   string `json:"CountryCode"`
}

// DriverList is a map of racing number (string) → DriverInfo.
// Topic: DriverList
type DriverList map[string]DriverInfo

// RaceControlMessage is a single race control event (flag, safety car, etc.).
type RaceControlMessage struct {
	Utc          string `json:"Utc"`
	Category     string `json:"Category"`     // "Flag", "SafetyCar", "Drs", etc.
	Message      string `json:"Message"`
	Flag         string `json:"Flag"`         // "GREEN", "YELLOW", "RED", etc.
	Scope        string `json:"Scope"`        // "Track", "Sector", "Driver"
	Sector       int    `json:"Sector"`
	RacingNumber string `json:"RacingNumber"` // Set when scope is "Driver"
}

// RaceControlMessages wraps the keyed map of race control events.
// Topic: RaceControlMessages
// Messages is keyed by a string index ("0", "1", ...) not an array.
type RaceControlMessages struct {
	Messages map[string]RaceControlMessage `json:"Messages"`
}

// TeamRadioCapture is a single radio transmission recording.
type TeamRadioCapture struct {
	Utc          string `json:"Utc"`
	RacingNumber string `json:"RacingNumber"`
	Path         string `json:"Path"` // Relative URL to the .mp3 file
}

// TeamRadio wraps the keyed map of radio captures.
// Topic: TeamRadio
type TeamRadio struct {
	Captures map[string]TeamRadioCapture `json:"Captures"`
}

// SessionCountry holds country metadata within a SessionMeeting.
type SessionCountry struct {
	Key  int    `json:"Key"`
	Code string `json:"Code"`
	Name string `json:"Name"`
}

// SessionCircuit holds circuit metadata within a SessionMeeting.
type SessionCircuit struct {
	Key       int    `json:"Key"`
	ShortName string `json:"ShortName"`
}

// SessionMeeting holds race weekend metadata.
type SessionMeeting struct {
	Key          int            `json:"Key"`
	Name         string         `json:"Name"`
	OfficialName string         `json:"OfficialName"`
	Location     string         `json:"Location"`
	Country      SessionCountry `json:"Country"`
	Circuit      SessionCircuit `json:"Circuit"`
}

// SessionInfo describes the current session.
// Topic: SessionInfo
type SessionInfo struct {
	Meeting   SessionMeeting `json:"Meeting"`
	Key       int            `json:"Key"`
	Type      string         `json:"Type"`       // "Practice", "Qualifying", "Race"
	Name      string         `json:"Name"`       // "Race", "Sprint", "Practice 1", etc.
	StartDate string         `json:"StartDate"`
	EndDate   string         `json:"EndDate"`
	GmtOffset string         `json:"GmtOffset"`
	Path      string         `json:"Path"`
}

// ---------------------------------------------------------------------
// Compressed topics: CarData.z and Position.z
// These are zlib-compressed, base64-encoded JSON blobs.
// After decompression the payloads match the structs below.
// ---------------------------------------------------------------------

// CarChannels holds the telemetry channel values for one car.
// Channel keys (as ints in JSON):
//
//	0  = RPM
//	2  = Speed (km/h)
//	3  = Gear (0=N)
//	4  = Throttle (0-100)
//	5  = Brake (0 or 1)
//	45 = DRS status
type CarChannels map[string]int

// CarEntry is a snapshot of all cars' telemetry at a single instant.
type CarEntry struct {
	Utc  string                 `json:"Utc"`
	Cars map[string]CarChannels `json:"Cars"` // keyed by racing number string
}

// CarDataPayload is the decompressed content of the CarData.z topic.
type CarDataPayload struct {
	Entries []CarEntry `json:"Entries"`
}

// PositionEntry holds the track position for one car at an instant.
type PositionEntry struct {
	Status string  `json:"Status"` // "OnTrack", "OffTrack", "Pit"
	X      float64 `json:"X"`
	Y      float64 `json:"Y"`
	Z      float64 `json:"Z"`
}

// PositionSnapshot holds all cars' positions at a single timestamp.
type PositionSnapshot struct {
	Timestamp string                    `json:"Timestamp"`
	Entries   map[string]PositionEntry  `json:"Entries"` // keyed by racing number string
}

// PositionPayload is the decompressed content of the Position.z topic.
type PositionPayload struct {
	Position []PositionSnapshot `json:"Position"`
}
