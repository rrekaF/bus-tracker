package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type ztmVehicle struct {
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
type ztmResponse struct {
	LastUpdate string
	Vehicles   []ztmVehicle
}

// The same as ztmVehicle but without
// routeShortName field because i will
// use routeShortName as a key to look up vehicles
type vehicle struct {
	Generated              string  `json:"generated"`
	TripId                 int     `json:"tripId"`
	RouteId                int     `json:"routeId"`
	Headsign               string  `json:"headsign"`
	VehicleCode            string  `json:"vehicleCode"`
	VehicleService         string  `json:"vehicleService"`
	VehicleId              int     `json:"vehicleId"`
	Speed                  int     `json:"speed"`
	Direction              int     `json:"direction"`
	Delay                  int     `json:"delay"`
	ScheduledTripStartTime string  `json:"scheduledTripStartTime"`
	Lat                    float64 `json:"lat"`
	Lon                    float64 `json:"lon"`
	GpsQuality             int     `json:"gpsQuality"`
}

// GLOBALS
var vehicles map[string][]vehicle
var lastUpdate string

func timeDifference(updatedAt string) time.Duration {
	// example date from ztm api -> "2024-04-03T14:03:12Z"
	// RFC3339 = "2006-01-02T15:04:05Z07:00" <- https://pkg.go.dev/time#pkg-constants
	date, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		log.Fatal("Error parsing date string: ", err)
	}
	// nanoseconds to seconds

	return time.Since(date) / 1000000000

}
func prepBuses(data ztmResponse) map[string][]vehicle {
	vehicles := make(map[string][]vehicle)

	// i want to have users choose a vehicle from a dropdown list.
	// route short name is the most important because that's what it says at bus and tram stops
	// e.g. '174' -> [ {"data": data}, {"data": data}, {"data": data} ]

	for _, b := range data.Vehicles {
		vehicles[b.RouteShortName] = append(vehicles[b.RouteShortName], vehicle{
			Generated:              b.Generated,
			TripId:                 b.TripId,
			RouteId:                b.RouteId,
			Headsign:               b.Headsign,
			VehicleCode:            b.VehicleCode,
			VehicleService:         b.VehicleService,
			VehicleId:              b.VehicleId,
			Speed:                  b.Speed,
			Direction:              b.Direction,
			Delay:                  b.Delay,
			ScheduledTripStartTime: b.ScheduledTripStartTime,
			Lat:                    b.Lat,
			Lon:                    b.Lon,
			GpsQuality:             b.GpsQuality,
		})
	}
	return vehicles
}

func fetchData() {
	resp, err := http.Get("https://ckan2.multimediagdansk.pl/gpsPositions?v=2")
	if err != nil {
		log.Fatal("Error fetching JSON: ", err)
	}
	log.Print("Fetched Buses\n")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading ztmResponse body: ", err)
	}

	var data ztmResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error unmarshaling json data: ", err)
	}

	// Format the ztm api response the way i want
	// v <- prepBuses(data)
	// t <- data.LastUpdate
	vehicles = prepBuses(data)
	lastUpdate = data.LastUpdate
	// return prepBuses(data), timeDifference(data.LastUpdate)
}

func fetchPeriodically() {
	for {
		// Fetch data every 5 seconds
		fetchData()
		time.Sleep(5 * time.Second)
	}
}

func main() {
	// dataCh := make(chan map[string][]vehicle, 100)
	// updatedAtCh := make(chan string, 100)

	go fetchPeriodically()

	http.HandleFunc("/api/vehicles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		jsonData, err := json.Marshal(vehicles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Error marshaling data: ", err)
		}
		log.Print("Serving vehicles\n")
		w.Write(jsonData)

	})

	http.HandleFunc("/api/updated", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		jsonData, err := json.Marshal(timeDifference(lastUpdate))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Error marshaling data: ", err)
		}
		log.Print("Serving lastUpdate\n")
		w.Write(jsonData)
	})

	log.Print("Listening on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
