#!/bin/sh
set -e

if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
  echo "Usage: sh scripts/migrations.test.sh ..."
  exit 0
fi

if [ "$1" = "--new" ] || [ "$1" = "-n" ]; then
  if [ -n "$2" ]; then
    goose -dir migration create "$2" sql
    exit 0
  fi
  echo "ERROR: forgot to add name"
  exit 1
fi

case "$1" in
  --status|-s) action=status ;;
  --up|-u) action=up; action_arg="$2" ;;
  --down|-d) action=down; action_arg="$2" ;;
  *) echo "Try: sh scripts/migrations.test.sh --help"; exit 1 ;;
esac

: "${DATABASE_USER:?missing}"
: "${DATABASE_PASSWORD:?missing}"
: "${DATABASE_DBNAME:?missing}"
: "${DATABASE_HOST:?missing}"
: "${DATABASE_PORT:?missing}"

goose -dir migration postgres \
  "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_DBNAME?sslmode=disable" \
  "$action" "$action_arg"