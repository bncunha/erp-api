package service

import (
	"context"
	"time"

	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/application/constants"
	"github.com/bncunha/erp-api/src/application/errors"
	"github.com/bncunha/erp-api/src/application/service/output"
	"github.com/bncunha/erp-api/src/domain"
)

var (
	ErrDashboardWidgetNotFound   = errors.New("Widget não encontrado")
	ErrDashboardInvalidPeriod    = errors.New("Período inválido")
	ErrDashboardResellerNotFound = errors.New("Revendedor não encontrado")
)

const (
	defaultLowStockThreshold = 0
	maxTopProducts           = 10
)

type DashboardService interface {
	ListWidgets(ctx context.Context) ([]output.DashboardWidgetItem, error)
	GetWidgetData(ctx context.Context, request request.DashboardWidgetDataRequest) (output.DashboardWidgetDataOutput, error)
}

type dashboardService struct {
	dashboardRepository domain.DashboardRepository
	userRepository      domain.UserRepository
}

func NewDashboardService(dashboardRepository domain.DashboardRepository, userRepository domain.UserRepository) DashboardService {
	return &dashboardService{
		dashboardRepository: dashboardRepository,
		userRepository:      userRepository,
	}
}

type widgetPeriod struct {
	From      time.Time
	To        time.Time
	FromDate  time.Time
	ToDate    time.Time
	FromLabel string
	ToLabel   string
}

type widgetInput struct {
	Period     widgetPeriod
	ResellerId *int64
	ProductId  *int64
}

type widgetDefinition struct {
	Enum        domain.DashboardWidgetEnum
	Type        domain.DashboardWidgetType
	Order       int
	Title       string
	Description string
	Roles       []domain.Role
	Handler     func(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error)
}

func (s *dashboardService) ListWidgets(ctx context.Context) ([]output.DashboardWidgetItem, error) {
	role, ok := ctx.Value(constants.ROLE_KEY).(string)
	if !ok || role == "" {
		return nil, ErrPermissionDenied
	}

	items := make([]output.DashboardWidgetItem, 0)
	for _, def := range s.widgetDefinitions() {
		if !s.roleAllowed(def.Roles, role) {
			continue
		}
		items = append(items, output.DashboardWidgetItem{
			Enum:        def.Enum,
			Type:        def.Type,
			Order:       def.Order,
			Title:       def.Title,
			Description: def.Description,
		})
	}

	return items, nil
}

func (s *dashboardService) GetWidgetData(ctx context.Context, request request.DashboardWidgetDataRequest) (output.DashboardWidgetDataOutput, error) {
	if err := request.Validate(); err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	role, ok := ctx.Value(constants.ROLE_KEY).(string)
	if !ok || role == "" {
		return output.DashboardWidgetDataOutput{}, ErrPermissionDenied
	}

	enum := domain.DashboardWidgetEnum(request.Enum)
	definition, ok := s.findWidget(enum)
	if !ok {
		return output.DashboardWidgetDataOutput{}, ErrDashboardWidgetNotFound
	}
	if !s.roleAllowed(definition.Roles, role) {
		return output.DashboardWidgetDataOutput{}, ErrPermissionDenied
	}

	period, err := s.parsePeriod(request.Period)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	resellerId, err := s.resolveResellerFilter(ctx, role, request.Filters)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	var productId *int64
	if request.Filters != nil {
		productId = request.Filters.ProductId
	}

	return definition.Handler(ctx, widgetInput{
		Period:     period,
		ResellerId: resellerId,
		ProductId:  productId,
	})
}

