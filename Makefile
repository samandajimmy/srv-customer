# --------
# Manifest
# --------
PROJECT_NAME:="PDS Service"
PROJECT_PKG:=repo.pegadaian.co.id/ms-pds/srv-customer
DOCKER_NAMESPACE:=artifactory.pegadaian.co.id:5443

# ---------------
# Command Aliases
# ---------------
GO_CMD:=go
GO_BUILD:=${GO_CMD} build
GO_MOD:=${GO_CMD} mod
GO_CLEAN:=${GO_CMD} clean
GO_GET:=${GO_CMD} get
DOCKER_CMD:=docker

# -----------------
# Project Variables
# -----------------
PROJECT_ROOT?=$(shell pwd)
PROJECT_WORKDIR?=${PROJECT_ROOT}
PROJECT_RESPONSES:=responses.yml
PROJECT_CONFIG:=.env
PROJECT_CONFIG_RELEASE:=.env
PROJECT_WEB_TEMPLATES=web/templates
PROJECT_WEB_STATIC=web/static
PROJECT_DOCKERFILE_DIR?=${PROJECT_ROOT}/deployments/customer-svc
OUTPUT_DIR:=${PROJECT_ROOT}/bin
DOCTOR_CMD:=${PROJECT_ROOT}/scripts/doctor.sh

# ---
# API
# ---
BINARY_NAME:=customer-svc
PROJECT_MAIN_PKG=cmd/${BINARY_NAME}
PROJECT_ENV_FILES:=$(addprefix ${PROJECT_ROOT}/,${PROJECT_CONFIG} ${PROJECT_RESPONSES})
PROJECT_ENV_FILES_RELEASE:=$(addprefix ${PROJECT_ROOT}/,${PROJECT_CONFIG_RELEASE} ${PROJECT_RESPONSES})

# ----------------------
# Debug Output Variables
# ----------------------
DEBUG_DIR:=${OUTPUT_DIR}/debug
DEBUG_BIN:=${DEBUG_DIR}/${BINARY_NAME}
DEBUG_ENV_FILES:=$(addprefix ${DEBUG_DIR}/,${PROJECT_CONFIG} ${PROJECT_RESPONSES})

# ------------------------
# Release Output Variables
# ------------------------
RELEASE_OUTPUT_DIR:=${OUTPUT_DIR}/release
RELEASE_ENV_APP_ENV?=1
RELEASE_ENV_LOG_LEVEL?=error
RELEASE_ENV_LOG_FORMAT?=console

# ----------------
# Docker Variables
# ----------------
CI_PROJECT_PATH ?= srv-customer
CI_COMMIT_REF_SLUG ?= local

IMAGE_APP ?= $(DOCKER_NAMESPACE)/$(CI_PROJECT_PATH)
IMAGE_APP_TAG ?= $(CI_COMMIT_REF_SLUG)

# -------------------
# Migration Variables
# -------------------
MIGRATION_TOOL_CMD:=flyway
MIGRATION_TOOL_CONF=flyway.conf

MIGRATION_DIR=${PROJECT_ROOT}/migrations
MIGRATION_SRC_UP?=${MIGRATION_DIR}/sql-up
MIGRATION_SRC_DOWN?=${MIGRATION_DIR}/sql-down
MIGRATION_CONFIG=${MIGRATION_DIR}/${MIGRATION_TOOL_CONF}

MIGRATION_SCRIPTS_DIR?=${PROJECT_ROOT}/scripts/migrations
MIGRATION_DOWN_CMD:=${MIGRATION_SCRIPTS_DIR}/flyway-undo.sh
MIGRATION_INIT_SERVER_CMD:=${MIGRATION_SCRIPTS_DIR}/init-server.sh
MIGRATION_INIT_CONFIG_CMD:=${MIGRATION_SCRIPTS_DIR}/init-config.sh
MIGRATION_CREATE_DB:=${MIGRATION_SCRIPTS_DIR}/pg-create-db.sh

# -----------
# API Version
# -----------
CI_COMMIT_TAG?=$$(git describe --tags $$(git rev-list --tags --max-count=1))
CI_COMMIT_SHA?=$$(git rev-parse HEAD)

# --------
# Commands
# --------

# Initialize CLI environment
-include ${PROJECT_CONFIG}
export

# ------------
# Common Rules
# ------------

## help: Show command help
.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "${PROJECT_NAME}":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## clean: Clean everything
.PHONY: clean
clean:
	@-echo "  > Deleting output dir..."
	@-rm -rf ${OUTPUT_DIR}
	@-echo "  > Done"

## doctor: Check for prerequisites
.PHONY: doctor
doctor: $(DOCTOR_CMD)
	@-echo "  > Checking dependencies..."
	@-${DOCTOR_CMD}

# ---------
# API Rules
# ---------

