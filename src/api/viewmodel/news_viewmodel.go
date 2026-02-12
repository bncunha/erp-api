package viewmodel

import (
	"time"

	"github.com/bncunha/erp-api/src/domain"
)

type NewsViewModel struct {
	Id          int64  `json:"id"`
	ContentHtml string `json:"content_html"`
	CreatedAt   string `json:"created_at"`
}

func ToNewsViewModel(news domain.News) NewsViewModel {
	return NewsViewModel{
		Id:          news.Id,
		ContentHtml: news.ContentHtml,
		CreatedAt:   news.CreatedAt.Format(time.RFC3339),
	}
}
