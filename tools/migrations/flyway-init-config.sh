#!/usr/bin/env sh

echo "\
  +--------------------------------------------------------------------------+
  | flyway-init-config: Generate flyway.conf file from Environment Variables |
  | Only support PostgreSQL                                                  |
  |                                                                          |
  | Author: Saggaf Arsyad <saggaf.arsyad@gmail.com>                          |
  +--------------------------------------------------------------------------+
"

export CONF_FILE=$1

main() {
    touch "$CONF_FILE"
    echo "\
flyway.url=jdbc:${DB_DRIVER}://${MIGRATION_DB_HOST}:${MIGRATION_DB_PORT}/${DB_NAME}
flyway.user=$MIGRATION_DB_USER
flyway.password=$MIGRATION_DB_PASS
" >>${CONF_FILE}
    echo "INFO: Done"
    exit 0
}

# ---------
# Functions
# ---------

get_env_with_prefix() {
    local ENV_KEY="${ENV_KEY_PREFIX}_$1"
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

# ----------
# Entrypoint
# ----------
load_env
main
