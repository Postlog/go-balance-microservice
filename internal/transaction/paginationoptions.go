package transaction

import (
	"errors"
	"fmt"
	"github.com/postlog/go-balance-microservice/internal/utils"
	"strconv"
	"strings"
)

type PaginationOptions struct {
	orderBy, orderDirection, limit, offset string
}

func NewPaginationOptions(limit, offset int, column, direction string) (*PaginationOptions, error) {
	opts := &PaginationOptions{
		orderBy:        "date",
		orderDirection: "asc",
		limit:          "NULL",
		offset:         "0",
	}

	var err error
	setError := func(someErr error) {
		if someErr != nil {
			err = someErr
		}
	}
	if column != "" {
		setError(opts.setOrdering(column))
	}
	if direction != "" {
		setError(opts.setOrderDirection(direction))
	}
	if limit != 0 {
		setError(opts.setLimit(limit))
	}
	if offset != 0 {
		setError(opts.setOffset(offset))
	}

	if err != nil {
		return nil, err
	}

	return opts, err
}

func (o *PaginationOptions) setLimit(limit int) error {
	if limit < 0 {
		return errors.New("limit cannot be less than 0")
	}

	o.limit = strconv.Itoa(limit)
	return nil
}

func (o *PaginationOptions) setOffset(offset int) error {
	if offset < 0 {
		return errors.New("offset cannot be less than 0")
	}

	o.offset = strconv.Itoa(offset)
	return nil
}

func (o *PaginationOptions) setOrdering(column string) error {
	column = strings.ToLower(column)
	if !utils.StringInCollection(column, "amount", "date") {
		return fmt.Errorf("column name must be equal \"amount\" or \"date\", not \"%s\"", column)
	}

	o.orderBy = column
	return nil
}

func (o *PaginationOptions) setOrderDirection(direction string) error {
	direction = strings.ToLower(direction)
	if !utils.StringInCollection(direction, "desc", "asc") {
		return fmt.Errorf("order direction must be equal \"asc\" or \"desc\", not \"%s\"", direction)
	}

	o.orderDirection = direction
	return nil
}

func (o *PaginationOptions) ToSQLClosure() string {
	return fmt.Sprintf("ORDER BY %s %s LIMIT %s OFFSET %s", o.orderBy, o.orderDirection, o.limit, o.offset)
}
