package main

import (
	"openf1_structs"
	"log"
	"os"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	Broker   = "tcp://mqtt.openf1.org:1883" // or ssl://mqtt.openf1.org:8883
	Topic    = "v1/car_data"                // Subscribe to live telemetry
	ClientID = "go-f1-dash-client"
	Token    = "YOUR_OPENF1_TOKEN_HERE"     // Required for MQTT
)

func main() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(Broker)
	opts.SetClientID(ClientID)
	
	opts.SetUsername("user")
	opts.SetPassword(Token)
}
