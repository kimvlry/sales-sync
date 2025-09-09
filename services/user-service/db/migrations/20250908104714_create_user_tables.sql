-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id uuid primary key,
    telegram_id text unique not null,
    name text not null,
    created_at timestamptz default now()
);

create table if not exists marketplace_accounts (
    id uuid primary key,
    user_id uuid references users(id) on delete cascade,
    marketplace_type text not null,
    account_id uuid,
    credentials jsonb,
    created_at timestamptz default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists marketplace_accounts;
drop table if exists users;
-- +goose StatementEnd
