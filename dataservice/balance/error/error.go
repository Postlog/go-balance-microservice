package error

import "errors"

var NotFoundErr = errors.New("user with specified id not found")
