package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	BASE = "http://api.openweathermap.org/data/2.5"
)

type Unit string

const (
	Unit_Kelvin     Unit = ""
	Unit_Fahrenheit Unit = "imperial"
	Unit_Celcius    Unit = "metric"
)

type ThreeHourForecast struct {
	Date        time.Time
	Temp        float64
	Humidity    int
	Description string
	Icon        string
	Wind        struct {
		Speed  float64
		Degree float64
	}
}

type Session struct {
	cli http.Client
	key string
}

func (s *Session) Init(key string) error {
	s.key = key
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.cli.Jar = jar
	return nil
}

func (s *Session) IsInitialized() bool {
	return s.cli.Jar != nil
}

func (s *Session) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "openweather api")
	return s.cli.Do(req)
}

func (s *Session) ThreeHourForecast(q string, u Unit) ([]ThreeHourForecast, error) {
	ru := fmt.Sprintf(BASE+"/forecast?APPID=%s&q=%s", s.key, url.QueryEscape(q))
	if u != Unit_Kelvin {
		ru += "&units=" + string(u)
	}
	res, err := s.request("GET", ru, nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		List []struct {
			Dt   int `json:"dt"`
			Main struct {
				Temp     float64 `json:"temp"`
				Humidity int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
			Wind struct {
				Speed float64 `json:"speed"`
				Deg   float64 `json:"deg"`
			} `json:"wind"`
		} `json:"list"`
	}
	err = json.NewDecoder(res.Body).Decode(&result)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	wfs := make([]ThreeHourForecast, 0, len(result.List))
	for _, e := range result.List {
		wfs = append(wfs, ThreeHourForecast{
			Date:        time.Unix(int64(e.Dt), 0),
			Temp:        e.Main.Temp,
			Humidity:    e.Main.Humidity,
			Description: e.Weather[0].Description,
			Icon: fmt.Sprintf("http://openweathermap.org/img/w/%s.png",
				e.Weather[0].Icon),
			Wind: struct {
				Speed  float64
				Degree float64
			}{
				Speed:  e.Wind.Speed,
				Degree: e.Wind.Deg,
			},
		})
	}
	return wfs, nil
}
