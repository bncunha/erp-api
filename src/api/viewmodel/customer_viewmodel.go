package viewmodel

import "github.com/bncunha/erp-api/src/domain"

type GetAllCustomersViewModel struct {
    Id          int64  `json:"id"`
    Name        string `json:"name"`
    PhoneNumber string `json:"phone_number"`
}

func ToCustomerViewModel(customers []domain.Customer) []GetAllCustomersViewModel {
	var viewmodel []GetAllCustomersViewModel = make([]GetAllCustomersViewModel, 0)
	for _, customer := range customers {
		viewmodel = append(viewmodel, GetAllCustomersViewModel{
			Id:          customer.Id,
			Name:        customer.Name,
			PhoneNumber: customer.PhoneNumber,
		})
	}
    return viewmodel
}

func ToGetCustomerViewModel(customer domain.Customer) GetAllCustomersViewModel {
    return GetAllCustomersViewModel{
        Id:          customer.Id,
        Name:        customer.Name,
        PhoneNumber: customer.PhoneNumber,
    }
}
