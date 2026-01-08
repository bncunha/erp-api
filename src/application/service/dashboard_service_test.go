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

type stubDashboardRepository struct {
	revenueResponses       []float64
	revenueErr             error
	revenueCalls           int
	revenueInputs          []domain.DashboardQueryInput
	salesCountResponses    []int64
	salesCountErr          error
	salesCountCalls        int
	salesCountInputs       []domain.DashboardQueryInput
	revenueByDay           []domain.DashboardTimeSeriesItem
	revenueByDayErr        error
	salesCountByDay        []domain.DashboardTimeSeriesItem
	salesCountByDayErr     error
	stockTotal             float64
	stockTotalErr          error
	stockInput             domain.DashboardStockQueryInput
	lowStock               []domain.DashboardLowStockItem
	lowStockErr            error
	lowStockInput          domain.DashboardStockQueryInput
	revenueByReseller      []domain.DashboardResellerSalesItem
	revenueByResellerErr   error
	revenueByResellerInput domain.DashboardQueryInput
	topProducts            []domain.DashboardProductSalesItem
	topProductsErr         error
	topProductsInput       domain.DashboardQueryInput
	topProductsLimit       int
}

func (s *stubDashboardRepository) GetRevenue(ctx context.Context, input domain.DashboardQueryInput) (float64, error) {
	s.revenueInputs = append(s.revenueInputs, input)
	if s.revenueErr != nil {
		return 0, s.revenueErr
	}
	if s.revenueCalls < len(s.revenueResponses) {
		resp := s.revenueResponses[s.revenueCalls]
		s.revenueCalls++
		return resp, nil
	}
	s.revenueCalls++
	return 0, nil
}

func (s *stubDashboardRepository) GetSalesCount(ctx context.Context, input domain.DashboardQueryInput) (int64, error) {
	s.salesCountInputs = append(s.salesCountInputs, input)
	if s.salesCountErr != nil {
		return 0, s.salesCountErr
	}
	if s.salesCountCalls < len(s.salesCountResponses) {
		resp := s.salesCountResponses[s.salesCountCalls]
		s.salesCountCalls++
		return resp, nil
	}
	s.salesCountCalls++
	return 0, nil
}

func (s *stubDashboardRepository) GetRevenueByDay(ctx context.Context, input domain.DashboardQueryInput) ([]domain.DashboardTimeSeriesItem, error) {
	if s.revenueByDayErr != nil {
		return nil, s.revenueByDayErr
	}
	return s.revenueByDay, nil
}

func (s *stubDashboardRepository) GetSalesCountByDay(ctx context.Context, input domain.DashboardQueryInput) ([]domain.DashboardTimeSeriesItem, error) {
	if s.salesCountByDayErr != nil {
		return nil, s.salesCountByDayErr
	}
	return s.salesCountByDay, nil
}

func (s *stubDashboardRepository) GetStockTotal(ctx context.Context, input domain.DashboardStockQueryInput) (float64, error) {
	s.stockInput = input
	if s.stockTotalErr != nil {
		return 0, s.stockTotalErr
	}
	return s.stockTotal, nil
}

func (s *stubDashboardRepository) GetLowStockProducts(ctx context.Context, input domain.DashboardStockQueryInput) ([]domain.DashboardLowStockItem, error) {
	s.lowStockInput = input
	if s.lowStockErr != nil {
		return nil, s.lowStockErr
	}
	return s.lowStock, nil
}

func (s *stubDashboardRepository) GetRevenueByReseller(ctx context.Context, input domain.DashboardQueryInput) ([]domain.DashboardResellerSalesItem, error) {
	s.revenueByResellerInput = input
	if s.revenueByResellerErr != nil {
		return nil, s.revenueByResellerErr
	}
	return s.revenueByReseller, nil
}

func (s *stubDashboardRepository) GetTopProductsByReseller(ctx context.Context, input domain.DashboardQueryInput, limit int) ([]domain.DashboardProductSalesItem, error) {
	s.topProductsInput = input
	s.topProductsLimit = limit
	if s.topProductsErr != nil {
		return nil, s.topProductsErr
	}
	return s.topProducts, nil
}

func newDashboardService(repo *stubDashboardRepository, userRepo domain.UserRepository) DashboardService {
	return &dashboardService{
		dashboardRepository: repo,
		userRepository:      userRepo,
	}
}

