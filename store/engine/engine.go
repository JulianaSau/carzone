package engine

import (
	"context"
	"database/sql"

	"github.com/JulianaSau/carzone/models"
)

type EngineStore struct {
	db *sql.DB
}

func new(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (s *EngineStore) EngineById(ctx context.Context, id string) (*models.Engine, error) {

}

func (s *EngineStore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.EngineRequest, error) {

}

func (s *EngineStore) UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (*models.EngineRequest, error) {

}

func (s *EngineStore) DeleteEngine(ctx context.Context, id string) error {

}
