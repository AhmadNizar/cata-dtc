#!/bin/bash

read -p "Enter database user [root]: " DB_USER
DB_USER=${DB_USER:-"root"}

read -s -p "Enter database password: " DB_PASSWORD
echo

echo "Running migrations with user: $DB_USER"
migrate -path internal/infrastructure/db/mysql/migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp(127.0.0.1:3306)/pokemondb" up