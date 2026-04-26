package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db   *sqlx.DB
	tx   *sqlx.Tx
	conn *sqlx.Conn
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

//
// ---------------- Base sqlx verbs ----------------
//

// Exec query (INSERT/UPDATE/DELETE) returning sql.Result
func (p *Postgres) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if p.tx != nil {
		return p.tx.ExecContext(ctx, query, args...)
	}
	return p.db.ExecContext(ctx, query, args...)
}

// MustExec query, panic jika error
func (p *Postgres) MustExec(ctx context.Context, query string, args ...any) sql.Result {
	if p.tx != nil {
		return p.tx.MustExecContext(ctx, query, args...)
	}
	return p.db.MustExecContext(ctx, query, args...)
}

// Queryx return sqlx.Rows (bisa StructScan)
func (p *Postgres) Queryx(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	if p.tx != nil {
		return p.tx.QueryxContext(ctx, query, args...)
	}
	return p.db.QueryxContext(ctx, query, args...)
}

// QueryRowx return sqlx.Row (bisa StructScan)
func (p *Postgres) QueryRowx(ctx context.Context, query string, args ...any) *sqlx.Row {
	if p.tx != nil {
		return p.tx.QueryRowxContext(ctx, query, args...)
	}
	return p.db.QueryRowxContext(ctx, query, args...)
}

// Get one row into struct/dest
func (p *Postgres) Get(ctx context.Context, dest any, query string, args ...any) error {
	if p.tx != nil {
		return p.tx.GetContext(ctx, dest, query, args...)
	}
	return p.db.GetContext(ctx, dest, query, args...)
}

// Select many rows into slice
func (p *Postgres) Select(ctx context.Context, dest any, query string, args ...any) error {
	if p.tx != nil {
		return p.tx.SelectContext(ctx, dest, query, args...)
	}
	return p.db.SelectContext(ctx, dest, query, args...)
}

//
// ---------------- Mini ORM chaining ----------------
//

type Query struct {
	db     *Postgres
	table  string
	where  string
	args   []any
	limit  string
	offset string
	order  string
	joins  string
	values map[string]any
}

func (p *Postgres) Table(name string) *Query {
	return &Query{
		db:     p,
		table:  name,
		values: make(map[string]any),
	}
}

func (q *Query) Where(cond string, args ...any) *Query {
	q.where = cond
	q.args = args
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = fmt.Sprintf(" LIMIT %d", limit)
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = fmt.Sprintf(" OFFSET %d", offset)
	return q
}

func (q *Query) Order(order string) *Query {
	q.order = " ORDER BY " + order
	return q
}

func (q *Query) Join(joinClause string) *Query {
	q.joins += " " + joinClause
	return q
}

// SELECT multiple
func (q *Query) Find(ctx context.Context, dest any) error {
	query := fmt.Sprintf("SELECT * FROM %s", q.table)
	if q.joins != "" {
		query += q.joins
	}
	if q.where != "" {
		query += " WHERE " + q.where
	}
	if q.order != "" {
		query += q.order
	}
	if q.limit != "" {
		query += q.limit
	}
	if q.offset != "" {
		query += q.offset
	}
	return q.db.Select(ctx, dest, query, q.args...)
}

// SELECT first row
func (q *Query) First(ctx context.Context, dest any) error {
	query := fmt.Sprintf("SELECT * FROM %s", q.table)
	if q.joins != "" {
		query += q.joins
	}
	if q.where != "" {
		query += " WHERE " + q.where
	}
	if q.order != "" {
		query += q.order
	}
	query += " LIMIT 1"
	return q.db.Get(ctx, dest, query, q.args...)
}

// INSERT
func (q *Query) Insert(ctx context.Context, values map[string]any) (int64, error) {
	if len(values) == 0 {
		return 0, errors.New("no values to insert")
	}

	cols := []string{}
	placeholders := []string{}
	args := []any{}
	i := 1
	for k, v := range values {
		cols = append(cols, k)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i)) // untuk PostgreSQL
		args = append(args, v)
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		q.table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	var id int64
	err := q.db.Get(ctx, &id, query, args...)
	return id, err
}

// UPDATE
func (q *Query) Update(ctx context.Context, values map[string]any) error {
	if q.where == "" {
		return errors.New("update without where is not allowed")
	}
	if len(values) == 0 {
		return errors.New("no values to update")
	}

	sets := []string{}
	args := []any{}
	i := 1
	for k, v := range values {
		sets = append(sets, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}
	args = append(args, q.args...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		q.table,
		strings.Join(sets, ", "),
		q.where,
	)
	_, err := q.db.Exec(ctx, query, args...)
	return err
}

// DELETE
func (q *Query) Delete(ctx context.Context) error {
	if q.where == "" {
		return errors.New("delete without where is not allowed")
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", q.table, q.where)
	_, err := q.db.Exec(ctx, query, q.args...)
	return err
}
