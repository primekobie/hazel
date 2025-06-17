-- +goose Up
-- +goose StatementBegin

DO
$$
    BEGIN
        CREATE TYPE task_status AS ENUM ('todo', 'started', 'complete');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

DO
$$
    BEGIN
        CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

CREATE TABLE IF NOT EXISTS tasks(
    id uuid NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    project_id uuid NOT NULL,
    status task_status DEFAULT 'todo' NOT NULL,
    priority task_priority DEFAULT 'low' NOT NULL,
    due TIMESTAMP,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    last_modified TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (project_id) REFERENCES projects (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS task_priority;
-- +goose StatementEnd
