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

	findRecordOrRaiseError(&ticketOption, "id", ctx)

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
	var quantity uint64
	var err error

	findRecordOrRaiseError(&ticketOption, "id", ctx)
	findRecordOrRaiseError(&user, "user_id", ctx)

	quantity, err = strconv.ParseUint(ctx.Param("quantity"), 10, 32)
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

func findRecordOrRaiseError[T any](record *T, param_identifier string, ctx *gin.Context) {
	err := models.DB.First(&record, "id = ?", ctx.Param(param_identifier)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
}

func main() {
	models.ConnectToDatabase()

	r := gin.Default()

	r.GET("/ticket_options/:id", getTicketOption)
	r.POST("/ticket_options", createTicketOption)
	r.POST("/ticket_options/:id/purchases", purchaseTickets)

	r.Run("localhost:3000")
}