func newDashboardRequest(enum domain.DashboardWidgetEnum) request.DashboardWidgetDataRequest {
	return request.DashboardWidgetDataRequest{
		Enum: string(enum),
		Period: request.DashboardWidgetPeriodRequest{
			From: "2026-01-01",
			To:   "2026-01-31",
		},
	}
}

func ctxWithRole(role domain.Role) context.Context {
	return context.WithValue(context.Background(), constants.ROLE_KEY, string(role))
}

func ctxWithRoleAndUser(role domain.Role, userId int64) context.Context {
	ctx := ctxWithRole(role)
	return context.WithValue(ctx, constants.USERID_KEY, float64(userId))
}

func TestDashboardServiceListWidgetsByRole(t *testing.T) {
	service := newDashboardService(&stubDashboardRepository{}, &stubUserRepository{})

	adminItems, err := service.ListWidgets(ctxWithRole(domain.UserRoleAdmin))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(adminItems) != 6 {
		t.Fatalf("expected 6 admin widgets, got %d", len(adminItems))
	}

	resellerItems, err := service.ListWidgets(ctxWithRole(domain.UserRoleReseller))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resellerItems) != 4 {
		t.Fatalf("expected 4 reseller widgets, got %d", len(resellerItems))
	}
}

func TestDashboardServiceGetWidgetDataFaturamento(t *testing.T) {
	repo := &stubDashboardRepository{
		revenueResponses: []float64{200, 100},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), newDashboardRequest(domain.DashboardWidgetFaturamento))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardCardData)
	if !ok {
		t.Fatalf("expected card data")
	}
	if data.Value != 200 || data.PreviousValue != 100 || data.DeltaPercent != 100 {
		t.Fatalf("unexpected card values: %+v", data)
	}
	if resp.Meta.Currency != "BRL" {
		t.Fatalf("expected BRL currency")
	}
}

func TestDashboardServiceGetWidgetDataTotalVendas(t *testing.T) {
	repo := &stubDashboardRepository{
		salesCountResponses: []int64{10, 5},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), newDashboardRequest(domain.DashboardWidgetTotalVendas))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardCardData)
	if !ok {
		t.Fatalf("expected card data")
	}
	if data.Value != 10 || data.PreviousValue != 5 {
		t.Fatalf("unexpected card values: %+v", data)
	}
}

func TestDashboardServiceGetWidgetDataProdutosEmEstoque(t *testing.T) {
	repo := &stubDashboardRepository{
		stockTotal: 45,
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), newDashboardRequest(domain.DashboardWidgetProdutosEmEstoque))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardCardData)
	if !ok {
		t.Fatalf("expected card data")
	}
	if data.Value != 45 || data.Unit != "UN" {
		t.Fatalf("unexpected card data: %+v", data)
	}
}

func TestDashboardServiceGetWidgetDataEstoqueBaixo(t *testing.T) {
	repo := &stubDashboardRepository{
		lowStock: []domain.DashboardLowStockItem{
			{ProductName: "A", Quantity: 0},
			{ProductName: "B", Quantity: 2},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), newDashboardRequest(domain.DashboardWidgetEstoqueBaixo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardTableData)
	if !ok {
		t.Fatalf("expected table data")
	}
	if len(data.Rows) != 2 || len(data.Columns) != 2 {
		t.Fatalf("unexpected table data: %+v", data)
	}
	if repo.lowStockInput.Threshold != 0 {
		t.Fatalf("expected threshold 0")
	}
}

func TestDashboardServiceGetWidgetDataFaturamentoNoTempo(t *testing.T) {
	repo := &stubDashboardRepository{
		revenueByDay: []domain.DashboardTimeSeriesItem{
			{Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Value: 10},
			{Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), Value: 20},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), newDashboardRequest(domain.DashboardWidgetFaturamentoNoTempo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardLineBarData)
	if !ok {
		t.Fatalf("expected line data")
	}
	if len(data.Labels) != 2 || data.Series[0].Values[1] != 20 {
		t.Fatalf("unexpected line data: %+v", data)
	}
}

func TestDashboardServiceGetWidgetDataVendasPorRevendedor(t *testing.T) {
	repo := &stubDashboardRepository{
		revenueByReseller: []domain.DashboardResellerSalesItem{
			{ResellerName: "A", Value: 100},
			{ResellerName: "B", Value: 50},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), newDashboardRequest(domain.DashboardWidgetVendasPorRevendedor))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardLineBarData)
	if !ok {
		t.Fatalf("expected bar data")
	}
	if len(data.Labels) != 2 || data.Series[0].Values[0] != 100 {
		t.Fatalf("unexpected bar data: %+v", data)
	}
}

