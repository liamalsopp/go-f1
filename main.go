package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	negotiateURL   = "https://livetiming.formula1.com/signalr/negotiate"
	connectURL     = "wss://livetiming.formula1.com/signalr/connect"
	hubName        = "Streaming"
	clientProtocol = "1.5"
	connectionData = `[{"name":"Streaming"}]`
	reconnectDelay = 5 * time.Second
)

// AllTopics is the list of F1 live timing topics to subscribe to.
var AllTopics = []string{
	"Heartbeat",
	"CarData.z",
	"Position.z",
	"ExtrapolatedClock",
	"TopThree",
	"RcmSeries",
	"TimingStats",
	"TimingAppData",
	"WeatherData",
	"TrackStatus",
	"DriverList",
	"RaceControlMessages",
	"SessionInfo",
	"SessionData",
	"LapCount",
	"TimingData",
	"TeamRadio",
	"AudioStreams",
	"ContentStreams",
}

// negotiateResp holds the fields we need from the SignalR negotiate response.
type negotiateResp struct {
	ConnectionToken string `json:"ConnectionToken"`
}

// signalRMsg is the top-level envelope for incoming SignalR messages.
type signalRMsg struct {
	C string    `json:"C"` // cursor / message ID
	M []feedMsg `json:"M"` // array of hub method invocations
}

// feedMsg represents a single hub message from the Streaming hub.
type feedMsg struct {
	H string            `json:"H"` // hub name
	M string            `json:"M"` // method name ("feed")
	A []json.RawMessage `json:"A"` // arguments: [topic, data, timestamp]
}

// subscribeReq is the message sent to the server to subscribe to topics.
type subscribeReq struct {
	H string     `json:"H"`
	M string     `json:"M"`
	A [][]string `json:"A"`
	I int        `json:"I"`
}

// negotiate performs the SignalR negotiate handshake and returns the
// connection token and any cookies to forward to the WebSocket.
func negotiate(authToken string) (token string, cookie string, err error) {
	encoded := url.QueryEscape(connectionData)
	reqURL := fmt.Sprintf("%s?connectionData=%s&clientProtocol=%s", negotiateURL, encoded, clientProtocol)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("build negotiate request: %w", err)
	}
	req.Header.Set("User-Agent", "BestHTTP")
	req.Header.Set("Accept-Encoding", "gzip,identity")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("negotiate request: %w", err)
	}
	defer resp.Body.Close()

	var neg negotiateResp
	if err := json.NewDecoder(resp.Body).Decode(&neg); err != nil {
		return "", "", fmt.Errorf("decode negotiate response: %w", err)
	}

	var parts []string
	for _, c := range resp.Cookies() {
		parts = append(parts, c.Name+"="+c.Value)
	}
	return neg.ConnectionToken, strings.Join(parts, "; "), nil
}

// connectWS dials the SignalR WebSocket endpoint using the negotiated token.
func connectWS(token, cookie string) (*websocket.Conn, error) {
	encoded := url.QueryEscape(connectionData)
	encodedToken := url.QueryEscape(token)

	wsURL := fmt.Sprintf("%s?clientProtocol=%s&transport=webSockets&connectionToken=%s&connectionData=%s",
		connectURL, clientProtocol, encodedToken, encoded)

	headers := http.Header{}
	headers.Set("User-Agent", "BestHTTP")
	headers.Set("Accept-Encoding", "gzip,identity")
	if cookie != "" {
		headers.Set("Cookie", cookie)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	return conn, err
}

// subscribe sends the topic subscription message over the open WebSocket.
func subscribe(conn *websocket.Conn) error {
	return conn.WriteJSON(subscribeReq{
		H: hubName,
		M: "Subscribe",
		A: [][]string{AllTopics},
		I: 1,
	})
}

// decompressZlib base64-decodes and zlib-inflates a compressed payload
// (used for CarData.z and Position.z topics).
func decompressZlib(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	r, err := zlib.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return nil, fmt.Errorf("zlib reader: %w", err)
	}
	defer r.Close()
	return io.ReadAll(r)
}

// handleMessage processes a single SignalR feed message, writing the payload
// to the console and log file.
func handleMessage(msg signalRMsg, logFile *os.File) {
	for _, feed := range msg.M {
		if feed.M != "feed" || len(feed.A) < 2 {
			continue
		}

		var topic string
		if err := json.Unmarshal(feed.A[0], &topic); err != nil {
			log.Printf("Failed to read topic name: %v", err)
			continue
		}

		var payload []byte
		if strings.HasSuffix(topic, ".z") {
			var compressed string
			if err := json.Unmarshal(feed.A[1], &compressed); err != nil {
				log.Printf("Failed to read compressed payload for %s: %v", topic, err)
				continue
			}
			var err error
			payload, err = decompressZlib(compressed)
			if err != nil {
				log.Printf("Failed to decompress %s: %v", topic, err)
				continue
			}
		} else {
			payload = feed.A[1]
		}

		output := fmt.Sprintf("[%s] %s\n", topic, string(payload))
		fmt.Print(output)
		if _, err := logFile.WriteString(output); err != nil {
			log.Printf("Error writing to file: %v", err)
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	authToken := os.Getenv("F1_AUTH_TOKEN") // Optional: F1 TV Pro bearer token

	logFile, err := os.OpenFile("f1_live_data.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer logFile.Close()
	fmt.Println("Logging data to f1_live_data.log")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		// Re-negotiate on each (re)connect attempt.
		log.Println("Negotiating connection...")
		token, cookie, err := negotiate(authToken)
		if err != nil {
			log.Printf("Negotiate failed: %v — retrying in %s", err, reconnectDelay)
			select {
			case <-quit:
				fmt.Println("\nInterrupted. Bye!")
				return
			case <-time.After(reconnectDelay):
				continue
			}
		}

		log.Println("Connecting to F1 live timing...")
		conn, err := connectWS(token, cookie)
		if err != nil {
			log.Printf("WebSocket connect failed: %v — retrying in %s", err, reconnectDelay)
			select {
			case <-quit:
				fmt.Println("\nInterrupted. Bye!")
				return
			case <-time.After(reconnectDelay):
				continue
			}
		}

		if err := subscribe(conn); err != nil {
			log.Printf("Subscribe failed: %v — retrying in %s", err, reconnectDelay)
			conn.Close()
			select {
			case <-quit:
				fmt.Println("\nInterrupted. Bye!")
				return
			case <-time.After(reconnectDelay):
				continue
			}
		}
		log.Printf("Subscribed to %d topics. Press Ctrl+C to stop.", len(AllTopics))

		// Read messages until the connection drops or the user quits.
		dropped := make(chan struct{})
		go func() {
			defer close(dropped)
			for {
				_, raw, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						log.Printf("WebSocket error: %v", err)
					}
					return
				}

				var msg signalRMsg
				if err := json.Unmarshal(raw, &msg); err != nil {
					log.Printf("Failed to parse SignalR message: %v", err)
					continue
				}
				handleMessage(msg, logFile)
			}
		}()

		select {
		case <-quit:
			fmt.Println("\nInterrupted! Disconnecting...")
			conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			conn.Close()
			fmt.Println("Bye!")
			return
		case <-dropped:
			conn.Close()
			log.Printf("Connection lost — reconnecting in %s...", reconnectDelay)
			select {
			case <-quit:
				fmt.Println("\nInterrupted. Bye!")
				return
			case <-time.After(reconnectDelay):
				// loop and reconnect
			}
		}
	}
}
