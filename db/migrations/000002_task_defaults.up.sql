alter table tasks
alter column due_date
set default (now() at time zone 'GMT+1'),
alter column completed
set default false;