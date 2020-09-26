package core

import (
	"github.com/pkg/errors"
)

var (
	ErrorUserNotFound      = errors.Errorf("User was not found")
	ErrorInvalidSort       = errors.Errorf("Invalid sort value")
	ErrorInsufficientFunds = errors.Errorf("User has insufficient funds")
)
