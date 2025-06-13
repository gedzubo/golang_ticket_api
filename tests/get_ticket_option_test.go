package tests

import (
	"fmt"
	"golang_ticket_api/handlers"
	"golang_ticket_api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTicketOption_Success(t *testing.T) {
	r := gin.Default()
	r.GET("/ticket_options/:id", handlers.GetTicketOption)

	db := setupTestDB()
	ticketOption := models.TicketOption{
		Name:       "First Ticket Option",
		Desc:       "This is the first ticket option",
		Allocation: 100,
	}
	db.Create(&ticketOption)
	fmt.Println("Ticket Option Created:", ticketOption.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ticket_options/"+ticketOption.ID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
