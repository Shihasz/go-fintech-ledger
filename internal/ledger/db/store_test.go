package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> Before: Account 1 = %d, Account 2 = %d\n", account1.Balance, account2.Balance)

	// Run 5 concurrent transfer transactions
	n := 5
	amount := int64(10)

	// Channels to communicate between goroutines and the main test thread
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Verify the results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check basic transfer details
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)

		// Verify the balances were updated correctly in the transaction
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		fmt.Printf(">> Tx %d: Account 1 = %d, Account 2 = %d\n", i+1, fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // Should be a multiple of the transfer amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balances
	updatedAccount1, _ := store.GetAccount(context.Background(), account1.ID)
	updatedAccount2, _ := store.GetAccount(context.Background(), account2.ID)

	fmt.Printf(">> After: Account 1 = %d, Account 2 = %d\n", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run 10 concurrent transfers (5 from A->B, 5 from B->A)
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		// Reverse direction for half of the transactions
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// Wait for all 10 goroutines to finish
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err) // If a deadlock occurred, this would fail or hang forever
	}

	// Because we sent 50 to B, and 50 back to A, the final balances should remain completely unchanged!
	updatedAccount1, _ := store.GetAccount(context.Background(), account1.ID)
	updatedAccount2, _ := store.GetAccount(context.Background(), account2.ID)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
