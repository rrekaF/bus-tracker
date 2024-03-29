package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type rawVehicle struct {
	Generated              string
	RouteShortName         string
	TripId                 int
	RouteId                int
	Headsign               string
	VehicleCode            string
	VehicleService         string
	VehicleId              int
	Speed                  int
	Direction              int
	Delay                  int
	ScheduledTripStartTime string
	Lat                    float64
	Lon                    float64
	GpsQuality             int
}
type response struct {
	LastUpdate string
	Vehicles   []rawVehicle
}

func fetchData() response {
	resp, err := http.Get("https://ckan2.multimediagdansk.pl/gpsPositions?v=2")
	if err != nil {
		log.Fatal("Error fetching JSON: ", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	var data response
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error unmarshaling json data: ", err)
	}
	prepBuses(data)
	return data
}

type vehicle struct {
	Generated              string
	TripId                 int
	RouteId                int
	Headsign               string
	VehicleCode            string
	VehicleService         string
	VehicleId              int
	Speed                  int
	Direction              int
	Delay                  int
	ScheduledTripStartTime string
	Lat                    float64
	Lon                    float64
	GpsQuality             int
}

// TODO
func prepBuses(data response) {
	// var vehicles map[string][]vehicle

	// I don't like this but i want to have users choose a vehicle from a dropdown list
	// for _, b := range data.Vehicles {
	// 	vehicles[b.RouteShortName] = append(vehicles[b.RouteShortName], vehicle{
	// 		Generated:              b.Generated,
	// 		TripId:                 b.TripId,
	// 		RouteId:                b.RouteId,
	// 		Headsign:               b.Headsign,
	// 		VehicleCode:            b.VehicleCode,
	// 		VehicleService:         b.VehicleService,
	// 		VehicleId:              b.VehicleId,
	// 		Speed:                  b.Speed,
	// 		Direction:              b.Direction,
	// 		Delay:                  b.Delay,
	// 		ScheduledTripStartTime: b.ScheduledTripStartTime,
	// 		Lat:                    b.Lat,
	// 		Lon:                    b.Lon,
	// 		GpsQuality:             b.GpsQuality,
	// 	})
	// }
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	fetchData()

	log.Print("Listening on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
