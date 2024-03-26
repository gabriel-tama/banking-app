package transaction

import (
	"context"
	"time"

	"github.com/gabriel-tama/banking-app/common/db"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Add(ctx context.Context, req *AddBalancePayload, userId int) error
	Reduce(ctx context.Context, req *ReduceBalancePayload, userId int) error
	Get(ctx context.Context, userId int) (*GetBalanceResponse, error)
	GetHistory(ctx context.Context, req *GetHistoryPayload, userId int) (*GetHistoryResponse, int, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Add(ctx context.Context, req *AddBalancePayload, userId int) error {
	var balanceId int
	err := d.db.StartTx(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"INSERT INTO balances (userId, balance, currencyCode) VALUES($1, $2, $3) ON CONFLICT (userId, currencyCode) DO UPDATE SET balance=balances.balance+$4 RETURNING balances.id",
			userId, req.AddBalance, req.Currency, req.AddBalance).Scan(&balanceId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO transactions (userId, amount, bankAccountNumber, bankName, transferProofImg, currencyCode, balanceId) VALUES($1, $2, $3, $4, $5, $6, $7)",
			userId, req.AddBalance, req.Sender, req.SenderBank, req.TransferProof, req.Currency, balanceId)
		return err
	})

	return err
}

func (d *dbRepository) Reduce(ctx context.Context, req *ReduceBalancePayload, userId int) error {
	var balanceId int
	err := d.db.StartTx(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			"UPDATE balances SET balance=balance-$2 WHERE userId=$1 AND currencyCode=$3 RETURNING balances.id",
			userId, req.AddBalance, req.Currency).Scan(&balanceId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO transactions (userId, amount, bankAccountNumber, bankName, currencyCode, balanceId) VALUES($1, $2, $3, $4, $5, $6)",
			userId, req.AddBalance, req.Sender, req.SenderBank, req.Currency, balanceId)
		return err
	})

	return err
}

func (d *dbRepository) Get(ctx context.Context, userId int) (*GetBalanceResponse, error) {
	rows, err := d.db.Pool.Query(ctx, "SELECT balance, currencyCode FROM balances WHERE userId=$1", userId)
	defer rows.Close()
	var res GetBalanceResponse

	for rows.Next() {
		var balance SingleBalanceRes
		err = rows.Scan(&balance.Amount, &balance.CurrencyCode)
		if err != nil {
			return nil, err
		}
		res = append(res, balance)
	}

	if err != nil {
		return nil, err
	}

	return &res, err
}

func (d *dbRepository) GetHistory(ctx context.Context, req *GetHistoryPayload, userId int) (*GetHistoryResponse, int, error) {
	stmt := `WITH UserTransactions AS (
            SELECT id, amount,  bankAccountNumber, bankName, COALESCE(transferProofImg,''), currencyCode, createdAt
            FROM transactions
            WHERE userId = $1
            ORDER BY createdAt DESC
            LIMIT $2
            OFFSET $3
        )
        SELECT *,
            (SELECT COUNT(*) FROM transactions WHERE userId = $1) AS total_count
        FROM UserTransactions;
        `
	rows, err := d.db.Pool.Query(ctx, stmt, userId, req.Limit, req.Offset)
	defer rows.Close()
	if err != nil {
		return nil, 0, err
	}
	var res GetHistoryResponse
	var total int
	for rows.Next() {
		var createdAt time.Time
		var history SingleHistoryRes

		err = rows.Scan(&history.TransactionId, &history.Balance, &history.SourceTransaction.BankAccountNumber, &history.SourceTransaction.BankName, &history.TransferProof, &history.CurrencyCode, &createdAt, &total)
		history.CreatedAt = createdAt.Unix()
		res = append(res, history)
		if err != nil {
			return nil, 0, err
		}
	}
	return &res, total, err

}
