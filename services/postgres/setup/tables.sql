DO $$
BEGIN

CREATE TABLE IF NOT EXISTS web.account(
    id serial primary key,
    name text not null,
    hash bytea not null,
    salt bytea not null
);

CREATE TABLE IF NOT EXISTS web.contact(
  id serial primary key,
  account_id int references web.account,
  email text not null,
  has_confirmed_email boolean
);

CREATE TABLE IF NOT EXISTS web.session(
  id serial primary key,
  account_id int references web.account,
  token bytea,
  active boolean,
  created_at timestamp default now(),
  completed_at timestamp
);

CREATE TABLE IF NOT EXISTS web.api(
    id serial primary key,
    account_id int references web.account,
    key text,
    secret text,
    salt text,
    deprecated bool
);

END;
$$;
