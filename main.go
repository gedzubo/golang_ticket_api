package main

import (
	"errors"
	"golang_ticket_api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getTicketOption(ctx *gin.Context) {
	var ticketOption models.TicketOption

	err := models.DB.First(&ticketOption, "id = ?", ctx.Param("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
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

func purchaseTickets(ctx *gin.Context) {
	var ticketOption models.TicketOption
	var user models.User

	ticketOptionErr := models.DB.First(&ticketOption, "id = ?", ctx.Param("id")).Error
	if ticketOptionErr != nil {
		if errors.Is(ticketOptionErr, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": ticketOptionErr})
			return
		}
	}

	userErr := models.DB.First(&user, "id = ?", ctx.Param("user_id")).Error
	if userErr != nil {
		if errors.Is(userErr, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": userErr})
			return
		}
	}

	quantity, err := strconv.ParseUint(ctx.Param("quantity"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide correct quantity"})
	}

	err = models.DB.Transaction(func(tx *gorm.DB) error {
		models.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&ticketOption)

		// check if we have enough available slots
		// update ticket option quantity
		// create a new purchase
		// create tickets based on quantity

		if err := tx.Create(&models.Purchase{Quantity: uint16(quantity)}).Error; err != nil {
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func main() {
	models.ConnectToDatabase()

	r := gin.Default()

	r.GET("/ticket_options/:id", getTicketOption)
	r.POST("/ticket_options", createTicketOption)
	r.POST("/ticket_options/:id/purchases", purchaseTickets)

	r.Run("localhost:3000")
}
