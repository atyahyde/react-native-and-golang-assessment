DELETE FROM tasks WHERE title IN (
    'Setup project repository',
    'Design database schema',
    'Implement task filtering',
    'Add Redis caching',
    'Write unit tests',
    'Fix duplicate title bug'
);