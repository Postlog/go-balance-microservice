package types

import "context"

type TransactionWrapper func(ctx context.Context, f func(ctx context.Context) error) error
