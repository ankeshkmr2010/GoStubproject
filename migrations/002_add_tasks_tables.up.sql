Create table if not exists tasks (
    id uuid,
    serial_number BIGSERIAL,
    name text not null,
    description text,
    status text not null default 'pending',
    priority integer not null default 1,
    created_by uuid not null,
    is_deleted boolean not null default false,
    task_data jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    primary key (id)
);

Create index index_tasks_assigned_to on tasks (created_by);


CREATE TABLE IF NOT EXISTS task_relationships (
    parent_id UUID NOT NULL,
    child_id UUID NOT NULL,
    relation_type TEXT DEFAULT 'depends_on', -- optional, if you want to describe the relation
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (parent_id)REFERENCES tasks(id)
        ON DELETE CASCADE,
    FOREIGN KEY (child_id) REFERENCES tasks(id)
        ON DELETE CASCADE
);