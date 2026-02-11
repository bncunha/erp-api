package repository

import (
	"context"
	"testing"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

func TestCreateManySaleItemWithEmptySlice(t *testing.T) {
	repo := &salesRepository{}
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(1))

	ids, err := repo.CreateManySaleItem(ctx, nil, domain.Sales{}, []domain.SalesItem{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected no ids, got %v", ids)
	}
}
