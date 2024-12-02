package models

type ParkingSpot struct {
	ID         int    `json:"id"`
	SpotNumber string `json:"spot_number"`
	IsOccupied bool   `json:"is_occupied"`
}

type Reservation struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CarNumber string `json:"car_number"`
	SpotID    int    `json:"spot_id"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
}
