package postgresql

import (
	"context"
	"fmt"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *entity.ActivityLog) error
	FindAll(ctx context.Context, schoolID *int64, limit, offset int) ([]*entity.ActivityLog, error)
	FindDeletions(ctx context.Context, schoolID *int64, limit, offset int) ([]*entity.ActivityLog, error)
}

type activityLogRepository struct {
	db *sqlx.DB
}

func NewActivityLogRepository(db *sqlx.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *entity.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (user_id, school_id, action, module, description, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	return r.db.QueryRowxContext(ctx, query,
		log.UserID, log.SchoolID, log.Action, log.Module, log.Description, log.IPAddress, log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *activityLogRepository) FindAll(ctx context.Context, schoolID *int64, limit, offset int) ([]*entity.ActivityLog, error) {
	var logs []*entity.ActivityLog
	query := `
		SELECT l.*, u.username as user_name
		FROM activity_logs l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if schoolID != nil {
		query += fmt.Sprintf(" AND l.school_id = $%d", argCount)
		args = append(args, *schoolID)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY l.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		return nil, err
	}
	if logs == nil {
		logs = []*entity.ActivityLog{}
	}
	return logs, nil
}

func (r *activityLogRepository) FindDeletions(ctx context.Context, schoolID *int64, limit, offset int) ([]*entity.ActivityLog, error) {
	var logs []*entity.ActivityLog
	query := `
		SELECT l.*, u.username as user_name
		FROM activity_logs l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE (l.action ILIKE '%delete%' OR l.action ILIKE '%remove%')
	`
	args := []interface{}{}
	argCount := 1

	if schoolID != nil {
		query += fmt.Sprintf(" AND l.school_id = $%d", argCount)
		args = append(args, *schoolID)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY l.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		return nil, err
	}
	if logs == nil {
		logs = []*entity.ActivityLog{}
	}
	return logs, nil
}
