// repository/pg_account_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/leesache/financial-api/model"

	"github.com/jackc/pgx/v5"
)

type AccountRepository interface {
	GetAccount(id int) (*model.Account, error)
	TransferFunds(fromID, toID int, amount float64) error
	GetTransactionHistory(accountID int) ([]*model.Transaction, error)
}

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// PGAccountRepository implements AccountRepository using pgx.
type PGAccountRepository struct {
	db *pgx.Conn
}

func NewPGAccountRepository(db *pgx.Conn) *PGAccountRepository {
	return &PGAccountRepository{db: db}
}

func (r *PGAccountRepository) GetAccount(id int) (*model.Account, error) {
	var account model.Account
	err := r.db.QueryRow(context.Background(), `
        SELECT id, name, balance FROM accounts WHERE id = $1
    `, id).Scan(&account.ID, &account.Name, &account.Balance)

	if err == pgx.ErrNoRows {
		return nil, ErrAccountNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (r *PGAccountRepository) TransferFunds(fromID, toID int, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	log.Printf("Starting transaction for fund transfer: fromID=%d, toID=%d, amount=%.2f", fromID, toID, amount)

	// Deduct from sender's account
	var fromBalance float64
	err = tx.QueryRow(ctx, `
        SELECT balance FROM accounts WHERE id = $1 FOR UPDATE
    `, fromID).Scan(&fromBalance)

	if err == pgx.ErrNoRows {
		return ErrAccountNotFound
	} else if err != nil {
		return fmt.Errorf("failed to get sender account (id=%d): %w", fromID, err)
	}

	if fromBalance < amount {
		return ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, `
        UPDATE accounts SET balance = balance - $1 WHERE id = $2
    `, amount, fromID)
	if err != nil {
		return fmt.Errorf("failed to update sender account (id=%d): %w", fromID, err)
	}

	// Add to receiver's account
	_, err = tx.Exec(ctx, `
        UPDATE accounts SET balance = balance + $1 WHERE id = $2
    `, amount, toID)
	if err != nil {
		return fmt.Errorf("failed to update receiver account (id=%d): %w", toID, err)
	}

	// Record the transaction in the transactions table
	_, err = tx.Exec(ctx, `
        INSERT INTO transactions (from_id, to_id, amount) VALUES ($1, $2, $3)
    `, fromID, toID, amount)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Transaction completed successfully: fromID=%d, toID=%d, amount=%.2f", fromID, toID, amount)
	return nil
}

func (r *PGAccountRepository) GetTransactionHistory(accountID int) ([]*model.Transaction, error) {
	rows, err := r.db.Query(context.Background(), `
        SELECT id, from_id, to_id, amount, created_at
        FROM transactions
        WHERE from_id = $1 OR to_id = $1
        ORDER BY created_at DESC
    `, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction history: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.FromID, &t.ToID, &t.Amount, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, &t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating over transaction rows: %w", rows.Err())
	}

	return transactions, nil
}
