package controller

import (
	_http "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bncunha/erp-api/src/api/http"
	request "github.com/bncunha/erp-api/src/api/requests"
	"github.com/bncunha/erp-api/src/api/viewmodel"
	helper "github.com/bncunha/erp-api/src/application/helpers"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/domain"
	"github.com/labstack/echo/v4"
)

type QuoteController struct {
	quoteService service.QuoteService
}

func NewQuoteController(quoteService service.QuoteService) *QuoteController {
	return &QuoteController{quoteService: quoteService}
}

func (c *QuoteController) Create(context echo.Context) error {
	var req request.CreateQuoteRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	quote, err := c.quoteService.Create(context.Request().Context(), req)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusCreated, viewmodel.ToQuoteDetailedViewModel(quote))
}

func (c *QuoteController) Update(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))
	var req request.UpdateQuoteRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	quote, err := c.quoteService.Update(context.Request().Context(), id, req)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusOK, viewmodel.ToQuoteDetailedViewModel(quote))
}

func (c *QuoteController) GetByID(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))
	quote, err := c.quoteService.GetByID(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusOK, viewmodel.ToQuoteDetailedViewModel(quote))
}

func (c *QuoteController) List(context echo.Context) error {
	input := domain.GetQuotesInput{
		Page:      parseIntOrDefault(context.QueryParam("page"), 1),
		PageSize:  parseIntOrDefault(context.QueryParam("page_size"), 10),
		SortBy:    context.QueryParam("sort_by"),
		SortOrder: context.QueryParam("sort_order"),
	}
	if statuses := parseStatusFilters(context); len(statuses) > 0 {
		input.Statuses = statuses
	}
	if rawCustomer := strings.TrimSpace(context.QueryParam("customer_id")); rawCustomer != "" {
		customerID := helper.ParseInt64(rawCustomer)
		input.CustomerID = &customerID
	}
	if raw := strings.TrimSpace(context.QueryParam("valid_until_start")); raw != "" {
		if value, err := time.Parse(time.DateOnly, raw); err == nil {
			input.ValidUntilStart = &value
		}
	}
	if raw := strings.TrimSpace(context.QueryParam("valid_until_end")); raw != "" {
		if value, err := time.Parse(time.DateOnly, raw); err == nil {
			input.ValidUntilEnd = &value
		}
	}
	if raw := strings.TrimSpace(context.QueryParam("created_at_start")); raw != "" {
		if value, err := time.Parse(time.DateOnly, raw); err == nil {
			input.CreatedAtStart = &value
		}
	}
	if raw := strings.TrimSpace(context.QueryParam("created_at_end")); raw != "" {
		if value, err := time.Parse(time.DateOnly, raw); err == nil {
			value = value.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			input.CreatedAtEnd = &value
		}
	}
	if raw := strings.TrimSpace(context.QueryParam("search")); raw != "" {
		input.Search = &raw
	}
	if strings.EqualFold(context.QueryParam("only_expired"), "true") {
		input.OnlyExpired = true
	}

	output, err := c.quoteService.List(context.Request().Context(), input)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusOK, viewmodel.ToQuoteListViewModel(output))
}

func (c *QuoteController) PatchStatus(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))
	var req request.PatchQuoteStatusRequest
	if err := context.Bind(&req); err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	quote, err := c.quoteService.PatchStatus(context.Request().Context(), id, req)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusOK, viewmodel.ToQuoteDetailedViewModel(quote))
}

func (c *QuoteController) Duplicate(context echo.Context) error {
	id := helper.ParseInt64(context.Param("id"))
	newID, err := c.quoteService.Duplicate(context.Request().Context(), id)
	if err != nil {
		return context.JSON(_http.StatusBadRequest, http.HandleError(err))
	}
	return context.JSON(_http.StatusCreated, viewmodel.ToQuoteDuplicateViewModel(newID))
}

func parseIntOrDefault(raw string, defaultValue int) int {
	if strings.TrimSpace(raw) == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func parseStatusFilters(context echo.Context) []domain.QuoteStatus {
	rawValues := append(context.QueryParams()["status"], context.QueryParams()["status[]"]...)
	if len(rawValues) == 0 {
		return nil
	}
	statuses := make([]domain.QuoteStatus, 0, len(rawValues))
	seen := map[domain.QuoteStatus]struct{}{}
	for _, raw := range rawValues {
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			value := domain.QuoteStatus(strings.TrimSpace(part))
			if value == "" {
				continue
			}
			if _, exists := seen[value]; exists {
				continue
			}
			seen[value] = struct{}{}
			statuses = append(statuses, value)
		}
	}
	return statuses
}
