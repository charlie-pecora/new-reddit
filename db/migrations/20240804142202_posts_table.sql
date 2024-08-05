-- +goose Up
-- +goose StatementBegin
create table posts (
  id bigserial primary key,
  user_id bigint references users(id),
  title varchar(1000) not null,
  content varchar,
  created timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table posts;
-- +goose StatementEnd
