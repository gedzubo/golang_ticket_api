package handlers

import (
	"golang_ticket_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTicketOption(ctx *gin.Context) {
	var ticketOptionInput models.TicketOptionInput

	if err := ctx.ShouldBindJSON(&ticketOptionInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketOption := models.TicketOption{
		Name:       ticketOptionInput.Name,
		Desc:       ticketOptionInput.Desc,
		Allocation: ticketOptionInput.Allocation,
	}
	models.DB.Create(&ticketOption)

	ctx.JSON(http.StatusOK, ticketOption)
}