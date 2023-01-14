package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SunsetAlert() {
	lat := os.Getenv("SUNSET_LATITUDE")
	lon := os.Getenv("SUNSET_LONGITUDE")
	fastDebug, err := strconv.ParseBool(os.Getenv("SUNSET_FAST_DEBUG"))
	if err != nil {
		fmt.Printf("SUNSET_FAST_DEBUG environment variable must be either 'true' or 'false'.\n")
		return
	}

	//Find sunset time
	sunsetTime := GetSunsetTime(lat, lon)
	fmt.Printf("The next sunset is at: %s\n", sunsetTime.String())
	//Then... let's say the program starts at noon. We'll need to sleep for a while.
	now := time.Now()
	//If sunset has not happened yet today, we can start our alert cycle
	if now.Before(sunsetTime) || fastDebug {
		//Then we sleep until an hour before sunset
		hourBeforeSunset := sunsetTime.Add(-time.Hour)
		if now.Before(hourBeforeSunset) || fastDebug {
			//Subtract a bit of sleep time to help the times work out
			sleepTime := hourBeforeSunset.Sub(now) - 30
			fmt.Printf("Sleeping program until one hour before sunset: %s\n", hourBeforeSunset.String())
			time.Sleep(sleepTime)
			//Then we call a yellow pulse for a bit
			fmt.Printf("Sending first warning pulse\n")
			SendWLEDPulse()
			fmt.Printf("First warning pulse complete\n")
		}
		//Then we sleep until a half hour before sunset
		halfHourBeforeSunset := sunsetTime.Add(-30 * time.Minute)
		if now.Before(halfHourBeforeSunset) || fastDebug {
			//Refresh Now
			now = time.Now()
			//Subtract a bit of sleep time to help the times work out
			sleepTime := halfHourBeforeSunset.Sub(now) - 30
			fmt.Printf("Sleeping program "+sleepTime.String()+" seconds, until one half hour before sunset: %s\n", halfHourBeforeSunset.String())
			time.Sleep(sleepTime)
			//Then we call a faster yellow pulse for a minute
			fmt.Printf("Sending second warning pulse\n")
			SendWLEDPulse()
			fmt.Printf("Second warning pulse complete\n")
		}
		//Then we sleep until 15 minutes before sunset
		quarterHourBeforeSunset := sunsetTime.Add(-15 * time.Minute)
		if now.Before(quarterHourBeforeSunset) || fastDebug {
			//Refresh Now
			now = time.Now()
			//Subtract a bit of sleep time to help the times work out
			sleepTime := quarterHourBeforeSunset.Sub(now) - 30
			fmt.Printf("Sleeping program "+sleepTime.String()+" seconds, until one quarter hour before sunset: %s\n", quarterHourBeforeSunset.String())
			time.Sleep(sleepTime)
			fmt.Printf("Sending third warning pulse\n")
			//Then we call an even faster yellow pulse for a minute
			SendWLEDPulse()
			fmt.Printf("Third warning pulse complete\n")
		}
	}

	//Then we sleep until noon tomorrow.
	nextNoon := getNextNoonTime(now)
	sleepTime := nextNoon.Sub(now)
	fmt.Printf("Sunset Alert complete. Sleeping %d seconds, until noon tomorrow.\n\n", sleepTime)

	time.Sleep(sleepTime)
}

func getNextNoonTime(now time.Time) (nextNoon time.Time) {
	tomorrow := now.Add(time.Hour * 24)
	nextNoon = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, tomorrow.Location())
	return
}

func GetSunsetTime(lat string, lon string) (sunsetTime time.Time) {
	now := time.Now()

	url := "https://api.sunrise-sunset.org/json?lat=" + lat + "&lng=" + lon
	resp, err := http.Get(url)
	println(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body := resp.Body
	var data SunriseSunset
	err = json.NewDecoder(body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return
	}

	layout := "3:04:05 PM"
	sunset, err := time.Parse(layout, data.Results.Sunset)
	if err != nil {
		fmt.Println(err)
	}
	sunset = time.Date(now.Year(), now.Month(), now.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), sunset.Location())

	return sunset
}

type SunriseSunset struct {
	Results struct {
		Sunrise                   string `json:"sunrise"`
		Sunset                    string `json:"sunset"`
		SolarNoon                 string `json:"solar_noon"`
		DayLength                 string `json:"day_length"`
		CivilTwilightBegin        string `json:"civil_twilight_begin"`
		CivilTwilightEnd          string `json:"civil_twilight_end"`
		NauticalTwilightBegin     string `json:"nautical_twilight_begin"`
		NauticalTwilightEnd       string `json:"nautical_twilight_end"`
		AstronomicalTwilightBegin string `json:"astronomical_twilight_begin"`
		AstronomicalTwilightEnd   string `json:"astronomical_twilight_end"`
	} `json:"results"`
	Status string `json:"status"`
}
