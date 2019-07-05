package domain

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsOngoing(t *testing.T) {
	trip := Trip{StartTime: pq.NullTime{Time: time.Now(), Valid: true}}
	assert.True(t, trip.Ongoing())
}

func createNullTime(time time.Time) pq.NullTime {
	return pq.NullTime{Time: time, Valid: true}
}

func TestActiveTimeStarted(t *testing.T) {
	startTime := time.Now().Add(-time.Duration(5) * time.Minute)
	trip := Trip{StartTime: createNullTime(startTime)}

	result, err := trip.ActiveTime()
	assert.Greater(t, result.Minutes(), float64(5))
	assert.Nil(t, err)
}

func TestActiveTimeEnded(t *testing.T) {
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(5) * time.Minute)
	trip := Trip{StartTime: createNullTime(startTime), EndTime: createNullTime(endTime)}

	result, err := trip.ActiveTime()
	assert.Equal(t, result.Minutes(), float64(5))
	assert.Nil(t, err)
}
