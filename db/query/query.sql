-- name: GetJob :one
select *
from jobs
where id = $1
limit 1
;

-- name: ListJobs :many
select *
from jobs
order by status
;

-- name: InsertJob :one
insert into jobs (
  status
) values ( $1 )
returning *;
