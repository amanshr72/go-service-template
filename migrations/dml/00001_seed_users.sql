-- +goose Up
INSERT INTO
    users (name, email, is_active)
VALUES
    ('Seed User', 'seed@test.com', true),
    ('John Doe', 'john.doe@test.com', true),
    ('Jane Smith', 'jane.smith@test.com', false),
    ('Alex Johnson', 'alex.johnson@test.com', true) ON CONFLICT (email) DO NOTHING;

-- +goose Down
DELETE FROM
    users
WHERE
    email IN (
        'seed@test.com',
        'john.doe@test.com',
        'jane.smith@test.com',
        'alex.johnson@test.com'
    );