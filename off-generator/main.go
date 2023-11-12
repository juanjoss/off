package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/juanjoss/off-generator/event"

	"github.com/go-co-op/gocron"
)

func main() {
	// register events
	eventProvider := event.NewEventProvider()
	eventProvider.RegisterEvent(&event.ProductOrder{})
	eventProvider.RegisterEvent(&event.UserRegistration{})

	// launch scheduler to generate random registered events
	s := gocron.NewScheduler(time.UTC)

	period, err := strconv.Atoi(os.Getenv("GENERATOR_RUN_PERIOD"))
	if err != nil {
		log.Fatal("unable to set the execution period, check the GENERATOR_RUN_PERIOD env variable.")
	}
	s.Every(period).Second().Do(eventProvider.SelectEvent)

	s.WaitForSchedule().StartBlocking()
}
