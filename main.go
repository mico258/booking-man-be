package main

import (
	"github.com/booking-man-be/lib/logger"
)

func main() {
	logClient := logger.New()
	logClient.Infof("Test format message for logger : %d", 43)

}
