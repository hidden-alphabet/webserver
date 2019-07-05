#!/usr/bin/env bash

case $ENV in
  production|*)
    BASE=/
  ;;
  development)
  BASE=$(git rev-parse --show-toplevel)
  ENVFILE=$BASE/.env

  if [[ -f $ENVFILE ]]; then
    echo "[!] Using configuration file: $ENVFILE"
    source $ENVFILE
    echo "[!] Waiting for database to start up"
    sleep 2
  fi

  if [[ -z $PG_HOST ]]; then
    echo "[!] PG_HOST must be set."
    exit -1
  fi

  if [[ -z $PG_DATABASE ]]; then
    echo "[!] PG_DATABASE must be set."
    exit -1
  fi

  if [[ -z $PG_USER ]]; then
    echo "[!] PG_USER must be set."
    exit -1
  fi

  if [[ -z $PG_PASSWORD ]]; then
    echo "[!] PG_PASSWORD must be set."
    exit -1
  fi

  if [[ ! -f $ENVFILE ]]; then
    cat <<-SH > $ENVFILE
# autogenerated
PG_HOST=$PG_HOST
PG_DATABASE=$PG_DATABASE
PG_USER=$PG_USER
PG_PASSWORD=$PG_PASSWORD
SH
  fi
  ;;
esac

CREATE_USER="CREATE USER $PG_USER WITH PASSWORD '$PG_PASSWORD' CREATEDB"
CREATE_DATABASE="CREATE DATABASE $PG_DATABASE WITH OWNER $PG_USER"

echo "[!] Waiting for postgres to start"
until psql -h $PG_HOST -p 5432 -U postgres -d postgres -c '\l' 2>/dev/null; do
  echo "[-] Waiting..."
  sleep 1
done

echo "[!] creating user '$PG_USER'"
psql -h $PG_HOST -p 5432 -U postgres -d postgres -c "$CREATE_USER;"

echo "[!] creating database '$PG_DATABASE'"
psql -h $PG_HOST -p 5432 -U postgres -d postgres -c "$CREATE_DATABASE;"

echo "[!] creating schema"
psql -h $PG_HOST -p 5432 -U $PG_USER -d $PG_DATABASE < $BASE/schema.sql

echo "[!] creating tables"
psql -h $PG_HOST -p 5432 -U $PG_USER -d $PG_DATABASE < $BASE/tables.sql

echo "[!] removing public access to database"
psql -h $PG_HOST -p 5432 -U $PG_USER -d $PG_DATABASE < $BASE/security.sql