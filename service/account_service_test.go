package service

import (
	"testing"

	"github.com/leesache/financial-api/mock"
	"github.com/leesache/financial-api/model"
	"github.com/stretchr/testify/assert"
)

func TestGetAccount(t *testing.T) {
	// Arrange
	mockRepo := mock.NewMockAccountRepository()
	mockRepo.Accounts[1] = &model.Account{ID: 1, Name: "Alice", Balance: 1000.0}

	service := NewAccountService(mockRepo)

	// Act
	account, err := service.GetAccount(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, "Alice", account.Name)
	assert.Equal(t, 1000.0, account.Balance)
}

func TestGetAccountNotFound(t *testing.T) {
	// Arrange
	mockRepo := mock.NewMockAccountRepository()
	service := NewAccountService(mockRepo)

	// Act
	account, err := service.GetAccount(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, "account not found", err.Error())
}

func TestTransferFunds(t *testing.T) {
	// Arrange
	mockRepo := mock.NewMockAccountRepository()
	mockRepo.Accounts[1] = &model.Account{ID: 1, Name: "Alice", Balance: 1000.0}
	mockRepo.Accounts[2] = &model.Account{ID: 2, Name: "Bob", Balance: 500.0}

	service := NewAccountService(mockRepo)

	// Act
	err := service.TransferFunds(1, 2, 450.0)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 550.0, mockRepo.Accounts[1].Balance)
	assert.Equal(t, 950.0, mockRepo.Accounts[2].Balance)
}

func TestTransferFundsInsufficientFunds(t *testing.T) {
	// Arrange
	mockRepo := mock.NewMockAccountRepository()
	mockRepo.Accounts[1] = &model.Account{ID: 1, Name: "Alice", Balance: 100.0}
	mockRepo.Accounts[2] = &model.Account{ID: 2, Name: "Bob", Balance: 500.0}

	service := NewAccountService(mockRepo)

	// Act
	err := service.TransferFunds(1, 2, 450.0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
}

func TestGetTransactionHistory(t *testing.T) {
	// Arrange
	mockRepo := mock.NewMockAccountRepository()
	mockRepo.Accounts[1] = &model.Account{ID: 1, Name: "Alice", Balance: 1000.0}
	mockRepo.Accounts[2] = &model.Account{ID: 2, Name: "Bob", Balance: 500.0}
	mockRepo.TransferFunds(1, 2, 450.0) // Perform a transfer to populate transaction history

	service := NewAccountService(mockRepo)

	// Act
	transactions, err := service.GetTransactionHistory(1)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, 1, transactions[0].FromID)
	assert.Equal(t, 2, transactions[0].ToID)
	assert.Equal(t, 450.0, transactions[0].Amount)
}