## setup: Make env from env example and grant permission.
.PHONY: setup
setup:
	@-echo "  > Creating env file..."
	@cp configs/.env-example .env
	@-echo "  > Fix scripts permission..."
	@chmod +x scripts/**/*.sh
	@chmod +x scripts/*.sh

## configure: Download dependencies
.PHONY: configure
configure: go.mod
	@-echo "  > Downloading dependencies..."
	@${GO_MOD} download
	@-echo "  > Done"

## serve: Run server in development mode
.PHONY: serve
serve: --dev-build ${DEBUG_ENV_FILES}
	@-echo "  > Starting Server...\n"
	@LOG_LEVEL=debug;LOG_FORMAT=console; ${DEBUG_BIN} -dir=${PROJECT_ROOT} -load-env-file

## vendor: Download dependencies to vendor folder
vendor: go.mod
	@-echo "  > Vendoring..."
	@${GO_MOD} vendor

## release: Compile binary for deployment.
.PHONY: release
release: vendor
	@-echo "  > Compiling for release..."
	@-echo "  >   Version: ${CI_COMMIT_TAG}"
	@-echo "  >   CommitHash: ${CI_COMMIT_SHA}"
	@CGO_ENABLED=0 GOOS=linux ${GO_BUILD} -a -v -mod=vendor \
		-ldflags "-X main.AppVersion=${CI_COMMIT_TAG} -X main.BuildHash=${CI_COMMIT_SHA}" \
		-o ${RELEASE_OUTPUT_DIR}/${BINARY_NAME} ${PROJECT_ROOT}/${PROJECT_MAIN_PKG}
	@-echo "  > Copying error codes..."
	@cp ${PROJECT_RESPONSES} ${RELEASE_OUTPUT_DIR}/
	@-echo "  > Output: $(RELEASE_OUTPUT_DIR)"

## image: Build a docker image from release
.PHONY: image
image:
	@-echo "  > Building image ${IMAGE_APP}:${IMAGE_APP_TAG}..."
	${DOCKER_CMD} build -t ${IMAGE_APP}:$(IMAGE_APP_TAG) \
		--build-arg ARG_PORT=${PORT} \
	    --progress plain -f ${PROJECT_DOCKERFILE_DIR}/Dockerfile .

## image-push: Push app image
.PHONY: image-push
image-push: image
	@-echo "  > Push image ${IMAGE_APP}:${IMAGE_APP_TAG} to Container Registry..."
	@${DOCKER_CMD} push ${IMAGE_APP}:${IMAGE_APP_TAG}

# ---------------
# Migration Rules
# ---------------

## db: Create Database
.PHONY: db
db: db-configure
	@${MIGRATION_CREATE_DB}

## db-configure: Generate a configuration for database migration tool
.PHONY: db-configure
db-configure: ${MIGRATION_CONFIG}
${MIGRATION_CONFIG}: $(PROJECT_CONFIG) $(MIGRATION_INIT_CONFIG_SCRIPT)
	@-echo "  > Removing ${MIGRATION_TOOL_CONF}..."
	@-rm ${MIGRATION_CONFIG}
	@-echo "  > Creating ${MIGRATION_TOOL_CONF}..."
	@-${MIGRATION_INIT_CONFIG_CMD} ${MIGRATION_CONFIG}

## db-status: Prints the details and status information about all the migrations.
.PHONY: db-status
db-status: db-configure
	@${MIGRATION_TOOL_CMD} info -configFiles=${MIGRATION_CONFIG} -locations=filesystem:${MIGRATION_SRC_UP}

## db-repair: Repair checksum
.PHONY: db-repair
db-repair: db-configure
	@${MIGRATION_TOOL_CMD} repair -configFiles=${MIGRATION_CONFIG} -locations=filesystem:${MIGRATION_SRC_UP}

## db-up: Upgrade database
.PHONY: db-up
db-up: db-configure
	@-echo "  > Running up scripts..."
	@${MIGRATION_TOOL_CMD} migrate -configFiles=${MIGRATION_CONFIG} -locations=filesystem:${MIGRATION_SRC_UP}

## db-down: (Experimental) undo to previous migration version
.PHONY: db-down
db-down: db-configure
	${MIGRATION_DOWN_CMD} ${MIGRATION_SRC_DOWN}

## db-clean: Clean database
.PHONY: db-clean
db-clean: db-configure --clean-prompt
	@-echo "  > Cleaning database..."
	@${MIGRATION_TOOL_CMD} clean -configFiles=${MIGRATION_CONFIG} -locations=filesystem:${MIGRATION_SRC_UP}

# -------------
# Private Rules
# -------------

.PHONY: --clean-release
--clean-release:
	@-echo "  > Cleaning ${RELEASE_OUTPUT_DIR}..."
	@rm -rf ${RELEASE_OUTPUT_DIR}

.PHONY: --dev-build
--dev-build:
	@-echo "  > Compiling..."
	@${GO_BUILD} -o ${DEBUG_BIN} ${PROJECT_ROOT}/${PROJECT_MAIN_PKG}
	@-echo "  > Output: ${DEBUG_BIN}"

.PHONY: --clean-prompt
--clean-prompt:
	@echo -n "Are you sure want to clean all data in database? [y/N] " && read ans && [ $${ans:-N} = y ]

${DEBUG_ENV_FILES}: $(PROJECT_ENV_FILES)
	@-echo "  > Copying environment files..."
	@-cp -R ${PROJECT_ENV_FILES} ${DEBUG_DIR}
