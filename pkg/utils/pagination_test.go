package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected PaginationParams
	}{
		{
			name:     "Valid parameters",
			page:     2,
			pageSize: 20,
			expected: PaginationParams{Page: 2, PageSize: 20},
		},
		{
			name:     "Page less than 1",
			page:     0,
			pageSize: 20,
			expected: PaginationParams{Page: 1, PageSize: 20},
		},
		{
			name:     "Negative page",
			page:     -5,
			pageSize: 20,
			expected: PaginationParams{Page: 1, PageSize: 20},
		},
		{
			name:     "PageSize less than 1",
			page:     2,
			pageSize: 0,
			expected: PaginationParams{Page: 2, PageSize: DefaultPageSize},
		},
		{
			name:     "PageSize greater than max",
			page:     2,
			pageSize: 150,
			expected: PaginationParams{Page: 2, PageSize: DefaultPageSize},
		},
		{
			name:     "Both invalid",
			page:     -1,
			pageSize: 150,
			expected: PaginationParams{Page: 1, PageSize: DefaultPageSize},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePagination(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		totalCount int64
		expected   PaginationResult
	}{
		{
			name:       "First page with exact division",
			page:       1,
			pageSize:   10,
			totalCount: 100,
			expected: PaginationResult{
				Page:       1,
				PageSize:   10,
				Offset:     0,
				TotalCount: 100,
				TotalPages: 10,
			},
		},
		{
			name:       "Second page",
			page:       2,
			pageSize:   10,
			totalCount: 100,
			expected: PaginationResult{
				Page:       2,
				PageSize:   10,
				Offset:     10,
				TotalCount: 100,
				TotalPages: 10,
			},
		},
		{
			name:       "Last page with remainder",
			page:       3,
			pageSize:   10,
			totalCount: 25,
			expected: PaginationResult{
				Page:       3,
				PageSize:   10,
				Offset:     20,
				TotalCount: 25,
				TotalPages: 3,
			},
		},
		{
			name:       "Empty result set",
			page:       1,
			pageSize:   10,
			totalCount: 0,
			expected: PaginationResult{
				Page:       1,
				PageSize:   10,
				Offset:     0,
				TotalCount: 0,
				TotalPages: 0,
			},
		},
		{
			name:       "Single item",
			page:       1,
			pageSize:   10,
			totalCount: 1,
			expected: PaginationResult{
				Page:       1,
				PageSize:   10,
				Offset:     0,
				TotalCount: 1,
				TotalPages: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePagination(tt.page, tt.pageSize, tt.totalCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "First page",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "Second page",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "Third page with different page size",
			page:     3,
			pageSize: 20,
			expected: 40,
		},
		{
			name:     "Invalid page (should normalize to 1)",
			page:     0,
			pageSize: 10,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOffset(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLimit(t *testing.T) {
	tests := []struct {
		name     string
		pageSize int
		expected int
	}{
		{
			name:     "Valid page size",
			pageSize: 20,
			expected: 20,
		},
		{
			name:     "Invalid page size (should use default)",
			pageSize: 0,
			expected: DefaultPageSize,
		},
		{
			name:     "Too large page size (should use default)",
			pageSize: 150,
			expected: DefaultPageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLimit(tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}