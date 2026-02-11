package service

import (
	"context"
	"errors"
	"testing"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

func TestSalesServiceCreateSales(t *testing.T) {
	useCase := &stubSalesUseCase{}
	repo := &stubSalesRepository{}
	service := NewSalesService(useCase, repo, &stubInventoryRepository{})

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(42))
	firstInstallment := time.Now().AddDate(0, 0, 10)
	installments := 3
	req := request.CreateSaleRequest{
		CustomerId: 99,
		Items: []request.CreateSaleRequestItems{{
			SkuId:    11,
			Quantity: 2,
		}},
		Payments: []request.CreateSaleRequestPayments{{
			PaymentType:          domain.PaymentTypeCreditStore,
			Value:                100,
			InstallmentsQuantity: &installments,
			FirstInstallmentDate: &firstInstallment,
		}},
	}

	if err := service.CreateSales(ctx, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if useCase.receivedInput.CustomerId != req.CustomerId {
		t.Fatalf("expected customer id %d, got %d", req.CustomerId, useCase.receivedInput.CustomerId)
	}
	if useCase.receivedInput.UserId != 42 {
		t.Fatalf("expected user id 42, got %d", useCase.receivedInput.UserId)
	}
	if len(useCase.receivedInput.Items) != 1 || useCase.receivedInput.Items[0].SkuId != req.Items[0].SkuId {
		t.Fatalf("unexpected items input: %+v", useCase.receivedInput.Items)
	}
	if len(useCase.receivedInput.Payments) != 1 {
		t.Fatalf("unexpected payments length: %d", len(useCase.receivedInput.Payments))
	}
	dates := useCase.receivedInput.Payments[0].Dates
	if len(dates) != installments {
		t.Fatalf("expected %d installments, got %d", installments, len(dates))
	}
	if !dates[0].DueDate.Equal(firstInstallment) {
		t.Fatalf("expected first due date %v, got %v", firstInstallment, dates[0].DueDate)
	}
	total := 0.0
	for _, d := range dates {
		total += d.InstallmentValue
	}
	if total != req.Payments[0].Value {
		t.Fatalf("expected total %.2f, got %.2f", req.Payments[0].Value, total)
	}
	if dates[0].InstallmentValue < dates[2].InstallmentValue {
		t.Fatalf("expected rounding remainder applied to first installment: %+v", dates)
	}
}

func TestSalesServiceCreateSalesValidationError(t *testing.T) {
	service := NewSalesService(&stubSalesUseCase{}, &stubSalesRepository{}, &stubInventoryRepository{})

	err := service.CreateSales(context.Background(), request.CreateSaleRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSalesServiceGetSales(t *testing.T) {
	repo := &stubSalesRepository{
		getSalesOutput: []output.GetSalesItemOutput{{Id: 1}},
	}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleAdmin))
	req := request.ListSalesRequest{UserId: []int64{7}}

	out, err := service.GetSales(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Sales) != 1 {
		t.Fatalf("expected 1 sale, got %d", len(out.Sales))
	}
	if len(repo.getSalesInput.UserId) != 1 || repo.getSalesInput.UserId[0] != 7 {
		t.Fatalf("expected repository to receive user id 7, got %+v", repo.getSalesInput.UserId)
	}
}

func TestSalesServiceGetSalesPermissionDenied(t *testing.T) {
	service := NewSalesService(&stubSalesUseCase{}, &stubSalesRepository{}, &stubInventoryRepository{})
	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, "")

	_, err := service.GetSales(ctx, request.ListSalesRequest{})
	if err != ErrPermissionDenied {
		t.Fatalf("expected permission denied error, got %v", err)
	}
}

func TestSalesServiceGetById(t *testing.T) {
	repo := &stubSalesRepository{
		saleByIdOutput: output.GetSaleByIdOutput{Id: 2},
		paymentsOutput: []output.GetSalesPaymentOutput{{PaymentType: domain.PaymentTypeCash}, {PaymentType: domain.PaymentTypeCash}},
		itemsOutput:    []output.GetItemsOutput{{Sku: domain.Sku{Id: 3}}},
	}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})

	sale, payments, items, returnsOutput, err := service.GetById(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sale.Id != 2 {
		t.Fatalf("expected sale id 2, got %d", sale.Id)
	}
	if len(payments) != 1 || len(payments[0].Installments) != 2 {
		t.Fatalf("expected grouped payments, got %+v", payments)
	}
	if len(items) != 1 || items[0].Sku.Id != 3 {
		t.Fatalf("unexpected items output: %+v", items)
	}
	if len(returnsOutput) != 0 {
		t.Fatalf("expected no returns, got %+v", returnsOutput)
	}
}

