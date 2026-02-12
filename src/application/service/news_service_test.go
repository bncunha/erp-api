package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/domain"
)

type stubNewsRepository struct {
	output         domain.News
	err            error
	lastTenantId   int64
	lastRole       domain.Role
	getLatestCalls int
}

func (s *stubNewsRepository) GetLatestVisible(ctx context.Context, tenantId int64, role domain.Role) (domain.News, error) {
	s.getLatestCalls++
	s.lastTenantId = tenantId
	s.lastRole = role
	if s.err != nil {
		return domain.News{}, s.err
	}
	return s.output, nil
}

func TestNewsServiceGetLatestSuccess(t *testing.T) {
	repo := &stubNewsRepository{
		output: domain.News{
			Id:          10,
			ContentHtml: "<h1>Nova funcionalidade</h1>",
			CreatedAt:   time.Now(),
		},
	}
	service := NewNewsService(repo)
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(7))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, string(domain.UserRoleAdmin))

	news, err := service.GetLatest(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if news.Id != 10 {
		t.Fatalf("expected id 10, got %d", news.Id)
	}
	if repo.lastTenantId != 7 || repo.lastRole != domain.UserRoleAdmin {
		t.Fatalf("expected tenant 7/admin, got %d/%s", repo.lastTenantId, repo.lastRole)
	}
}

func TestNewsServiceGetLatestWithFloatTenant(t *testing.T) {
	repo := &stubNewsRepository{output: domain.News{Id: 1}}
	service := NewNewsService(repo)
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, float64(3))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, string(domain.UserRoleReseller))

	_, err := service.GetLatest(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.lastTenantId != 3 || repo.lastRole != domain.UserRoleReseller {
		t.Fatalf("expected tenant 3/reseller, got %d/%s", repo.lastTenantId, repo.lastRole)
	}
}

func TestNewsServiceGetLatestNotFound(t *testing.T) {
	repo := &stubNewsRepository{err: domain.ErrNewsNotFound}
	service := NewNewsService(repo)
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(1))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, string(domain.UserRoleAdmin))

	_, err := service.GetLatest(ctx)
	if !errors.Is(err, domain.ErrNewsNotFound) {
		t.Fatalf("expected ErrNewsNotFound, got %v", err)
	}
}

func TestNewsServiceGetLatestInvalidTenant(t *testing.T) {
	repo := &stubNewsRepository{}
	service := NewNewsService(repo)
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, "invalid")
	ctx = context.WithValue(ctx, constants.ROLE_KEY, string(domain.UserRoleAdmin))

	_, err := service.GetLatest(ctx)
	if err == nil || err.Error() != "tenant id invalido" {
		t.Fatalf("expected tenant id invalido, got %v", err)
	}
}

func TestNewsServiceGetLatestInvalidRole(t *testing.T) {
	repo := &stubNewsRepository{}
	service := NewNewsService(repo)
	ctx := context.WithValue(context.Background(), constants.TENANT_KEY, int64(1))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, "")

	_, err := service.GetLatest(ctx)
	if err == nil || err.Error() != "role invalida" {
		t.Fatalf("expected role invalida, got %v", err)
	}
}
