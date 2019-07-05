package postgres

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/phillebaba/bike-tracker/pkg/domain"
)

func Init() domain.ServiceContext {
	// DB Connection
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=mysecretpassword sslmode=disable")

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}

	//defer db.Close()

	db.AutoMigrate(&domain.Trip{})
	db.AutoMigrate(&domain.Checkin{})
	db.AutoMigrate(&domain.Checkpoint{})

	serviceContext := domain.ServiceContext{
		TripService:       TripService{db},
		CheckinService:    CheckinService{db},
		CheckpointService: CheckpointService{db},
	}

	return serviceContext
}