func (s *dashboardService) widgetDefinitions() []widgetDefinition {
	return []widgetDefinition{
		{
			Enum:        domain.DashboardWidgetFaturamento,
			Type:        domain.DashboardWidgetTypeCard,
			Order:       3,
			Title:       "Faturamento",
			Description: "Total vendido no período",
			Roles:       []domain.Role{domain.UserRoleAdmin},
			Handler:     s.handleFaturamento,
		},
		{
			Enum:        domain.DashboardWidgetTotalVendas,
			Type:        domain.DashboardWidgetTypeCard,
			Order:       2,
			Title:       "Total de vendas",
			Description: "Quantidade de vendas no período",
			Roles:       []domain.Role{domain.UserRoleAdmin},
			Handler:     s.handleTotalVendas,
		},
		{
			Enum:        domain.DashboardWidgetProdutosEmEstoque,
			Type:        domain.DashboardWidgetTypeCard,
			Order:       1,
			Title:       "Produtos em estoque",
			Description: "Quantidade total em estoque",
			Roles:       []domain.Role{domain.UserRoleAdmin},
			Handler:     s.handleProdutosEmEstoque,
		},
		{
			Enum:        domain.DashboardWidgetEstoqueBaixo,
			Type:        domain.DashboardWidgetTypeTable,
			Order:       4,
			Title:       "Estoque baixo",
			Description: "Produtos com estoque abaixo do mínimo",
			Roles:       []domain.Role{domain.UserRoleAdmin},
			Handler:     s.handleEstoqueBaixo,
		},
		{
			Enum:        domain.DashboardWidgetFaturamentoNoTempo,
			Type:        domain.DashboardWidgetTypeLine,
			Order:       3,
			Title:       "Faturamento no tempo",
			Description: "Evolução diária do faturamento",
			Roles:       []domain.Role{domain.UserRoleAdmin},
			Handler:     s.handleFaturamentoNoTempo,
		},
		{
			Enum:        domain.DashboardWidgetVendasPorRevendedor,
			Type:        domain.DashboardWidgetTypeBar,
			Order:       6,
			Title:       "Vendas por revendedor",
			Description: "Ranking por faturamento no período",
			Roles:       []domain.Role{domain.UserRoleAdmin},
			Handler:     s.handleVendasPorRevendedor,
		},
		{
			Enum:        domain.DashboardWidgetMeuFaturamento,
			Type:        domain.DashboardWidgetTypeCard,
			Order:       1,
			Title:       "Meu faturamento",
			Description: "Total vendido no período",
			Roles:       []domain.Role{domain.UserRoleReseller},
			Handler:     s.handleMeuFaturamento,
		},
		{
			Enum:        domain.DashboardWidgetMinhasVendas,
			Type:        domain.DashboardWidgetTypeCard,
			Order:       2,
			Title:       "Minhas vendas",
			Description: "Quantidade de vendas no período",
			Roles:       []domain.Role{domain.UserRoleReseller},
			Handler:     s.handleMinhasVendas,
		},
		{
			Enum:        domain.DashboardWidgetMinhasVendasNoTempo,
			Type:        domain.DashboardWidgetTypeLine,
			Order:       3,
			Title:       "Minhas vendas no tempo",
			Description: "Evolução diária das vendas",
			Roles:       []domain.Role{domain.UserRoleReseller},
			Handler:     s.handleMinhasVendasNoTempo,
		},
		{
			Enum:        domain.DashboardWidgetMeusProdutosMaisVendidos,
			Type:        domain.DashboardWidgetTypeBar,
			Order:       4,
			Title:       "Meus produtos mais vendidos",
			Description: "Top produtos por quantidade vendida",
			Roles:       []domain.Role{domain.UserRoleReseller},
			Handler:     s.handleMeusProdutosMaisVendidos,
		},
	}
}

func (s *dashboardService) findWidget(enum domain.DashboardWidgetEnum) (widgetDefinition, bool) {
	for _, def := range s.widgetDefinitions() {
		if def.Enum == enum {
			return def, true
		}
	}
	return widgetDefinition{}, false
}

func (s *dashboardService) roleAllowed(roles []domain.Role, role string) bool {
	for _, allowed := range roles {
		if string(allowed) == role {
			return true
		}
	}
	return false
}

func (s *dashboardService) parsePeriod(period request.DashboardWidgetPeriodRequest) (widgetPeriod, error) {
	fromDate, err := time.Parse(time.DateOnly, period.From)
	if err != nil {
		return widgetPeriod{}, ErrDashboardInvalidPeriod
	}
	toDate, err := time.Parse(time.DateOnly, period.To)
	if err != nil {
		return widgetPeriod{}, ErrDashboardInvalidPeriod
	}
	if fromDate.After(toDate) {
		return widgetPeriod{}, ErrDashboardInvalidPeriod
	}

	from := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 0, 0, 0, 0, time.UTC)
	to = to.Add(24*time.Hour - time.Nanosecond)

	return widgetPeriod{
		From:      from,
		To:        to,
		FromDate:  fromDate,
		ToDate:    toDate,
		FromLabel: period.From,
		ToLabel:   period.To,
	}, nil
}

func (s *dashboardService) resolveResellerFilter(ctx context.Context, role string, filters *request.DashboardWidgetFiltersRequest) (*int64, error) {
	if role == string(domain.UserRoleAdmin) {
		if filters == nil || filters.ResellerId == nil {
			return nil, nil
		}
		resellerId := *filters.ResellerId
		user, err := s.userRepository.GetById(ctx, resellerId)
		if err != nil || user.Role != string(domain.UserRoleReseller) {
			return nil, ErrDashboardResellerNotFound
		}
		return &resellerId, nil
	}

	userIdValue, ok := ctx.Value(constants.USERID_KEY).(float64)
	if !ok {
		return nil, ErrPermissionDenied
	}
	userId := int64(userIdValue)
	return &userId, nil
}

