#!/usr/bin/env bash

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

if [[ -f /.dockerenv ]]; then
  echo "[!] Waiting for database to start up"
  sleep 2
fi

CREATE_USER="CREATE USER $PG_USER WITH PASSWORD '$PG_PASSWORD' CREATEDB"

if [[ $PG_PASSWORD -eq default ]]; then
  TOMORROW=$(date -v +1d +'%Y-%m-%d')
  CREATE_USER="$CREATE_USER VALID UNTIL '$TOMORROW'"
fi

# h/t https://stackoverflow.com/a/4774063/8738498
SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"

echo "[!] creating user '$PG_USER'"
psql -h $PG_HOST -p 5432 -U postgres -d postgres -c "$CREATE_USER;"

echo "[!] creating database '$PG_DATABASE'"
psql -h $PG_HOST -p 5432 -U postgres -d postgres -c "CREATE DATABASE $PG_DATABASE WITH OWNER $PG_USER;"

echo "[!] creating schema"
psql -h $PG_HOST -p 5432 -U $PG_USER -d $PG_DATABASE < $SCRIPTPATH/schema.sql

echo "[!] creating tables"
psql -h $PG_HOST -p 5432 -U $PG_USER -d $PG_DATABASE < $SCRIPTPATH/tables.sql

echo "[!] removing public access to database"
psql -h $PG_HOST -p 5432 -U $PG_USER -d $PG_DATABASE < $SCRIPTPATH/security.sql
