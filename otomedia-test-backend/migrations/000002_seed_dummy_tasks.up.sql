INSERT INTO tasks (title, description, status, assignee, created_at, updated_at) VALUES
('Setup project repository', 'Inisialisasi repo & CI', 'done', 'Andi', NOW() - INTERVAL 10 DAY, NOW() - INTERVAL 9 DAY),
('Design database schema', 'Rancang skema tabel tasks', 'done', 'Budi', NOW() - INTERVAL 9 DAY, NOW() - INTERVAL 8 DAY),
('Implement task filtering', 'Filter status/keyword/assignee', 'in_progress', 'Citra', NOW() - INTERVAL 5 DAY, NOW() - INTERVAL 1 DAY),
('Add Redis caching', 'Cache GET /api/tasks 60 detik', 'in_progress', 'Dewi', NOW() - INTERVAL 4 DAY, NOW() - INTERVAL 1 DAY),
('Write unit tests', 'Test update, search, cache invalidation', 'todo', 'Eka', NOW() - INTERVAL 2 DAY, NOW() - INTERVAL 2 DAY),
('Fix duplicate title bug', 'Harus return 409 bukan 500', 'todo', 'Andi', NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY);