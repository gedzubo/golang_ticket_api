package handlers

import (
	"errors"
	"fmt"
	"golang_ticket_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseRequest struct {
	Quantity uint64 `json:"quantity" binding:"required,min=1"`
	UserID   string `json:"user_id" binding:"required"`
}

func PurchaseTickets(ctx *gin.Context) {
	ticketOptionID := ctx.Param("id")
	if ticketOptionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ticket option ID is required"})
		return
	}

	var req PurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	quantity := req.Quantity
	userID := req.UserID

	var ticketOption models.TicketOption
	if err := models.DB.First(&ticketOption, "id = ?", ticketOptionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket option with provided ID does not exist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ticket option"})
		}
		return
	}

	var user models.User
	if err := models.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user with provided ID does not exist"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		}
		return
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		var lockedTicketOption models.TicketOption
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&lockedTicketOption, "id = ?", ticketOptionID).Error; err != nil {
			return fmt.Errorf("failed to lock ticket option: %w", err)
		}

		if lockedTicketOption.Allocation < quantity {
			return fmt.Errorf("insufficient tickets available: requested %d, available %d",
				quantity, lockedTicketOption.Allocation)
		}

		if err := tx.Model(&lockedTicketOption).
			Update("allocation", lockedTicketOption.Allocation-quantity).Error; err != nil {
			return fmt.Errorf("failed to update ticket allocation: %w", err)
		}

		purchase := models.Purchase{
			Quantity:       quantity,
			UserId:         user.ID,
			TicketOptionId: ticketOption.ID,
		}
		if err := tx.Create(&purchase).Error; err != nil {
			return fmt.Errorf("failed to create purchase: %w", err)
		}

		tickets := make([]models.Ticket, quantity)
		for i := range tickets {
			tickets[i] = models.Ticket{
				PurchaseId:     purchase.ID,
				TicketOptionId: ticketOption.ID,
			}
		}

		if err := tx.CreateInBatches(&tickets, 100).Error; err != nil {
			return fmt.Errorf("failed to create tickets: %w", err)
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tickets purchased successfully"})
}
