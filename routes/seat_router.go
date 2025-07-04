package routes

import (
	"flyhorizons-bookingservice/services/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterSeatRoutes(router *gin.Engine, seatService interfaces.SeatService) {
	router.GET("/bookings/seats/:flightCode", func(ctx *gin.Context) {
		flightCode := ctx.Param("flightCode")
		seats, err := seatService.GetByFlightCode(flightCode)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, seats)
	})
}
