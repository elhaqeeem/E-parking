package routes

import (
	"project/controllers"

	"github.com/gin-contrib/cors" // Import package CORS
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Konfigurasi CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true                                             // Izinkan semua origin (untuk development, sesuaikan untuk produksi)
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTION"}  // Izinkan metode HTTP
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // Izinkan header tertentu

	// Apply CORS middleware
	router.Use(cors.New(corsConfig))

	// Routes untuk tempat parkir
	router.GET("/parking-spots", controllers.GetAllParkingSpots)
	router.POST("/parking-spots", controllers.CreateParkingSpot)
	router.PUT("/parking-spots/:id", controllers.EditParkingSpot)
	router.DELETE("/parking-spots/:id", controllers.DeleteParkingSpot)

	// Routes untuk reservasi parkir

	router.POST("/book-spot", controllers.BookParkingSpot)
	router.PUT("/reservation/edit", controllers.EditBookParkingSpot)            // Endpoint untuk mengedit reservasi
	router.DELETE("/reservation/delete/:id", controllers.DeleteBookParkingSpot) // Endpoint untuk menghapus reservasi
	router.GET("/reservations", controllers.GetReservations)

	return router
}
