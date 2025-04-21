package handlers

import (
	"golang_ticket_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTicketOption(ctx *gin.Context) {
	var ticketOption models.TicketOption

	err := models.DB.First(&ticketOption, "id = ?", ctx.Param("id")).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket option with provided ID does not exist"})
		return
	}

	ctx.JSON(http.StatusOK, ticketOption)
}
