// account_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leesache/financial-api/service"
)

// AccountHandler defines the handler for account-related HTTP requests.
type AccountHandler struct {
	service service.AccountService
}

// NewAccountHandler initializes a new AccountHandler with the given service.
func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// GetAccount handles GET requests to retrieve an account by ID.
func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.service.GetAccount(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// TransferFunds handles POST requests to transfer funds between accounts.
func (h *AccountHandler) TransferFunds(c *gin.Context) {
	var input struct {
		FromID int     `json:"from_id"`
		ToID   int     `json:"to_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.TransferFunds(input.FromID, input.ToID, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer completed successfully"})
}