func TestDashboardServiceGetWidgetDataMeusProdutosMaisVendidos(t *testing.T) {
	repo := &stubDashboardRepository{
		topProducts: []domain.DashboardProductSalesItem{
			{ProductName: "Camisa", Quantity: 5},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRoleAndUser(domain.UserRoleReseller, 7), newDashboardRequest(domain.DashboardWidgetMeusProdutosMaisVendidos))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardLineBarData)
	if !ok {
		t.Fatalf("expected bar data")
	}
	if data.Series[0].Values[0] != 5 || repo.topProductsLimit != 10 {
		t.Fatalf("unexpected bar data or limit: %+v", data)
	}
}

func TestDashboardServiceGetWidgetDataMinhasVendasNoTempo(t *testing.T) {
	repo := &stubDashboardRepository{
		salesCountByDay: []domain.DashboardTimeSeriesItem{
			{Date: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), Value: 2},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	resp, err := service.GetWidgetData(ctxWithRoleAndUser(domain.UserRoleReseller, 7), newDashboardRequest(domain.DashboardWidgetMinhasVendasNoTempo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardLineBarData)
	if !ok {
		t.Fatalf("expected line data")
	}
	if data.Series[0].Values[0] != 2 {
		t.Fatalf("unexpected line data: %+v", data)
	}
}

func TestDashboardServiceResellerIgnoresResellerFilter(t *testing.T) {
	repo := &stubDashboardRepository{
		revenueResponses: []float64{10, 5},
	}
	service := newDashboardService(repo, &stubUserRepository{})

	ctx := ctxWithRole(domain.UserRoleReseller)
	ctx = context.WithValue(ctx, constants.USERID_KEY, float64(7))

	req := newDashboardRequest(domain.DashboardWidgetMeuFaturamento)
	req.Filters = &request.DashboardWidgetFiltersRequest{ResellerId: func() *int64 { v := int64(99); return &v }()}

	if _, err := service.GetWidgetData(ctx, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.revenueInputs) == 0 || repo.revenueInputs[0].ResellerId == nil || *repo.revenueInputs[0].ResellerId != 7 {
		t.Fatalf("expected reseller id from context")
	}
}

func TestDashboardServiceGetWidgetDataPermissionDenied(t *testing.T) {
	service := newDashboardService(&stubDashboardRepository{}, &stubUserRepository{})

	_, err := service.GetWidgetData(ctxWithRole(domain.UserRoleReseller), newDashboardRequest(domain.DashboardWidgetFaturamento))
	if err == nil {
		t.Fatalf("expected permission error")
	}
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("expected ErrPermissionDenied")
	}
}

func TestDashboardServiceGetWidgetDataInvalidEnum(t *testing.T) {
	service := newDashboardService(&stubDashboardRepository{}, &stubUserRepository{})

	req := newDashboardRequest(domain.DashboardWidgetFaturamento)
	req.Enum = "INVALID"
	_, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), req)
	if err == nil {
		t.Fatalf("expected widget not found")
	}
	if !errors.Is(err, ErrDashboardWidgetNotFound) {
		t.Fatalf("expected ErrDashboardWidgetNotFound")
	}
}

func TestDashboardServiceAdminResellerValidation(t *testing.T) {
	userRepo := &stubUserRepository{getById: domain.User{Id: 8, Role: string(domain.UserRoleAdmin)}}
	service := newDashboardService(&stubDashboardRepository{}, userRepo)

	req := newDashboardRequest(domain.DashboardWidgetFaturamento)
	resellerId := int64(8)
	req.Filters = &request.DashboardWidgetFiltersRequest{ResellerId: &resellerId}

	_, err := service.GetWidgetData(ctxWithRole(domain.UserRoleAdmin), req)
	if err == nil {
		t.Fatalf("expected reseller validation error")
	}
	if !errors.Is(err, ErrDashboardResellerNotFound) {
		t.Fatalf("expected ErrDashboardResellerNotFound")
	}
}

