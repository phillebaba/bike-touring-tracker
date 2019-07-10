package main

import (
	"flag"
	"log"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
	"github.com/phillebaba/bike-touring-tracker/pkg/http"
	"github.com/phillebaba/bike-touring-tracker/pkg/postgres"
)

func main() {
	host := flag.String("postgres-host", "localhost", "host name for postgres")
	port := flag.Int("postgres-port", 5432, "port for postgres")
	dbname := flag.String("postgres-dbname", "postgres", "name of database")
	user := flag.String("postgres-username", "postgres", "user for postgres")
	password := flag.String("postgres-password", "password", "password for postgres")
	flag.Parse()

	log.Println("Starting Bike Tracker")

	serviceContext := postgres.Init(*host, *port, *user, *password, *dbname)

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
