#!/bin/bash
## createNewAppDBAndUser.sh
## version : 1.0.0
## script to create a postgresql user and database based on the given app name argument
echo "NUM ARGS : " $#
if [ $# -eq 1 ]
then
  APP_NAME=${1,,}
else
  echo "## You must give the name of the app as first argument"
  echo "## this app name will be used to create a local postgres user and db "
  echo "usage : sudo $0 myapp"
  exit 1
fi
cd /tmp || exit 1
# generate a random password of 32 chars with chars selected in alphanumeric and some special chars
#DB_PASSWORD=`tr -dc '_+=()A-Z-a-z-0-9' < /dev/urandom | fold -w32 | head -n1`
#in this case i prefer to generate it with openssl, no user will enter this password manually
DB_PASSWORD=$(openssl rand -base64 32)
echo "## Will try to create postgres user "
echo "## username       : ${APP_NAME}"
echo "## password       : ${DB_PASSWORD}"
CREATE_USER="psql -c \"CREATE USER ${APP_NAME} WITH PASSWORD '${DB_PASSWORD}';\""
echo "about to run : ${CREATE_USER}"
su -c "${CREATE_USER}"   postgres
echo "## Will try to create database ${APP_NAME} with owner=${APP_NAME}"
su -c "createdb -O ${APP_NAME} ${APP_NAME}" postgres
# uncomment next line to add postgis extension to the db
#su -c "psql -c 'CREATE EXTENSION postgis;' ${APP_NAME}" postgres
cd -
# https://www.freedesktop.org/software/systemd/man/systemd.service.html
echo "## Will prepare a systemd unit conf file in current directory: ${APP_NAME}.conf"
cat  > "${APP_NAME}".conf << EOS
[Service]
Environment="PORT=8080"
Environment="SERVERIP=localhost"
# a way to indicate which storage to use for now one of (memory|postgres)
Environment="DB_DRIVER=postgres"
Environment="DB_HOST=127.0.0.1"
Environment="DB_PORT=5433"
Environment="DB_NAME=${APP_NAME}"
Environment="DBUSER=${APP_NAME}"
Environment="DBPASS=${DB_PASSWORD}"
# in dev env it can be ok to disable SSL mode but in prod it is another story
# it depends on various factor. is your service (go) running in the same host as the db (localhost ?)
# if not, is the network between your server and your db trusted ?? read the doc and ask your security officer:
# https://www.postgresql.org/docs/11/libpq-ssl.html#LIBPQ-SSL-PROTECTION
Environment="DB_SSL_MODE=disable"
EOS



~
~
