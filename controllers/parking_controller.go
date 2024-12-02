package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"project/config"
	"project/models"

	"github.com/gin-gonic/gin"
)

// CreateParkingSpot - Menambahkan tempat parkir baru
func CreateParkingSpot(c *gin.Context) {
	var spot models.ParkingSpot
	if err := c.ShouldBindJSON(&spot); err != nil {
		log.Printf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Pastikan spot_number tidak kosong dan unique
	if spot.SpotNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spot number is required"})
		return
	}

	// Insert parking spot into the database
	_, err := config.DB.Exec("INSERT INTO parking_spots (spot_number, is_occupied) VALUES (?, ?)", spot.SpotNumber, false)
	if err != nil {
		log.Printf("Error creating parking spot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create parking spot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking spot created successfully"})
}

// EditParkingSpot - Mengubah tempat parkir yang ada
func EditParkingSpot(c *gin.Context) {
	var spot models.ParkingSpot
	if err := c.ShouldBindJSON(&spot); err != nil {
		log.Printf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validasi apakah spot ada
	var existingSpot models.ParkingSpot
	err := config.DB.QueryRow("SELECT id, spot_number, is_occupied FROM parking_spots WHERE id = ?", spot.ID).Scan(&existingSpot.ID, &existingSpot.SpotNumber, &existingSpot.IsOccupied)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking spot not found"})
		return
	} else if err != nil {
		log.Printf("Error checking spot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check parking spot"})
		return
	}

	// Update parking spot details
	_, err = config.DB.Exec("UPDATE parking_spots SET spot_number = ?, is_occupied = ? WHERE id = ?", spot.SpotNumber, spot.IsOccupied, spot.ID)
	if err != nil {
		log.Printf("Error updating parking spot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking spot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking spot updated successfully"})
}

// DeleteParkingSpot - Menghapus tempat parkir
func DeleteParkingSpot(c *gin.Context) {
	spotID := c.Param("id")
	if spotID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spot ID is required"})
		return
	}

	// Check if the spot exists before trying to delete
	var existingSpot models.ParkingSpot
	err := config.DB.QueryRow("SELECT id FROM parking_spots WHERE id = ?", spotID).Scan(&existingSpot.ID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking spot not found"})
		return
	} else if err != nil {
		log.Printf("Error checking spot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check parking spot"})
		return
	}

	// Delete the parking spot
	_, err = config.DB.Exec("DELETE FROM parking_spots WHERE id = ?", spotID)
	if err != nil {
		log.Printf("Error deleting parking spot: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete parking spot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking spot deleted successfully"})
}

// GetAllParkingSpots - Menampilkan semua tempat parkir
func GetAllParkingSpots(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, spot_number, is_occupied FROM parking_spots")
	if err != nil {
		log.Printf("Error fetching parking spots: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parking spots"})
		return
	}
	defer rows.Close()

	var spots []models.ParkingSpot
	for rows.Next() {
		var spot models.ParkingSpot
		if err := rows.Scan(&spot.ID, &spot.SpotNumber, &spot.IsOccupied); err != nil {
			log.Printf("Error parsing parking spot: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse parking spot"})
			return
		}
		spots = append(spots, spot)
	}

	c.JSON(http.StatusOK, spots)
}

// BookParkingSpot - Memesan tempat parkir
func BookParkingSpot(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		log.Printf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validasi dan format waktu
	startTime, err := time.Parse(time.RFC3339, reservation.StartTime)
	if err != nil {
		log.Printf("Invalid start_time format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
		return
	}
	formattedTime := startTime.Format("2006-01-02 15:04:05")

	// Validasi apakah tempat parkir tersedia
	var isOccupied bool
	err = config.DB.QueryRow("SELECT is_occupied FROM parking_spots WHERE id = ?", reservation.SpotID).Scan(&isOccupied)
	if err == sql.ErrNoRows {
		log.Printf("Spot ID %d not found", reservation.SpotID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking spot not found"})
		return
	} else if err != nil {
		log.Printf("Error checking spot status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check parking spot availability"})
		return
	}

	if isOccupied {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parking spot is already occupied"})
		return
	}

	// Proses pemesanan
	_, err = config.DB.Exec("INSERT INTO reservations (name, car_number, spot_id, start_time, duration) VALUES (?, ?, ?, ?, ?)",
		reservation.Name, reservation.CarNumber, reservation.SpotID, formattedTime, reservation.Duration)
	if err != nil {
		log.Printf("Error creating reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	_, err = config.DB.Exec("UPDATE parking_spots SET is_occupied = TRUE WHERE id = ?", reservation.SpotID)
	if err != nil {
		log.Printf("Error updating parking spot status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking spot status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking spot reserved successfully"})
}

// GetReservations - Menampilkan semua reservasi
func GetReservations(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, car_number, spot_id, start_time, duration FROM reservations")
	if err != nil {
		log.Printf("Error fetching reservations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var res models.Reservation
		if err := rows.Scan(&res.ID, &res.Name, &res.CarNumber, &res.SpotID, &res.StartTime, &res.Duration); err != nil {
			log.Printf("Error parsing reservation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reservation"})
			return
		}

		// Format waktu jika diperlukan
		res.StartTime = strings.Replace(res.StartTime, "T", " ", 1)
		res.StartTime = strings.TrimSuffix(res.StartTime, "Z")

		reservations = append(reservations, res)
	}

	c.JSON(http.StatusOK, reservations)
}

// EditBookParkingSpot - Mengedit reservasi parkir
func EditBookParkingSpot(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		log.Printf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validasi reservasi yang ada
	var existingReservation models.Reservation
	err := config.DB.QueryRow("SELECT id, spot_id FROM reservations WHERE id = ?", reservation.ID).Scan(&existingReservation.ID, &existingReservation.SpotID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	} else if err != nil {
		log.Printf("Error checking reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check reservation"})
		return
	}

	// Update reservasi
	_, err = config.DB.Exec(
		"UPDATE reservations SET name = ?, car_number = ?, spot_id = ?, start_time = ?, duration = ? WHERE id = ?",
		reservation.Name, reservation.CarNumber, reservation.SpotID, reservation.StartTime, reservation.Duration, reservation.ID,
	)
	if err != nil {
		log.Printf("Error updating reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation updated successfully"})
}

// DeleteBookParkingSpot - Menghapus reservasi parkir
func DeleteBookParkingSpot(c *gin.Context) {
	reservationID := c.Param("id")
	if reservationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation ID is required"})
		return
	}

	// Validasi apakah reservasi ada
	var existingReservation models.Reservation
	err := config.DB.QueryRow("SELECT id, spot_id FROM reservations WHERE id = ?", reservationID).Scan(&existingReservation.ID, &existingReservation.SpotID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	} else if err != nil {
		log.Printf("Error checking reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check reservation"})
		return
	}

	// Hapus reservasi dari database
	_, err = config.DB.Exec("DELETE FROM reservations WHERE id = ?", reservationID)
	if err != nil {
		log.Printf("Error deleting reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reservation"})
		return
	}

	// Perbarui status tempat parkir terkait menjadi tidak terpakai (tidak ditempati)
	_, err = config.DB.Exec("UPDATE parking_spots SET is_occupied = FALSE WHERE id = ?", existingReservation.SpotID)
	if err != nil {
		log.Printf("Error updating parking spot status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking spot status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted successfully"})
}
