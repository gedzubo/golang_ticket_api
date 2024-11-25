package main

import (
	"golang_ticket_api/models"

	"github.com/gin-gonic/gin"
)

func getTicketOption(ctx *gin.Context) {
	var ticketOption models.TicketOption

	result := models.DB.First(&ticketOption, ctx.Param("id"))
	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error})
		return
	}

	ctx.JSON(200, ticketOption)
}

func main() {
	models.ConnectToDatabase()

	r := gin.Default()

	r.GET("/ticket_options/:id", getTicketOption)

	r.Run("localhost:3000")
}
