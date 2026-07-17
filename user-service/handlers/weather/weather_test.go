package weather

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func setupMockTransport(resp *http.Response, err error) func() {
	originalTransport := http.DefaultTransport
	http.DefaultTransport = &mockTransport{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return resp, err
		},
	}
	return func() {
		http.DefaultTransport = originalTransport
	}
}

func TestGetWeather(t *testing.T) {
	origAPIKey := os.Getenv("WEATHER_API_KEY")
	defer os.Setenv("WEATHER_API_KEY", origAPIKey)

	t.Run("пустой API ключ", func(t *testing.T) {
		os.Unsetenv("WEATHER_API_KEY")

		data, err := GetWeather("Moscow")

		if err != nil {
			t.Errorf("ожидалась nil ошибка, получено: %v", err)
		}
		if data != nil {
			t.Errorf("ожидались nil данные при пустом ключе, получено: %+v", data)
		}
	})

	t.Run("успешный запрос", func(t *testing.T) {
		os.Setenv("WEATHER_API_KEY", "test-key")
		defer setupMockTransport(nil, nil)()

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"main": {"temp": 22.5},
				"weather": [{"main": "Clouds"}]
			}`)),
		}
		mockResp.Header.Set("Content-Type", "application/json")
		defer setupMockTransport(mockResp, nil)()

		data, err := GetWeather("Moscow")

		if err != nil {
			t.Fatalf("ожидалась nil ошибка, получено: %v", err)
		}
		if data == nil {
			t.Fatal("ожидались данные, получено nil")
		}
		if data.Temp != 22.5 {
			t.Errorf("ожидалась температура 22.5, получено: %f", data.Temp)
		}
		if data.Condition != "Облачно" {
			t.Errorf("ожидалось условие 'Облачно', получено: %s", data.Condition)
		}
	})

	t.Run("город не найден (404)", func(t *testing.T) {
		os.Setenv("WEATHER_API_KEY", "test-key")

		mockResp := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{"message": "city not found"}`)),
		}
		defer setupMockTransport(mockResp, nil)()

		data, err := GetWeather("UnknownCity12345")

		if err != nil {
			t.Errorf("ожидалась nil ошибка для 404, получено: %v", err)
		}
		if data != nil {
			t.Errorf("ожидались nil данные для 404, получено: %+v", data)
		}
	})

	t.Run("ошибка сервера (500)", func(t *testing.T) {
		os.Setenv("WEATHER_API_KEY", "test-key")

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(`{"message": "internal error"}`)),
		}
		defer setupMockTransport(mockResp, nil)()

		data, err := GetWeather("Moscow")

		if err == nil {
			t.Errorf("ожидалась ошибка для 500 статуса, получено nil")
		}
		if data != nil {
			t.Errorf("ожидались nil данные для 500, получено: %+v", data)
		}
	})

	t.Run("ошибка сети", func(t *testing.T) {
		os.Setenv("WEATHER_API_KEY", "test-key")

		netErr := errors.New("network timeout")
		defer setupMockTransport(nil, netErr)()

		data, err := GetWeather("Moscow")

		if err == nil {
			t.Errorf("ожидалась ошибка сети, получено nil")
		}
		if data != nil {
			t.Errorf("ожидались nil данные при ошибке сети, получено: %+v", data)
		}
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		os.Setenv("WEATHER_API_KEY", "test-key")

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`this is not json`)),
		}
		defer setupMockTransport(mockResp, nil)()

		data, err := GetWeather("Moscow")

		if err == nil {
			t.Errorf("ожидалась ошибка декодирования, получено nil")
		}
		if data != nil {
			t.Errorf("ожидались nil данные при невалидном JSON, получено: %+v", data)
		}
	})

	t.Run("пустой массив weather", func(t *testing.T) {
		os.Setenv("WEATHER_API_KEY", "test-key")

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"main": {"temp": 15.0},
				"weather": []
			}`)),
		}
		mockResp.Header.Set("Content-Type", "application/json")
		defer setupMockTransport(mockResp, nil)()

		data, err := GetWeather("Moscow")

		if err != nil {
			t.Fatalf("ожидалась nil ошибка, получено: %v", err)
		}
		if data.Condition != "unknown" {
			t.Errorf("ожидалось условие 'unknown' при пустом массиве, получено: %s", data.Condition)
		}
	})
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Clear", "Ясно"},
		{"Clouds", "Облачно"},
		{"Rain", "Дождь"},
		{"Drizzle", "Дождь"},
		{"Thunderstorm", "Дождь"},
		{"Mist", "Туман"},
		{"Fog", "Туман"},
		{"Haze", "Туман"},
		{"Snow", "Переменная облачность"}, // default case
		{"", "Переменная облачность"},     // default case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalize(tt.input)
			if result != tt.expected {
				t.Errorf("normalize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
