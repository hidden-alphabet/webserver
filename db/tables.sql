DO $$
BEGIN

CREATE TABLE IF NOT EXISTS web.user(
    id serial primary key,
    name text not null,
    email text not null,
    hash bytea not null,
    salt bytea not null
);

CREATE TABLE IF NOT EXISTS web.activated(
  id serial primary key,
  user_id int references web.user,
  active boolean
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
