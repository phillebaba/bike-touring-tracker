package postgres

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
)

func Init(host string, port int, user string, password string, dbname string) domain.ServiceContext {
	options := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", options)

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
