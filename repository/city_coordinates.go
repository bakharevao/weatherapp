package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"weatherapp/types"
)

var ErrNotFound = errors.New("not found")

type CityCoordinatesRepository interface {
	GetCoordinates(ctx context.Context, city string) (types.Coordinates, error)
	ListCities(ctx context.Context) ([]string, error)
	UpsertCoordinates(ctx context.Context, city string, coords types.Coordinates) error
}

type cityCoordinatesRepository struct {
	db *sql.DB
}

func NewCityCoordinatesRepository(db *sql.DB) CityCoordinatesRepository {
	return &cityCoordinatesRepository{db: db}
}

func (r *cityCoordinatesRepository) GetCoordinates(ctx context.Context, city string) (types.Coordinates, error) {
	var c types.Coordinates
	err := r.db.QueryRowContext(ctx, `
		SELECT latitude, longitude
		FROM city_coordinates
		WHERE city = ?`,
		strings.ToLower(city),
	).Scan(&c.Lat, &c.Lon)
	if errors.Is(err, sql.ErrNoRows) {
		return types.Coordinates{}, ErrNotFound
	} else if err != nil {
		return types.Coordinates{}, fmt.Errorf("failed to get coordinates: %w", err)
	}
	return c, nil
}

func (r *cityCoordinatesRepository) ListCities(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT city FROM city_coordinates ORDER BY city`)
	if err != nil {
		return nil, fmt.Errorf("failed to list cities: %w", err)
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, city)
	}
	return cities, rows.Err()
}

func (r *cityCoordinatesRepository) UpsertCoordinates(ctx context.Context, city string, coords types.Coordinates) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO city_coordinates (city, latitude, longitude)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			latitude = VALUES(latitude),
			longitude = VALUES(longitude)`,
		strings.ToLower(city), coords.Lat, coords.Lon,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert coordinates: %w", err)
	}
	return nil
}
