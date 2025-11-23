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

func IsForeignKeyViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503"
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

func ParseDuplicatedMessage(title string, err error) error {
	pqError, ok := err.(*pq.Error)
	if !ok {
		return errors.New(title + " já cadastrado!")
	}
	if strings.Contains(pqError.Detail, "phone_number") {
		return errors.New(title + " já cadastrado com este telefone!")
	}
	if strings.Contains(pqError.Detail, "email") {
		return errors.New(title + " já cadastrado com este email!")
	}
	return errors.New(title + " já cadastrado!")
}
