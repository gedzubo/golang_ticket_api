package main

import (
	"errors"
	"golang_ticket_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getTicketOption(ctx *gin.Context) {
	var ticketOption models.TicketOption

	result := models.DB.First(&ticketOption, "id = ?", ctx.Param("id"))
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
			return
		}
	}

	ctx.JSON(http.StatusOK, ticketOption)
}

func createTicketOption(ctx *gin.Context) {
	var ticketOption models.TicketOption

	if err := ctx.ShouldBindJSON(&ticketOption); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ticketOption)
}

func main() {
	models.ConnectToDatabase()

	r := gin.Default()

	r.GET("/ticket_options/:id", getTicketOption)
	r.POST("/ticket_options", createTicketOption)

	r.Run("localhost:3000")
}
