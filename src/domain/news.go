package domain

import (
	"errors"
	"time"
)

var ErrNewsNotFound = errors.New("Novidade não encontrada")

type News struct {
	Id          int64
	ContentHtml string
	CreatedAt   time.Time
}
