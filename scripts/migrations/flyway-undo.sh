#!/usr/bin/env bash

echo \
"> Flyway undo script for Community Edition
> Author: Saggaf Arsyad <saggaf.arsyad@gmail.com>
-------------------------------------------------"

# Init DSN
DSN="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
SRC_DOWN_DIR=$1

# If SSL Mode set to false, then set option
if [[ -z $SSL_MODE || $SSL_MODE == "false" ]]; then
  echo "DEBUG: Non SSL Mode"
  DSN+="?sslmode=disable"
fi

# Init query
LATEST_VERSION_QUERY=$(cat <<EOF
  SELECT version
  FROM flyway_schema_history 
  ORDER BY installed_rank 
  DESC LIMIT 1;
EOF
)

# Get latest version from schema history
echo "INFO: Retrieving latest version in migration history..."
LATEST_VERSION=$(usql -t -c "${LATEST_VERSION_QUERY}" "${DSN}" | tail -2)

if [[ -z ${LATEST_VERSION} || ${LATEST_VERSION} == "version" ]]; then
  echo "ERROR: cannot undo. migration has not been started"
  exit 5
fi

# Trim whitespace
LATEST_VERSION=$(echo ${LATEST_VERSION} | xargs)
echo "INFO: Latest migration version: ${LATEST_VERSION}"

# Find undo file with prefix
UNDO_FILE=$(find ${SRC_DOWN_DIR} -type f -name U${LATEST_VERSION}__*)

if [[ -z ${UNDO_FILE} ]]; then
  echo "WARN: Undo file for this version is not available. (version=${LATEST_VERSION})"
else
  echo "DEBUG: Undo script file: ${UNDO_FILE}"

  # Execute file
  echo "INFO: Running undo script..."
  usql -f ${UNDO_FILE} ${DSN}

  # If not success, return
  RESULT=$?
  if [[ ${RESULT} -ne 0 ]]; then
    echo "ERROR: error occurred while executing undo file. (error=${RESULT})"
    exit 3
  fi
fi

# Delete migrate history for latest version
echo "INFO: Removing migration history..."

DELETE_HISTORY_PROCEDURE_QUERY=$(cat <<EOF
  DO \$\$
    DECLARE
        var_latest_rev int;
    BEGIN
        SELECT installed_rank into var_latest_rev from flyway_schema_history order by installed_rank desc limit 1;
        delete from flyway_schema_history where installed_rank = var_latest_rev;
  END \$\$;
EOF
)

usql -c "${DELETE_HISTORY_PROCEDURE_QUERY}" "${DSN}"

# If not success, return
RESULT=$?
if [[ ${RESULT} -ne 0 ]]; then
  echo "ERROR: error occurred while removing migration history. (error=${RESULT})"
  exit 4
fi

echo "INFO: Done"
exit 0
