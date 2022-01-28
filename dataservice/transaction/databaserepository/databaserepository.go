package databaserepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/postlog/go-balance-microservice/dataservice/models"
	"github.com/postlog/go-balance-microservice/pkg/database"
	"github.com/postlog/go-balance-microservice/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// repository implements Repository interface
type repository struct {
	db database.Database
}

func New(db database.Database) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, t models.Transaction) error {
	query := "INSERT INTO transaction (sender_id, receiver_id, amount, description, date) VALUES ($1, $2, $3, $4, $5)"
	return r.db.Exec(ctx, query, t.SenderId, t.ReceiverId, t.Amount, t.Description, t.Date)
}

func (r *repository) Get(
	ctx context.Context,
	userId uuid.UUID,
	count, startFrom int, orderBy, orderDirection string,
) ([]models.Transaction, error) {
	closure, err := convertToSQLClosure(count, startFrom, orderBy, orderDirection)
	if err != nil {
		return nil, err
	}
	query := "SELECT sender_id, receiver_id, amount, description, date FROM transaction " +
		"WHERE (sender_id=$1 OR receiver_id=$1) " + closure

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Transaction{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

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

		transactions = append(transactions, models.Transaction{
			SenderId: senderId, ReceiverId: receiverId, Amount: amount, Description: description, Date: date,
		})
	}

	return transactions, nil
}

func convertToSQLClosure(limit, offset int, column, direction string) (string, error) {
	o := struct {
		orderBy, orderDirection, limit, offset string
	}{
		"date", "asc", "NULL", "0",
	}

	if limit != 0 {
		if limit < 0 {
			return "", errors.New("limit cannot be less than 0")
		}
		o.limit = strconv.Itoa(limit)
	}

	if offset != 0 {
		if offset < 0 {
			return "", errors.New("offset cannot be less than 0")
		}
		o.offset = strconv.Itoa(offset)
	}

	column = strings.ToLower(column)
	if column != "" {
		if !utils.StringInCollection(column, "amount", "date") {
			return "", fmt.Errorf("column name must be equal \"amount\" or \"date\", not \"%s\"", column)
		}
		o.orderBy = column
	}

	direction = strings.ToLower(direction)
	if direction != "" {
		if !utils.StringInCollection(direction, "amount", "date") {
			return "", fmt.Errorf("order direction must be equal \"asc\" or \"desc\", not \"%s\"", direction)
		}
		o.orderDirection = direction
	}

	return fmt.Sprintf("ORDER BY %s %s LIMIT %s OFFSET %s", o.orderBy, o.orderDirection, o.limit, o.offset), nil
}
