package main

import (
	"log"

	"github.com/phillebaba/bike-tracker/pkg/domain"
	"github.com/phillebaba/bike-tracker/pkg/http"
	"github.com/phillebaba/bike-tracker/pkg/postgres"
)

func main() {
	log.Println("Starting Bike Tracker")

	serviceContext := postgres.Init()

	result := serviceContext.CheckpointService.List()
	if len(result) == 0 {
		log.Println("Creating data")

		trip := domain.Trip{
			Name:        "Göteborg-Gotland-Göteborg",
			Description: "Trip to Gotland, around the island and back",
		}
		serviceContext.TripService.Add(&trip)

		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.674778, Lng: 11.931406}) // Skäpplandsgatan
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.758292, Lng: 16.636712}) // Västervik
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.635028, Lng: 18.286178}) // Visby
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.716103, Lng: 18.793730}) // Slite
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.631040, Lng: 18.277280}) // Visby
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.263863, Lng: 16.456182}) // Oskarshamn
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.699514, Lng: 11.952062}) // Göteborg
	}

	http.Run(serviceContext)
}
