package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/bncunha/erp-api/src/domain"
)

func TestNewsRepositoryGetLatestVisibleSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewNewsRepository(db)
	now := time.Now()
	queryRegex := regexp.QuoteMeta(`
SELECT n.id, n.content_html, n.created_at
FROM news n
JOIN news_roles nr ON nr.news_id = n.id
WHERE (n.tenant_id = $1 OR n.tenant_id IS NULL)
  AND nr.role = $2
ORDER BY n.created_at DESC, n.id DESC
LIMIT 1`)

	rows := sqlmock.NewRows([]string{"id", "content_html", "created_at"}).
		AddRow(int64(5), "<h1>Atualização</h1>", now)
	mock.ExpectQuery(queryRegex).
		WithArgs(int64(2), domain.UserRoleAdmin).
		WillReturnRows(rows)

	news, err := repo.GetLatestVisible(context.Background(), 2, domain.UserRoleAdmin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if news.Id != 5 || news.ContentHtml != "<h1>Atualização</h1>" {
		t.Fatalf("unexpected news: %+v", news)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestNewsRepositoryGetLatestVisibleUnexpectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewNewsRepository(db)
	queryRegex := regexp.QuoteMeta(`
SELECT n.id, n.content_html, n.created_at
FROM news n
JOIN news_roles nr ON nr.news_id = n.id
WHERE (n.tenant_id = $1 OR n.tenant_id IS NULL)
  AND nr.role = $2
ORDER BY n.created_at DESC, n.id DESC
LIMIT 1`)

	mock.ExpectQuery(queryRegex).
		WithArgs(int64(2), domain.UserRoleReseller).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetLatestVisible(context.Background(), 2, domain.UserRoleReseller)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestNewsRepositoryGetLatestVisibleReturnsDomainNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := &newsRepository{db: db}
	queryRegex := regexp.QuoteMeta(`
SELECT n.id, n.content_html, n.created_at
FROM news n
JOIN news_roles nr ON nr.news_id = n.id
WHERE (n.tenant_id = $1 OR n.tenant_id IS NULL)
  AND nr.role = $2
ORDER BY n.created_at DESC, n.id DESC
LIMIT 1`)

	mock.ExpectQuery(queryRegex).
		WithArgs(int64(2), domain.UserRoleReseller).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content_html", "created_at"}))

	_, err = repo.GetLatestVisible(context.Background(), 2, domain.UserRoleReseller)
	if err == nil || err != domain.ErrNewsNotFound {
		t.Fatalf("expected ErrNewsNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
