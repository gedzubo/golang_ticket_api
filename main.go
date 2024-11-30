package main

import (
	"fmt"
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket option with provided ID does not exist"})
		return
	}

	ctx.JSON(http.StatusOK, ticketOption)
}

func createTicketOption(ctx *gin.Context) {
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

func purchaseTickets(ctx *gin.Context) {
	var ticketOption models.TicketOption
	var user models.User
	var quantity uint64
	var err error

	err = models.DB.First(&ticketOption, "id = ?", ctx.Param("id")).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket option with provided ID does not exist"})
		return
	}

	err = models.DB.First(&user, "id = ?", ctx.Param("user_id")).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user with provided ID does not exist"})
		return
	}

	quantity, err = strconv.ParseUint(ctx.Param("quantity"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid quantity"})
	}

	err = models.DB.Transaction(func(tx *gorm.DB) error {
		models.DB.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&ticketOption)

		currentQuantity := ticketOption.Allocation

		if currentQuantity < quantity {
			return fmt.Errorf("We don't have enough tickets to complete your purchase")
		}

		tx.Model(&ticketOption).Update("Allocation", currentQuantity-quantity)

		purchase := models.Purchase{Quantity: quantity, UserId: user.ID, TicketOptionId: ticketOption.ID}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		for i := 1; i <= int(quantity); i++ {
			if err := tx.Create(&models.Ticket{PurchaseId: purchase.ID, TicketOptionId: ticketOption.ID}).Error; err != nil {
				return err
			}
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
