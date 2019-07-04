DO $$
BEGIN

CREATE TABLE IF NOT EXISTS web.user(
    id serial primary key,
    name text not null,
    email text not null,
    hash bytea not null,
    salt bytea not null
);

CREATE TABLE IF NOT EXISTS web.meta(
  id serial primary key,
  user_id int references web.user,
  is_active boolean,
  email_confirmed boolean,
  email_confirmation_path text,
  created_at timestamp default now(),
  last_active_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS web.session(
  id serial primary key,
  user_id int references web.user,
  created_at timestamp default now(),
  active boolean,
  token bytea
);

CREATE TABLE IF NOT EXISTS web.resource(
    id serial primary key,
    user_id int references web.user,
    name text,
    endpoint text
);

CREATE TABLE IF NOT EXISTS web.api(
    id serial primary key,
    user_id int references web.user,
    key text,
    secret text,
    salt text,
    deprecated bool
);

END;
$$;
