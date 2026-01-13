package domain

type DashboardWidgetEnum string

const (
	DashboardWidgetFaturamento              DashboardWidgetEnum = "FATURAMENTO"
	DashboardWidgetTotalVendas              DashboardWidgetEnum = "TOTAL_VENDAS"
	DashboardWidgetProdutosEmEstoque        DashboardWidgetEnum = "PRODUTOS_EM_ESTOQUE"
	DashboardWidgetEstoqueBaixo             DashboardWidgetEnum = "ESTOQUE_BAIXO"
	DashboardWidgetFaturamentoNoTempo       DashboardWidgetEnum = "FATURAMENTO_NO_TEMPO"
	DashboardWidgetVendasPorRevendedor      DashboardWidgetEnum = "VENDAS_POR_REVENDEDOR"
	DashboardWidgetProdutosMaisVendidos     DashboardWidgetEnum = "PRODUTOS_MAIS_VENDIDOS"
	DashboardWidgetMeuFaturamento           DashboardWidgetEnum = "MEU_FATURAMENTO"
	DashboardWidgetMinhasVendas             DashboardWidgetEnum = "MINHAS_VENDAS"
	DashboardWidgetMinhasVendasNoTempo      DashboardWidgetEnum = "MINHAS_VENDAS_NO_TEMPO"
	DashboardWidgetMeusProdutosMaisVendidos DashboardWidgetEnum = "MEUS_PRODUTOS_MAIS_VENDIDOS"
)

type DashboardWidgetType string

const (
	DashboardWidgetTypeCard  DashboardWidgetType = "CARD"
	DashboardWidgetTypeBar   DashboardWidgetType = "BAR"
	DashboardWidgetTypeLine  DashboardWidgetType = "LINE"
	DashboardWidgetTypePie   DashboardWidgetType = "PIE"
	DashboardWidgetTypeTable DashboardWidgetType = "TABLE"
)
