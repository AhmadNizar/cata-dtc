#!/bin/bash

mysql -h 127.0.0.1 -u root -p"${DB_PASSWORD}" < scripts/create_db.sql