func TestSalesServiceChangePaymentStatus(t *testing.T) {
	now := time.Now()
	repo := &stubSalesRepository{
		paymentDateBySaleAndPaymentId: domain.SalesPaymentDates{Id: 5},
	}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})

	req := request.ChangePaymentStatusRequest{
		Status: string(domain.PaymentStatusPaid),
		Date:   now,
	}

	err := service.ChangePaymentStatus(context.Background(), 10, 20, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.changePaymentStatusCalledWith.id != 20 || repo.changePaymentStatusCalledWith.status != domain.PaymentStatus(req.Status) {
		t.Fatalf("unexpected status change args: %+v", repo.changePaymentStatusCalledWith)
	}
	if repo.changePaymentDateCalledWith.date == nil || !repo.changePaymentDateCalledWith.date.Equal(now) {
		t.Fatalf("expected payment date to be set: %+v", repo.changePaymentDateCalledWith)
	}
}

func TestSalesServiceChangePaymentStatusValidationError(t *testing.T) {
	service := NewSalesService(&stubSalesUseCase{}, &stubSalesRepository{}, &stubInventoryRepository{})

	err := service.ChangePaymentStatus(context.Background(), 1, 2, request.ChangePaymentStatusRequest{})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSalesServiceGroupPaymentsByPaymentType(t *testing.T) {
	service := &salesService{}
	payments := []output.GetSalesPaymentOutput{
		{PaymentType: domain.PaymentTypeCash, InstallmentNumber: 1},
		{PaymentType: domain.PaymentTypeCash, InstallmentNumber: 2},
		{PaymentType: domain.PaymentTypeCreditCard, InstallmentNumber: 1},
	}

	groups := service.groupPaymentsByPaymentType(payments)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if len(groups[0].Installments)+len(groups[1].Installments) != 3 {
		t.Fatalf("unexpected grouping result: %+v", groups)
	}
}

func TestSalesServiceCalculateTotalValueRemainder(t *testing.T) {
	service := &salesService{}
	out := service.calculateTotalValue(10.0, 3)
	if len(out) != 3 {
		t.Fatalf("expected 3 installments")
	}
	if out[0] != 3.34 || out[1] != 3.33 || out[2] != 3.33 {
		t.Fatalf("unexpected installments: %+v", out)
	}
}

func TestSalesServiceGetSalesNonAdminUsesContextUser(t *testing.T) {
	repo := &stubSalesRepository{
		getSalesOutput: []output.GetSalesItemOutput{},
	}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})

	ctx := context.WithValue(context.Background(), constants.ROLE_KEY, string(domain.UserRoleReseller))
	ctx = context.WithValue(ctx, constants.USERID_KEY, float64(55))

	if _, err := service.GetSales(ctx, request.ListSalesRequest{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.getSalesInput.UserId) != 1 || repo.getSalesInput.UserId[0] != 55 {
		t.Fatalf("expected reseller user id propagated, got %+v", repo.getSalesInput.UserId)
	}
}

func TestSalesServiceChangePaymentStatusUsesNilDateWhenNotPaid(t *testing.T) {
	repo := &stubSalesRepository{
		paymentDateBySaleAndPaymentId: domain.SalesPaymentDates{},
	}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})

	req := request.ChangePaymentStatusRequest{
		Status: string(domain.PaymentStatusPending),
	}

	if err := service.ChangePaymentStatus(context.Background(), 1, 2, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.changePaymentDateCalledWith.date != nil {
		t.Fatalf("expected nil date when status not paid")
	}
}

func TestSalesServiceCreateSalesUsesDefaultDueDate(t *testing.T) {
	useCase := &stubSalesUseCase{}
	service := NewSalesService(useCase, &stubSalesRepository{}, &stubInventoryRepository{})

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(9))
	req := request.CreateSaleRequest{
		CustomerId: 1,
		Items: []request.CreateSaleRequestItems{{
			SkuId:    1,
			Quantity: 1,
		}},
		Payments: []request.CreateSaleRequestPayments{{
			PaymentType: domain.PaymentTypeCash,
			Value:       12,
		}},
	}

	if err := service.CreateSales(ctx, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(useCase.receivedInput.Payments) != 1 || len(useCase.receivedInput.Payments[0].Dates) != 1 {
		t.Fatalf("expected single installment: %+v", useCase.receivedInput.Payments)
	}
	due := useCase.receivedInput.Payments[0].Dates[0].DueDate
	if time.Since(due) > time.Second {
		t.Fatalf("expected due date to default to current time, got %v", due)
	}
}

