package mock

import (
	"errors"
	"time"

	"github.com/leesache/financial-api/model"
)

type MockAccountRepository struct {
	Accounts     map[int]*model.Account
	Transactions []*model.Transaction
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		Accounts:     make(map[int]*model.Account),
		Transactions: []*model.Transaction{},
	}
}

func (m *MockAccountRepository) GetAccount(id int) (*model.Account, error) {
	account, exists := m.Accounts[id]
	if !exists {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func (m *MockAccountRepository) TransferFunds(fromID, toID int, amount float64) error {
	fromAccount, exists := m.Accounts[fromID]
	if !exists || fromAccount.Balance < amount {
		return errors.New("insufficient funds")
	}

	toAccount, exists := m.Accounts[toID]
	if !exists {
		return errors.New("receiver account not found")
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount

	var mockTimestamp = time.Date(2023, 10, 7, 12, 0, 0, 0, time.UTC)

	m.Transactions = append(m.Transactions, &model.Transaction{
		FromID:    fromID,
		ToID:      toID,
		Amount:    amount,
		CreatedAt: mockTimestamp,
	})

	return nil
}

func (m *MockAccountRepository) GetTransactionHistory(accountID int) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	for _, t := range m.Transactions {
		if t.FromID == accountID || t.ToID == accountID {
			transactions = append(transactions, t)
		}
	}
	return transactions, nil
}
