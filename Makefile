#!/usr/bin/make

DC=docker-compose -f ./build/docker-compose.yml

.SILENT:
.ONESHELL:
.PHONY: go-run prod-up dev-up dev-build ps stop logs

go-run:
	go run ./cmd/customizer/

prod-up:
	${DC} up -d --build --no-cache mongo redis customizer

dev-up:
	${DC} up -d mongo redis

dev-build:
	${DC} up -d --build --no-cache mongo redis

ps:
	${DC} ps

stop:
	${DC} stop

logs:
	${DC} logs ${s}

.DEFAULT_GOAL := go-run