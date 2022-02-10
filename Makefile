# --------
# Manifest
# --------
PROJECT_NAME := "PDS Service"
PROJECT_PKG := repo.pegadaian.co.id/ms-pds/srv-customer
DOCKER_NAMESPACE := artifactory.pegadaian.co.id:5443

# -----------------
# Project Variables
# -----------------
PROJECT_ROOT ?= $(shell pwd)
PROJECT_WORKDIR ?= ${PROJECT_ROOT}
PROJECT_RESPONSES := responses.yml
PROJECT_CONFIG := .env
PROJECT_CONFIG_RELEASE := .env
PROJECT_WEB_TEMPLATES = web/templates
PROJECT_WEB_STATIC = web/static
PROJECT_DOCKERFILE_DIR ?= ${PROJECT_ROOT}/deployments/pds-svc
OUTPUT_DIR := ${PROJECT_ROOT}/bin
DOCTOR_CMD := ${PROJECT_ROOT}/scripts/doctor.sh

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

# Project Directories, configurable value using env
PROJECT_CONFIG?=$(PROJECT_WORKDIR)/.env
DEPLOYMENT_DIR=$(PROJECT_ROOT)/deployments
SCRIPTS_DIR=$(PROJECT_ROOT)/tools

# -------------------
# Migration Variables
# -------------------
MIGRATION_TOOL_CMD:=flyway
MIGRATION_TOOL_CONF=flyway.conf

MIGRATION_DIR=$(PROJECT_ROOT)/migrations
MIGRATION_SRC_UP?=$(MIGRATION_DIR)/sql-up
MIGRATION_SRC_DOWN?=$(MIGRATION_DIR)/sql-down
MIGRATION_CONFIG=$(MIGRATION_DIR)/$(MIGRATION_TOOL_CONF)

MIGRATION_SCRIPTS_DIR?=$(SCRIPTS_DIR)/migrations
MIGRATION_DOWN_CMD:=$(MIGRATION_SCRIPTS_DIR)/flyway-undo.sh
MIGRATION_INIT_CONFIG_CMD:=$(MIGRATION_SCRIPTS_DIR)/flyway-init-config.sh
MIGRATION_CREATE_DB:=$(MIGRATION_SCRIPTS_DIR)/pg-create-db.sh

# -----------
# API Version
# -----------
CI_COMMIT_TAG?=$$(git describe --tags $$(git rev-list --tags --max-count=1))
CI_COMMIT_SHA?=$$(git rev-parse HEAD)

# ------------------------------------------------------------------------------------------------------------------ #
# NEW CONFIG                                                                                                         #
# ------------------------------------------------------------------------------------------------------------------ #

# Working directory
PROJECT_ROOT ?= $(shell pwd)
PROJECT_WORKDIR ?= $(PROJECT_ROOT)
TEMP_DIR ?= $(PROJECT_WORKDIR)/.tmp

# Manifest
PROJECT_NAME := "PDS Customer Service"
PROJECT_SLUG ?= pds
PROJECT_PKG := repo.pegadaian.co.id/ms-pds/srv-customer
PROJECT_MAIN_DIR ?= $(PROJECT_ROOT)/cmd/customer
PROJECT_RELEASE_BIN ?= $(TEMP_DIR)/release/customer-svc
PROJECT_RELEASE_VERSION ?= $$(git describe --tags $$(git rev-list --tags --max-count=1))
PROJECT_RELEASE_BUILD_SIGNATURE ?= $$(git rev-parse HEAD)
DOCKER_NAMESPACE := artifactory.pegadaian.co.id:5443

# Docker variables
DOCKER_IMAGE_TAG ?= $(PROJECT_RELEASE_VERSION)
DOCKER_IMAGE_NAME ?= $(DOCKER_NAMESPACE)
DOCKER_IMAGE ?= $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_RELEASE_FILE := $(PROJECT_ROOT)/build/svc/Dockerfile
DOCKER_ARG_BUILD_VERSION ?= $(PROJECT_RELEASE_VERSION)
DOCKER_ARG_BUILD_SIGNATURE ?= $(PROJECT_RELEASE_BUILD_SIGNATURE)

# Project Directories, configurable value using env
PROJECT_CONFIG?=$(PROJECT_WORKDIR)/.env
DEPLOYMENT_DIR=$(PROJECT_ROOT)/deployments
PROTO_SRC_DIR=$(PROJECT_ROOT)/api/proto
PROTO_OUT_DIR=$(PROJECT_ROOT)/internal/dto
SCRIPTS_DIR=$(PROJECT_ROOT)/tools

# Commands
DOCTOR_CMD:=$(SCRIPTS_DIR)/doctor.sh

# Temp Variables
__LOCAL__DOCKER_COMPOSE=docker-compose.local.yml
__LOCAL__CMD_DOCKER_COMPOSE:=docker-compose -f $(__LOCAL__DOCKER_COMPOSE) --env-file $(PROJECT_CONFIG)

# Migration Variables
MIGRATION_TOOL_CMD:=flyway
MIGRATION_TOOL_CONF=flyway.conf

