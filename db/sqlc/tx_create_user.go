package db

import "context"

type CreateUserTransactionParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTransactionResult struct {
	User User
}

func (store *SQLStore) CreateUserTransaction(ctx context.Context, args CreateUserTransactionParams) (CreateUserTransactionResult, error) {
	var result CreateUserTransactionResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, args.CreateUserParams)
		if err != nil {
			return err
		}

		return args.AfterCreate(result.User)
	})

	return result, err
}
