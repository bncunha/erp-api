package viewmodel

import (
	"testing"
	"time"

	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

func TestToSaleByIdViewModelIncludesSkuIdOnItems(t *testing.T) {
	sale := output.GetSaleByIdOutput{
		Id:            1,
		Code:          "V-1",
		Date:          time.Now(),
		TotalValue:    100,
		SellerName:    "Seller",
		CustomerName:  "Customer",
		ReceivedValue: 50,
		FutureRevenue: 50,
		PaymentStatus: domain.PaymentStatusPending,
	}

	itemsOutput := []output.GetItemsOutput{
		{
			Sku: domain.Sku{
				Id:    42,
				Code:  "SKU-42",
				Price: 25,
				Product: domain.Product{
					Name: "T-Shirt",
				},
				Color: "Blue",
				Size:  "M",
			},
			Quantity: 2,
		},
	}

	viewModel := ToSaleByIdViewModel(sale, nil, itemsOutput, nil)
	if len(viewModel.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(viewModel.Items))
	}
	if viewModel.Items[0].SkuId != 42 {
		t.Fatalf("expected sku_id 42, got %d", viewModel.Items[0].SkuId)
	}
}
