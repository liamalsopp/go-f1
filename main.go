package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type Driver struct {
	DriverNumber int
	DriverName   string
	Team         string
	Country      string
}

type Race struct {
	RaceName    string
	Location    string
	Date        time.Time
	Winner      Driver
	SecondPlace Driver
	ThirdPlace  Driver
}

func main() {
	ctx := context.Background()
	router := gin.Default()
	router.GET("/albums")

	router.Run("localhost:8080")

}
