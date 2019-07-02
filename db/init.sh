#!/usr/bin/env bash

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

CREATE_USER="CREATE USER $PG_USER WITH PASSWORD '$PG_PASSWORD' CREATEDB"

if [[ $PG_PASSWORD -eq default ]]; then
  TOMORROW=$(date -v +1d +'%Y-%m-%d')
  CREATE_USER="$CREATE_USER VALID UNTIL '$TOMORROW'"
fi

echo "[!] creating user '$PG_USER'"
psql -h 0.0.0.0 -p 5432 -U postgres -d postgres -c "$CREATE_USER;" 

echo "[!] creating database '$PG_DATABASE'"
psql -h 0.0.0.0 -p 5432 -U postgres -d postgres -c "CREATE DATABASE $PG_DATABASE WITH OWNER $PG_USER;"

echo "[!] creating schema"
psql -h 0.0.0.0 -p 5432 -U $PG_USER -d $PG_DATABASE < schema.sql

echo "[!] creating tables"
psql -h 0.0.0.0 -p 5432 -U $PG_USER -d $PG_DATABASE < tables.sql

echo "[!] removing public access to database"
psql -h 0.0.0.0 -p 5432 -U $PG_USER -d $PG_DATABASE < security.sql
