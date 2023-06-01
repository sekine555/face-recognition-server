#!/usr/bin/env bash

mysql -u testUser -pPassword123 sun < "/docker-entrypoint-initdb.d/1_create_database.sql"
mysql -u testUser -pPassword123 sun < "/docker-entrypoint-initdb.d/2_ddl.sql"
mysql -u testUser -pPassword123 sun < "/docker-entrypoint-initdb.d/3_dml.sql"
