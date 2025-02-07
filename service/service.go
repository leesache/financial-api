package service

import (
	"errors"

	"github.com/leesache/financial-api/model"
	"github.com/leesache/financial-api/repository"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// AccountService defines the interface for account-related business logic.
type AccountService interface {
	GetAccount(id int) (*model.Account, error)
	TransferFunds(fromID, toID int, amount float64) error
	GetTransactionHistory(accountID int) ([]*model.Transaction, error)
}

// accountService implements AccountService using a repository.
type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

// GetAccount retrieves an account by ID.
func (s *accountService) GetAccount(id int) (*model.Account, error) {
	account, err := s.repo.GetAccount(id)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

// TransferFunds transfers funds between two accounts.
func (s *accountService) TransferFunds(fromID, toID int, amount float64) error {
	// Validate input
	if fromID <= 0 || toID <= 0 {
		return errors.New("invalid account IDs")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Delegate to the repository
	err := s.repo.TransferFunds(fromID, toID, amount)
	if err != nil {
		return err
	}
	return nil
}

// GetTransactionHistory retrieves the transaction history for a specific account.
func (s *accountService) GetTransactionHistory(accountID int) ([]*model.Transaction, error) {
	if accountID <= 0 {
		return nil, errors.New("invalid account ID")
	}

	transactions, err := s.repo.GetTransactionHistory(accountID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
