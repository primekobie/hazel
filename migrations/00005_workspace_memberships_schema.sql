-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workspace_memberships(
    workspace_id uuid NOT NULL,
    user_id uuid NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (workspace_id, user_id),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workspace_memberships;
-- +goose StatementEnd
