package output

import "github.com/bncunha/erp-api/src/domain"

type DashboardWidgetItem struct {
	Enum           domain.DashboardWidgetEnum `json:"enum"`
	Type           domain.DashboardWidgetType `json:"type"`
	Order          int                        `json:"order"`
	Title          string                     `json:"title,omitempty"`
	Description    string                     `json:"description,omitempty"`
	DefaultPeriod  *DashboardWidgetPeriod     `json:"default_period,omitempty"`
	DefaultFilters *DashboardWidgetFilters    `json:"default_filters,omitempty"`
}

type DashboardWidgetFilters struct {
	ResellerId *int64 `json:"reseller_id,omitempty"`
	ProductId  *int64 `json:"product_id,omitempty"`
}

type DashboardWidgetPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DashboardWidgetMeta struct {
	Title      string                `json:"title"`
	Period     DashboardWidgetPeriod `json:"period"`
	ShowPeriod bool                  `json:"show_period"`
	Currency   string                `json:"currency,omitempty"`
}

type DashboardWidgetDataOutput struct {
	Enum domain.DashboardWidgetEnum `json:"enum"`
	Type domain.DashboardWidgetType `json:"type"`
	Meta DashboardWidgetMeta        `json:"meta"`
	Data any                        `json:"data"`
}

type DashboardCardData struct {
	Value         float64 `json:"value"`
	Unit          string  `json:"unit,omitempty"`
	DeltaPercent  float64 `json:"delta_percent,omitempty"`
	PreviousValue float64 `json:"previous_value,omitempty"`
}

type DashboardSeries struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

type DashboardLineBarData struct {
	Labels []string          `json:"labels"`
	Series []DashboardSeries `json:"series"`
}

type DashboardPieData struct {
	Labels []string  `json:"labels"`
	Values []float64 `json:"values"`
}

type DashboardTableColumn struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type DashboardTableData struct {
	Columns []DashboardTableColumn `json:"columns"`
	Rows    []map[string]any       `json:"rows"`
}
