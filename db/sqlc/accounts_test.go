package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/santinofajardo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func deleteAccountByID(accountID int64) {
	testQueries.DeleteAccount(context.Background(), accountID)
}

func createRandomAccount(t *testing.T) Account { // This function doesn't has the 'Test' prefix so will doesn't run with the tests
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	// Tests
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreate(t *testing.T) {
	account := createRandomAccount(t)
	deleteAccountByID(account.ID)

}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	defer deleteAccountByID(account1.ID)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	// Tests
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	defer deleteAccountByID(account1.ID)

	account2, err := testQueries.UpdateAccount(
		context.Background(),
		UpdateAccountParams{account1.ID, util.RandomMoney()})

	// Tests
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	// Tests
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{5, 5})

	// Tests
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		deleteAccountByID(account.ID)
	}
}
