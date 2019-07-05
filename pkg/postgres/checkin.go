package postgres

import (
	"log"

	"github.com/golang/geo/s2"
	"github.com/jinzhu/gorm"
	"math"
	"time"

	"github.com/phillebaba/bike-tracker/pkg/domain"
)

type CheckinService struct {
	DB *gorm.DB
}

func (c CheckinService) List() []domain.Checkin {
	checkins := []domain.Checkin{}
	c.DB.Find(&checkins)
	return checkins
}

func (c CheckinService) Delete(id int) {
	var checkin domain.Checkin
	c.DB.First(&checkin, id)
	c.DB.Delete(&checkin)
}

func (c CheckinService) Register(lat float64, lng float64, precision int, description string) {
	var trip domain.Trip
	c.DB.First(&trip)

	// Need to start trip if not started
	if !trip.Ongoing() {
		c.DB.Model(&trip).Update("StartTime", time.Now())
	}

	cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
	level := 15 - precision
	parentCellID := cellID.Parent(level)
	latLng := parentCellID.LatLng()
	approxArea := s2.CellFromCellID(parentCellID).AverageArea() * math.Pow(10, 13)
	radius := math.Sqrt(approxArea / math.Pi)

	log.Println(approxArea)
	log.Println(int(radius))
	log.Println(parentCellID.Level())

	checkIn := domain.Checkin{
		Trip:        trip,
		Time:        time.Now(),
		Lat:         latLng.Lat.Degrees(),
		Lng:         latLng.Lng.Degrees(),
		Radius:      int(radius),
		Description: description,
	}

	c.DB.Create(&checkIn)
}
