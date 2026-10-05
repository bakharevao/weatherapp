package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"weatherapp/model"
	"weatherapp/repository"
	"weatherapp/types"
	"weatherapp/utils"
)

type WeatherProvider interface {
	FetchWeather(ctx context.Context, city string) (types.StoreData, error)
}

type openMeteoProvider struct {
	baseURL         string
	cityCoordinates repository.CityCoordinatesRepository
}

func NewOpenMeteoProvider(baseURL string, cityCoordinates repository.CityCoordinatesRepository) WeatherProvider {
	return &openMeteoProvider{baseURL: baseURL, cityCoordinates: cityCoordinates}
}

func (p *openMeteoProvider) FetchWeather(ctx context.Context, city string) (types.StoreData, error) {
	coords, err := p.cityCoordinates.GetCoordinates(ctx, city)
	if err != nil {
		return types.StoreData{}, err
	}

	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current_weather=true", p.baseURL, coords.Lat, coords.Lon)
	body, err := utils.ApiCall(ctx, url)
	if err != nil {
		return types.StoreData{}, err
	}

	var resp model.OpenMeteoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return types.StoreData{}, err
	}

	return types.StoreData{
		Name:        city,
		Latitude:    float32(coords.Lat),
		Longitude:   float32(coords.Lon),
		TempC:       float32(resp.CurrentWeather.Temperature),
		LastUpdated: resp.CurrentWeather.Time,
	}, nil
}
