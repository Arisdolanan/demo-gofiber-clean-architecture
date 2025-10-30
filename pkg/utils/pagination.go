package utils

import (
	"math"
	"strconv"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// PaginationParams holds common pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// PaginationResult holds pagination calculation results
type PaginationResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Offset     int   `json:"offset"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

// DefaultPageSize is the default number of items per page
const DefaultPageSize = 10

// MaxPageSize is the maximum allowed page size
const MaxPageSize = 100

// NormalizePagination validates and normalizes pagination parameters
func NormalizePagination(page, pageSize int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize int, totalCount int64) PaginationResult {
	params := NormalizePagination(page, pageSize)
	offset := (params.Page - 1) * params.PageSize
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))

	return PaginationResult{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Offset:     offset,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

// CreatePageMetadata creates a PageMetadata for response
func CreatePageMetadata(page, pageSize int, totalCount int64) response.PageMetadata {
	pagination := CalculatePagination(page, pageSize, totalCount)
	return response.PageMetadata{
		Page:      pagination.Page,
		Size:      pagination.PageSize,
		TotalItem: totalCount,
		TotalPage: int64(pagination.TotalPages),
	}
}

// CreatePageResponse creates a paginated response
func CreatePageResponse[T any](data []T, page, pageSize int, totalCount int64) response.PageResponse[T] {
	return response.PageResponse[T]{
		Data:         data,
		PageMetadata: CreatePageMetadata(page, pageSize, totalCount),
	}
}

// GetOffset calculates the offset for database queries
func GetOffset(page, pageSize int) int {
	params := NormalizePagination(page, pageSize)
	return (params.Page - 1) * params.PageSize
}

// GetLimit returns the normalized page size
func GetLimit(pageSize int) int {
	params := NormalizePagination(1, pageSize)
	return params.PageSize
}

// QueryPaginationParams represents query parameter parsing results
type QueryPaginationParams struct {
	Page     int
	PageSize int
	Limit    int
	Offset   int
}

// ParsePaginationQuery parses pagination parameters from Fiber context
// Supports both page/pageSize and limit/offset approaches
func ParsePaginationQuery(ctx *fiber.Ctx) QueryPaginationParams {
	pageStr := ctx.Query("page", "1")
	pageSizeStr := ctx.Query("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return QueryPaginationParams{
		Page:     page,
		PageSize: pageSize,
		Limit:    limit,
		Offset:   offset,
	}
}

// ParseInt64FromParam safely parses int64 from route parameter
func ParseInt64FromParam(ctx *fiber.Ctx, paramName string) (int64, error) {
	paramStr := ctx.Params(paramName)
	return strconv.ParseInt(paramStr, 10, 64)
}

// ParseIntFromQuery safely parses int from query parameter with default
func ParseIntFromQuery(ctx *fiber.Ctx, paramName string, defaultValue int) int {
	paramStr := ctx.Query(paramName)
	if paramStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil {
		return defaultValue
	}

	return value
}
