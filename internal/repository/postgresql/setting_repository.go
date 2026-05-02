package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type SettingRepository interface {
	FindAll(ctx context.Context, schoolID int64) ([]*entity.Setting, error)
	FindByGroup(ctx context.Context, schoolID int64, groupName string) ([]*entity.Setting, error)
	FindByKey(ctx context.Context, schoolID int64, key string) (*entity.Setting, error)
	Upsert(ctx context.Context, schoolID int64, settings []entity.SettingItem, updatedBy int64) error
}

type settingRepository struct {
	*BaseRepository[entity.Setting]
	db *sqlx.DB
}

func NewSettingRepository(db *sqlx.DB) SettingRepository {
	return &settingRepository{
		BaseRepository: NewBaseRepository[entity.Setting](db, "settings"),
		db:             db,
	}
}

func (r *settingRepository) FindAll(ctx context.Context, schoolID int64) ([]*entity.Setting, error) {
	return r.BaseRepository.FindAll(ctx, "school_id = $1 AND deleted_at IS NULL ORDER BY group_name, setting_key", schoolID)
}

func (r *settingRepository) FindByGroup(ctx context.Context, schoolID int64, groupName string) ([]*entity.Setting, error) {
	return r.BaseRepository.FindAll(ctx, "school_id = $1 AND group_name = $2 AND deleted_at IS NULL ORDER BY setting_key", schoolID, groupName)
}

func (r *settingRepository) FindByKey(ctx context.Context, schoolID int64, key string) (*entity.Setting, error) {
	return r.BaseRepository.FindOne(ctx, "school_id = $1 AND setting_key = $2 AND deleted_at IS NULL", schoolID, key)
}

func (r *settingRepository) Upsert(ctx context.Context, schoolID int64, settings []entity.SettingItem, updatedBy int64) error {
	if len(settings) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Use an upsert strategy for Postgres
	query := `
		INSERT INTO settings (school_id, setting_key, setting_value, group_name, description, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (school_id, setting_key) 
		DO UPDATE SET 
			setting_value = EXCLUDED.setting_value,
			group_name = EXCLUDED.group_name,
			description = COALESCE(EXCLUDED.description, settings.description),
			updated_at = EXCLUDED.updated_at,
			updated_by = EXCLUDED.updated_by
	`

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range settings {
		_, err := stmt.ExecContext(ctx, schoolID, s.SettingKey, s.SettingValue, s.GroupName, s.Description, now, now, updatedBy, updatedBy)
		if err != nil {
			return fmt.Errorf("failed to upsert setting key %s for school %d: %w", s.SettingKey, schoolID, err)
		}
	}

	return tx.Commit()
}
