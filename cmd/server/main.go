package main

import (
	"flag"
	"log"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
	"github.com/phillebaba/bike-touring-tracker/pkg/http"
	"github.com/phillebaba/bike-touring-tracker/pkg/postgres"
)

func main() {
	adminPassword := flag.String("admin-password", "password", "password for admin user")
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
			Name:        "Gotland Runt",
			Description: "Trip to Gotland, around the island and back",
		}
		serviceContext.TripService.Add(&trip)

		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.681960, Lng: 11.946025}) // Skäpplandsgatan
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.263863, Lng: 16.456182}) // Oskarshamn
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.635028, Lng: 18.286178}) // Visby
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.974794, Lng: 19.322754}) // Fårö
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.421742, Lng: 18.910230}) // Herrvik
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 56.922863, Lng: 18.132187}) // Hoborg
		serviceContext.CheckpointService.Add(&domain.Checkpoint{TripID: trip.ID, Lat: 57.699514, Lng: 11.952062}) // Göteborg
	}

	http.Run(*adminPassword, serviceContext)
}
