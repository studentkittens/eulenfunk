package ui

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	owm "github.com/briandowns/openweathermap"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/util"
)

func celsius(c float64) string {
	if c > 100 {
		c = 99.0
	}

	return fmt.Sprintf("%d৹C", int(c))
}

func degToDirection(deg int) string {
	switch {
	case deg > 315 || deg <= 45:
		return "↑"
	case deg <= 135:
		return "→"
	case deg <= 225:
		return "↓"
	case deg <= 315:
		return "←"
	default:
		return "o"
	}
}

func weatherForecast() (*owm.ForecastWeatherData, error) {
	w, err := owm.NewForecast("C", "DE")
	if err != nil {
		log.Printf("Failed to instance new forecast: %v", err)
		return nil, err
	}

	// That's Augsburg:
	err = w.DailyByCoordinates(
		&owm.Coordinates{
			Latitude:  48.3830555,
			Longitude: 10.8830555,
		},
		3, // 3 days of forecast
	)

	if err != nil {
		fmt.Printf("Failed to download by coordinate")
		return nil, err
	}

	return w, nil
}

func toScreen(w *owm.ForecastWeatherData, p *owm.ForecastWeatherList, dayOff, width int) []string {
	top := w.City.Name

	now := time.Now()
	date := fmt.Sprintf(
		"%d.%d.%d",
		now.Day()+dayOff,
		now.Month(),
		now.Year()-2000,
	)

	top += strings.Repeat(" ", width-len(top)-len(date))
	top += date

	status := "No weather today."
	if len(p.Weather) > 0 {
		status = util.Center(p.Weather[0].Description, width, ' ')
	}

	humidity := p.Humidity
	if humidity >= 100 {
		humidity = 99
	}

	stats := fmt.Sprintf(
		"R%5.1f%% %2d%% %s%4.1fm/s",
		p.Rain,
		humidity,
		degToDirection(p.Deg),
		p.Speed,
	)

	temps := fmt.Sprintf(
		"%s %s %s %s",
		celsius(p.Temp.Morn),
		celsius(p.Temp.Day),
		celsius(p.Temp.Eve),
		celsius(p.Temp.Night),
	)

	return []string{
		top,
		status,
		stats,
		temps,
	}
}

func errorScreen(width int) [][]string {
	return [][]string{{
		strings.Repeat("=", width),
		"Sorry, no weather...",
		"Please see the log.",
		strings.Repeat("=", width),
	}}
}

func downloadData(width int) [][]string {
	w, err := weatherForecast()
	if err != nil {
		log.Printf("Failed to retrieve forecast: %v", err)
		return errorScreen(width)
	}

	screens := [][]string{}
	for dayOff, p := range w.List {
		screens = append(screens, toScreen(w, &p, dayOff, width))
	}

	return screens
}

func displayWeather(lw *display.LineWriter, screen []string) {
	for idx, line := range screen {
		if err := lw.Line("weather", idx, line); err != nil {
			log.Printf("Failed to display weather widget: %v", err)
		}

		log.Printf("weather: %02d: %s", idx, line)
	}
}

func init() {
	// The OWM API is a bit weird:
	// They expect the API Key in the OWM_API_KEY env var.
	if err := os.Setenv("OWM_API_KEY", "7e8a8d42af13c734b8960a714e966c5c"); err != nil {
		log.Printf("Failed to set OWM_API_KEY env var (huh?): %v", err)
	}
}

// RunWeather displays a weather forecast in the "weather" window.
func RunWeather(lw *display.LineWriter, width int, ctx context.Context) {
	switchTicker := time.NewTicker(10 * time.Second)
	updateTicker := time.NewTicker(30 * time.Minute)

	screens := downloadData(width)
	screenIdx := 0

	if len(screens) > 0 {
		displayWeather(lw, screens[0])
	}

	for {
		select {
		// Update the data:
		case <-updateTicker.C:
			screens = downloadData(width)
		// Toggle through:
		case <-switchTicker.C:
			if screenIdx < len(screens) {
				displayWeather(lw, screens[screenIdx])
				screenIdx = (screenIdx + 1) % len(screens)
			}

		// Watch for aborts:
		case <-ctx.Done():
			return
		}
	}
}
