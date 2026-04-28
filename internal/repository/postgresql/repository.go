package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/internal/infrastructure/database"
	"github.com/jmoiron/sqlx"
)

// BaseRepository provides common CRUD operations for all entities
type BaseRepository[T any] struct {
	postgres *database.Postgres
	table    string
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T any](db *sqlx.DB, table string) *BaseRepository[T] {
	return &BaseRepository[T]{
		postgres: database.NewPostgres(db),
		table:    table,
	}
}

// Create inserts a new entity into the database
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	values := r.structToMap(entity)

	// Remove ID if it's zero value (new entity)
	if id, exists := values["id"]; exists {
		if reflect.ValueOf(id).IsZero() {
			delete(values, "id")
		}
	}

	// Handle audit fields for entities that have them
	now := time.Now()
	if _, hasCreatedAt := values["created_at"]; hasCreatedAt {
		values["created_at"] = now
	}
	if _, hasUpdatedAt := values["updated_at"]; hasUpdatedAt {
		values["updated_at"] = now
	}

	// Handle audit trail fields if context contains user ID
	if userID := ctx.Value("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			if _, hasCreatedBy := values["created_by"]; hasCreatedBy {
				values["created_by"] = &userIDVal
			}
			if _, hasUpdatedBy := values["updated_by"]; hasUpdatedBy {
				values["updated_by"] = &userIDVal
			}
		}
	}

	id, err := r.postgres.Table(r.table).Insert(ctx, values)
	if err != nil {
		return err
	}

	// Set the ID back to the entity
	r.setEntityID(entity, id)

	return nil
}

// Update updates an existing entity in the database
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T, whereClause string, args ...any) error {
	values := r.structToMap(entity)

	// Remove fields that shouldn't be updated
	delete(values, "id")
	delete(values, "created_at")
	delete(values, "created_by")

	// Remove zero values to allow partial updates and avoid overwriting with defaults
	for k, v := range values {
		if reflect.ValueOf(v).IsZero() {
			delete(values, k)
		}
	}

	// Handle audit fields for entities that have them
	if _, hasUpdatedAt := values["updated_at"]; hasUpdatedAt {
		values["updated_at"] = time.Now()
	}

	// Handle audit trail fields if context contains user ID
	if userID := ctx.Value("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			if _, hasUpdatedBy := values["updated_by"]; hasUpdatedBy {
				values["updated_by"] = &userIDVal
			}
		}
	}

	return r.postgres.Table(r.table).
		Where(whereClause, args...).
		Update(ctx, values)
}

// Delete permanently removes an entity from the database
func (r *BaseRepository[T]) Delete(ctx context.Context, whereClause string, args ...any) error {
	return r.postgres.Table(r.table).
		Where(whereClause, args...).
		Delete(ctx)
}

// SoftDelete marks an entity as deleted without removing it
func (r *BaseRepository[T]) SoftDelete(ctx context.Context, whereClause string, args ...any) error {
	values := map[string]any{
		"deleted_at": time.Now(),
	}

	// Handle audit trail fields if context contains user ID
	if userID := ctx.Value("user_id"); userID != nil {
		if userIDVal, ok := userID.(int64); ok {
			values["deleted_by"] = &userIDVal
		}
	}

	return r.postgres.Table(r.table).
		Where(whereClause, args...).
		Update(ctx, values)
}

// FindByID retrieves an entity by its ID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id int64) (*T, error) {
	var entity T
	err := r.postgres.Table(r.table).
		Where("id = $1 AND deleted_at IS NULL", id).
		First(ctx, &entity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

// FindOne retrieves the first entity matching the where clause
func (r *BaseRepository[T]) FindOne(ctx context.Context, whereClause string, args ...any) (*T, error) {
	var entity T
	err := r.postgres.Table(r.table).
		Where(whereClause, args...).
		First(ctx, &entity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

// FindAll retrieves all entities matching the where clause
func (r *BaseRepository[T]) FindAll(ctx context.Context, whereClause string, args ...any) ([]*T, error) {
	var entities []*T
	query := r.postgres.Table(r.table)

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Find(ctx, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// FindAllWithPagination retrieves entities with pagination
func (r *BaseRepository[T]) FindAllWithPagination(ctx context.Context, limit, offset int, whereClause string, args ...any) ([]*T, error) {
	var entities []*T
	query := r.postgres.Table(r.table)

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(ctx, &entities)

	if err != nil {
		return nil, err
	}

	return entities, nil
}

// Count counts entities matching the where clause
func (r *BaseRepository[T]) Count(ctx context.Context, whereClause string, args ...any) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.table)

	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	err := r.postgres.Get(ctx, &count, query, args...)
	return count, err
}

// Exists checks if an entity exists matching the where clause
func (r *BaseRepository[T]) Exists(ctx context.Context, whereClause string, args ...any) (bool, error) {
	count, err := r.Count(ctx, whereClause, args...)
	return count > 0, err
}

// GetContext executes a query and scans the result into a destination
func (r *BaseRepository[T]) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.postgres.Get(ctx, dest, query, args...)
}

// SelectContext executes a query and scans the results into a slice
func (r *BaseRepository[T]) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.postgres.Select(ctx, dest, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row
func (r *BaseRepository[T]) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return r.postgres.QueryRowx(ctx, query, args...)
}

// ExecContext executes a query without returning any rows
func (r *BaseRepository[T]) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.postgres.Exec(ctx, query, args...)
}

// Helper methods

// structToMap converts a struct to a map for database operations
func (r *BaseRepository[T]) structToMap(entity *T) map[string]any {
	values := make(map[string]any)
	r.mapFields(reflect.ValueOf(entity).Elem(), values)
	return values
}

// mapFields recursively maps struct fields to a map based on db tags
func (r *BaseRepository[T]) mapFields(v reflect.Value, values map[string]any) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("db")

		// If it's an embedded struct and has no db tag, recurse
		if fieldType.Anonymous && field.Kind() == reflect.Struct && tag == "" {
			r.mapFields(field, values)
			continue
		}

		if tag != "" && tag != "-" {
			// Use the tag as provided in the db tag (usually already snake_case)
			values[tag] = field.Interface()
		}
	}
}

// getEntityFieldPtr gets a pointer to a field in the entity by name
func (r *BaseRepository[T]) getEntityFieldPtr(entity *T, fieldName string) (interface{}, bool) {
	v := reflect.ValueOf(entity).Elem()
	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		return nil, false
	}

	// Return a pointer to the field
	return field.Addr().Interface(), true
}

// setEntityID sets the ID field of an entity
func (r *BaseRepository[T]) setEntityID(entity *T, id int64) {
	v := reflect.ValueOf(entity).Elem()
	idField := v.FieldByName("ID")

	if idField.IsValid() && idField.CanSet() {
		idField.SetInt(id)
	}
}

// toSnakeCase converts camelCase to snake_case
func (r *BaseRepository[T]) toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
