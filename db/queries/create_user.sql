-- name: CreateOrUpdateUser :one
insert into users (sub, name) values (
  @Sub,
  @Name
)
on conflict (sub) do update
set name = @Name
returning id, sub, name
;
