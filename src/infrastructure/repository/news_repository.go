package repository

import (
	"context"
	"database/sql"

	"github.com/bncunha/erp-api/src/domain"
)

type newsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) domain.NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) GetLatestVisible(ctx context.Context, tenantId int64, role domain.Role) (domain.News, error) {
	var news domain.News
	query := `
SELECT n.id, n.content_html, n.created_at
FROM news n
JOIN news_roles nr ON nr.news_id = n.id
WHERE (n.tenant_id = $1 OR n.tenant_id IS NULL)
  AND nr.role = $2
ORDER BY n.created_at DESC, n.id DESC
LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, tenantId, role).Scan(&news.Id, &news.ContentHtml, &news.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return news, domain.ErrNewsNotFound
		}
		return news, err
	}

	return news, nil
}
