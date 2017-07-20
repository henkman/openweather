package openweather

import (
	"fmt"
	"testing"
)

const (
	API_KEY = "your key"
)

func TestThreeHourForecast(t *testing.T) {
	var s Session
	s.Init(API_KEY)
	wfs, err := s.ThreeHourForecast("Mettmann, Germany", Unit_Celcius)
	if err != nil {
		t.Fatal(err)
	}
	for _, wf := range wfs {
		fmt.Println(wf)
	}
}

func TestDailyForecast(t *testing.T) {
	var s Session
	s.Init(API_KEY)
	wfs, err := s.DailyForecast("Mettmann, Germany", Unit_Celcius)
	if err != nil {
		t.Fatal(err)
	}
	for _, wf := range wfs {
		fmt.Println(wf)
	}
}
