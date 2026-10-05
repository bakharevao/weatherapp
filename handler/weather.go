package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"weatherapp/service"
	"weatherapp/utils"
)

func GetWeatherHandler(weatherService service.WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {
		city := c.Param("city")
		if city == "" {
			utils.NewError(c, http.StatusBadRequest, fmt.Errorf("city is required"))
			return nil
		}

		ctx := c.Request().Context()
		data, err := weatherService.AddWeather(ctx, city)
		if err != nil {
			utils.NewError(c, http.StatusBadGateway, err)
			return nil
		}
		return c.JSON(http.StatusOK, data)
	}
}
