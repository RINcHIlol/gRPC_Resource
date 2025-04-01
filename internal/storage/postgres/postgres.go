package postgres

import (
	"ResourceService/internal/domain/models"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Data(ctx context.Context, resourceName string) (models.Data, error) {
	const op = "postgres.Data"
	var info string
	var success bool

	query := "SELECT info, access from data where data_name = $1"

	err := s.db.QueryRowContext(ctx, query, resourceName).Scan(&info, &success)
	if err != nil {
		return models.Data{}, fmt.Errorf("%s: %w", op, err)
	}

	data := models.Data{Info: info, Success: success}
	return data, nil
}
