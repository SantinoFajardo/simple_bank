package db

import "context"

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
