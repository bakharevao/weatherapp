package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"weatherapp/model"
	"weatherapp/types"
	"weatherapp/utils"
)

type GeocodingProvider interface {
	Resolve(ctx context.Context, city string) (types.Coordinates, error)
}

type geocodingProvider struct {
	baseURL string
}

func NewGeocodingProvider(baseURL string) GeocodingProvider {
	return &geocodingProvider{baseURL: baseURL}
}

func (g *geocodingProvider) Resolve(ctx context.Context, city string) (types.Coordinates, error) {
	reqURL := fmt.Sprintf("%s?name=%s&count=1", g.baseURL, url.QueryEscape(city))
	body, err := utils.ApiCall(ctx, reqURL)
	if err != nil {
		return types.Coordinates{}, err
	}

	var resp model.GeocodingResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return types.Coordinates{}, err
	}

	if len(resp.Results) == 0 {
		return types.Coordinates{}, fmt.Errorf("no coordinates found for %q", city)
	}

	return types.Coordinates{
		Lat: resp.Results[0].Latitude,
		Lon: resp.Results[0].Longitude,
	}, nil
}
