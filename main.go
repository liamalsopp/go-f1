package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

const (
	ClientID = "go-f1-dash-client"
)

var (
	Broker    string
	Token     string
	Username  string
	AllTopics = []string{
		"v1/car_data",
		"v1/position",
		"v1/weather",
		"v1/session_status",
		"v1/laps",
		"v1/intervals",
		"v1/pit",
		"v1/race_control",
		"v1/stints",
		"v1/team_radio",
		"v1/drivers",
		"v1/location",
		"v1/position",
		"v1/overtakes",
	}
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	Broker = os.Getenv("MQTT_BROKER")
	Token = os.Getenv("OPENF1_TOKEN")
	Username = os.Getenv("MQTT_USER")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(Broker)
	opts.SetClientID(ClientID)

	opts.SetUsername(Username)
	opts.SetPassword(Token)

	logFile, err := os.OpenFile("f1_live_data.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ Could not open log file: %v", err)
	}
	defer logFile.Close()
	fmt.Println("📁 Logging data to f1_live_data.log")

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		output := fmt.Sprintf("[%s] %s\n", msg.Topic(), payload)

		// Print to Console
		fmt.Print(output)

		// Write to File
		if _, err := logFile.WriteString(output); err != nil {
			log.Printf("Error writing to file: %v", err)
		}

		// Optional: If you want to parse specific topics:

		if msg.Topic() == "v1/car_data" {
			var data CarData // Defined in your openf1_types.go
			if err := json.Unmarshal(msg.Payload(), &data); err == nil {
				log.Printf("Failed to save car data: %v", err)
			}
		}
	})
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	filters := make(map[string]byte)
	for _, t := range AllTopics {
		filters[t] = 0 // We set QoS 0 inside the map for every topic
	}

	if token := client.SubscribeMultiple(filters, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Subscribe failed: %v", token.Error())
	}
	fmt.Println("🌊 Subscribed to firehose (#). Press Ctrl+C to stop.")

	// Create a channel that listens for system signals (Ctrl+C)
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	// The program will pause here until a signal is received
	<-keepAlive

	fmt.Println("\n🛑 Interrupted! Disconnecting...")
	client.Disconnect(250)
	fmt.Println("Bye!")
}
