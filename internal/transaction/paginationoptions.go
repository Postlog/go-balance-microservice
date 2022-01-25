package transaction

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PaginationOptions struct {
	orderBy, orderDirection, limit, offset string
}

type PaginationOption func(*PaginationOptions) error

func NewPaginationOptions(opts ...PaginationOption) (*PaginationOptions, error) {
	pg := &PaginationOptions{
		orderBy:        "date",
		orderDirection: "asc",
		limit:          "NULL",
		offset:         "0",
	}

	for _, opt := range opts {
		err := opt(pg)
		if err != nil {
			return nil, err
		}
	}
	return pg, nil
}

func (o *PaginationOptions) ToSQLClosure() string {
	orderClosure := ""
	if o.orderBy != "" {
		orderClosure = fmt.Sprintf("ORDER BY %s %s", o.orderBy, o.orderDirection)
	}

	return fmt.Sprintf("%s LIMIT %s OFFSET %s", orderClosure, o.limit, o.offset)
}

func WithOrdering(columnName string) PaginationOption {
	columnName = strings.ToLower(columnName)
	return func(pg *PaginationOptions) error {
		if columnName != "amount" && columnName != "date" {
			return fmt.Errorf("column name must be equal \"amount\" or \"date\", not \"%s\"", columnName)
		}
		pg.orderBy = columnName
		return nil
	}
}

func WithDirection(dir string) PaginationOption {
	dir = strings.ToLower(dir)
	return func(pg *PaginationOptions) error {
		if dir != "asc" && dir != "desc" {
			return fmt.Errorf("order direction must be equal \"asc\" or \"desc\", not \"%s\"", dir)
		}
		pg.orderDirection = dir
		return nil
	}
}

func WithLimit(limit int) PaginationOption {
	return func(pg *PaginationOptions) error {
		if limit < 0 {
			return errors.New("limit cannot be less than 0")
		}
		pg.limit = strconv.Itoa(limit)
		return nil
	}
}

func WithOffset(offset int) PaginationOption {
	return func(pg *PaginationOptions) error {
		if offset < 0 {
			return errors.New("offset cannot be less than 0")
		}
		pg.offset = strconv.Itoa(offset)
		return nil
	}
}
