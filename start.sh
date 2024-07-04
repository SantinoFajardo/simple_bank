#!/bin/sh

set -e
echo "run db migration"
source app.env
./migrate -path ./migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"