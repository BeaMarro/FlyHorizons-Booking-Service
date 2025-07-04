package routes

import (
	"flyhorizons-bookingservice/models"
	"flyhorizons-bookingservice/services/errors"
	"flyhorizons-bookingservice/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handles the booking CRUD functionality
func RegisterBookingRoutes(router *gin.Engine, bookingService interfaces.BookingService, authMiddleware interfaces.GatewayAuthMiddleware) {
	// Public routes
	router.DELETE("/bookings/:ID", func(ctx *gin.Context) {
		bookingIDString := ctx.Param("ID")

		// Convert string to int
		bookingID, err := strconv.Atoi(bookingIDString)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bookingID"})
			return
		}

		success, err := bookingService.DeleteByBookingID(bookingID)
		if err != nil {
			if _, ok := err.(*errors.BookingNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()}) // 404 Not Found
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		// Uses success to confirm the deletion
		if success {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Booking deleted successfully",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete Booking, but no error has occurred",
			})
		}
	})

	router.POST("/bookings", func(ctx *gin.Context) {
		var booking models.Booking

		if err := ctx.ShouldBindJSON(&booking); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postBooking, err := bookingService.Create(booking)
		if err != nil {
			// 409 Conflict
			if _, ok := err.(*errors.BookingExistsError); ok {
				ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
				return
			}
			// 500 Internal Server Error
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, postBooking)
	})

	bookingGroup := router.Group("/bookings")
	bookingGroup.Use(authMiddleware.GatewayAuthMiddleware())

	// Protected routes
	// Can only be accessible by the logged in user (userID)
	bookingGroup.GET("/", func(ctx *gin.Context) {
		userIDRaw, _ := ctx.Get("user_id")

		userID, ok := userIDRaw.(int)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not a string"})
			return
		}

		bookings := bookingService.GetByUserID(userID)

		// Check that the bookings userID matches the logged in user
		for _, booking := range bookings {
			if userID != booking.UserID {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot access the bookings belonging to another user"})
				return
			}
		}

		ctx.JSON(http.StatusOK, bookings)
	})

	bookingGroup.PUT("/", func(ctx *gin.Context) {
		userIDRaw, _ := ctx.Get("user_id")

		userID, ok := userIDRaw.(int)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not a string"})
			return
		}

		var booking models.Booking
		// Convert JSON to a Booking object
		if err := ctx.ShouldBindJSON(&booking); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the booking userID matches the logged in user
		if userID != booking.UserID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot access the bookings belonging to another user"})
			return
		}

		put_booking, err := bookingService.Update(booking)
		if err != nil {
			if _, ok := err.(*errors.BookingNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, put_booking)
	})
}
