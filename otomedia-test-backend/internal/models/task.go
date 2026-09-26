package models

import "time"

// TaskStatus merepresentasikan status task yang valid.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

// Valid mengecek apakah status termasuk salah satu nilai yang diperbolehkan.
func (s TaskStatus) Valid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	}
	return false
}

// Task merepresentasikan satu baris di tabel `tasks`.
// DeletedAt tidak diekspos ke JSON karena task yang soft-deleted tidak pernah
// dikembalikan ke client (Task 4 bug fix: hide soft-deleted tasks).
type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Assignee    string     `json:"assignee"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}