func TestDashboardServiceParsePeriodInvalid(t *testing.T) {
	service := newDashboardService(&stubDashboardRepository{}, &stubUserRepository{})

	if _, err := service.(*dashboardService).parsePeriod(request.DashboardWidgetPeriodRequest{From: "invalid", To: "2026-01-01"}); err == nil {
		t.Fatalf("expected invalid period error")
	}
}

func TestDashboardServiceResolveResellerFilterAdminValid(t *testing.T) {
	userRepo := &stubUserRepository{
		getByIdResponses: map[int64]domain.User{12: {Id: 12, Role: string(domain.UserRoleReseller)}},
	}
	service := newDashboardService(&stubDashboardRepository{}, userRepo)
	resellerId := int64(12)

	id, err := service.(*dashboardService).resolveResellerFilter(ctxWithRole(domain.UserRoleAdmin), string(domain.UserRoleAdmin), &request.DashboardWidgetFiltersRequest{ResellerId: &resellerId})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || *id != 12 {
		t.Fatalf("expected reseller id 12")
	}
}

func TestDashboardServiceResolveResellerFilterResellerMissingUser(t *testing.T) {
	service := newDashboardService(&stubDashboardRepository{}, &stubUserRepository{})

	_, err := service.(*dashboardService).resolveResellerFilter(ctxWithRole(domain.UserRoleReseller), string(domain.UserRoleReseller), nil)
	if err == nil {
		t.Fatalf("expected permission error")
	}
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("expected ErrPermissionDenied")
	}
}

func TestDashboardServiceHandleMinhasVendas(t *testing.T) {
	repo := &stubDashboardRepository{
		salesCountResponses: []int64{12, 6},
	}
	service := newDashboardService(repo, &stubUserRepository{}).(*dashboardService)
	period, _ := service.parsePeriod(request.DashboardWidgetPeriodRequest{From: "2026-01-01", To: "2026-01-31"})
	resellerId := int64(3)

	resp, err := service.handleMinhasVendas(context.Background(), widgetInput{Period: period, ResellerId: &resellerId})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardCardData)
	if !ok {
		t.Fatalf("expected card data")
	}
	if data.Value != 12 || data.PreviousValue != 6 || data.DeltaPercent != 100 {
		t.Fatalf("unexpected card data: %+v", data)
	}
}

func TestDashboardServiceHandleMinhasVendasNoTempo(t *testing.T) {
	repo := &stubDashboardRepository{
		salesCountByDay: []domain.DashboardTimeSeriesItem{
			{Date: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), Value: 2},
			{Date: time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC), Value: 4},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{}).(*dashboardService)
	period, _ := service.parsePeriod(request.DashboardWidgetPeriodRequest{From: "2026-01-01", To: "2026-01-31"})
	resellerId := int64(3)

	resp, err := service.handleMinhasVendasNoTempo(context.Background(), widgetInput{Period: period, ResellerId: &resellerId})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardLineBarData)
	if !ok {
		t.Fatalf("expected line data")
	}
	if len(data.Labels) != 2 || data.Series[0].Values[1] != 4 {
		t.Fatalf("unexpected line data: %+v", data)
	}
}

func TestDashboardServiceHandleMeusProdutosMaisVendidos(t *testing.T) {
	repo := &stubDashboardRepository{
		topProducts: []domain.DashboardProductSalesItem{
			{ProductName: "Camisa", Quantity: 7},
			{ProductName: "Calca", Quantity: 2},
		},
	}
	service := newDashboardService(repo, &stubUserRepository{}).(*dashboardService)
	period, _ := service.parsePeriod(request.DashboardWidgetPeriodRequest{From: "2026-01-01", To: "2026-01-31"})
	resellerId := int64(3)

	resp, err := service.handleMeusProdutosMaisVendidos(context.Background(), widgetInput{Period: period, ResellerId: &resellerId})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, ok := resp.Data.(output.DashboardLineBarData)
	if !ok {
		t.Fatalf("expected bar data")
	}
	if len(data.Labels) != 2 || data.Series[0].Values[0] != 7 {
		t.Fatalf("unexpected bar data: %+v", data)
	}
}
