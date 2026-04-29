#!/bin/sh

if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo -n "Usage: sh scripts/migrations.sh [OPTIONS]
OPTIONS:
    --help,   -h         - show this text and exit
    --new,    -n NAME    - generate new migration file
                             with given name
    --status, -s         - get migrations status
    --up,     -u VERSION - apply all available migrations.
                             If a version is given,
                             migrate up to a specific version
    --down,   -d VERSION - roll back a single migration
                             from the current version.
                             If a version is given,
                             roll back migrations
                             to a specific version
"
    exit
elif [ "$1" = "--new" ] || [ "$1" = "-n" ]; then
    if [ -n "$2" ]; then
        goose -dir migration create $2 sql
        exit
    else
        echo "ERROR: forgot to add name"
        exit
    fi

elif [ "$1" = "--status" ] || [ "$1" = "-s" ]; then
    action=status
elif [ "$1" = "--up" ] || [ "$1" = "-u" ]; then
    action=up
    if [ -n "$2" ]; then
        action=up-to
        action_arg="$2"
    fi
elif [ "$1" = "--down" ] || [ "$1" = "-d" ]; then
    action=down
    if [ -n "$2" ]; then
        action=down-to
        action_arg="$2"
    fi
else
    echo "Try: sh scripts/migrations.sh --help"
    exit
fi

DATABASE_USER=$(grep -E "DATABASE_USER" .env | grep -Eo "[^=]+$")
DATABASE_PASSWORD=$(grep -E "DATABASE_PASSWORD" .env | grep -Eo "[^=]+$")
DATABASE_DBNAME=$(grep -E "DATABASE_DBNAME" .env | grep -Eo "[^=]+$")
DATABASE_HOST=$(grep -E "DATABASE_HOST" .env | grep -Eo "[^=]+$")
DATABASE_PORT=$(grep -E "DATABASE_PORT" .env | grep -Eo "[^=]+$")

goose -dir migration postgres "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_DBNAME?sslmode=disable" "$action" $action_arg