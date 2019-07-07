package postgres

import (
	"github.com/jinzhu/gorm"

	"github.com/phillebaba/bike-touring-tracker/pkg/domain"
)

type CheckpointService struct {
	DB *gorm.DB
}

func (c CheckpointService) List() []domain.Checkpoint {
	checkpoints := []domain.Checkpoint{}
	c.DB.Find(&checkpoints)
	return checkpoints
}

func (c CheckpointService) Add(checkpoint *domain.Checkpoint) {
	c.DB.Create(checkpoint)
}
