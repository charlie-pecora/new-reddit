#!/usr/bin/env bash
podman run -d --name postgres \
	--replace \
	-p "5432:5432" \
	-e "POSTGRES_PASSWORD=test" \
	-e "POSTGRES_USER=test" \
	-e "POSTGRES_DB=test" \
	postgres:latest
