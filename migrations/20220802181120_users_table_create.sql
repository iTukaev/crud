-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users (
  name          varchar(30) NOT NULL CONSTRAINT name_right CHECK ( name ~ '^[A-Za-z0-9_\.]+$' ) PRIMARY KEY,
  password      varchar(30) NOT NULL,
  email         varchar(50) NOT NULL UNIQUE CONSTRAINT email_right CHECK(email ~ '^.*@[A-Za-z0-9\-_\.]*$'),
  full_name     varchar(255) NOT NULL,
  created_at    integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
