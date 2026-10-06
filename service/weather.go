package service

import (
	"context"
	"errors"
	"time"
	"weatherapp/provider"
	"weatherapp/repository"
	"weatherapp/types"
)

type weatherService struct {
	weatherProvider           provider.WeatherProvider
	geocodingProvider         provider.GeocodingProvider
	cityCoordinatesRepository repository.CityCoordinatesRepository
	contextTimeout            time.Duration
}

type WeatherService interface {
	AddWeather(ctx context.Context, city string) (*types.StoreData, error)
	ListCities(ctx context.Context) ([]string, error)
}

func NewWeatherService(weatherProvider provider.WeatherProvider, geocodingProvider provider.GeocodingProvider, cityCoordinatesRepository repository.CityCoordinatesRepository, timeout time.Duration) WeatherService {
	return &weatherService{
		weatherProvider:           weatherProvider,
		geocodingProvider:         geocodingProvider,
		cityCoordinatesRepository: cityCoordinatesRepository,
		contextTimeout:            timeout,
	}
}

func (w *weatherService) AddWeather(ctx context.Context, city string) (*types.StoreData, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	_, err := w.cityCoordinatesRepository.GetCoordinates(ctx, city)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}

		coords, err := w.geocodingProvider.Resolve(ctx, city)
		if err != nil {
			return nil, err
		}

		if err := w.cityCoordinatesRepository.UpsertCoordinates(ctx, city, coords); err != nil {
			return nil, err
		}
	}

	store, err := w.weatherProvider.FetchWeather(ctx, city)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (w *weatherService) ListCities(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.cityCoordinatesRepository.ListCities(ctx)
}
