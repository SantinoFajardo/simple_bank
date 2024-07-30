package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTransaction(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run concurrent transfer transactions
	n := 5
	amount := int64(10)

	// Create chanels to store the errors and results between the goroutines
	errs := make(chan error)
	results := make(chan TransferTransactionResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			// Fill the channels when the result and error is taken
			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		// Check tranfers records
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, -fromEntry.Amount, amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntrie(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntrie(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		// Check balances
		diff1 := account1.Balance - fromAccount.Balance // 50 - 40 = 10
		diff2 := toAccount.Balance - account2.Balance   // 30 + 40 = 10
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, account1.Balance-amount*int64(n), updatedAccount.Balance)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance+amount*int64(n), updatedAccount2.Balance)
}

func TestTransferTransactionDeadLock(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run concurrent transfer transactions
	n := 10
	amount := int64(10)

	// Create chanels to store the errors and results between the goroutines
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := testStore.TransferTransaction(context.Background(), TransferTransactionParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})

			// Fill the channels when the result and error is taken
			errs <- err
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	// check the final updated balances
	updatedAccount, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, account1.Balance, updatedAccount.Balance)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
