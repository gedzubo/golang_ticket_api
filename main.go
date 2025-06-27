package main

import (
	"golang_ticket_api/handlers"
	"golang_ticket_api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectToDatabase()

	r := gin.Default()

	r.GET("/ticket_options/:id", handlers.GetTicketOption)
	r.POST("/ticket_options", handlers.CreateTicketOption)
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	r.Run("localhost:3000")
}
