package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang_ticket_api/handlers"
	"golang_ticket_api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type PurchaseRequest struct {
	Quantity uint64 `json:"quantity"`
	UserID   string `json:"user_id"`
}

func TestPurchaseTickets_Success(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()
	
	user := models.User{
		Username: "testuser",
	}
	db.Create(&user)

	ticketOption := models.TicketOption{
		Name:       "Concert Ticket",
		Desc:       "General admission",
		Allocation: 100,
	}
	db.Create(&ticketOption)

	purchaseReq := PurchaseRequest{
		Quantity: 5,
		UserID:   user.ID,
	}
	jsonData, _ := json.Marshal(purchaseReq)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/ticket_options/%s/purchases", ticketOption.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Tickets purchased successfully", response["message"])

	// Verify ticket allocation was reduced
	var updatedTicketOption models.TicketOption
	db.First(&updatedTicketOption, "id = ?", ticketOption.ID)
	assert.Equal(t, uint64(95), updatedTicketOption.Allocation)

	// Verify purchase was created
	var purchase models.Purchase
	db.First(&purchase, "user_id = ? AND ticket_option_id = ?", user.ID, ticketOption.ID)
	assert.Equal(t, uint64(5), purchase.Quantity)

	// Verify tickets were created
	var tickets []models.Ticket
	db.Find(&tickets, "purchase_id = ?", purchase.ID)
	assert.Equal(t, 5, len(tickets))
}

func TestPurchaseTickets_TicketOptionNotFound(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()
	
	user := models.User{
		Username: "testuser",
	}
	db.Create(&user)

	purchaseReq := PurchaseRequest{
		Quantity: 5,
		UserID:   user.ID,
	}
	jsonData, _ := json.Marshal(purchaseReq)

	w := httptest.NewRecorder()
	url := "/ticket_options/invalid-id/purchases"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ticket option with provided ID does not exist", response["error"])
}

func TestPurchaseTickets_UserNotFound(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()

	ticketOption := models.TicketOption{
		Name:       "Concert Ticket",
		Desc:       "General admission",
		Allocation: 100,
	}
	db.Create(&ticketOption)

	purchaseReq := PurchaseRequest{
		Quantity: 5,
		UserID:   "invalid-user-id",
	}
	jsonData, _ := json.Marshal(purchaseReq)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/ticket_options/%s/purchases", ticketOption.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user with provided ID does not exist", response["error"])
}

func TestPurchaseTickets_InvalidJSON(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()

	ticketOption := models.TicketOption{
		Name:       "Concert Ticket",
		Desc:       "General admission",
		Allocation: 100,
	}
	db.Create(&ticketOption)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/ticket_options/%s/purchases", ticketOption.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid request body")
}

func TestPurchaseTickets_ZeroQuantity(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()
	
	user := models.User{
		Username: "testuser",
	}
	db.Create(&user)

	ticketOption := models.TicketOption{
		Name:       "Concert Ticket",
		Desc:       "General admission",
		Allocation: 100,
	}
	db.Create(&ticketOption)

	purchaseReq := PurchaseRequest{
		Quantity: 0,
		UserID:   user.ID,
	}
	jsonData, _ := json.Marshal(purchaseReq)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/ticket_options/%s/purchases", ticketOption.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid request body")
}

func TestPurchaseTickets_InsufficientTickets(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()
	
	user := models.User{
		Username: "testuser",
	}
	db.Create(&user)

	ticketOption := models.TicketOption{
		Name:       "Concert Ticket",
		Desc:       "General admission",
		Allocation: 5,
	}
	db.Create(&ticketOption)

	purchaseReq := PurchaseRequest{
		Quantity: 10,
		UserID:   user.ID,
	}
	jsonData, _ := json.Marshal(purchaseReq)

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/ticket_options/%s/purchases", ticketOption.ID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "insufficient tickets available")
}

func TestPurchaseTickets_MissingFields(t *testing.T) {
	r := gin.Default()
	r.POST("/ticket_options/:id/purchases", handlers.PurchaseTickets)

	db := setupTestDB()

	ticketOption := models.TicketOption{
		Name:       "Concert Ticket",
		Desc:       "General admission",
		Allocation: 100,
	}
	db.Create(&ticketOption)

	tests := []struct {
		name        string
		reqBody     string
		expectedMsg string
	}{
		{
			name:        "Missing user_id",
			reqBody:     `{"quantity": 5}`,
			expectedMsg: "invalid request body",
		},
		{
			name:        "Missing quantity",
			reqBody:     `{"user_id": "test-user-id"}`,
			expectedMsg: "invalid request body",
		},
		{
			name:        "Empty user_id",
			reqBody:     `{"quantity": 5, "user_id": ""}`,
			expectedMsg: "invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/ticket_options/%s/purchases", ticketOption.ID)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(tt.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], tt.expectedMsg)
		})
	}
}