package tests

import (
	"bytes"
	"encoding/json"
	"golang_ticket_api/handlers"
	"golang_ticket_api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicketOption_Success(t *testing.T) {
	setupTestDB()

	r := gin.Default()
	r.POST("/ticket_options", handlers.CreateTicketOption)

	ticketOptionInput := models.TicketOptionInput{
		Name:       "Test Ticket",
		Desc:       "Test ticket description",
		Allocation: 50,
	}

	jsonData, _ := json.Marshal(ticketOptionInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ticket_options", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TicketOption
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Ticket", response.Name)
	assert.Equal(t, "Test ticket description", response.Desc)
	assert.Equal(t, uint64(50), response.Allocation)
}

func TestCreateTicketOption_InvalidJSON(t *testing.T) {
	setupTestDB()

	r := gin.Default()
	r.POST("/ticket_options", handlers.CreateTicketOption)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ticket_options", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTicketOption_MissingFields(t *testing.T) {
	setupTestDB()

	r := gin.Default()
	r.POST("/ticket_options", handlers.CreateTicketOption)

	ticketOptionInput := models.TicketOptionInput{
		Name: "Test Ticket",
	}

	jsonData, _ := json.Marshal(ticketOptionInput)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ticket_options", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TicketOption
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Ticket", response.Name)
	assert.Equal(t, "", response.Desc)
	assert.Equal(t, uint64(0), response.Allocation)
}
