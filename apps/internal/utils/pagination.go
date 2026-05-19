package utils

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

// PaginationParams represents the pagination request parameters
type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginationMeta contains metadata about the paginated result
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    any            `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// GetPaginationParams extracts and validates pagination parameters from the request query string.
// Defaults: page=1, limit=10. Limit is capped at MaxLimit (100).
func GetPaginationParams(c *fiber.Ctx) PaginationParams {
	page := c.QueryInt("page", DefaultPage)
	limit := c.QueryInt("limit", DefaultLimit)

	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// Offset computes the database offset for the given pagination params.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// BuildPaginationMeta creates pagination metadata from the given params and total count.
func BuildPaginationMeta(params PaginationParams, total int) PaginationMeta {
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}

	return PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// PaginateSuccess returns a paginated JSON success response.
// Usage:
//
//	params := GetPaginationParams(c)
//	total, err := query.Count(ctx)
//	items, err := query.Offset(params.Offset()).Limit(params.Limit).All(ctx)
//	return PaginateSuccess(c, "users retrieved", items, params, total)
func PaginateSuccess(c *fiber.Ctx, message string, data any, params PaginationParams, total int) error {
	resp := PaginatedResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    data,
		Meta:    BuildPaginationMeta(params, total),
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
