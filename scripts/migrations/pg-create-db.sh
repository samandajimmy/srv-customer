#!/usr/bin/env bash

echo \
"> Create Postgres Database Script
> Author: Saggaf Arsyad <saggaf.arsyad@gmail.com>
-------------------------------------------------"

# Init DSN
DSN="postgres://${DB_USER}:${DB_PASS}@${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${DB_DEFAULT:=postgres}"
SRC_DOWN_DIR=$1

# If SSL Mode set to false, then set option
if [[ -z $SSL_MODE || $SSL_MODE == "false" ]]; then
  echo "DEBUG: Non SSL Mode"
  DSN+="?sslmode=disable"
fi

# Init query
CREATE_DB_QUERY=$(cat <<EOF
  CREATE DATABASE ${DB_NAME}
  ENCODING = 'UTF8'
  TABLESPACE = pg_default
  CONNECTION LIMIT = -1;
EOF
)

# Execute query
echo "INFO: Creating database..."
usql -t -c "${CREATE_DB_QUERY}" "${DSN}"

# If not success, return
RESULT=$?
if [[ ${RESULT} -ne 0 ]]; then
  echo "ERROR: error occurred while creating database. (error=${RESULT})"
  exit 1
fi

echo "INFO: Done"
exit 0
