package main

import "time"

// ---------------------------------------------------------------------
// OpenF1 Data Structures
// Source: https://openf1.org/#api-endpoints
// ---------------------------------------------------------------------

// CarData represents telemetry data (Speed, RPM, etc.)
// Endpoint: /v1/car_data
type CarData struct {
	Date         time.Time `json:"date"`          // UTC timestamp
	DriverNumber int       `json:"driver_number"` // Driver number (e.g., 44)
	RPM          int       `json:"rpm"`           // Engine RPM
	Speed        int       `json:"speed"`         // Speed in km/h
	NGear        int       `json:"n_gear"`        // Current gear (0 = Neutral/Reverse)
	Throttle     int       `json:"throttle"`      // Throttle percentage (0-100)
	Brake        int       `json:"brake"`         // Brake (0 or 100 usually, sometimes granular)
	DRS          int       `json:"drs"`           // DRS status (0-14, see docs)
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// Driver represents driver information for a session
// Endpoint: /v1/drivers
type Driver struct {
	DriverNumber  int    `json:"driver_number"`
	BroadcastName string `json:"broadcast_name"` // e.g., "L HAMILTON"
	FullName      string `json:"full_name"`      // e.g., "Lewis Hamilton"
	NameAcronym   string `json:"name_acronym"`   // e.g., "HAM"
	TeamName      string `json:"team_name"`      // e.g., "Mercedes"
	TeamColour    string `json:"team_colour"`    // Hex color e.g., "00D2BE"
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	HeadshotURL   string `json:"headshot_url"`
	CountryCode   string `json:"country_code"`
	SessionKey    int    `json:"session_key"`
	MeetingKey    int    `json:"meeting_key"`
}

// Interval represents the live gap between drivers
// Endpoint: /v1/intervals
type Interval struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	GapToLeader  float64   `json:"gap_to_leader"` // Seconds (can be null/0)
	Interval     float64   `json:"interval"`      // Gap to car ahead
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// Lap represents detailed lap timing data
// Endpoint: /v1/laps
type Lap struct {
	DateStart       time.Time `json:"date_start"`
	DriverNumber    int       `json:"driver_number"`
	LapNumber       int       `json:"lap_number"`
	LapDuration     float64   `json:"lap_duration"` // Total lap time in seconds
	IsPitOutLap     bool      `json:"is_pit_out_lap"`
	DurationSector1 float64   `json:"duration_sector_1"`
	DurationSector2 float64   `json:"duration_sector_2"`
	DurationSector3 float64   `json:"duration_sector_3"`
	SegmentsSector1 []int     `json:"segments_sector_1"` // Mini-sector colors
	SegmentsSector2 []int     `json:"segments_sector_2"`
	SegmentsSector3 []int     `json:"segments_sector_3"`
	I1Speed         int       `json:"i1_speed"` // Speed at Interval 1
	I2Speed         int       `json:"i2_speed"` // Speed at Interval 2
	STSpeed         int       `json:"st_speed"` // Speed Trap speed
	SessionKey      int       `json:"session_key"`
	MeetingKey      int       `json:"meeting_key"`
}

// Location represents the X, Y, Z coordinates of a car on track
// Endpoint: /v1/location
type Location struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	X            float64   `json:"x"`
	Y            float64   `json:"y"`
	Z            float64   `json:"z"`
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// Meeting represents a race weekend (Grand Prix)
// Endpoint: /v1/meetings
type Meeting struct {
	MeetingKey          int       `json:"meeting_key"`
	MeetingName         string    `json:"meeting_name"`          // e.g., "Bahrain Grand Prix"
	MeetingOfficialName string    `json:"meeting_official_name"` // Full sponsored name
	Location            string    `json:"location"`              // e.g., "Sakhir"
	CountryKey          int       `json:"country_key"`
	CountryCode         string    `json:"country_code"`
	CountryName         string    `json:"country_name"`
	CircuitKey          int       `json:"circuit_key"`
	CircuitShortName    string    `json:"circuit_short_name"`
	DateStart           time.Time `json:"date_start"`
	GMTOffset           string    `json:"gmt_offset"`
	Year                int       `json:"year"`
}

// Pit represents a pit stop event
// Endpoint: /v1/pit
type Pit struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	LapNumber    int       `json:"lap_number"`
	PitDuration  float64   `json:"pit_duration"` // Time stationary or total pit lane time
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// Position represents a driver's position on track
// Endpoint: /v1/position
type Position struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	Position     int       `json:"position"`
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// RaceControl represents flags (Yellow, Red, etc.) and messages
// Endpoint: /v1/race_control
type RaceControl struct {
	Date         time.Time `json:"date"`
	LapNumber    int       `json:"lap_number"`
	Category     string    `json:"category"` // "Flag", "SafetyCar", etc.
	Flag         string    `json:"flag"`     // "GREEN", "YELLOW", "DOUBLE YELLOW", "RED"
	Scope        string    `json:"scope"`    // "Track", "Sector", "Driver"
	Sector       int       `json:"sector"`   // Sector number (if scope is Sector)
	Message      string    `json:"message"`  // Display message
	DriverNumber int       `json:"driver_number"`
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// Session represents a specific session (FP1, Quali, Race)
// Endpoint: /v1/sessions
type Session struct {
	SessionKey       int       `json:"session_key"`
	SessionName      string    `json:"session_name"` // "Practice 1", "Race"
	SessionType      string    `json:"session_type"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
	GMTOffset        string    `json:"gmt_offset"`
	MeetingKey       int       `json:"meeting_key"`
	Location         string    `json:"location"`
	CountryKey       int       `json:"country_key"`
	CountryCode      string    `json:"country_code"`
	CountryName      string    `json:"country_name"`
	CircuitKey       int       `json:"circuit_key"`
	CircuitShortName string    `json:"circuit_short_name"`
	Year             int       `json:"year"`
}

// Stint represents a tyre stint
// Endpoint: /v1/stints
type Stint struct {
	StintNumber    int    `json:"stint_number"`
	DriverNumber   int    `json:"driver_number"`
	LapStart       int    `json:"lap_start"`
	LapEnd         int    `json:"lap_end"`
	Compound       string `json:"compound"`          // "SOFT", "MEDIUM", "HARD"
	TyreAgeAtStart int    `json:"tyre_age_at_start"` // Laps old
	SessionKey     int    `json:"session_key"`
	MeetingKey     int    `json:"meeting_key"`
}

// TeamRadio represents a radio communication URL
// Endpoint: /v1/team_radio
type TeamRadio struct {
	Date         time.Time `json:"date"`
	DriverNumber int       `json:"driver_number"`
	RecordingURL string    `json:"recording_url"` // .mp3 link
	SessionKey   int       `json:"session_key"`
	MeetingKey   int       `json:"meeting_key"`
}

// Weather represents weather data
// Endpoint: /v1/weather
type Weather struct {
	Date             time.Time `json:"date"`
	AirTemperature   float64   `json:"air_temperature"`
	TrackTemperature float64   `json:"track_temperature"`
	Humidity         float64   `json:"humidity"`
	Pressure         float64   `json:"pressure"`
	Rainfall         int       `json:"rainfall"` // 0 = No rain, 1 = Rain
	WindDirection    int       `json:"wind_direction"`
	WindSpeed        float64   `json:"wind_speed"`
	SessionKey       int       `json:"session_key"`
	MeetingKey       int       `json:"meeting_key"`
}
