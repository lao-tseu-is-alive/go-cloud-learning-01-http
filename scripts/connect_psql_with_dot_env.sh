#!/bin/bash
if test -f ".env"; then
  export $(cat .env | sed 's/#.*//g' | xargs)
  eval $(egrep -v '^#' .env | xargs)  psql "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME"
else
  echo "Your env file .env was not found !"
fi
