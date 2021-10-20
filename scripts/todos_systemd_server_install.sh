#!/bin/bash
## todos_systemd_server_install.sh
## Version : 1.0.0
## Example script to allow deployment of your Go Server as a systemd unit (service) in recent Linux OS
## must be used with the files : 'todos.conf' (for env variables) and 'todos.service' for the systemd unit service
## this script is provided as is, you must adapt it to your needs
# defining a name for your systemd unit, must be unique and not used on your target server
GONAME="todos"
SERVER_PATH=/usr/local/bin/${GONAME}Server
SERVICE_ENV_CONFIG="${GONAME}.conf"
#check if service definition was given as  first argument
SERVICE_DEFINITION="${GONAME}.service"
echo "NUM ARGS : " $#
if [ $# -eq 1 ]
then
  SERVICE_DEFINITION=$1
fi

# check if the service definition file exist
if [ ! -f "${SERVICE_DEFINITION}" ]
then
    echo "The '${GONAME}' Service definition file : '${SERVICE_DEFINITION}'  was not found."
    echo "Please pass the path to your service definition as first argument "
    echo "or adapt this scripts accordingly, BEFORE running this script."
    exit 1
fi
# check if the service configuration file exist
if [ ! -f "${SERVICE_ENV_CONFIG}" ]
then
    echo "The '${GONAME}' Service configuration file : '${SERVICE_ENV_CONFIG}'  was not found."
    echo "place it on the same dir or adapt this scripts accordingly, BEFORE running this script."
    exit 1
fi
if [ ! -f "${SERVER_PATH}" ]
then
    echo "The ${GONAME} Server binary is not present on the usual path : '${SERVER_PATH}'."
    echo "Please copy your latest service binary file in : ${SERVER_PATH}"
    echo "or adapt this scripts accordingly. BEFORE running this script."
    exit 1
fi
echo "creating group $GONAME"
groupadd --system $GONAME
# let's check that the group was created
grep $GONAME /etc/group
echo "creating user $GONAME and adding it to group $GONAME"
useradd -M -r -s /sbin/nologin $GONAME -g $GONAME
#let's check that the user was created by displaying it with the id command
id $GONAME
chown $GONAME:$GONAME ${SERVER_PATH}
echo "Your server script is now owned, by the ${GONAME} user and group :"
lsa ${SERVER_PATH}
echo "Will try to copy your service definition : ${SERVICE_DEFINITION} to /etc/systemd/system/"
mkdir -p /etc/systemd/system/${GONAME}.service.d/
cp "${SERVICE_DEFINITION}" /etc/systemd/system/
# copy the env definition for this service
mkdir -p /etc/systemd/system/${GONAME}.service.d/
cp "${SERVICE_ENV_CONFIG}" /etc/systemd/system/${GONAME}.service.d/
mkdir /var/lib/$GONAME
chown -R $GONAME:$GONAME /var/lib/$GONAME
echo "you can edit and check with : vim /etc/systemd/system/${GONAME}.service"
mkdir /var/log/$GONAME
chown -R $GONAME:$GONAME /var/log/$GONAME
chmod -R 775 /var/log/$GONAME
chmod -R 775 /var/lib/$GONAME
systemctl status "${GONAME}.service"
systemctl enable "${GONAME}.service"
systemctl start "${GONAME}.service"
systemctl status "${GONAME}.service"
echo "Your ${GONAME} Systemd Unit was deployed and enabled as a service"
echo "To check (and follow) the logs you can as usual use : journalctl -u ${GONAME}.service -f"
echo "And you can (stop|start|restart  as usual use : systemctl stop ${GONAME}.service"

