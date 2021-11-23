#!/usr/bin/env bash

# TODO: Support MySQL database

export FOLDER=${1}
export FILE_NAME=${2}

# copy file to tmp
docker cp ${PWD}/${FOLDER}/${FILE_NAME} ${DB_CONTAINER_NAME}:/tmp/${FILE_NAME}

# exec file
docker exec -i ${DB_CONTAINER_NAME} psql \
        -U ${DB_USER} \
        -w ${DB_PASS} \
        -d ${DB_NAME} \
        -f /tmp/${FILE_NAME}
