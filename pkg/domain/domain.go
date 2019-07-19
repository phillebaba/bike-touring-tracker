package domain

import (
	"errors"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type ServiceContext struct {
	TripService       TripService
	CheckpointService CheckpointService
	CheckinService    CheckinService
}

type Trip struct {
	gorm.Model
	Name        string
	Description string
	StartTime   pq.NullTime
	EndTime     pq.NullTime
	Checkpoints []Checkpoint
	Checkins    []Checkin
	Distance    int64
}

func (t Trip) IsFirst(checkpoint Checkpoint) bool {
	return checkpoint == t.Checkpoints[0]
}

func (t Trip) IsLast(checkpoint Checkpoint) bool {
	return checkpoint == t.Checkpoints[len(t.Checkpoints)-1]
}

func (t Trip) Ongoing() bool {
	return t.StartTime.Valid && !t.EndTime.Valid
}

func (t Trip) Ended() bool {
	return t.EndTime.Valid
}

func (t Trip) ActiveTime() (time.Duration, error) {
	if t.Ongoing() {
		return time.Since(t.StartTime.Time), nil
	} else if t.Ended() {
		return t.EndTime.Time.Sub(t.StartTime.Time), nil
	}

	return 0, errors.New("Trip not started")
}

func (t Trip) ActiveTimeFormatted() (string, error) {
	duration, err := t.ActiveTime()
	if err != nil {
		return "", err
	}

	return durafmt.ParseShort(duration).String(), nil
}

type TripService interface {
	List() []Trip
	Add(trip *Trip)
}

type Checkin struct {
	gorm.Model
	Name        string
	Description string
	Time        time.Time
	Lat         float64
	Lng         float64
	Radius      int

	Trip   Trip
	TripID uint
}

func (c Checkin) TimeSinceFormatted() string {
	return humanize.Time(c.Time)
}

type CheckinService interface {
	List() []Checkin
	Delete(id int) error
	Register(name string, lat float64, lng float64, precision int, description string)
}

type CheckpointType int

type Checkpoint struct {
	gorm.Model
	Lat    float64
	Lng    float64
	TripID uint
}

type CheckpointService interface {
	List() []Checkpoint
	Add(checkpoint *Checkpoint)
}