MIGRATION_DIR=$(PROJECT_ROOT)/migrations
MIGRATION_SRC_UP?=$(MIGRATION_DIR)/sql-up
MIGRATION_SRC_DOWN?=$(MIGRATION_DIR)/sql-down
MIGRATION_CONFIG=$(MIGRATION_DIR)/$(MIGRATION_TOOL_CONF)

MIGRATION_SCRIPTS_DIR?=$(SCRIPTS_DIR)/migrations
MIGRATION_DOWN_CMD:=$(MIGRATION_SCRIPTS_DIR)/flyway-undo.sh
MIGRATION_INIT_CONFIG_CMD:=$(MIGRATION_SCRIPTS_DIR)/flyway-init-config.sh
MIGRATION_CREATE_DB:=$(MIGRATION_SCRIPTS_DIR)/pg-create-db.sh

# ----
# Init
# ----
-include ${PROJECT_CONFIG}
export

# ---------------
# Common Commands
# ---------------

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

## configure: Configure project
.PHONY: configure
configure: --permit-exec --copy-env vendor
	@-echo "  > Configure: Done"

# --------------------------
# Local Development Commands
# --------------------------

## servers: Create servers for local development
.PHONY: servers
servers:
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) build
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) up -d db-pg13

## serve: Run server in development mode
.PHONY: serve
serve: --dev-build ${DEBUG_ENV_FILES}
	@-echo "  > Starting Server...\n"
	@LOG_LEVEL=debug;LOG_FORMAT=console; ${DEBUG_BIN} -dir=${PROJECT_ROOT}

## serve: Serve locally with live reloading and attach container stdout. For development purpose only
.PHONY: serve-dbg
serve-dbg:
	@-echo "  > Press CTRL+C to end"
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) up ${CTX_SLUG}-svc

## serve-rebuild: Re-build and serve with live reloading. For development purpose only
.PHONY: serve-rebuild
serve-rebuild:
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) up --build

## serve-log: Print last 100 lines of log from local server
.PHONY: serve-log
serve-log:
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) logs --tail 100

## serve-ps: See local docker process
.PHONY: serve-ps
serve-ps:
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) ps

## stop-serve: Stop local server
.PHONY: stop-serve
stop-serve:
	@-echo "  > Stopping..."
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) stop

## clean-vendor: Delete vendor
.PHONY: clean-vendor
clean-vendor:
	@-echo "  > Delete vendor..."
	@-rm -rf vendor

## clean-local: Clean local
.PHONY: clean-local
clean-local: clean-vendor
	@-echo "  > Deleting docker-compose for development..."
	@-$(__LOCAL__CMD_DOCKER_COMPOSE) down
	@-echo "  > Delete .tmp files"
	@-rm -rf .tmp

## vendor: Download dependencies to vendor
.PHONY: vendor
vendor:
	@-echo "  > Vendoring..."
	@go mod vendor
	@-echo "  > Vendoring: Done"

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

# ------------------
# Database Migration
# ------------------

## db: Create Database
.PHONY: db
db: db-config
	@$(MIGRATION_CREATE_DB)

## db-config: Generate a configuration file for database migration tool
.PHONY: db-config
db-config: $(MIGRATION_CONFIG)
$(MIGRATION_CONFIG): $(PROJECT_CONFIG) $(MIGRATION_INIT_CONFIG_SCRIPT)
	@-echo "  > Removing $(MIGRATION_TOOL_CONF)..."
	@-rm $(MIGRATION_CONFIG)
	@-echo "  > Creating $(MIGRATION_TOOL_CONF)..."
	@-$(MIGRATION_INIT_CONFIG_CMD) $(MIGRATION_CONFIG)

## db-status: Prints the details and status information about all the migrations.
.PHONY: db-status
db-status: db-config
	@$(MIGRATION_TOOL_CMD) info -configFiles=$(MIGRATION_CONFIG) -locations=filesystem:$(MIGRATION_SRC_UP)

## db-repair: Repair checksum
.PHONY: db-repair
db-repair: db-config
	@$(MIGRATION_TOOL_CMD) repair -configFiles=$(MIGRATION_CONFIG) -locations=filesystem:$(MIGRATION_SRC_UP)

## db-up: Upgrade database
.PHONY: db-up
db-up: db-config
	@-echo "  > Running up scripts..."
	@$(MIGRATION_TOOL_CMD) migrate -configFiles=$(MIGRATION_CONFIG) -locations=filesystem:$(MIGRATION_SRC_UP)

## db-down: (Experimental) undo to previous migration version
.PHONY: db-down
db-down: db-config
	$(MIGRATION_DOWN_CMD) $(MIGRATION_SRC_DOWN)

# -------------
# Private Rules
# -------------
.PHONY: --copy-env
--copy-env:
	@-echo "  > Copy .env (did not overwrite existing file)..."
	@-cp -n $(PROJECT_ROOT)/configs/.example.env $(PROJECT_CONFIG)

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

.PHONY: --permit-exec
--permit-exec: $(shell find $(SCRIPTS_DIR) -type f -name "*.sh")
	@-echo "  > Set executable permission to scripts..."
	@-chmod +x $(SCRIPTS_DIR)/**/*.sh
	@-chmod +x $(SCRIPTS_DIR)/*.sh
