package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

type stubNewsService struct {
	output domain.News
	err    error
}

func (s *stubNewsService) GetLatest(ctx context.Context) (domain.News, error) {
	if s.err != nil {
		return domain.News{}, s.err
	}
	return s.output, nil
}

func TestNewsControllerGetLatestSuccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := NewNewsController(&stubNewsService{
		output: domain.News{
			Id:          1,
			ContentHtml: "<h1>Teste</h1>",
			CreatedAt:   time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC),
		},
	})

	if err := controller.GetLatest(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestNewsControllerGetLatestNoContent(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := NewNewsController(&stubNewsService{err: domain.ErrNewsNotFound})
	if err := controller.GetLatest(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestNewsControllerGetLatestBadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/news/latest", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := NewNewsController(&stubNewsService{err: errors.New("fail")})
	if err := controller.GetLatest(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