func TestSalesServiceCreateSalesUseCaseError(t *testing.T) {
	expected := errors.New("boom")
	useCase := &stubSalesUseCase{err: expected}
	service := NewSalesService(useCase, &stubSalesRepository{}, &stubInventoryRepository{})

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(1))
	installments := 1
	req := request.CreateSaleRequest{
		CustomerId: 1,
		Items:      []request.CreateSaleRequestItems{{SkuId: 1, Quantity: 1}},
		Payments: []request.CreateSaleRequestPayments{{
			PaymentType:          domain.PaymentTypeCreditStore,
			Value:                10,
			InstallmentsQuantity: &installments,
			FirstInstallmentDate: func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
		}},
	}

	if err := service.CreateSales(ctx, req); err != expected {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestSalesServiceGetByIdErrors(t *testing.T) {
	repo := &stubSalesRepository{saleByIdErr: errors.New("fail")}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})
	if _, _, _, _, err := service.GetById(context.Background(), 1); err == nil {
		t.Fatalf("expected error")
	}
	if !repo.getSaleByIdCalled {
		t.Fatalf("expected sale by id to be called")
	}
	if repo.getPaymentsCalled || repo.getItemsCalled {
		t.Fatalf("expected no further repository calls")
	}

	repo = &stubSalesRepository{saleByIdOutput: output.GetSaleByIdOutput{Id: 1}, paymentsErr: errors.New("payments")}
	service = NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})
	if _, _, _, _, err := service.GetById(context.Background(), 1); err == nil {
		t.Fatalf("expected payments error")
	}
	if !repo.getPaymentsCalled || repo.getItemsCalled {
		t.Fatalf("expected payments call only")
	}

	repo = &stubSalesRepository{
		saleByIdOutput: output.GetSaleByIdOutput{Id: 1},
		paymentsOutput: []output.GetSalesPaymentOutput{},
		itemsErr:       errors.New("items"),
	}
	service = NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})
	if _, _, _, _, err := service.GetById(context.Background(), 1); err == nil {
		t.Fatalf("expected items error")
	}
	if !repo.getItemsCalled {
		t.Fatalf("expected items to be requested")
	}
}

func TestSalesServiceChangePaymentStatusErrors(t *testing.T) {
	repo := &stubSalesRepository{paymentDateErr: errors.New("dates")}
	service := NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})
	req := request.ChangePaymentStatusRequest{Status: string(domain.PaymentStatusPending)}
	if err := service.ChangePaymentStatus(context.Background(), 1, 2, req); err == nil {
		t.Fatalf("expected error from payment date lookup")
	}
	if repo.changePaymentStatusCalledWith.id != 0 {
		t.Fatalf("expected change status not called")
	}

	repo = &stubSalesRepository{paymentDateBySaleAndPaymentId: domain.SalesPaymentDates{}, changePaymentStatusErr: errors.New("status")}
	service = NewSalesService(&stubSalesUseCase{}, repo, &stubInventoryRepository{})
	if err := service.ChangePaymentStatus(context.Background(), 1, 2, req); err == nil {
		t.Fatalf("expected status change error")
	}
}

func TestSalesServiceCreateReturnResellerUsesOwnInventory(t *testing.T) {
	useCase := &stubSalesUseCase{}
	inventoryRepo := &stubInventoryRepository{
		getByUser: domain.Inventory{Id: 77},
	}
	service := NewSalesService(useCase, &stubSalesRepository{}, inventoryRepo)

	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(42))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, string(domain.UserRoleReseller))

	req := request.CreateSalesReturnRequest{
		ReturnerName: "Cliente",
		Reason:       "Defeito",
		Items: []request.CreateSalesReturnItemRequest{
			{SkuId: 1, Quantity: 1},
		},
	}

	if err := service.CreateReturn(ctx, 10, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if useCase.receivedReturnInput.InventoryDestinationId != 77 {
		t.Fatalf("expected destination inventory from reseller, got %d", useCase.receivedReturnInput.InventoryDestinationId)
	}
}

func TestSalesServiceCreateReturnAdminRequiresInventory(t *testing.T) {
	service := NewSalesService(&stubSalesUseCase{}, &stubSalesRepository{}, &stubInventoryRepository{})
	ctx := context.WithValue(context.Background(), constants.USERID_KEY, float64(1))
	ctx = context.WithValue(ctx, constants.ROLE_KEY, string(domain.UserRoleAdmin))

	req := request.CreateSalesReturnRequest{
		ReturnerName: "Cliente",
		Reason:       "Troca",
		Items: []request.CreateSalesReturnItemRequest{
			{SkuId: 1, Quantity: 1},
		},
	}

	if err := service.CreateReturn(ctx, 1, req); err == nil {
		t.Fatalf("expected admin inventory validation error")
	}
}
