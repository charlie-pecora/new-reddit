-- +goose Up
-- +goose StatementBegin
create table users (
  id bigserial primary key,
  sub varchar(20) unique not null,
  name varchar(300) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
