-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task_assignments(
    task_id uuid NOT NULL,
    user_id uuid NOT NULL,
    PRIMARY KEY (task_id, user_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_assignments;
-- +goose StatementEnd
