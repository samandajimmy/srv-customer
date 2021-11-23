#!/bin/sh

export CONF_FILE=${1}

touch $CONF_FILE

echo "flyway.url=jdbc:${DB_DRIVER}://${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${DB_NAME}" >> $CONF_FILE
echo "flyway.user=$DB_USER" >> $CONF_FILE
echo "flyway.password=$DB_PASS" >> $CONF_FILE