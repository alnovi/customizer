#!/usr/bin/make

DC_PROD=docker-compose -f ./docker-compose.yml
DC_DEV=docker-compose -f ./docker-compose.dev.yml

.SILENT:
.ONESHELL:
.PHONY: go-run up-prod build-prod up-dev build-dev ps stop exec logs

#################################################
# PROD ##########################################
#################################################

prod-up:
	${DC_PROD} up -d

prod-build:
	${DC_PROD} build --no-cache

prod-ps:
	${DC_PROD} ps

prod-stop:
	${DC_PROD} stop

#################################################
# DEV ###########################################
#################################################

dev-up:
	${DC_DEV} up -d

dev-build:
	${DC_DEV} build --no-cache

dev-ps:
	${DC_DEV} ps

dev-stop:
	${DC_DEV} stop

dev-go-run:
	${DC_DEV} exec golang go run ./cmd/customizer

dev-node-install:
	${DC_DEV} exec node yarn install

dev-node-build:
	${DC_DEV} exec node yarn build

.DEFAULT_GOAL := dev-up