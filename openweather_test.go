package openweather

import (
	"fmt"
	"testing"
)

func TestForecast(t *testing.T) {
	var s Session
	s.Init("apikey")
	wfs, err := s.ThreeHourForecast("Mettmann, Germany", Unit_Celcius)
	if err != nil {
		t.Fatal(err)
	}
	for _, wf := range wfs {
		fmt.Println(wf)
	}
}
