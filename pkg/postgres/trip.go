package postgres

import (
	"github.com/jinzhu/gorm"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
)

type TripService struct {
	DB *gorm.DB
}

func (c TripService) List() []domain.Trip {
	trips := []domain.Trip{}
	c.DB.Preload("Checkins").Preload("Checkpoints").Find(&trips)
	return trips
}

func (c TripService) Add(trip *domain.Trip) {
	c.DB.Create(trip)
}
