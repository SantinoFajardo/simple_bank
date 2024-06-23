package db

import (
	"context"
	"testing"
	"time"

	"github.com/santinofajardo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntrie(t *testing.T) Entry {
	account := createRandomAccount(t)
	args := CreateEntrieParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	defer deleteAccountByID(account.ID)

	entrie, err := testQueries.CreateEntrie(context.Background(), args)

	// Tests
	require.NoError(t, err)
	require.NotEmpty(t, entrie)
	require.Equal(t, entrie.Amount, args.Amount)
	require.Equal(t, args.AccountID, entrie.AccountID)
	require.NotZero(t, entrie.ID)
	require.NotZero(t, entrie.CreatedAt)

	return entrie
}

func TestCreateEntrie(t *testing.T) {
	createRandomAccount(t)
}

func TestGetEntrie(t *testing.T) {
	entrie1 := createRandomEntrie(t)
	entrie2, err := testQueries.GetEntrie(context.Background(), entrie1.ID)

	// Tests
	require.NoError(t, err)
	require.NotEmpty(t, entrie2)
	require.Equal(t, entrie1.AccountID, entrie2.AccountID)
	require.Equal(t, entrie1.Amount, entrie2.Amount)
	require.WithinDuration(t, entrie1.CreatedAt, entrie2.CreatedAt, time.Second)
}
