package postgresql

import (
	"context"
	"database/sql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
)

type BackupRepository interface {
	FindAll(ctx context.Context, schoolID int64) ([]*entity.BackupRecord, error)
	Create(ctx context.Context, schoolID int64, backup *entity.BackupRecord) error
	FindByFilename(ctx context.Context, schoolID int64, filename string) (*entity.BackupRecord, error)
}

type backupRepository struct {
	db *sql.DB
}

func NewBackupRepository(db *sql.DB) BackupRepository {
	return &backupRepository{db: db}
}

func (r *backupRepository) FindAll(ctx context.Context, schoolID int64) ([]*entity.BackupRecord, error) {
	query := "SELECT id, filename, storage_path, size_bytes, status, created_at, school_id FROM backups WHERE school_id = $1 ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	backups := make([]*entity.BackupRecord, 0)
	for rows.Next() {
		b := &entity.BackupRecord{}
		if err := rows.Scan(&b.ID, &b.Filename, &b.StoragePath, &b.SizeBytes, &b.Status, &b.CreatedAt, &b.SchoolID); err != nil {
			return nil, err
		}
		backups = append(backups, b)
	}
	return backups, nil
}

func (r *backupRepository) Create(ctx context.Context, schoolID int64, b *entity.BackupRecord) error {
	query := "INSERT INTO backups (filename, storage_path, size_bytes, status, school_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at"
	return r.db.QueryRowContext(ctx, query, b.Filename, b.StoragePath, b.SizeBytes, b.Status, schoolID).Scan(&b.ID, &b.CreatedAt)
}

func (r *backupRepository) FindByFilename(ctx context.Context, schoolID int64, filename string) (*entity.BackupRecord, error) {
	query := "SELECT id, filename, storage_path, size_bytes, status, created_at, school_id FROM backups WHERE filename = $1 AND school_id = $2"
	b := &entity.BackupRecord{}
	err := r.db.QueryRowContext(ctx, query, filename, schoolID).Scan(&b.ID, &b.Filename, &b.StoragePath, &b.SizeBytes, &b.Status, &b.CreatedAt, &b.SchoolID)
	if err != nil {
		return nil, err
	}
	return b, nil
}
