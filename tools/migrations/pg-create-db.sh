#!/usr/bin/env bash

echo "\
  +-------------------------------------------------+
  | pg-create-db: Create Postgres Database Script   |
  |                                                 |
  | Author: Saggaf Arsyad <saggaf.arsyad@gmail.com> |
  +-------------------------------------------------+
"

main() {
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
}

prepare_query() {
    CREATE_DB_QUERY="CREATE DATABASE ${DB_NAME} ENCODING = 'UTF8' TABLESPACE = pg_default CONNECTION LIMIT = -1;"
}

# ---------
# Functions
# ---------

get_env_with_prefix() {
    local ENV_KEY="${ENV_KEY_PREFIX}_${1}"
    echo ${!ENV_KEY}
}

load_env() {
    # Check if env has prefix
    if [[ -n ${ENV_KEY_PREFIX} ]]; then
        echo "  > ENV_KEY_PREFIX is set: ${ENV_KEY_PREFIX}"
        DB_DRIVER=$(get_env_with_prefix "DB_DRIVER")
        DB_MIGRATION_HOST=$(get_env_with_prefix "DB_MIGRATION_HOST")
        DB_MIGRATION_PORT=$(get_env_with_prefix "DB_MIGRATION_PORT")
        DB_NAME=$(get_env_with_prefix "DB_NAME")
        DB_USER=$(get_env_with_prefix "DB_USER")
        DB_PASS=$(get_env_with_prefix "DB_PASS")
    fi
}

init_usql_dsn() {
    # Normalize driver name
    if [[ ${DB_DRIVER} == "postgresql" ]]; then
        DB_DRIVER=postgres
    fi

    # Init DSN
    DSN="${DB_DRIVER}://${DB_USER}:${DB_PASS}@${DB_MIGRATION_HOST}:${DB_MIGRATION_PORT}/${DB_NAME}"

    # If SSL Mode set to false, then set option
    if [[ -z $SSL_MODE || $SSL_MODE == "false" ]]; then
        echo "DEBUG: Non SSL Mode"
        DSN+="?sslmode=disable"
    fi
}

# ----------
# Entrypoint
# ----------
load_env
init_usql_dsn
prepare_query
main
