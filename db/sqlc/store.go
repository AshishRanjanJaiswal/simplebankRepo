package db

import (
	"context"
	"database/sql"
	"fmt"
)

// store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // Starts the transaction
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q) // Callback method
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: % v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer from one account to other
//It creates the transfer record, add account enteries, and update accounts, balance within a singles database transaction

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entery   `json:"from_entry"`
	ToEntry     Entery   `json:"to_entry"`
}

// var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account enteries, and updatte accounts'balance
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err = store.execTx(ctx, func(q *Queries) error {
		var err error
		// txName := ctx.Value(txKey)
		// fmt.Println(txName, "Create Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, "Create Entry 1")
		result.FromEntry, err = q.CreateEnteries(ctx, CreateEnteriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}
		// fmt.Println(txName, "Create Entry 2")
		result.ToEntry, err = q.CreateEnteries(ctx, CreateEnteriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO:updatae account's balance

		// fmt.Println(txName, "Get Account For Update 1")

		// fmt.Println(txName, "Update account 1")
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, arg.Amount)
		}

		// fmt.Println(txName, "Get Account For Update 2")

		// fmt.Println(txName, "Update account 2")

		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
