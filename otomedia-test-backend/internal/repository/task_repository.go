package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"otomedia/task-managment/internal/dto"
	"otomedia/task-managment/internal/models"
)

// TaskRepository adalah abstraksi akses data supaya TaskService bisa diuji
// dengan mock, tanpa perlu koneksi MySQL sungguhan.
type TaskRepository interface {
	List(ctx context.Context, filter dto.ListFilter) ([]models.Task, int64, error)
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	// FindActiveByTitle mencari task aktif (belum di-soft-delete) dengan judul
	// yang sama, dikecualikan task dengan ID excludeID. Dipakai untuk validasi
	// duplicate title saat create/update (Task 4: harus 409, bukan 500).
	FindActiveByTitle(ctx context.Context, title string, excludeID int64) (*models.Task, error)
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	SoftDelete(ctx context.Context, id int64) error
}

type MySQLTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(db *sql.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

func (r *MySQLTaskRepository) List(ctx context.Context, f dto.ListFilter) ([]models.Task, int64, error) {
	// deleted_at IS NULL selalu ada di WHERE -> soft-deleted task tidak pernah
	// ikut terhitung/terambil (Task 4 bug fix).
	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}

	if f.Status != "" {
		where = append(where, "status = ?")
		args = append(args, f.Status)
	}
	if f.Assignee != "" {
		where = append(where, "assignee = ?")
		args = append(args, f.Assignee)
	}
	if f.Keyword != "" {
		where = append(where, "(title LIKE ? OR description LIKE ?)")
		kw := "%" + f.Keyword + "%"
		args = append(args, kw, kw)
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if f.Sort == "oldest" {
		order = "created_at ASC"
	}

	limit := f.Limit
	page := f.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(
		`SELECT id, title, description, status, assignee, created_at, updated_at
		 FROM tasks WHERE %s ORDER BY %s LIMIT ? OFFSET ?`,
		whereClause, order,
	)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Assignee, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *MySQLTaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `SELECT id, title, description, status, assignee, created_at, updated_at
	          FROM tasks WHERE id = ? AND deleted_at IS NULL`
	var t models.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Assignee, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *MySQLTaskRepository) FindActiveByTitle(ctx context.Context, title string, excludeID int64) (*models.Task, error) {
	query := `SELECT id FROM tasks WHERE title = ? AND deleted_at IS NULL AND id != ? LIMIT 1`
	var id int64
	err := r.db.QueryRowContext(ctx, query, title, excludeID).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &models.Task{ID: id}, nil
}

func (r *MySQLTaskRepository) Create(ctx context.Context, t *models.Task) error {
	query := `INSERT INTO tasks (title, description, status, assignee, created_at, updated_at)
	          VALUES (?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Status, t.Assignee)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (r *MySQLTaskRepository) Update(ctx context.Context, t *models.Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, status = ?, assignee = ?, updated_at = NOW()
	          WHERE id = ? AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Status, t.Assignee, t.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MySQLTaskRepository) SoftDelete(ctx context.Context, id int64) error {
	query := `UPDATE tasks SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
