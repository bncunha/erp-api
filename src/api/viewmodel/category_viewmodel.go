package viewmodel

import "github.com/bncunha/erp-api/src/domain"

type GetCategoryViewModel struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func ToGetCategoryViewModel(category domain.Category) GetCategoryViewModel {
	return GetCategoryViewModel{
		Id:   category.Id,
		Name: category.Name,
	}
}