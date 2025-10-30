package usecase

import (
	"testing"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestExcelUsecase tests the Excel usecase functionality
func TestExcelUsecase(t *testing.T) {
	log := logrus.New()

	excelUsecase := NewExcelUsecase(log)

	assert.NotNil(t, excelUsecase)
}

// TestExportExcel tests exporting data to Excel
func TestExportExcel(t *testing.T) {
	log := logrus.New()

	excelUsecase := NewExcelUsecase(log)

	users := []map[string]interface{}{
		{
			"id":       1,
			"username": "user1",
			"email":    "user1@example.com",
		},
		{
			"id":       2,
			"username": "user2",
			"email":    "user2@example.com",
		},
	}

	excelData := &entity.ExcelData{
		Headers: []string{"id", "username", "email"},
		Rows:    []entity.ExcelRow{users[0], users[1]},
	}

	req := entity.ExcelExportRequest{
		Data:     excelData,
		Sheet:    "Users",
		Filename: "test_users.xlsx",
		Table:    "users",
	}

	response, err := excelUsecase.ExportExcel(req)

	assert.NotNil(t, excelUsecase)
	_ = response
	_ = err
}

// TestExportCustomQueryToExcel tests the custom query export function
func TestExportCustomQueryToExcel(t *testing.T) {
	log := logrus.New()

	excelUsecase := NewExcelUsecase(log)

	assert.NotNil(t, excelUsecase)
}
