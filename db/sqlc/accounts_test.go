package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/santinofajardo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func deleteAccountByID(accountID int64) {
	testStore.DeleteAccount(context.Background(), accountID)
}

func createRandomAccount(t *testing.T) Account { // This function doesn't has the 'Test' prefix so will doesn't run with the tests
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)

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

	account2, err := testStore.GetAccount(context.Background(), account1.ID)

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

	account2, err := testStore.UpdateAccount(
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

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	// Tests
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	accounts, err := testStore.ListAccounts(context.Background(), ListAccountsParams{lastAccount.Owner, 5, 0})

	// Tests
	require.NoError(t, err)
	require.Len(t, accounts, 1)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
		deleteAccountByID(account.ID)
	}
}
