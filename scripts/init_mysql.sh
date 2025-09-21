#!/bin/bash

mysql -h mysql -u root -p"${DB_PASSWORD}" < scripts/create_db.sql