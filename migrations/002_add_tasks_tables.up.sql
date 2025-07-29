Create table if not exists tasks (
    id uuid primary key,
    name text not null,
    description text,
    status text not null default 'pending',
    priority integer not null default 1,
    due_date timestamptz,
    assigned_to uuid not null,
    created_by uuid not null,
    is_deleted boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
Create index index_tasks_assigned_to on tasks (assigned_to);

Create table if not exists task_comments (
    id uuid primary key,
    task_id uuid not null references tasks(id) on delete cascade,
    comment text not null,
    created_by uuid not null,
    is_deleted boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
Create index task_comments_index on task_comments (task_id);

Create table if not exists parent_child_tasks (
    parent_task_id uuid not null references tasks(id) on delete cascade,
    child_task_id uuid not null references tasks(id) on delete cascade,
    primary key (parent_task_id, child_task_id)
);