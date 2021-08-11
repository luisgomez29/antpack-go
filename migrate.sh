#!/bin/sh
# Database URL
DSN=postgres://postgres:$DBPWD@localhost:5432/antpack-go?sslmode=disable

echo "Running migrations"
c:/go-migrate/migrate.exe -source file://migrations/versions -database $DSN $1 $2