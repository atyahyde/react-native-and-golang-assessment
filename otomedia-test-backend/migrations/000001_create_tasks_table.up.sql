CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status ENUM('todo', 'in_progress', 'done') NOT NULL DEFAULT 'todo',
    assignee VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,

    INDEX idx_tasks_status (status),
    INDEX idx_tasks_assignee (assignee),
    INDEX idx_tasks_deleted_at (deleted_at),
    INDEX idx_tasks_created_at (created_at),
    -- index gabungan untuk query list yang selalu memfilter deleted_at IS NULL
    -- lalu ORDER BY created_at
    INDEX idx_tasks_deleted_created (deleted_at, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Catatan: judul task tidak diberi UNIQUE constraint di level database karena
-- MySQL tidak mendukung partial/filtered unique index (unik hanya untuk baris
-- yang deleted_at IS NULL). Validasi duplicate title untuk task aktif
-- dilakukan di application layer (lihat TaskRepository.FindActiveByTitle),
-- supaya title task yang sudah di-soft-delete tetap bisa dipakai ulang.