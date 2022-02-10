#!/usr/bin/env bash

echo "\
  +--------------------------------------------------------------+
  | flyway-undo.sh: Undo workaround for Flyway Community Edition |
  | Only support PostgreSQL                                      |
  |                                                              |
  | Author: Saggaf Arsyad <saggaf.arsyad@gmail.com>              |
  +--------------------------------------------------------------+
"

# Load Arguments
SRC_DOWN_DIR=$1

main() {
    # Get latest version from schema history
    echo "INFO: Retrieving latest version in migration history..."
    LATEST_VERSION=$(usql -t -c "${LATEST_VERSION_QUERY}" "${DSN}" | tail -2 | xargs)

    if [[ $LATEST_VERSION == "Connected with driver"* ]]; then
        echo "ERROR: cannot undo. migration has not been started"
        exit 5
    fi

    echo "INFO: Latest migration version: ${LATEST_VERSION}"

    # Find undo file with prefix
    UNDO_FILE=$(find ${SRC_DOWN_DIR} -type f -name "U${LATEST_VERSION}__"*)

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

    usql -c "${DELETE_HISTORY_PROCEDURE_QUERY}" "${DSN}"

    # If not success, return
    RESULT=$?
    if [[ ${RESULT} -ne 0 ]]; then
        echo "ERROR: error occurred while removing migration history. (error=${RESULT})"
        exit 4
    fi

    echo "INFO: Done"
    exit 0
}

prepare_query() {
    LATEST_VERSION_QUERY="SELECT version FROM flyway_schema_history WHERE version IS NOT NULL ORDER BY installed_rank DESC LIMIT 1"

    DELETE_HISTORY_PROCEDURE_QUERY="\
      DO
      \$\$
        DECLARE
            var_latest_rev int;
        BEGIN
            SELECT installed_rank into var_latest_rev from flyway_schema_history where version is not null order by installed_rank desc limit 1;
            -- Remove repeatable script that is executed after migration history
            delete from flyway_schema_history where version is NULL AND installed_rank > var_latest_rev;
            -- Remove migration history
            delete from flyway_schema_history where installed_rank = var_latest_rev;
        END
      \$\$;
    "
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
        echo "DEBUG: ENV_KEY_PREFIX is set: ${ENV_KEY_PREFIX}"
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
    DSN="${DB_DRIVER}://${MIGRATION_DB_USER}:${MIGRATION_DB_PASS}@${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${DB_NAME}"

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