func (s *dashboardService) handleFaturamento(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	current, err := s.dashboardRepository.GetRevenue(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	previous, err := s.getPreviousRevenue(ctx, input)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	return s.buildCardResponse(domain.DashboardWidgetFaturamento, domain.DashboardWidgetTypeCard, "Faturamento do mês", input.Period, current, previous, "BRL"), nil
}

func (s *dashboardService) handleMeuFaturamento(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	current, err := s.dashboardRepository.GetRevenue(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	previous, err := s.getPreviousRevenue(ctx, input)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	return s.buildCardResponse(domain.DashboardWidgetMeuFaturamento, domain.DashboardWidgetTypeCard, "Meu faturamento", input.Period, current, previous, "BRL"), nil
}

func (s *dashboardService) handleTotalVendas(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	current, err := s.dashboardRepository.GetSalesCount(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	previous, err := s.getPreviousSalesCount(ctx, input)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	return s.buildCardResponse(domain.DashboardWidgetTotalVendas, domain.DashboardWidgetTypeCard, "Total de vendas", input.Period, float64(current), float64(previous), ""), nil
}

func (s *dashboardService) handleMinhasVendas(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	current, err := s.dashboardRepository.GetSalesCount(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	previous, err := s.getPreviousSalesCount(ctx, input)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	return s.buildCardResponse(domain.DashboardWidgetMinhasVendas, domain.DashboardWidgetTypeCard, "Minhas vendas", input.Period, float64(current), float64(previous), ""), nil
}

func (s *dashboardService) handleProdutosEmEstoque(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	total, err := s.dashboardRepository.GetStockTotal(ctx, domain.DashboardStockQueryInput{
		ResellerId: input.ResellerId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	data := output.DashboardCardData{
		Value: total,
		Unit:  "UN",
	}

	return output.DashboardWidgetDataOutput{
		Enum: domain.DashboardWidgetProdutosEmEstoque,
		Type: domain.DashboardWidgetTypeCard,
		Meta: s.buildMeta("Produtos em estoque", input.Period, "", false),
		Data: data,
	}, nil
}

func (s *dashboardService) handleEstoqueBaixo(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	items, err := s.dashboardRepository.GetLowStockProducts(ctx, domain.DashboardStockQueryInput{
		ResellerId: input.ResellerId,
		Threshold:  defaultLowStockThreshold,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rows = append(rows, map[string]any{
			"product": item.ProductName,
			"qty":     item.Quantity,
		})
	}

	data := output.DashboardTableData{
		Columns: []output.DashboardTableColumn{
			{Key: "product", Label: "Produto"},
			{Key: "qty", Label: "Qtd"},
		},
		Rows: rows,
	}

	return output.DashboardWidgetDataOutput{
		Enum: domain.DashboardWidgetEstoqueBaixo,
		Type: domain.DashboardWidgetTypeTable,
		Meta: s.buildMeta("Estoque baixo", input.Period, "", true),
		Data: data,
	}, nil
}

func (s *dashboardService) handleFaturamentoNoTempo(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	series, err := s.dashboardRepository.GetRevenueByDay(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	labels := make([]string, 0, len(series))
	values := make([]float64, 0, len(series))
	for _, item := range series {
		labels = append(labels, item.Date.Format(time.DateOnly))
		values = append(values, item.Value)
	}

	data := output.DashboardLineBarData{
		Labels: labels,
		Series: []output.DashboardSeries{
			{Name: "Faturamento", Values: values},
		},
	}

	return output.DashboardWidgetDataOutput{
		Enum: domain.DashboardWidgetFaturamentoNoTempo,
		Type: domain.DashboardWidgetTypeLine,
		Meta: s.buildMeta("Faturamento no tempo", input.Period, "BRL", true),
		Data: data,
	}, nil
}

func (s *dashboardService) handleMinhasVendasNoTempo(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	series, err := s.dashboardRepository.GetSalesCountByDay(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	labels := make([]string, 0, len(series))
	values := make([]float64, 0, len(series))
	for _, item := range series {
		labels = append(labels, item.Date.Format(time.DateOnly))
		values = append(values, item.Value)
	}

	data := output.DashboardLineBarData{
		Labels: labels,
		Series: []output.DashboardSeries{
			{Name: "Vendas", Values: values},
		},
	}

	return output.DashboardWidgetDataOutput{
		Enum: domain.DashboardWidgetMinhasVendasNoTempo,
		Type: domain.DashboardWidgetTypeLine,
		Meta: s.buildMeta("Minhas vendas no tempo", input.Period, "", true),
		Data: data,
	}, nil
}

func (s *dashboardService) handleVendasPorRevendedor(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	items, err := s.dashboardRepository.GetRevenueByReseller(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	labels := make([]string, 0, len(items))
	values := make([]float64, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.ResellerName)
		values = append(values, item.Value)
	}

	data := output.DashboardLineBarData{
		Labels: labels,
		Series: []output.DashboardSeries{
			{Name: "Faturamento", Values: values},
		},
	}

	return output.DashboardWidgetDataOutput{
		Enum: domain.DashboardWidgetVendasPorRevendedor,
		Type: domain.DashboardWidgetTypeBar,
		Meta: s.buildMeta("Vendas por revendedor", input.Period, "BRL", true),
		Data: data,
	}, nil
}

func (s *dashboardService) handleMeusProdutosMaisVendidos(ctx context.Context, input widgetInput) (output.DashboardWidgetDataOutput, error) {
	items, err := s.dashboardRepository.GetTopProductsByReseller(ctx, domain.DashboardQueryInput{
		From:       input.Period.From,
		To:         input.Period.To,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	}, maxTopProducts)
	if err != nil {
		return output.DashboardWidgetDataOutput{}, err
	}

	labels := make([]string, 0, len(items))
	values := make([]float64, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.ProductName)
		values = append(values, item.Quantity)
	}

	data := output.DashboardLineBarData{
		Labels: labels,
		Series: []output.DashboardSeries{
			{Name: "Quantidade", Values: values},
		},
	}

	return output.DashboardWidgetDataOutput{
		Enum: domain.DashboardWidgetMeusProdutosMaisVendidos,
		Type: domain.DashboardWidgetTypeBar,
		Meta: s.buildMeta("Meus produtos mais vendidos", input.Period, "", true),
		Data: data,
	}, nil
}

func (s *dashboardService) getPreviousRevenue(ctx context.Context, input widgetInput) (float64, error) {
	prevFrom, prevTo := s.previousPeriod(input.Period)
	return s.dashboardRepository.GetRevenue(ctx, domain.DashboardQueryInput{
		From:       prevFrom,
		To:         prevTo,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
}

func (s *dashboardService) getPreviousSalesCount(ctx context.Context, input widgetInput) (int64, error) {
	prevFrom, prevTo := s.previousPeriod(input.Period)
	return s.dashboardRepository.GetSalesCount(ctx, domain.DashboardQueryInput{
		From:       prevFrom,
		To:         prevTo,
		ResellerId: input.ResellerId,
		ProductId:  input.ProductId,
	})
}

func (s *dashboardService) previousPeriod(period widgetPeriod) (time.Time, time.Time) {
	days := int(period.ToDate.Sub(period.FromDate).Hours()/24) + 1
	prevFrom := period.FromDate.AddDate(0, 0, -days)
	prevTo := period.FromDate.AddDate(0, 0, -1)
	prevTo = time.Date(prevTo.Year(), prevTo.Month(), prevTo.Day(), 0, 0, 0, 0, time.UTC).Add(24*time.Hour - time.Nanosecond)
	prevFrom = time.Date(prevFrom.Year(), prevFrom.Month(), prevFrom.Day(), 0, 0, 0, 0, time.UTC)
	return prevFrom, prevTo
}

func (s *dashboardService) buildMeta(title string, period widgetPeriod, currency string, showPeriod bool) output.DashboardWidgetMeta {
	return output.DashboardWidgetMeta{
		Title:      title,
		ShowPeriod: showPeriod,
		Period: output.DashboardWidgetPeriod{
			From: period.FromLabel,
			To:   period.ToLabel,
		},
		Currency: currency,
	}
}

func (s *dashboardService) buildCardResponse(enum domain.DashboardWidgetEnum, widgetType domain.DashboardWidgetType, title string, period widgetPeriod, value float64, previous float64, unit string) output.DashboardWidgetDataOutput {
	delta := 0.0
	if previous > 0 {
		delta = ((value - previous) / previous) * 100
	}

	data := output.DashboardCardData{
		Value:         value,
		Unit:          unit,
		DeltaPercent:  delta,
		PreviousValue: previous,
	}

	return output.DashboardWidgetDataOutput{
		Enum: enum,
		Type: widgetType,
		Meta: s.buildMeta(title, period, unit, true),
		Data: data,
	}
}
