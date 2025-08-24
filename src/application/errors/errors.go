package errors

import (
	"errors"
	"strings"

	"github.com/lib/pq"
)

func New(message string) error {
	return errors.New(message)
}

func IsUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func IsNoRowsFinded(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}

func IsDuplicated(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
