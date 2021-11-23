#!/bin/sh

# TODO: Support MySQL database

# TODO: Allow overriding from .env
export DB_CONTAINER_IMAGE=postgres:11-alpine

docker run -d --name ${DB_CONTAINER_NAME} \
  -p ${DB_PORT}:5432 \
  -e "POSTGRES_USER=$DB_USER" \
  -e "POSTGRES_PASSWORD=$DB_PASS" \
  -e "POSTGRES_DB=$DB_NAME" \
  $DB_CONTAINER_IMAGE
