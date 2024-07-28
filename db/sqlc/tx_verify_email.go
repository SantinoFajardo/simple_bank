package db

import (
	"context"
	"database/sql"
	"fmt"
)

type VerifyEmailTransactionParams struct {
	SecretCode string
	EmailId    int64
}

type VerifyEmailTransactionResult struct {
	VerifyEmail VerifyEmail
	User        User
}

func (store *SQLStore) VerifyEmailTransaction(ctx context.Context, args VerifyEmailTransactionParams) (VerifyEmailTransactionResult, error) {
	var result VerifyEmailTransactionResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         args.EmailId,
			SecretCode: args.SecretCode,
		})
		if err != nil {
			return fmt.Errorf("error updating verify email: %w", err)
		}

		user, err := q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}
		result.User = user

		return nil
	})

	return result, err
}
