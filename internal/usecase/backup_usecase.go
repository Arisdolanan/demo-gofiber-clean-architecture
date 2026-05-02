package usecase

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/entity"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/repository/postgresql"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/sirupsen/logrus"
)

type BackupUseCase interface {
	ListBackups(ctx context.Context, schoolID int64) ([]*entity.BackupRecord, error)
	CreateBackup(ctx context.Context, schoolID int64) (*entity.BackupRecord, error)
	RestoreBackup(ctx context.Context, schoolID int64, filename string) error
}

type backupUseCase struct {
	repo      postgresql.BackupRepository
	log       *logrus.Logger
	backupDir string
}

func NewBackupUseCase(repo postgresql.BackupRepository, log *logrus.Logger) BackupUseCase {
	backupDir, _ := filepath.Abs("./storage/backup")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		_ = os.MkdirAll(backupDir, 0755)
	}

	return &backupUseCase{
		repo:      repo,
		log:       log,
		backupDir: backupDir,
	}
}

func (uc *backupUseCase) ListBackups(ctx context.Context, schoolID int64) ([]*entity.BackupRecord, error) {
	return uc.repo.FindAll(ctx, schoolID)
}

func (uc *backupUseCase) CreateBackup(ctx context.Context, schoolID int64) (*entity.BackupRecord, error) {
	dbConfig := configuration.GetPostgresConfig()
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("backup_%d_%s.sql", schoolID, timestamp)
	storagePath := filepath.Join(uc.backupDir, filename)

	// Use docker exec to run pg_dump inside the PostgreSQL container
	// Output is piped via stdout to a local file
	uc.log.Infof("Creating database backup using docker exec for container: postgres_container")
	uc.log.Infof("Target file: %s", storagePath)

	cmd := exec.Command("docker", "exec",
		"-e", fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password),
		"postgres_container",
		"pg_dump",
		"-h", "localhost",
		"-U", dbConfig.Username,
		"-d", dbConfig.DBName,
	)

	// Capture stdout output and write directly to the file
	outFile, err := os.Create(storagePath)
	if err != nil {
		uc.log.Errorf("Failed to create output file: %v", err)
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		uc.log.Errorf("Error running docker exec pg_dump: %v, stderr: %s", err, stderr.String())
		_ = os.Remove(storagePath) // cleanup empty file
		return nil, fmt.Errorf("backup failed: %s", stderr.String())
	}

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		uc.log.Errorf("Backup command finished but file not found at: %s", storagePath)
		return nil, fmt.Errorf("backup file was not created at %s", storagePath)
	}

	info, _ := os.Stat(storagePath)

	backupRecord := &entity.BackupRecord{
		Filename:    filename,
		StoragePath: storagePath,
		SizeBytes:   info.Size(),
		Status:      "success",
		SchoolID:    schoolID,
	}

	if err := uc.repo.Create(ctx, schoolID, backupRecord); err != nil {
		uc.log.Errorf("Error saving backup metadata: %v", err)
		return nil, err
	}

	return backupRecord, nil
}

func (uc *backupUseCase) RestoreBackup(ctx context.Context, schoolID int64, filename string) error {
	dbConfig := configuration.GetPostgresConfig()
	
	record, err := uc.repo.FindByFilename(ctx, schoolID, filename)
	if err != nil {
		return fmt.Errorf("backup record not found: %w", err)
	}

	if _, err := os.Stat(record.StoragePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found on disk: %s", record.StoragePath)
	}

	// Open the backup file and pipe it into docker exec psql via stdin
	inFile, err := os.Open(record.StoragePath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer inFile.Close()

	cmd := exec.Command("docker", "exec", "-i",
		"postgres_container",
		"psql",
		"-U", dbConfig.Username,
		"-d", dbConfig.DBName,
	)

	cmd.Stdin = inFile
	var stderr strings.Builder
	cmd.Stderr = &stderr

	uc.log.Infof("Restoring database from %s via docker exec psql", record.StoragePath)

	if err := cmd.Run(); err != nil {
		uc.log.Errorf("Error restoring backup: %v, stderr: %s", err, stderr.String())
		return fmt.Errorf("failed to restore backup: %s", stderr.String())
	}

	uc.log.Infof("Database restored successfully from %s", filename)
	return nil
}
