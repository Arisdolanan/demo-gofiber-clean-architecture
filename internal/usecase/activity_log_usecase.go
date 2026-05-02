package usecase

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/sirupsen/logrus"
)

type ActivityLogUsecase interface {
	GetActivities(ctx context.Context, schoolID *int64, page, pageSize int) ([]*entity.ActivityLog, error)
	GetDeletions(ctx context.Context, schoolID *int64, page, pageSize int) ([]*entity.ActivityLog, error)
	LogActivity(ctx context.Context, log *entity.ActivityLog) error
}

type activityLogUsecase struct {
	repo postgresql.ActivityLogRepository
	log  *logrus.Logger
}

func NewActivityLogUsecase(repo postgresql.ActivityLogRepository, log *logrus.Logger) ActivityLogUsecase {
	return &activityLogUsecase{
		repo: repo,
		log:  log,
	}
}

func (uc *activityLogUsecase) GetActivities(ctx context.Context, schoolID *int64, page, pageSize int) ([]*entity.ActivityLog, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return uc.repo.FindAll(ctx, schoolID, pageSize, offset)
}

func (uc *activityLogUsecase) GetDeletions(ctx context.Context, schoolID *int64, page, pageSize int) ([]*entity.ActivityLog, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return uc.repo.FindDeletions(ctx, schoolID, pageSize, offset)
}

func (uc *activityLogUsecase) LogActivity(ctx context.Context, log *entity.ActivityLog) error {
	err := uc.repo.Create(ctx, log)
	if err != nil {
		uc.log.Errorf("Failed to create activity log in repository: %v", err)
		return err
	}
	return nil
}
