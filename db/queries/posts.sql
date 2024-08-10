-- name: ListPosts :many
select
  p.id,
  p.user_id,
  u.name,
  p.title,
  p.created
from posts p
join users u on u.id = p.user_id
order by created desc
limit 100
;

-- name: CreatePosts :one
insert into posts (user_id, title, content) values (
  (select id from users where sub = @Sub),
  @Title,
  @Content
)
returning id, user_id, title, content, created
;
