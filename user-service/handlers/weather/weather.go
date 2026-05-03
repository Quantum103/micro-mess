package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WeatherData struct {
	Temp      float64 `json:"temp"`
	Condition string  `json:"condition"`
}

func GetWeather(city string) (*WeatherData, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	// ДЛЯ РЕКРУТЕРОВ - если нет API, то просто не будет показывать, но приложение не сломает.
	// можно вставить свой API с openweathermap
	if apiKey == "" {
		return nil, nil
	}
	baseURL := "https://api.openweathermap.org/data/2.5/weather"
	params := url.Values{}
	params.Set("q", city)
	params.Set("appid", apiKey)
	params.Set("units", "metric")
	params.Set("lang", "ru")

	fullURL := baseURL + "?" + params.Encode()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API ошибка: %d", resp.StatusCode)
	}

	var result struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка декодирования: %w", err)
	}

	condition := "unknown"
	if len(result.Weather) > 0 {
		condition = normalize(result.Weather[0].Main)
	}

	return &WeatherData{
		Temp:      result.Main.Temp,
		Condition: condition,
	}, nil
}

func normalize(own string) string {
	switch own {
	case "Clear":
		return "Ясно"
	case "Clouds":
		return "Облачно"
	case "Rain", "Drizzle", "Thunderstorm":
		return "Дождь"
	case "Mist", "Fog", "Haze":
		return "Туман"
	default:
		return "Переменная облачность"
	}
}
