package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTransaction(ctx context.Context, args TransferTransactionParams) (TransferTransactionResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{db: db, Queries: New(db)}
}

func moveMoney(
	ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	return
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v\nRollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTransactionParams contains the input parameters of the transfer transaction
type TransferTransactionParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTransactionResult is the result of the transfer transaction
type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTransaction perform a money transfer from one account to the other.
// It creates a transfer record, add accounts entries, and update accounts balance within a single database transaction
func (store *SQLStore) TransferTransaction(ctx context.Context, args TransferTransactionParams) (TransferTransactionResult, error) {
	var result TransferTransactionResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountId,
			ToAccountID:   args.ToAccountId,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntrie(ctx, CreateEntrieParams{
			AccountID: args.FromAccountId,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntrie(ctx, CreateEntrieParams{
			AccountID: args.ToAccountId,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		if args.FromAccountId < args.ToAccountId {
			result.FromAccount, result.ToAccount, err = moveMoney(ctx, q, args.FromAccountId, -args.Amount, args.ToAccountId, args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = moveMoney(ctx, q, args.ToAccountId, args.Amount, args.FromAccountId, -args.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
