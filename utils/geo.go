package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	googleMapsAPIURL string = "https://maps.googleapis.com/maps/api/geocode/json"
)

type GeocodingResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
	}
	Status string `json:"status"`
}

type Location struct {
	Address string
	MapLink string
	Country string
}

func GenerateMapLinkFromLatLng(lat float64, lng float64) string {
	return fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, lng)
}

func GetShopLocation(lat float64, lng float64) (*Location, error) {
	location := new(Location)
	url := fmt.Sprintf("%s?latlng=%f,%f&key=%s", googleMapsAPIURL, lat, lng, os.Getenv("GOOGLE_MAPS_API_KEY"))

	res, err := http.Get(url)
	if err != nil {
		return location, fmt.Errorf("failed to call google maps geocode API: %v", err)
	}
	defer res.Body.Close()

	var geoRes GeocodingResponse
	if err := json.NewDecoder(res.Body).Decode(&geoRes); err != nil {
		return location, fmt.Errorf("failed to decode API response: %v", err)
	}

	if geoRes.Status != "OK" || len(geoRes.Results) == 0 {
		return location, fmt.Errorf("no results found for given coordinates: %v", err)
	}

	location.Address = geoRes.Results[0].FormattedAddress

	// Yes very ugly way to get country, for some reason could not quite get it directly from google maps api
	addressParts := strings.Split(location.Address, ", ")
	location.Country = addressParts[len(addressParts)-1]

	// Super scuffed but problem is that google maps return Myanmar as Myanmar (Burma)
	if strings.Contains(location.Country, "Myanmar") {
		location.Country = "Myanmar"
	}

	location.MapLink = fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, lng)

	return location, nil

}
