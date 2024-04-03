package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func fetchData() (map[string][]vehicle, string) {
	resp, err := http.Get("https://ckan2.multimediagdansk.pl/gpsPositions?v=2")
	if err != nil {
		log.Fatal("Error fetching JSON: ", err)
	}

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
	return prepBuses(data), data.LastUpdate
}
func prepBuses(data ztmResponse) map[string][]vehicle {
	vehicles := make(map[string][]vehicle)

	// I don't like this but i want to have users choose a vehicle from a dropdown list.
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

func main() {
	data, updatedAt := fetchData()
	http.HandleFunc("/api/vehicles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Error marshaling data: ", err)
		}
		w.Write(jsonData)

	})

	http.HandleFunc("/api/updated", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(updatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal("Error marshaling data: ", err)
		}
		w.Write(jsonData)

	})

	fetchData()

	log.Print("Listening on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
