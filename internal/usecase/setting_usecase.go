package usecase

import (
	"context"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type SettingUseCase interface {
	GetAllSettings(ctx context.Context, schoolID int64) ([]*entity.Setting, error)
	GetSettingsByGroup(ctx context.Context, schoolID int64, groupName string) ([]*entity.Setting, error)
	GetSettingByKey(ctx context.Context, schoolID int64, key string) (*entity.Setting, error)
	UpdateSettings(ctx context.Context, schoolID int64, req entity.SettingUpdateRequest, updatedBy int64) error
}

type settingUseCase struct {
	settingRepo postgresql.SettingRepository
	log         *logrus.Logger
	validate    *validator.Validate
}

func NewSettingUseCase(
	settingRepo postgresql.SettingRepository,
	log *logrus.Logger,
	validate *validator.Validate,
) SettingUseCase {
	return &settingUseCase{
		settingRepo: settingRepo,
		log:         log,
		validate:    validate,
	}
}

func (uc *settingUseCase) GetAllSettings(ctx context.Context, schoolID int64) ([]*entity.Setting, error) {
	settings, err := uc.settingRepo.FindAll(ctx, schoolID)
	if err != nil {
		uc.log.Errorf("Error getting all settings for school %d: %v", schoolID, err)
		return nil, err
	}
	return settings, nil
}

func (uc *settingUseCase) GetSettingsByGroup(ctx context.Context, schoolID int64, groupName string) ([]*entity.Setting, error) {
	settings, err := uc.settingRepo.FindByGroup(ctx, schoolID, groupName)
	if err != nil {
		uc.log.Errorf("Error getting settings by group %s for school %d: %v", groupName, schoolID, err)
		return nil, err
	}
	return settings, nil
}

func (uc *settingUseCase) GetSettingByKey(ctx context.Context, schoolID int64, key string) (*entity.Setting, error) {
	setting, err := uc.settingRepo.FindByKey(ctx, schoolID, key)
	if err != nil {
		uc.log.Errorf("Error getting setting by key %s for school %d: %v", key, schoolID, err)
		return nil, err
	}
	return setting, nil
}

func (uc *settingUseCase) UpdateSettings(ctx context.Context, schoolID int64, req entity.SettingUpdateRequest, updatedBy int64) error {
	if err := uc.validate.Struct(req); err != nil {
		return err
	}

	if err := uc.settingRepo.Upsert(ctx, schoolID, req.Settings, updatedBy); err != nil {
		uc.log.Errorf("Error updating settings for school %d: %v", schoolID, err)
		return err
	}

	uc.log.Infof("Settings updated successfully for school %d by user %d", schoolID, updatedBy)
	return nil
}
