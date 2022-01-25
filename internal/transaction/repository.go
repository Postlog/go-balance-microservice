package transaction

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/internal/database"
	"time"
)

type Transaction struct {
	SenderId    uuid.NullUUID `json:"senderId"`
	ReceiverId  uuid.NullUUID `json:"receiverId"`
	Amount      float64       `json:"amount"`
	Description string        `json:"description"`
	Date        time.Time     `json:"date"`
}

type Repository interface {
	Get(ctx context.Context, userId uuid.UUID, pg *PaginationOptions) ([]Transaction, error)
	Create(ctx context.Context, t Transaction) error
}

func NewRepository(db *database.Database) Repository {
	return &postgresRepository{db}
}

type postgresRepository struct {
	db *database.Database
}

var NoTransactionsErr = errors.New("user with specified id has no transactions")

func (r *postgresRepository) Create(ctx context.Context, t Transaction) error {
	query := "INSERT INTO transaction (sender_id, receiver_id, amount, description, date) VALUES ($1, $2, $3, $4, $5)"
	return r.db.Exec(ctx, query, t.SenderId, t.ReceiverId, t.Amount, t.Description, t.Date)
}

func (r *postgresRepository) Get(ctx context.Context, userId uuid.UUID, pg *PaginationOptions) ([]Transaction, error) {
	query := "SELECT sender_id, receiver_id, amount, description, date " +
		"FROM transaction WHERE (sender_id=$1 OR receiver_id=$1) " + pg.ToSQLClosure()
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NoTransactionsErr
		}
		return nil, err
	}

	var transactions []Transaction

	for rows.Next() {
		var (
			senderId, receiverId uuid.NullUUID
			amount               float64
			description          string
			date                 time.Time
		)
		err = rows.Scan(&senderId, &receiverId, &amount, &description, &date)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, Transaction{
			SenderId: senderId, ReceiverId: receiverId, Amount: amount, Description: description, Date: date,
		})
	}

	return transactions, nil
}
