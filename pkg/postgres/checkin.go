package postgres

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/golang/geo/s2"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
)

type CheckinService struct {
	DB *gorm.DB
}

func (c CheckinService) List() []domain.Checkin {
	checkins := []domain.Checkin{}
	c.DB.Find(&checkins)
	return checkins
}

func (c CheckinService) Delete(id int) error {
	var checkin domain.Checkin
	c.DB.Preload("Trip").First(&checkin, id)

	// Reset trip if last checkin
	var count int
	c.DB.Table("checkins").Where("trip_id = ?", checkin.Trip.ID).Count(&count)
	if count == 1 {
		c.DB.Model(&checkin.Trip).Updates(map[string]interface{}{"Distance": 0, "StartTime": pq.NullTime{}})
		c.DB.Unscoped().Delete(&checkin)
		return nil
	}

	// Check if checkin being delted is the first
	var firstCheckin domain.Checkin
	c.DB.First(&firstCheckin)
	if checkin.ID != firstCheckin.ID {
		return errors.New("Need to delete Checkins in order")
	}

	// Calculate distance to remove from the trip
	var prevCheckin domain.Checkin
	c.DB.Where("id < ?", id).First(&prevCheckin)
	distance := int64(distance(checkin.Lat, checkin.Lng, prevCheckin.Lat, prevCheckin.Lng))
	totalDistance := checkin.Trip.Distance - distance

	c.DB.Model(&checkin.Trip).Update("Distance", totalDistance)
	c.DB.Unscoped().Delete(&checkin)

	return nil
}

func (c CheckinService) Register(name string, lat float64, lng float64, precision int, description string) {
	var trip domain.Trip
	c.DB.First(&trip)

	// Need to start trip if not started
	isOngoing := trip.Ongoing()
	if !isOngoing {
		c.DB.Model(&trip).Update("StartTime", time.Now())
	}

	cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
	level := 15 - precision
	parentCellID := cellID.Parent(level)
	latLng := parentCellID.LatLng()
	approxArea := s2.CellFromCellID(parentCellID).AverageArea() * math.Pow(10, 13)
	radius := math.Sqrt(approxArea / math.Pi)

	if isOngoing {
		var prevCheckin domain.Checkin
		c.DB.Where("trip_id = ?", trip.ID).First(&prevCheckin)
		distance := distance(lat, lng, prevCheckin.Lat, prevCheckin.Lng)
		log.Println(distance)
		totalDistance := trip.Distance + int64(distance)
		c.DB.Model(&trip).Update("Distance", totalDistance)
	}

	checkIn := domain.Checkin{
		Trip:        trip,
		Name:        name,
		Description: description,
		Time:        time.Now(),
		Lat:         latLng.Lat.Degrees(),
		Lng:         latLng.Lng.Degrees(),
		Radius:      int(radius),
	}

	c.DB.Create(&checkIn)
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	return dist * 1.609344
}
