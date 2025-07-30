Create table if not exists tasks (
    id uuid primary key,
    name text not null,
    description text,
    status text not null default 'pending',
    priority integer not null default 1,
    created_by uuid not null,
    is_deleted boolean not null default false,
    version integer not null default 1,
    task_data jsonb not null default '{}'::jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
Create index index_tasks_assigned_to on tasks (created_by);

Create table if not exists parent_child_tasks (
    parent_task_id uuid not null references tasks(id) on delete cascade,
    child_task_id uuid not null references tasks(id) on delete cascade,
    primary key (parent_task_id, child_task_id)
);