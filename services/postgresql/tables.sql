DO $$
BEGIN

CREATE TABLE IF NOT EXISTS user.account(
    id serial primary key,
    name text not null,
    email text not null,
    hash bytea not null,
    salt bytea not null
);

CREATE TABLE IF NOT EXISTS user.meta(
  id serial primary key,
  account_id int references user.account,
  is_active boolean,
  email_confirmed boolean,
  email_confirmation_path text,
  account_created_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS user.session(
  id serial primary key,
  account_id int references user.account,
  created_at timestamp default now(),
  active boolean,
  token bytea
);

CREATE TABLE IF NOT EXISTS user.resource(
    id serial primary key,
    account_id int references user.account,
    name text,
    endpoint text
);

CREATE TABLE IF NOT EXISTS user.api(
    id serial primary key,
    account_id int references user.account,
    key text,
    secret text,
    salt text,
    deprecated bool
);

END;
$$;
