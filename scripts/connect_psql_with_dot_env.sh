#!/bin/bash
eval $(egrep -v '^#' .env | xargs)  psql postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSL_MODE
