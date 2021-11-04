#!/bin/bash
#let's clean up docker db for unit test
make db-docker-delete
make db-docker-init-data
set -o allexport; source .env; set +o allexport; go test  ./... -v
